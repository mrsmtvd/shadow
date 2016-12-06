package mail

import (
	"strings"
	"sync"
	"time"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
	"github.com/rs/xlog"
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
	config *config.Config
	logger xlog.Logger
	open   bool
	dialer *gomail.Dialer
	closer gomail.SendCloser
	queue  chan *mailTask
}

func (r *Mail) GetName() string {
	return "mail"
}

func (r *Mail) GetConfigVariables() []config.ConfigVariable {
	return []config.ConfigVariable{
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
			Value: 25,
			Usage: "SMTP port",
		},
		{
			Key:   "mail.from.address",
			Value: "",
			Usage: "Mail from address",
		},
		{
			Key:   "mail.from.name",
			Value: "",
			Usage: "Mail from name",
		},
	}
}

func (r *Mail) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	r.config = resourceConfig.(*config.Config)

	resourceLogger, err := a.GetResource("logger")
	if err != nil {
		return err
	}
	r.logger = resourceLogger.(*logger.Logger).Get(r.GetName())

	return nil
}

func (r *Mail) Run(wg *sync.WaitGroup) error {
	r.open = false
	r.dialer = gomail.NewDialer(
		r.config.GetString("mail.smtp.host"),
		r.config.GetInt("mail.smtp.port"),
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
						r.logger.Error("Dialer close failed", xlog.F{
							"error": err.Error(),
						})
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
			r.logger.Error("Dialer dial failed", xlog.F{"error": err.Error()})
			task.result <- err
		} else {
			r.logger.Debug("Dialer open success")
			r.open = true
		}
	}

	if r.open {
		if len(task.message.GetHeader("From")) == 0 {
			task.message.SetAddressHeader("From", r.config.GetString("mail.from.address"), r.config.GetString("mail.from.name"))
		}

		if err = gomail.Send(r.closer, task.message); err != nil {
			if strings.Contains(err.Error(), "4.4.2") {
				r.logger.Debug("SMTP server response timeout exceeded", xlog.F{
					"message": task.message,
					"error":   err.Error(),
				})

				r.open = false
				r.execute(task)
			} else {
				r.logger.Error(err.Error(), xlog.F{"message": task.message})
				task.result <- err
			}
		} else {
			r.logger.Debug("Send message success", xlog.F{"message": task.message})
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

	r.logger.Debug("Send new message to queue", xlog.F{"message": message})

}

func (r *Mail) SendAndReturn(message *gomail.Message) error {
	task := &mailTask{
		message: message,
		result:  make(chan error),
	}
	r.queue <- task

	return <-task.result
}
