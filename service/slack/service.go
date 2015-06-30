package slack

import (
	"flag"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
	"github.com/nlopes/slack"
)

const (
	keepAlive = 20 * time.Second
	connectDelay = 5 * time.Second
)

type ServiceSlackCommands interface {
	GetSlackCommands() []SlackCommand
}

type SlackService struct {
	application *shadow.Application
	config      *resource.Config
	logger      *logrus.Entry

	Commands map[string]SlackCommand
	Rtm      *slack.SlackWS
	Bot      *slack.UserDetails

	api             *slack.Slack
	sender          chan slack.OutgoingMessage
	receiver        chan slack.SlackEvent
	connectAttempts int
}

func (s *SlackService) GetName() string {
	return "slack"
}

func (s *SlackService) Init(a *shadow.Application) error {
	s.application = a

	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	s.config = resourceConfig.(*resource.Config)

	resourceLogger, err := a.GetResource("logger")
	if err != nil {
		return err
	}
	s.logger = resourceLogger.(*resource.Logger).Get(s.GetName())

	s.Commands = map[string]SlackCommand{}

	return nil
}

func (s *SlackService) Run(wg *sync.WaitGroup) (err error) {
	s.api = slack.New(s.config.GetString("slack-token"))
	s.api.SetDebug(s.config.GetBool("debug"))

	s.sender = make(chan slack.OutgoingMessage)
	s.receiver = make(chan slack.SlackEvent)

	for _, service := range s.application.GetServices() {
		if serviceCast, ok := service.(ServiceSlackCommands); ok {
			for _, command := range serviceCast.GetSlackCommands() {
				if err := s.RegisterCommand(command, service.(shadow.Service)); err != nil {
					s.logger.Errorf("Error register slack command %s", command.GetName())
					// ignore error
				}
			}
		}
	}

	if err = s.connect(); err != nil {
		return err
	}

	return nil
}

func (s *SlackService) RegisterCommand(command SlackCommand, service shadow.Service) error {
	name := command.GetName()

	if _, ok := s.Commands[name]; ok {
		return errors.Newf("There are two command mapped to %s!", name)
	} else {
		command.Init(service, s.application)
		s.Commands[name] = command
	}

	return nil
}

func (s *SlackService) handleCommand(m *slack.MessageEvent) {
	// ignore self messages
	if m.UserId == s.Bot.Id {
		return
	}

	f := flag.NewFlagSet("slack", flag.ExitOnError)
	f.Parse(strings.Split(m.Text, " "))
	name := f.Arg(0)
	args := f.Args()[1:]

	botTag := "<@" + s.Bot.Id + ">"
	appeal := name == s.Bot.Name || strings.HasPrefix(name, s.Bot.Name+":") || strings.HasPrefix(name, botTag) || strings.HasPrefix(name, botTag+":")
	name = strings.ToLower(name)

	// Not direct message
	ch, _ := s.Rtm.GetChannelInfo(m.ChannelId)
	if ch != nil {
		if !appeal {
			if _, ok := s.Commands["hello"]; ok && strings.Contains(m.Text, botTag) {
				name = "hello"
			} else {
				return
			}
		}
	}

	// ignore bot name
	if appeal {
		name = strings.ToLower(f.Arg(1))

		if len(args) > 1 {
			args = args[1:]
		}
	}

	var (
		ok      bool
		command SlackCommand
	)

	if name != "" {
		command, ok = s.Commands[name]
		if !ok {
			if _, ok := s.Commands["unknown"]; ok {
				command = s.Commands["unknown"]
				s.logger.Warnf("unknown command: %s (%s)", name, m.Text)
			} else {
				return
			}
		}
	} else if _, ok := s.Commands["hello"]; ok {
		command = s.Commands["hello"]
	} else {
		return
	}

	// check permissions
	if ch != nil && !command.AllowChannel() {
		return
	}

	if ch == nil && !command.AllowDirectMessage() {
		return
	}

	command.Run(m, args...)
}

func (s *SlackService) connect() (err error) {
	s.connectAttempts = s.connectAttempts + 1

	s.Rtm, err = s.api.StartRTM("", "http://localhost/")

	go func(wsAPI *slack.SlackWS, ch chan slack.SlackEvent) {
		defer func() {
			if r := recover(); r != nil {
				// TODO: reconnect
				panic(r)
			}
		}()

		wsAPI.HandleIncomingEvents(ch)
	}(s.Rtm, s.receiver)

	go s.keepalive()

	go func(wsAPI *slack.SlackWS, chSender chan slack.OutgoingMessage) {
		for {
			select {
			case msg := <-chSender:
				wsAPI.SendMessage(&msg)
			}
		}
	}(s.Rtm, s.sender)

	s.Bot = s.Rtm.GetInfo().User
	s.connectAttempts = 0
	s.logger.Infof("Connect slack as %s", s.Bot.Name)

	go func() {
		for {
			select {
			case msg := <-s.receiver:
				switch msg.Data.(type) {
				case *slack.MessageEvent:
					s.handleCommand(msg.Data.(*slack.MessageEvent))
				default:
					if s.config.GetBool("debug") {
						s.logger.Warnf("Unexpected: %v", msg.Data)
					}
				}
			}
		}
	}()

	return nil
}

func (s *SlackService) keepalive() {
	ticker := time.NewTicker(keepAlive)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.Rtm.Ping(); err != nil {
				s.logger.Errorf("Ping error: %s ", err.Error())
			}
		}
	}
}
