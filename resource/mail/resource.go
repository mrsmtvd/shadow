package mail

import (
	"strings"
	"sync"
	"time"

	kit "github.com/go-kit/kit/metrics"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
	"github.com/kihamo/shadow/resource/metrics"
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

type Resource struct {
	config  *config.Resource
	metrics *metrics.Resource
	logger  xlog.Logger

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

	if a.HasResource("logger") {
		resourceLogger, _ := a.GetResource("logger")
		r.logger = resourceLogger.(*logger.Resource).Get(r.GetName())
	}

	if a.HasResource("metrics") {
		resourceMetrics, _ := a.GetResource("metrics")
		r.metrics = resourceMetrics.(*metrics.Resource)
	}

	return nil
}

func (r *Resource) Run(wg *sync.WaitGroup) error {
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

		var metricTotal kit.Counter
		if r.metrics != nil {
			metricTotal = r.metrics.NewCounter(MetricMailTotal)
		}

		for {
			select {
			case task, ok := <-r.queue:
				if !ok {
					return
				}

				err := r.execute(task)
				if metricTotal != nil {
					if err != nil {
						metricTotal.With("result", "failed").Add(1)
					} else {
						metricTotal.With("result", "success").Add(1)
					}
				}

			case <-time.After(mailDaemonTimeOut):
				if r.open {
					if err := r.closer.Close(); err != nil && !strings.Contains(err.Error(), "4.4.2") {
						r.logger.Error("Dialer close failed", xlog.F{"error": err.Error()})
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

func (r *Resource) execute(task *mailTask) error {
	var err error

	if !r.open {
		if r.closer, err = r.dialer.Dial(); err != nil {
			r.logger.Error("Dialer dial failed", xlog.F{"error": err.Error()})
			task.result <- err

			return err
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
				return r.execute(task)
			} else {
				r.logger.Error(err.Error(), xlog.F{"message": task.message})
				task.result <- err

				return err
			}
		} else {
			r.logger.Debug("Send message success", xlog.F{"message": task.message})
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

	r.logger.Debug("Send new message to queue", xlog.F{"message": message})

}

func (r *Resource) SendAndReturn(message *gomail.Message) error {
	task := &mailTask{
		message: message,
		result:  make(chan error),
	}
	r.queue <- task

	return <-task.result
}
