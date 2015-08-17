package slack

import (
	"flag"
	"strings"
	"sync"

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
	mutex       sync.RWMutex

	Commands map[string]SlackCommand
	Rtm      *slack.RTM
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

func (s *SlackService) Run(wg *sync.WaitGroup) (err error) {
	token := s.config.GetString("slack.token")
	if token == "" {
		s.logger.Error("Slack token is empty")
		return nil
	}

	api := slack.New(s.config.GetString("slack.token"))
	api.SetDebug(s.config.GetBool("debug"))

	if s.Rtm = api.NewRTM(); err != nil {
		return err
	}

	for _, service := range s.application.GetServices() {
		if serviceCast, ok := service.(ServiceSlackCommands); ok {
			for _, command := range serviceCast.GetSlackCommands() {
				logEntry := s.logger.WithFields(logrus.Fields{
					"command": command.GetName(),
					"service": service.GetName(),
				})

				if !command.IsActive() {
					logEntry.Debug("Ignore disable command")
					continue
				}

				if err := s.RegisterCommand(command, service.(shadow.Service)); err != nil {
					logEntry.WithField("error", err.Error()).Error("Error register slack command")
					// ignore error
				} else {
					logEntry.Debug("Register command")
				}
			}
		}
	}

	go s.Rtm.ManageConnection()
	go s.process()

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

func (s *SlackService) process() {
	for {
		select {
		case msg := <-s.Rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				s.handleCommand(ev)

			case *slack.ConnectedEvent:
				s.mutex.Lock()
				s.logger.Info("Connected success")
				s.mutex.Unlock()

				s.Bot = s.Rtm.GetInfo().User

			case *slack.ConnectingEvent:
				// Ignore hello

			case *slack.HelloEvent:
				// Ignore hello

			case *slack.LatencyReport:
				// Ignore latency report

			case *slack.PresenceChangeEvent:
				// Ignore presence change

			case *slack.UserTypingEvent:
				// Ignore user typing

			case *slack.AckMessage:
				// Ignore ack message

			case *slack.DisconnectedEvent:
				s.mutex.Lock()
				s.logger.Info("Disconnected")
				s.mutex.Unlock()

			case *slack.InvalidAuthEvent:
				s.mutex.Lock()
				s.logger.Error("Invalid auth")
				s.mutex.Unlock()

			case *slack.AckErrorEvent:
				s.mutex.Lock()
				s.logger.Errorf("Ack error: %v", ev.Error())
				s.mutex.Unlock()

			case *slack.ConnectionErrorEvent:
				s.mutex.Lock()
				s.logger.Errorf("Connecting error: %v", ev.Error())
				s.mutex.Unlock()

			case *slack.IncomingEventError:
				s.mutex.Lock()
				s.logger.Errorf("Incomming error: %v", ev.Error())
				s.mutex.Unlock()

			case *slack.OutgoingErrorEvent:
				s.mutex.Lock()
				s.logger.Errorf("Outgoing error: %v", ev.Error())
				s.mutex.Unlock()

			case *slack.RTMError:
				s.mutex.Lock()
				s.logger.Errorf("RTM error: %v", ev.Error())
				s.mutex.Unlock()

			case *slack.SlackErrorEvent:
				s.mutex.Lock()
				s.logger.Errorf("Slack error: %v", ev.Error())
				s.mutex.Unlock()

			default:
				if s.config.GetBool("debug") {
					s.mutex.Lock()
					s.logger.Warnf("Unexpected: %v %v", msg.Type, msg.Data)
					s.mutex.Unlock()
				}
			}
		}
	}
}

func (s *SlackService) handleCommand(m *slack.MessageEvent) {
	// ignore self messages
	if m.User == s.Bot.ID {
		return
	}

	f := flag.NewFlagSet("slack", flag.ExitOnError)
	f.Parse(strings.Split(m.Text, " "))
	name := f.Arg(0)
	args := f.Args()[1:]

	botTag := "<@" + s.Bot.ID + ">"
	appeal := name == s.Bot.Name || strings.HasPrefix(name, s.Bot.Name+":") || strings.HasPrefix(name, botTag) || strings.HasPrefix(name, botTag+":")
	name = strings.ToLower(name)

	// Not direct message
	ch, _ := s.Rtm.GetChannelInfo(m.Channel)
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
		} else {
			args = []string{}
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
