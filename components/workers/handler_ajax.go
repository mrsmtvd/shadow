package workers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kihamo/go-workers/dispatcher"
	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow/components/dashboard"
)

// easyjson:json
type ajaxHandlerResponseSuccess struct {
	Result string `json:"result"`
}

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
	Status  int64                    `json:"status"`
	Task    *ajaxHandlerResponseTask `json:"task"`
}

// easyjson:json
type ajaxHandlerResponseListener struct {
	Name              string     `json:"name"`
	CreatedAt         time.Time  `json:"created_at"`
	LastTaskSuccessAt *time.Time `json:"last_task_success_at"`
	LastTaskFailedAt  *time.Time `json:"last_task_failed_at"`
	CountTaskSuccess  uint64     `json:"count_task_success"`
	CountTaskFailed   uint64     `json:"count_task_failed"`
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

func (h *AjaxHandler) actionStats(w http.ResponseWriter, r *http.Request) {
	workersList := []ajaxHandlerResponseWorker{}
	workersWait := 0
	workersBusy := 0

	for _, wrk := range h.component.dispatcher.GetWorkers() {
		switch wrk.GetStatus() {
		case worker.WorkerStatusWait:
			workersWait++
		case worker.WorkerStatusBusy:
			workersBusy++
		}

		workerInfo := ajaxHandlerResponseWorker{
			Id:      wrk.GetId(),
			Created: wrk.GetCreatedAt(),
			Status:  wrk.GetStatus(),
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
		item := l.(dispatcher.ListenerItem)

		listenersList = append(listenersList, ajaxHandlerResponseListener{
			Name:              item.GetName(),
			CreatedAt:         item.GetCreatedAt(),
			LastTaskSuccessAt: item.GetLastTaskSuccessAt(),
			LastTaskFailedAt:  item.GetLastTaskFailedAt(),
			CountTaskSuccess:  item.GetCountTaskSuccess(),
			CountTaskFailed:   item.GetCountTaskFailed(),
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
}

func (h *AjaxHandler) actionWorkersReset(w http.ResponseWriter, r *http.Request) {
	workers := h.component.dispatcher.GetWorkers()
	checkId := r.FormValue("id")

	go func() {
		for _, worker := range workers {
			if checkId == "" || worker.GetId() == checkId {
				worker.Reset()
				h.component.logger.Infof("Reseted worker #%s", worker.GetId())

				if checkId != "" {
					break
				}
			}
		}
	}()

	h.SendJSON(ajaxHandlerResponseSuccess{
		Result: "success",
	}, w)
}

func (h *AjaxHandler) actionWorkersKill(w http.ResponseWriter, r *http.Request) {
	workers := h.component.dispatcher.GetWorkers()
	checkId := r.FormValue("id")

	for _, worker := range workers {
		if checkId == "" || worker.GetId() == checkId {
			h.component.RemoveWorker(worker)

			if checkId != "" {
				break
			}
		}
	}

	h.SendJSON(ajaxHandlerResponseSuccess{
		Result: "success",
	}, w)
}

func (h *AjaxHandler) actionWorkersAdd(w http.ResponseWriter, r *http.Request) {
	count := r.FormValue("count")
	if count != "" {
		if c, err := strconv.Atoi(count); err == nil {
			for i := 0; i < c; i++ {
				h.component.AddWorker()
			}
		}
	}

	h.SendJSON(ajaxHandlerResponseSuccess{
		Result: "success",
	}, w)
}

func (h *AjaxHandler) actionListenersRemove(w http.ResponseWriter, r *http.Request) {
	listeners := h.component.dispatcher.GetListeners()
	checkId := r.FormValue("id")

	for _, listener := range listeners {
		if checkId == "" || listener.GetName() == checkId {
			if listener.GetName() != h.component.getDefaultListenerName() {
				h.component.RemoveListener(listener)
			}

			if checkId != "" {
				break
			}
		}
	}

	h.SendJSON(ajaxHandlerResponseSuccess{
		Result: "success",
	}, w)
}

func (h *AjaxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("action") {
	case "stats":
		h.actionStats(w, r)

	case "listeners-remove":
		if h.IsPost(r) {
			h.actionListenersRemove(w, r)
		} else {
			h.MethodNotAllowed(w, r)
		}

	case "workers-reset":
		if h.IsPost(r) {
			h.actionWorkersReset(w, r)
		} else {
			h.MethodNotAllowed(w, r)
		}

	case "workers-kill":
		if h.IsPost(r) {
			h.actionWorkersKill(w, r)
		} else {
			h.MethodNotAllowed(w, r)
		}

	case "workers-add":
		if h.IsPost(r) {
			h.actionWorkersAdd(w, r)
		} else {
			h.MethodNotAllowed(w, r)
		}

	default:
		h.NotFound(w, r)
	}
}
