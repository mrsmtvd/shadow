package handlers

import (
	"strconv"
	"time"

	"github.com/kihamo/go-workers/dispatcher"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/workers"
)

type HasDispatcher interface {
	workers.Component
	GetDispatcher() *dispatcher.Dispatcher
}

// easyjson:json
type managerHandlerResponseSuccess struct {
	Result string `json:"result"`
}

// easyjson:json
type managerHandlerResponseTask struct {
	Id        string      `json:"id"`
	Name      string      `json:"name"`
	Status    int64       `json:"status"`
	Priority  int64       `json:"priority"`
	Attempts  int64       `json:"attempts"`
	LastError interface{} `json:"last_error"`
	Created   time.Time   `json:"created"`
}

// easyjson:json
type managerHandlerResponseWorker struct {
	Id      string                      `json:"id"`
	Created time.Time                   `json:"created"`
	Status  int64                       `json:"status"`
	Task    *managerHandlerResponseTask `json:"task"`
}

// easyjson:json
type managerHandlerResponseListener struct {
	Name              string     `json:"name"`
	CreatedAt         time.Time  `json:"created_at"`
	LastTaskSuccessAt *time.Time `json:"last_task_success_at"`
	LastTaskFailedAt  *time.Time `json:"last_task_failed_at"`
	CountTaskSuccess  uint64     `json:"count_task_success"`
	CountTaskFailed   uint64     `json:"count_task_failed"`
}

// easyjson:json
type managerHandlerResponse struct {
	Tasks []managerHandlerResponseTask `json:"tasks_wait"`

	Workers      []managerHandlerResponseWorker `json:"workers"`
	WorkersCount int                            `json:"workers_count"`
	WorkersWait  int                            `json:"workers_wait"`
	WorkersBusy  int                            `json:"workers_busy"`

	Listeners      []managerHandlerResponseListener `json:"listeners"`
	ListenersCount int                              `json:"listeners_count"`
}

type ManagerHandler struct {
	dashboard.Handler

	Component HasDispatcher
}

func (h *ManagerHandler) fillResponseTask(t task.Tasker) *managerHandlerResponseTask {
	return &managerHandlerResponseTask{
		Id:        t.GetId(),
		Name:      t.GetName(),
		Status:    t.GetStatus(),
		Priority:  t.GetPriority(),
		Attempts:  t.GetAttempts(),
		LastError: t.GetLastError(),
		Created:   t.GetCreatedAt(),
	}
}

func (h *ManagerHandler) actionStats(w *dashboard.Response, _ *dashboard.Request) {
	tasksList := []managerHandlerResponseTask{}

	for _, t := range h.Component.GetDispatcher().GetTasks() {
		switch t.GetStatus() {
		case task.TaskStatusWait, task.TaskStatusRepeatWait:
			tasksList = append(tasksList, *h.fillResponseTask(t))
		}
	}

	workersList := []managerHandlerResponseWorker{}
	workersWait := 0
	workersBusy := 0

	for _, wrk := range h.Component.GetDispatcher().GetWorkers() {
		switch wrk.GetStatus() {
		case worker.WorkerStatusWait:
			workersWait++
		case worker.WorkerStatusBusy:
			workersBusy++
		}

		workerInfo := managerHandlerResponseWorker{
			Id:      wrk.GetId(),
			Created: wrk.GetCreatedAt(),
			Status:  wrk.GetStatus(),
		}

		t := wrk.GetTask()
		if t != nil {
			workerInfo.Task = h.fillResponseTask(t)

			workersList = append([]managerHandlerResponseWorker{workerInfo}, workersList...)
		} else {
			workersList = append(workersList, workerInfo)
		}
	}

	listenersList := []managerHandlerResponseListener{}
	listenersCount := 0

	for _, l := range h.Component.GetDispatcher().GetListeners() {
		listenersCount++
		item := l.(dispatcher.ListenerItem)

		listenersList = append(listenersList, managerHandlerResponseListener{
			Name:              item.GetName(),
			CreatedAt:         item.GetCreatedAt(),
			LastTaskSuccessAt: item.GetLastTaskSuccessAt(),
			LastTaskFailedAt:  item.GetLastTaskFailedAt(),
			CountTaskSuccess:  item.GetCountTaskSuccess(),
			CountTaskFailed:   item.GetCountTaskFailed(),
		})
	}

	stats := managerHandlerResponse{
		Tasks: tasksList,

		Workers:      workersList,
		WorkersCount: len(workersList),
		WorkersWait:  workersWait,
		WorkersBusy:  workersBusy,

		Listeners:      listenersList,
		ListenersCount: listenersCount,
	}

	w.SendJSON(stats)
}

