package frontend

import (
	"time"

	"github.com/kihamo/shadow"
)

const (
	maxAlertsInList = 50
	cleatTime       = time.Minute * 15
)

var (
	alertsList  []*alert
	alertsQueue chan *alert
)

type alert struct {
	Icon    string
	Message string
	Date    time.Time
}

func (a *alert) DateAsMessage() string {
	return shadow.DateSinceAsMessage(a.Date)
}

func (s *FrontendService) initAlerts() {
	alertsList = []*alert{}
	alertsQueue = make(chan *alert)

	go func() {
		ticker := time.NewTicker(cleatTime)

		for {
			select {
			case a := <-alertsQueue:
				alertsList = append([]*alert{a}, alertsList...)

			case <-ticker.C:
				if len(alertsList) > maxAlertsInList {
					alertsList = alertsList[:maxAlertsInList]
				}
			}
		}
	}()
}

func (s *FrontendService) SendAlert(message string, icon string) {
	alertsQueue <- &alert{
		Icon:    icon,
		Message: message,
		Date:    time.Now(),
	}
}

func (s *FrontendService) GetAlerts() []*alert {
	return alertsList
}
