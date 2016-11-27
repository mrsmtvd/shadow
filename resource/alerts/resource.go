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

type Alerts struct {
	mutex  sync.RWMutex
	alerts []*Alert
	queue  chan *Alert
}

func (r *Alerts) GetName() string {
	return "alerts"
}

func (r *Alerts) Init(a *shadow.Application) error {
	r.alerts = make([]*Alert, 0)
	r.queue = make(chan *Alert)

	return nil
}

func (r *Alerts) Run(wg *sync.WaitGroup) error {
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(ClearTime)

		for {
			select {
			case alert := <-r.queue:
				r.mutex.Lock()
				r.alerts = append([]*Alert{alert}, r.alerts...)
				r.mutex.Unlock()

			case <-ticker.C:
				r.mutex.Lock()
				if len(r.alerts) > MaxInList {
					r.alerts = r.alerts[:MaxInList]
				}
				r.mutex.Unlock()
			}
		}
	}()

	return nil
}

func (r *Alerts) Send(title string, message string, icon string) {
	r.queue <- NewAlert(title, message, icon, time.Now())
}

func (a *Alerts) GetAlerts() []*Alert {
	a.mutex.Lock()
	a.mutex.Unlock()

	alerts := make([]*Alert, len(a.alerts))
	copy(alerts, a.alerts)

	return alerts
}
