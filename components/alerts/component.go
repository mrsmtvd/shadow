package alerts

import (
	"sync"
	"time"

	"github.com/kihamo/shadow"
)

const (
	MaxInList = 50
	ClearTime = time.Minute * 15
)

type Component struct {
	application *shadow.Application

	mutex  sync.RWMutex
	alerts []*Alert
	queue  chan *Alert
}

func (c *Component) GetName() string {
	return "alerts"
}

func (c *Component) GetVersion() string {
	return "1.0.0"
}

func (c *Component) Init(a *shadow.Application) error {
	c.alerts = make([]*Alert, 0)
	c.queue = make(chan *Alert)

	c.application = a

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
				c.alerts = append([]*Alert{alert}, c.alerts...)
				c.mutex.Unlock()

				if metricAlertsTotal != nil {
					metricAlertsTotal.Add(1)
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

func (c *Component) GetAlerts() []*Alert {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	alerts := make([]*Alert, len(c.alerts))
	copy(alerts, c.alerts)

	return alerts
}
