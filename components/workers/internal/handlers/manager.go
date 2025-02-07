package handlers

import (
	"strconv"
	"time"

	ws "github.com/kihamo/go-workers"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/i18n"
	"github.com/mrsmtvd/shadow/components/workers"
)

// easyjson:json
type managerHandlerResponseSuccess struct {
	Result string `json:"result"`
}

// easyjson:json
type managerHandlerItemWorker struct {
	ID      string                  `json:"id"`
	Created time.Time               `json:"created"`
	Status  string                  `json:"status"`
	Locked  bool                    `json:"locked"`
	Task    *managerHandlerItemTask `json:"task"`
}

// easyjson:json
type managerHandlerItemTask struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Priority       int64         `json:"priority"`
	Repeats        int64         `json:"repeats"`
	RepeatInterval time.Duration `json:"repeat_interval"`
	Timeout        time.Duration `json:"timeout"`
	CreatedAt      time.Time     `json:"created_at"`
	StartedAt      *time.Time    `json:"started_at"`
	Status         string        `json:"status"`
	Locked         bool          `json:"locked"`
	Attempts       int64         `json:"attempts"`
	AllowStartAt   *time.Time    `json:"allow_start_at"`
	FirstStartedAt *time.Time    `json:"first_started_at"`
	LastStartedAt  *time.Time    `json:"last_started_at"`
}

// easyjson:json
type managerHandlerItemListener struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Locked       bool              `json:"locked"`
	Events       map[string]string `json:"events"`
	Fires        int64             `json:"fires"`
	FirstFiredAt *time.Time        `json:"first_fired_at"`
	LastFiredAt  *time.Time        `json:"last_fired_at"`
}

type ManagerHandler struct {
	dashboard.Handler

	component workers.Component
}

func NewManagerHandler(component workers.Component) *ManagerHandler {
	return &ManagerHandler{
		component: component,
	}
}

