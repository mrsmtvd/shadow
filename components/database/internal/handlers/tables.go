package handlers

import (
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/database"
	"github.com/mrsmtvd/shadow/components/database/storage"
)

type TablesHandler struct {
	dashboard.Handler

	component database.Component
}

func NewTablesHandler(component database.Component) *TablesHandler {
	return &TablesHandler{
		component: component,
	}
}

func (h *TablesHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), "tables", map[string]interface{}{
		"tables": h.component.Storage().(*storage.SQL).Tables(),
	})
}
