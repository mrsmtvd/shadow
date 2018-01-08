package handlers

import (
	"strconv"
	"time"

	ws "github.com/kihamo/go-workers"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/workers"
)

// easyjson:json
type managerHandlerResponseSuccess struct {
	Result string `json:"result"`
}

// easyjson:json
type managerHandlerStats struct {
	TasksCount          uint64                       `json:"tasks_count"`
	TasksWaitingCount   uint64                       `json:"tasks_waiting_count"`
	WorkersCount        uint64                       `json:"workers_count"`
	WorkersWaitingCount uint64                       `json:"workers_waiting_count"`
	ListenersCount      uint64                       `json:"listeners_count"`
	Workers             []managerHandlerItemWorker   `json:"workers"`
	Tasks               []managerHandlerItemTask     `json:"tasks"`
	Listeners           []managerHandlerItemListener `json:"listeners"`
}

// easyjson:json
type managerHandlerItemWorker struct {
	Id      string                  `json:"id"`
	Created time.Time               `json:"created"`
	Status  string                  `json:"status"`
	Task    *managerHandlerItemTask `json:"task"`
}

// easyjson:json
type managerHandlerItemTask struct {
	Id             string        `json:"id"`
	Name           string        `json:"name"`
	Priority       int64         `json:"priority"`
	Repeats        int64         `json:"repeats"`
	RepeatInterval time.Duration `json:"repeat_interval"`
	Timeout        time.Duration `json:"timeout"`
	CreatedAt      time.Time     `json:"created_at"`
	StartedAt      *time.Time    `json:"started_at"`
	Status         string        `json:"status"`
	Attempts       int64         `json:"attempts"`
	AllowStartAt   *time.Time    `json:"allow_start_at"`
	FirstStartedAt *time.Time    `json:"first_started_at"`
	LastStartedAt  *time.Time    `json:"last_started_at"`
}

// easyjson:json
type managerHandlerItemListener struct {
	Id           string           `json:"id"`
	Name         string           `json:"name"`
	Locked       bool             `json:"locked"`
	Events       map[string]int64 `json:"events"`
	Fires        int64            `json:"fires"`
	FirstFiredAt *time.Time       `json:"first_fired_at"`
	LastFiredAt  *time.Time       `json:"last_fired_at"`
}

type ManagerHandler struct {
	dashboard.Handler

	Component workers.Component
}

func (h *ManagerHandler) isLocked(id string) bool {
	for _, listenerId := range h.Component.GetLockedListeners() {
		if id == listenerId {
			return true
		}
	}

	return false
}

