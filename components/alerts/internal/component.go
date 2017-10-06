package internal

import (
	"sync"
	"time"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/alerts"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/metrics"
)

const (
	MaxInList = 50
	ClearTime = time.Minute * 15
)

type Component struct {
	application shadow.Application

	mutex  sync.RWMutex
	alerts []alerts.Alert
	queue  chan alerts.Alert
	routes []dashboard.Route
}

func (c *Component) GetName() string {
	return alerts.ComponentName
}

func (c *Component) GetVersion() string {
	return alerts.ComponentVersion
}

func (c *Component) GetDependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name: dashboard.ComponentName,
		},
		{
			Name: metrics.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.alerts = make([]alerts.Alert, 0)
	c.queue = make(chan alerts.Alert)

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(ClearTime)

		for {
			select {
			case alert := <-c.queue:
				c.mutex.Lock()
				c.alerts = append([]alerts.Alert{alert}, c.alerts...)
				c.mutex.Unlock()

				if metricAlertsTotal != nil {
					metricAlertsTotal.Inc()
				}

			case <-ticker.C:
				c.mutex.Lock()
				if len(c.alerts) > MaxInList {
					c.alerts = c.alerts[:MaxInList]
				}
				c.mutex.Unlock()
			}
		}
	}()

	return nil
}

func (c *Component) Send(title string, message string, icon string) {
	c.queue <- NewAlert(title, message, icon, time.Now())
}

func (c *Component) GetAlerts() []alerts.Alert {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	list := make([]alerts.Alert, 0, len(c.alerts))
	for _, a := range c.alerts {
		list = append(list, a)
	}

	return list
}
