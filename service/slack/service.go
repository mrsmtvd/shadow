package slack

import (
	"flag"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
	"github.com/nlopes/slack"
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

func (s *SlackService) Run() error {
	api := slack.New(s.config.GetString("slack-token"))
	api.SetDebug(s.config.GetBool("debug"))

	sender := make(chan slack.OutgoingMessage)
	receiver := make(chan slack.SlackEvent)

	var err error
	s.Rtm, err = api.StartRTM("", "http://localhost/")

	for _, service := range s.application.GetServices() {
		if serviceCast, ok := service.(ServiceSlackCommands); ok {
			for _, command := range serviceCast.GetSlackCommands() {
				if err := s.RegisterCommand(command, service.(shadow.Service)); err != nil {
					s.logger.Errorf("Error register slack command " + command.GetName())
				}
			}
		}
	}

	if err != nil {
		return err
	}

	go s.Rtm.HandleIncomingEvents(receiver)
	go s.Rtm.Keepalive(20 * time.Second)

	go func(wsAPI *slack.SlackWS, chSender chan slack.OutgoingMessage) {
		for {
			select {
			case msg := <-chSender:
				wsAPI.SendMessage(&msg)
			}
		}
	}(s.Rtm, sender)

	s.Bot = s.Rtm.GetInfo().User
	s.logger.Infof("Connect slack as %s", s.Bot.Name)

	go func() {
		for {
			select {
			case msg := <-receiver:
				switch msg.Data.(type) {
				case *slack.MessageEvent:
					s.handleCommand(msg.Data.(*slack.MessageEvent))
				default:
					if s.config.GetBool("debug") {
						s.logger.Warnf("Unexpected: %v\n", msg.Data)
					}
				}
			}
		}
	}()

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
