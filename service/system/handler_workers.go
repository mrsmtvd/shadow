package system

import (
	"runtime"

	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow/resource"
	"github.com/kihamo/shadow/service/frontend"
)

type WorkersHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *WorkersHandler) Handle() {
	if h.IsAjax() {
		resourceWorkers, _ := h.Application.GetResource("workers")
		dispatcher := resourceWorkers.(*resource.Workers).GetDispatcher()

		members := dispatcher.GetWorkers().GetItems()
		workersList := []map[string]interface{}{}
		workersWait := 0
		workersBusy := 0

		for _, member := range members {
			switch member.GetStatus() {
			case worker.WorkerStatusWait:
				workersWait += 1
			case worker.WorkerStatusBusy:
				workersBusy += 1
			}

			workersList = append(workersList, map[string]interface{}{
				"id":      member.GetId(),
				"created": member.GetCreatedAt(),
			})
		}

		stats := map[string]interface{}{
			"workers":       workersList,
			"workers_count": len(members),
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
