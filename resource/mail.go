package resource

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kihamo/shadow"
	"gopkg.in/gomail.v2"
)

const (
	mailDaemonTimeOut = 5 * time.Minute
)

type Mail struct {
	config *Config
	logger *logrus.Entry
	queue  chan *gomail.Message
}

func (r *Mail) GetName() string {
	return "mail"
}

func (r *Mail) GetConfigVariables() []ConfigVariable {
	return []ConfigVariable{
		{
			Key:   "mail.smtp.username",
			Value: "",
			Usage: "SMTP username",
		},
		{
			Key:   "mail.smtp.password",
			Value: "",
			Usage: "SMTP password",
		},
		{
			Key:   "mail.smtp.host",
			Value: "",
			Usage: "SMTP host",
		},
		{
			Key:   "mail.smtp.port",
			Value: int64(465),
			Usage: "SMTP port",
		},
		{
			Key:   "mail.from",
			Value: "",
			Usage: "Mail from",
		},
	}
}

func (r *Mail) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	r.config = resourceConfig.(*Config)

	resourceLogger, err := a.GetResource("logger")
	if err != nil {
		return err
	}
	r.logger = resourceLogger.(*Logger).Get(r.GetName())

	return nil
}

func (r *Mail) Run(wg *sync.WaitGroup) error {
	r.queue = make(chan *gomail.Message)

	go func() {
		defer wg.Done()

		var (
			closer gomail.SendCloser
			err    error
		)

		open := false
		dialer := gomail.NewDialer(
			r.config.GetString("mail.smtp.host"),
			int(r.config.GetInt64("mail.smtp.port")),
			r.config.GetString("mail.smtp.username"),
			r.config.GetString("mail.smtp.password"),
		)

		for {
			select {
			case message, ok := <-r.queue:
				if !ok {
					return
				}

				if !open {
					if closer, err = dialer.Dial(); err != nil {
						r.logger.WithField("error", err).Panicf("Dialer dial failed", err.Error())
						open = false
					} else {
						r.logger.Debug("Dialer open success")
						open = true
					}
				}

				if len(message.GetHeader("From")) == 0 {
					message.SetHeader("From", r.config.GetString("mail.from"))
				}

				if err := gomail.Send(closer, message); err != nil {
					r.logger.WithField("message", message).Error(err.Error())
				} else {
					r.logger.WithField("message", message).Debug("Send message success")
				}

			case <-time.After(mailDaemonTimeOut):
				if open {
					if err := closer.Close(); err != nil {
						r.logger.WithField("error", err).Panicf("Dialer close failed", err.Error())
					} else {
						r.logger.Debug("Dialer close success")
					}
					open = false
				}
			}
		}
	}()

	return nil
}

func (r *Mail) Send(message *gomail.Message) {
	r.logger.WithField("message", message).Debug("Send new message to queue")
	r.queue <- message
}
