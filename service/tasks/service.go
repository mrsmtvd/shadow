package tasks

import (
	"runtime"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
)

type TasksService struct {
	application *shadow.Application
	dispatcher  *Dispatcher
	logger      *logrus.Entry
}

func (s *TasksService) GetName() string {
	return "tasks"
}

func (s *TasksService) Init(a *shadow.Application) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	s.application = a
	s.dispatcher = NewDispatcher()

	resourceLogger, err := a.GetResource("logger")
	if err != nil {
		return err
	}
	s.logger = resourceLogger.(*resource.Logger).Get(s.GetName())

	return nil
}

func (s *TasksService) Run(wg *sync.WaitGroup) error {
	s.logger.Info("Start tasks manager")

	return nil
}

func (s *TasksService) AddTask(fn func(...interface{}) (bool, time.Duration), args ...interface{}) {
	s.dispatcher.AddTask(fn, args...)
}
