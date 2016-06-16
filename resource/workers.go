package resource

import (
	"sync"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/shadow"
)

type Workers struct {
	config     *Config
	dispatcher *workers.Dispatcher
}

func (r *Workers) GetName() string {
	return "workers"
}

func (r *Workers) GetConfigVariables() []ConfigVariable {
	return []ConfigVariable{
		ConfigVariable{
			Key:   "workers.count",
			Value: 2,
			Usage: "Default workers count",
		},
	}
}

func (r *Workers) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	r.config = resourceConfig.(*Config)

	r.dispatcher = workers.NewDispatcher()

	return nil
}

func (r *Workers) Run(wg *sync.WaitGroup) (err error) {
	for i := 1; i <= int(r.config.GetInt64("workers.count")); i++ {
		r.dispatcher.AddWorker()
	}

	go func() {
		defer wg.Done()

		r.dispatcher.Run()
	}()

	return nil
}

func (r *Workers) GetDispatcher() *workers.Dispatcher {
	return r.dispatcher
}
