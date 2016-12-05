package system

import (
	"runtime"

	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow/resource/workers"
	"github.com/kihamo/shadow/service/frontend"
)

type WorkersHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *WorkersHandler) Handle() {
	if h.IsAjax() {
		resourceWorkers, _ := h.Application.GetResource("workers")
		dispatcher := resourceWorkers.(*workers.Workers)

		workersList := []map[string]interface{}{}
		workersWait := 0
		workersBusy := 0

		for _, w := range dispatcher.GetWorkers() {
			switch w.GetStatus() {
			case worker.WorkerStatusWait:
				workersWait += 1
			case worker.WorkerStatusBusy:
				workersBusy += 1
			}

			workerInfo := map[string]interface{}{
				"id":      w.GetId(),
				"created": w.GetCreatedAt(),
			}

			t := w.GetTask()
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

		h.SendJSON(stats)
		return
	}

	h.SetTemplate("workers.tpl.html")
	h.SetPageTitle("Workers")
	h.SetPageHeader("Workers")
}
