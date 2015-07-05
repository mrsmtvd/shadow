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
	keepAlive    = 10 * time.Second
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

	mutex           sync.RWMutex
	connected       bool
	api             *slack.Slack
	senderChannel   chan slack.OutgoingMessage
	receiverChannel chan slack.SlackEvent
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
    token := s.config.GetString("slack-token")
    if token == "" {
        s.logger.Error("Slack token is empty")
        return nil
    }

	s.api = slack.New(s.config.GetString("slack-token"))
	s.api.SetDebug(s.config.GetBool("debug"))

	s.senderChannel = make(chan slack.OutgoingMessage)
	s.receiverChannel = make(chan slack.SlackEvent)

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

	go s.sender()
	go s.receiver()
	go s.connect()
	go s.keepalive()

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

func (s *SlackService) connect() {
	var err error

	defer func() {
		if r := recover(); r != nil {
			s.mutex.Lock()
			s.connected = false
			s.logger.Errorf("Panic: %v", r)
			s.logger.Debug("Connect closed. Sleep ", time.Duration(connectDelay).String())
			s.mutex.Unlock()

			time.Sleep(connectDelay)
			go s.connect()
		}
	}()

	s.Rtm, err = s.api.StartRTM("", "http://localhost/")
	if err != nil {
		panic(err)
	}

	s.Bot = s.Rtm.GetInfo().User

	s.mutex.Lock()
	s.connected = true
	s.logger.Infof("Connect slack as %s", s.Bot.Name)
	s.mutex.Unlock()

	s.Rtm.HandleIncomingEvents(s.receiverChannel)
}

func (s *SlackService) sender() {
	for {
		select {
		case msg := <-s.senderChannel:
			s.Rtm.SendMessage(&msg)
		}
	}
}

func (s *SlackService) receiver() {
	for {
		select {
		case msg := <-s.receiverChannel:
			switch msg.Data.(type) {
			case slack.HelloEvent:
			// Ignore hello
			case *slack.MessageEvent:
				s.handleCommand(msg.Data.(*slack.MessageEvent))
			default:
				if s.config.GetBool("debug") {
					// Ignore pong
					if _, ok := msg.Data.(slack.LatencyReport); !ok {
						s.mutex.Lock()
						s.logger.Warnf("Unexpected: %v", msg.Data)
						s.mutex.Unlock()
					}
				}
			}
		}
	}
}

func (s *SlackService) keepalive() {
	ticker := time.NewTicker(keepAlive)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mutex.RLock()
			if s.connected {
				s.mutex.RUnlock()

				if err := s.Rtm.Ping(); err != nil {
					s.mutex.Lock()
					s.logger.Errorf("Ping error: %s ", err.Error())
					s.mutex.Unlock()
				}
			} else {
				s.mutex.RUnlock()
			}
		}
	}
}
