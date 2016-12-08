package alerts

import (
	"sync"
	"time"

	kit "github.com/go-kit/kit/metrics"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/metrics"
)

const (
	MaxInList = 50
	ClearTime = time.Minute * 15
)

type Resource struct {
	metrics *metrics.Resource

	mutex  sync.RWMutex
	alerts []*Alert
	queue  chan *Alert
}

func (r *Resource) GetName() string {
	return "alerts"
}

func (r *Resource) Init(a *shadow.Application) error {
	r.alerts = make([]*Alert, 0)
	r.queue = make(chan *Alert)

	if a.HasResource("metrics") {
		resourceMetrics, _ := a.GetResource("metrics")
		r.metrics = resourceMetrics.(*metrics.Resource)
	}

	return nil
}

func (r *Resource) Run(wg *sync.WaitGroup) error {
	go func() {
		defer wg.Done()

		var metricTotal kit.Counter
		if r.metrics != nil {
			metricTotal = r.metrics.NewCounter(MetricAlertsTotal)
		}

		ticker := time.NewTicker(ClearTime)

		for {
			select {
			case alert := <-r.queue:
				r.mutex.Lock()
				r.alerts = append([]*Alert{alert}, r.alerts...)
				r.mutex.Unlock()

				if metricTotal != nil {
					metricTotal.Add(1)
				}

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

func (r *Resource) Send(title string, message string, icon string) {
	r.queue <- NewAlert(title, message, icon, time.Now())
}

func (a *Resource) GetAlerts() []*Alert {
	a.mutex.Lock()
	a.mutex.Unlock()

	alerts := make([]*Alert, len(a.alerts))
	copy(alerts, a.alerts)

	return alerts
}
