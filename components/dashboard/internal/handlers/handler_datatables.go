package handlers

import (
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
)

const (
	DataTablesMessageCtx = "datatables"
)

type DataTablesHandler struct {
	dashboard.Handler
}

func (h *DataTablesHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	locale := i18n.Locale(r.Context())

	translate := map[string]interface{}{
		"processing":     locale.Translate(dashboard.ComponentName, "Processing...", DataTablesMessageCtx),
		"search":         locale.Translate(dashboard.ComponentName, "Search:", DataTablesMessageCtx),
		"lengthMenu":     locale.Translate(dashboard.ComponentName, "Show _MENU_ entries", DataTablesMessageCtx),
		"info":           locale.Translate(dashboard.ComponentName, "Showing _START_ to _END_ of _TOTAL_ entries", DataTablesMessageCtx),
		"infoEmpty":      locale.Translate(dashboard.ComponentName, "Showing 0 to 0 of 0 entries", DataTablesMessageCtx),
		"infoFiltered":   locale.Translate(dashboard.ComponentName, "(filtered from _MAX_ total entries)", DataTablesMessageCtx),
		"infoPostFix":    "",
		"loadingRecords": locale.Translate(dashboard.ComponentName, "Loading...", DataTablesMessageCtx),
		"zeroRecords":    locale.Translate(dashboard.ComponentName, "No matching records found", DataTablesMessageCtx),
		"emptyTable":     locale.Translate(dashboard.ComponentName, "No data available in table", DataTablesMessageCtx),
		"paginate": map[string]interface{}{
			"first":    locale.Translate(dashboard.ComponentName, "First", DataTablesMessageCtx),
			"previous": locale.Translate(dashboard.ComponentName, "Previous", DataTablesMessageCtx),
			"next":     locale.Translate(dashboard.ComponentName, "Next", DataTablesMessageCtx),
			"last":     locale.Translate(dashboard.ComponentName, "Last", DataTablesMessageCtx),
		},
		"aria": map[string]interface{}{
			"sortAscending":  locale.Translate(dashboard.ComponentName, ": activate to sort column ascending", DataTablesMessageCtx),
			"sortDescending": locale.Translate(dashboard.ComponentName, ": activate to sort column descending", DataTablesMessageCtx),
		},
	}

	w.Header().Set("Cache-Control", "max-age=315360000, private, immutable")
	if err := w.SendJSON(translate); err != nil {
		h.InternalError(w, r, err)
	}
}