func (h *ManagerHandler) actionStats(w *dashboard.Response, r *dashboard.Request) {
	stats := struct {
		Draw     int         `json:"draw"`
		Total    int         `json:"recordsTotal"`
		Filtered int         `json:"recordsFiltered"`
		Data     interface{} `json:"data"`
		Error    string      `json:"error,omitempty"`
	}{}

	stats.Draw = 1

	switch r.URL().Query().Get("entity") {
	case "listeners":
		list := make([]managerHandlerItemListener, 0, 0)

		for _, item := range h.Component.GetListeners() {
			listener := managerHandlerItemListener{
				Id:     item.Id(),
				Name:   item.Name(),
				Locked: h.isLocked(item.Id()),
			}

			if md := h.Component.GetListenerMetadata(item.Id()); md != nil {
				listener.Fires = md[ws.ListenerMetadataFires].(int64)
				listener.FirstFiredAt = md[ws.ListenerMetadataFirstFiredAt].(*time.Time)
				listener.LastFiredAt = md[ws.ListenerMetadataLastFireAt].(*time.Time)

				events := md[ws.ListenerMetadataEventIds].([]ws.EventId)
				listener.Events = make(map[string]int64, len(events))
				for _, eventId := range events {
					listener.Events[eventId.String()] = eventId.Int64()
				}
			}

			list = append(list, listener)
		}

		stats.Data = list
		stats.Total = len(list)

	case "workers":
		list := make([]managerHandlerItemWorker, 0, 0)

		for _, item := range h.Component.GetWorkers() {
			data := managerHandlerItemWorker{
				Id:      item.Id(),
				Created: item.CreatedAt(),
			}

			if md := h.Component.GetWorkerMetadata(item.Id()); md != nil {
				data.Status = md[ws.WorkerMetadataStatus].(ws.Status).String()

				if task := md[ws.WorkerMetadataTask]; task != nil {
					item := task.(ws.Task)

					data.Task = &managerHandlerItemTask{
						Id:             item.Id(),
						Name:           item.Name(),
						Priority:       item.Priority(),
						Repeats:        item.Repeats(),
						RepeatInterval: item.RepeatInterval(),
						Timeout:        item.Timeout(),
						CreatedAt:      item.CreatedAt(),
						StartedAt:      item.StartedAt(),
					}

					if taskMD := h.Component.GetTaskMetadata(item.Id()); taskMD != nil {
						data.Task.Status = taskMD[ws.TaskMetadataStatus].(ws.Status).String()
						data.Task.Attempts = taskMD[ws.TaskMetadataAttempts].(int64)
						data.Task.AllowStartAt = taskMD[ws.TaskMetadataAllowStartAt].(*time.Time)
						data.Task.FirstStartedAt = taskMD[ws.TaskMetadataFirstStartedAt].(*time.Time)
						data.Task.LastStartedAt = taskMD[ws.TaskMetadataLastStartedAt].(*time.Time)
					}
				}
			}

			list = append(list, data)
		}

		stats.Data = list
		stats.Total = len(list)

	case "tasks":
		list := make([]managerHandlerItemTask, 0, 0)

		for _, item := range h.Component.GetTasks() {
			data := managerHandlerItemTask{
				Id:             item.Id(),
				Name:           item.Name(),
				Priority:       item.Priority(),
				Repeats:        item.Repeats(),
				RepeatInterval: item.RepeatInterval(),
				Timeout:        item.Timeout(),
				CreatedAt:      item.CreatedAt(),
				StartedAt:      item.StartedAt(),
			}

			if md := h.Component.GetTaskMetadata(item.Id()); md != nil {
				data.Status = md[ws.TaskMetadataStatus].(ws.Status).String()
				data.Attempts = md[ws.TaskMetadataAttempts].(int64)
				data.AllowStartAt = md[ws.TaskMetadataAllowStartAt].(*time.Time)
				data.FirstStartedAt = md[ws.TaskMetadataFirstStartedAt].(*time.Time)
				data.LastStartedAt = md[ws.TaskMetadataLastStartedAt].(*time.Time)
			}

			list = append(list, data)
		}

		stats.Data = list
		stats.Total = len(list)

	default:
		h.NotFound(w, r)
		return
	}

	stats.Filtered = stats.Total

	w.SendJSON(stats)
}

func (h *ManagerHandler) actionListenersRemove(w *dashboard.Response, r *dashboard.Request) {
	checkId := r.Original().PostFormValue("id")

	if checkId != "" && !h.isLocked(checkId) {
		for _, listener := range h.Component.GetListeners() {
			if listener.Id() == checkId {
				checkEvents := r.Original().PostForm["events[]"]

				if len(checkEvents) != 0 {
					for _, eventId := range checkEvents {
						if i, err := strconv.ParseInt(eventId, 10, 64); err == nil {
							h.Component.RemoveListenerByEvent(ws.EventId(i), listener)
						}
					}
				} else {
					h.Component.RemoveListener(listener)
				}

				break
			}
		}
	}

	w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})
}

func (h *ManagerHandler) actionTasksRemove(w *dashboard.Response, r *dashboard.Request) {
	checkId := r.Original().FormValue("id")

	for _, task := range h.Component.GetTasks() {
		if checkId == "" || task.Id() == checkId {
			h.Component.RemoveTask(task)

			if checkId != "" {
				break
			}
		}
	}

	w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})
}

func (h *ManagerHandler) actionWorkerRemove(w *dashboard.Response, r *dashboard.Request) {
	checkId := r.Original().FormValue("id")

	for _, worker := range h.Component.GetWorkers() {
		if checkId == "" || worker.Id() == checkId {
			h.Component.RemoveWorker(worker)

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
				h.Component.AddSimpleWorker()
			}
		}
	}

	w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})
}

func (h *ManagerHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if !r.IsAjax() {
		h.Render(r.Context(), h.Component.GetName(), "manager", map[string]interface{}{})
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

	case "tasks-remove":
		if r.IsPost() {
			h.actionTasksRemove(w, r)
		} else {
			h.MethodNotAllowed(w, r)
		}

	case "workers-remove":
		if r.IsPost() {
			h.actionWorkerRemove(w, r)
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