func (h *ManagerHandler) actionListenersRemove(w *dashboard.Response, r *dashboard.Request) {
	listeners := h.Component.GetDispatcher().GetListeners()
	checkId := r.Original().FormValue("id")

	for _, listener := range listeners {
		if checkId == "" || listener.GetName() == checkId {
			if listener.GetName() != h.Component.GetDefaultListenerName() {
				h.Component.RemoveListener(listener)
			}

			if checkId != "" {
				break
			}
		}
	}

	w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})
}

func (h *ManagerHandler) actionTaskRemove(w *dashboard.Response, r *dashboard.Request) {
	removeId := r.Original().FormValue("id")

	if removeId != "" {
		h.Component.RemoveTaskById(removeId)
	}

	w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})
}

func (h *ManagerHandler) actionWorkersReset(w *dashboard.Response, r *dashboard.Request) {
	workers := h.Component.GetDispatcher().GetWorkers()
	checkId := r.Original().FormValue("id")

	go func() {
		for _, w := range workers {
			if checkId == "" || w.GetId() == checkId {
				w.Reset()
				r.Logger().Infof("Reseted worker #%s", w.GetId())

				if checkId != "" {
					break
				}
			}
		}
	}()

	w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})
}

func (h *ManagerHandler) actionWorkersKill(w *dashboard.Response, r *dashboard.Request) {
	workers := h.Component.GetDispatcher().GetWorkers()
	checkId := r.Original().FormValue("id")

	for _, w := range workers {
		if checkId == "" || w.GetId() == checkId {
			h.Component.RemoveWorker(w)

			if checkId != "" {
				break
			}
		}
	}

	w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})
}

func (h *ManagerHandler) actionWorkersAdd(w *dashboard.Response, r *dashboard.Request) {
	count := r.Original().FormValue("count")
	if count != "" {
		if c, err := strconv.Atoi(count); err == nil {
			for i := 0; i < c; i++ {
				h.Component.AddWorker()
			}
		}
	}

	w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})
}

func (h *ManagerHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if !r.IsAjax() {
		h.Render(r.Context(), h.Component.GetName(), "manager", map[string]interface{}{
			"defaultListenerName": h.Component.GetDefaultListenerName(),
		})
		return
	}

	switch r.URL().Query().Get("action") {
	case "stats":
		h.actionStats(w, r)

	case "listeners-remove":
		if r.IsPost() {
			h.actionListenersRemove(w, r)
		} else {
			h.MethodNotAllowed(w, r)
		}

	case "task-remove":
		if r.IsPost() {
			h.actionTaskRemove(w, r)
		} else {
			h.MethodNotAllowed(w, r)
		}

	case "workers-reset":
		if r.IsPost() {
			h.actionWorkersReset(w, r)
		} else {
			h.MethodNotAllowed(w, r)
		}

	case "workers-kill":
		if r.IsPost() {
			h.actionWorkersKill(w, r)
		} else {
			h.MethodNotAllowed(w, r)
		}

	case "workers-add":
		if r.IsPost() {
			h.actionWorkersAdd(w, r)
		} else {
			h.MethodNotAllowed(w, r)
		}

	default:
		h.NotFound(w, r)
	}
}
