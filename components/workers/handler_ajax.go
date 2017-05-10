package workers

import (
	"net/http"
	"time"

	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow/components/dashboard"
)

// easyjson:json
type ajaxHandlerResponseTask struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Status  int64     `json:"status"`
	Created time.Time `json:"created"`
}

// easyjson:json
type ajaxHandlerResponseWorker struct {
	Id      string                   `json:"id"`
	Created time.Time                `json:"created"`
	Task    *ajaxHandlerResponseTask `json:"task"`
}

// easyjson:json
type ajaxHandlerResponseListener struct {
	Name string `json:"name"`
}

// easyjson:json
type ajaxHandlerResponse struct {
	Workers      []ajaxHandlerResponseWorker `json:"workers"`
	WorkersCount int                         `json:"workers_count"`
	WorkersWait  int                         `json:"workers_wait"`
	WorkersBusy  int                         `json:"workers_busy"`

	Listeners      []ajaxHandlerResponseListener `json:"listeners"`
	ListenersCount int                           `json:"listeners_count"`
}

type AjaxHandler struct {
	dashboard.Handler

	component *Component
}

func (h *AjaxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	workersList := []ajaxHandlerResponseWorker{}
	workersWait := 0
	workersBusy := 0

	for _, wrk := range h.component.dispatcher.GetWorkers().GetItems() {
		switch wrk.GetStatus() {
		case worker.WorkerStatusWait:
			workersWait++
		case worker.WorkerStatusBusy:
			workersBusy++
		}

		workerInfo := ajaxHandlerResponseWorker{
			Id:      wrk.GetId(),
			Created: wrk.GetCreatedAt(),
		}

		t := wrk.GetTask()
		if t != nil {
			workerInfo.Task = &ajaxHandlerResponseTask{
				Id:      t.GetId(),
				Name:    t.GetName(),
				Status:  t.GetStatus(),
				Created: t.GetCreatedAt(),
			}
		}

		workersList = append(workersList, workerInfo)
	}

	listenersList := []ajaxHandlerResponseListener{}
	listenersCount := 0

	for _, l := range h.component.dispatcher.GetListeners() {
		listenersCount++

		listenersList = append(listenersList, ajaxHandlerResponseListener{
			Name: l.GetName(),
		})
	}

	stats := ajaxHandlerResponse{
		Workers:      workersList,
		WorkersCount: len(workersList),
		WorkersWait:  workersWait,
		WorkersBusy:  workersBusy,

		Listeners:      listenersList,
		ListenersCount: listenersCount,
	}

	h.SendJSON(stats, w)
	return
}