func (h *ManagerHandler) isLocked(id string) bool {
	for _, listener := range h.component.LockedListeners() {
		if id == listener.Id() {
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
		listListeners := h.component.GetListeners()
		list := make([]managerHandlerItemListener, 0, len(listListeners))

		for _, item := range listListeners {
			listener := managerHandlerItemListener{
				ID:     item.Id(),
				Name:   item.Name(),
				Locked: h.isLocked(item.Id()),
			}

			md := h.component.GetListenerMetadata(item.Id())
			if md == nil {
				continue
			}

			listener.Fires = md[ws.ListenerMetadataFires].(int64)
			listener.FirstFiredAt = md[ws.ListenerMetadataFirstFiredAt].(*time.Time)
			listener.LastFiredAt = md[ws.ListenerMetadataLastFireAt].(*time.Time)

			events := md[ws.ListenerMetadataEvents].([]ws.Event)
			listener.Events = make(map[string]string, len(events))

			for _, event := range events {
				listener.Events[event.Id()] = event.Name()
			}

			list = append(list, listener)
		}

		stats.Data = list
		stats.Total = len(list)

	case "workers":
		listWorkers := h.component.GetWorkers()
		list := make([]managerHandlerItemWorker, 0, len(listWorkers))
		locale := i18n.Locale(r.Context())

		for _, item := range listWorkers {
			data := managerHandlerItemWorker{
				ID:      item.Id(),
				Created: item.CreatedAt(),
			}

			if md := h.component.GetWorkerMetadata(item.Id()); md != nil {
				data.Status = locale.Translate(workers.ComponentName, md[ws.WorkerMetadataStatus].(ws.Status).String(), "worker")
				data.Locked = md[ws.WorkerMetadataLocked].(bool)

				if task := md[ws.WorkerMetadataTask]; task != nil {
					item := task.(ws.Task)

					data.Task = &managerHandlerItemTask{
						ID:             item.Id(),
						Name:           item.Name(),
						Priority:       item.Priority(),
						Repeats:        item.Repeats(),
						RepeatInterval: item.RepeatInterval(),
						Timeout:        item.Timeout(),
						CreatedAt:      item.CreatedAt(),
						StartedAt:      item.StartedAt(),
					}

					if taskMD := h.component.GetTaskMetadata(item.Id()); taskMD != nil {
						data.Task.Status = taskMD[ws.TaskMetadataStatus].(ws.Status).String()
						data.Task.Locked = taskMD[ws.TaskMetadataLocked].(bool)
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
		listTasks := h.component.GetTasks()
		list := make([]managerHandlerItemTask, 0, len(listTasks))
		locale := i18n.Locale(r.Context())

		for _, item := range listTasks {
			data := managerHandlerItemTask{
				ID:             item.Id(),
				Name:           item.Name(),
				Priority:       item.Priority(),
				Repeats:        item.Repeats(),
				RepeatInterval: item.RepeatInterval(),
				Timeout:        item.Timeout(),
				CreatedAt:      item.CreatedAt(),
				StartedAt:      item.StartedAt(),
			}

			if md := h.component.GetTaskMetadata(item.Id()); md != nil {
				data.Status = locale.Translate(workers.ComponentName, md[ws.TaskMetadataStatus].(ws.Status).String(), "task")
				data.Locked = md[ws.TaskMetadataLocked].(bool)
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

	if err := w.SendJSON(stats); err != nil {
		h.InternalError(w, r, err)
	}
}

func (h *ManagerHandler) actionListenersRemove(w *dashboard.Response, r *dashboard.Request) {
	checkID := r.Original().PostFormValue("id")

	if checkID != "" && !h.isLocked(checkID) {
		for _, listener := range h.component.GetListeners() {
			if listener.Id() == checkID {
				checkEvents := r.Original().PostForm["events[]"]

				if len(checkEvents) != 0 {
					md := h.component.GetListenerMetadata(listener.Id())
					if md == nil {
						continue
					}

					mdEvents := md[ws.ListenerMetadataEvents].([]ws.Event)
					events := make(map[string]ws.Event, len(mdEvents))

					for _, event := range mdEvents {
						events[event.Id()] = event
					}

					for _, eventID := range checkEvents {
						event, ok := events[eventID]
						if !ok {
							continue
						}

						h.component.RemoveListenerByEvent(event, listener)
					}
				} else {
					h.component.RemoveListener(listener)
				}

				break
			}
		}
	}

	err := w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})

	if err != nil {
		h.InternalError(w, r, err)
	}
}

func (h *ManagerHandler) actionTasksRemove(w *dashboard.Response, r *dashboard.Request) {
	checkID := r.Original().FormValue("id")

	for _, task := range h.component.GetTasks() {
		if checkID == "" || task.Id() == checkID {
			h.component.RemoveTask(task)

			if checkID != "" {
				break
			}
		}
	}

	err := w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})

	if err != nil {
		h.InternalError(w, r, err)
	}
}

func (h *ManagerHandler) actionWorkerRemove(w *dashboard.Response, r *dashboard.Request) {
	checkID := r.Original().FormValue("id")

	for _, worker := range h.component.GetWorkers() {
		if checkID == "" || worker.Id() == checkID {
			h.component.RemoveWorker(worker)

			if checkID != "" {
				break
			}
		}
	}

	err := w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})

	if err != nil {
		h.InternalError(w, r, err)
	}
}

func (h *ManagerHandler) actionWorkersAdd(w *dashboard.Response, r *dashboard.Request) {
	count := r.Original().FormValue("count")
	if count != "" {
		if c, err := strconv.Atoi(count); err == nil {
			for i := 0; i < c; i++ {
				h.component.AddSimpleWorker()
			}
		}
	}

	err := w.SendJSON(managerHandlerResponseSuccess{
		Result: "success",
	})

	if err != nil {
		h.InternalError(w, r, err)
	}
}

func (h *ManagerHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if !r.IsAjax() {
		h.Render(r.Context(), "manager", nil)
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
