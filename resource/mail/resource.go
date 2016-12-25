package mail

import (
	"strings"
	"sync"
	"time"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
	"gopkg.in/gomail.v2"
)

const (
	mailDaemonTimeOut = 5 * time.Minute
)

type mailTask struct {
	message *gomail.Message
	result  chan error
}

type Resource struct {
	application *shadow.Application
	config      *config.Resource
	logger      logger.Logger

	mutex  sync.RWMutex
	open   bool
	dialer *gomail.Dialer
	closer gomail.SendCloser
	queue  chan *mailTask
}

func (r *Resource) GetName() string {
	return "mail"
}

func (r *Resource) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	r.config = resourceConfig.(*config.Resource)

	r.application = a

	return nil
}

func (r *Resource) Run(wg *sync.WaitGroup) error {
	r.logger = logger.NewOrNop(r.GetName(), r.application)

	r.open = false
	r.queue = make(chan *mailTask)
	r.initDialer(
		r.config.GetString(ConfigMailSmtpHost),
		r.config.GetInt(ConfigMailSmtpPort),
		r.config.GetString(ConfigMailSmtpUsername),
		r.config.GetString(ConfigMailSmtpPassword),
	)

	go func() {
		defer wg.Done()

		for {
			select {
			case task, ok := <-r.queue:
				if !ok {
					return
				}

				err := r.execute(task)
				if metricMailTotal != nil {
					if err != nil {
						metricMailTotal.With("result", "failed").Add(1)
					} else {
						metricMailTotal.With("result", "success").Add(1)
					}
				}

			case <-time.After(mailDaemonTimeOut):
				if r.open {
					if err := r.closer.Close(); err != nil && !strings.Contains(err.Error(), "4.4.2") {
						r.logger.Error("Dialer close failed", map[string]interface{}{"error": err.Error()})
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

func (r *Resource) initDialer(host string, port int, username, password string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.dialer = gomail.NewDialer(host, port, username, password)
	r.open = false
}

func (r *Resource) execute(task *mailTask) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var err error

	if !r.open {
		if r.closer, err = r.dialer.Dial(); err != nil {
			r.logger.Error("Dialer dial failed", map[string]interface{}{"error": err.Error()})
			task.result <- err

			return err
		} else {
			r.logger.Debug("Dialer open success")
			r.open = true
		}
	}

	if r.open {
		if len(task.message.GetHeader("From")) == 0 {
			task.message.SetAddressHeader("From", r.config.GetString(ConfigMailFromAddress), r.config.GetString(ConfigMailFromName))
		}

		if err = gomail.Send(r.closer, task.message); err != nil {
			if strings.Contains(err.Error(), "4.4.2") {
				r.logger.Debug("SMTP server response timeout exceeded", map[string]interface{}{
					"message": task.message,
					"error":   err.Error(),
				})

				r.open = false
				return r.execute(task)
			} else {
				r.logger.Error(err.Error(), map[string]interface{}{"message": task.message})
				task.result <- err

				return err
			}
		} else {
			r.logger.Debug("Send message success", map[string]interface{}{"message": task.message})
			task.result <- nil
		}
	}

	return nil
}

func (r *Resource) Send(message *gomail.Message) {
	task := &mailTask{
		message: message,
		result:  make(chan error),
	}
	r.queue <- task

	r.logger.Debug("Send new message to queue", map[string]interface{}{"message": message})

}

func (r *Resource) SendAndReturn(message *gomail.Message) error {
	task := &mailTask{
		message: message,
		result:  make(chan error),
	}
	r.queue <- task

	return <-task.result
}
