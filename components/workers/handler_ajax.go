package workers

import (
	"net/http"
	"runtime"

	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow/components/dashboard"
)

type AjaxHandler struct {
	dashboard.Handler

	component *Component
}

func (h *AjaxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	workersList := []map[string]interface{}{}
	workersWait := 0
	workersBusy := 0

	for _, wrk := range h.component.GetWorkers() {
		switch wrk.GetStatus() {
		case worker.WorkerStatusWait:
			workersWait += 1
		case worker.WorkerStatusBusy:
			workersBusy += 1
		}

		workerInfo := map[string]interface{}{
			"id":      wrk.GetId(),
			"created": wrk.GetCreatedAt(),
		}

		t := wrk.GetTask()
		if t != nil {
			workerInfo["task"] = map[string]interface{}{
				"id":      t.GetId(),
				"name":    t.GetName(),
				"status":  t.GetStatus(),
				"created": t.GetCreatedAt(),
			}
		}

		workersList = append(workersList, workerInfo)
	}

	stats := map[string]interface{}{
		"workers":       workersList,
		"workers_count": len(workersList),
		"workers_wait":  workersWait,
		"workers_busy":  workersBusy,
		"goroutines":    runtime.NumGoroutine(),
	}

	h.SendJSON(stats, w)
	return
}
