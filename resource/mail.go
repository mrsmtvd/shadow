package resource

import (
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kihamo/shadow"
	"gopkg.in/gomail.v2"
)

const (
	mailDaemonTimeOut = 5 * time.Minute
)

type mailTask struct {
	message *gomail.Message
	result  chan error
}

type Mail struct {
	config *Config
	logger *logrus.Entry
	open   bool
	dialer *gomail.Dialer
	closer gomail.SendCloser
	queue  chan *mailTask
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
			Value: int64(25),
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
	r.open = false
	r.dialer = gomail.NewDialer(
		r.config.GetString("mail.smtp.host"),
		int(r.config.GetInt64("mail.smtp.port")),
		r.config.GetString("mail.smtp.username"),
		r.config.GetString("mail.smtp.password"),
	)
	r.queue = make(chan *mailTask)

	go func() {
		defer wg.Done()

		for {
			select {
			case task, ok := <-r.queue:
				if !ok {
					return
				}

				r.execute(task)

			case <-time.After(mailDaemonTimeOut):
				if r.open {
					if err := r.closer.Close(); err != nil && !strings.Contains(err.Error(), "4.4.2") {
						r.logger.WithField("error", err).Error("Dialer close failed", err.Error())
					} else {
						r.logger.Debug("Dialer close success")
					}

					r.open = false
				}
			}
		}
	}()

	return nil
}

func (r *Mail) execute(task *mailTask) {
	var err error

	if !r.open {
		if r.closer, err = r.dialer.Dial(); err != nil {
			r.logger.WithField("error", err).Error("Dialer dial failed", err.Error())
			task.result <- err
		} else {
			r.logger.Debug("Dialer open success")
			r.open = true
		}
	}

	if r.open {
		if len(task.message.GetHeader("From")) == 0 {
			task.message.SetHeader("From", r.config.GetString("mail.from"))
		}

		if err = gomail.Send(r.closer, task.message); err != nil {
			if strings.Contains(err.Error(), "4.4.2") {
				r.logger.WithFields(logrus.Fields{
					"message": task.message,
					"error":   err.Error(),
				}).Debug("SMTP server response timeout exceeded")

				r.open = false
				r.execute(task)
			} else {
				r.logger.WithField("message", task.message).Error(err.Error())
				task.result <- err
			}
		} else {
			r.logger.WithField("message", task.message).Debug("Send message success")
			task.result <- nil
		}
	}
}

func (r *Mail) Send(message *gomail.Message) {
	task := &mailTask{
		message: message,
		result:  make(chan error),
	}
	r.queue <- task

	r.logger.WithField("message", message).Debug("Send new message to queue")

}

func (r *Mail) SendAndReturn(message *gomail.Message) error {
	task := &mailTask{
		message: message,
		result:  make(chan error),
	}
	r.queue <- task

	return <-task.result
}
