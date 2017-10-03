package handlers

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/components/grpc"
	"github.com/kihamo/shadow/components/metrics"
)

type ComponentsHandler struct {
	dashboard.Handler

	Application shadow.Application
}

func (h *ComponentsHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
	contextComponents := []map[string]interface{}{}

	components, _ := h.Application.GetComponents()
	for _, cmp := range components {
		row := map[string]interface{}{
			"name":                    cmp.GetName(),
			"version":                 cmp.GetVersion(),
			"dependencies":            []string{},
			"has_config_variables":    false,
			"has_config_watchers":     false,
			"has_dashboard_menu":      false,
			"has_dashboard_routes":    false,
			"has_dashboard_templates": false,
			"has_database_migrations": false,
			"has_grpc_server":         false,
			"has_metrics":             false,
		}

		if deps, ok := cmp.(shadow.ComponentDependency); ok {
			row["dependencies"] = deps.GetDependencies()
		}

		if _, ok := cmp.(config.HasVariables); ok {
			row["has_config_variables"] = true
		}

		if _, ok := cmp.(config.HasWatchers); ok {
			row["has_config_watchers"] = true
		}

		if _, ok := cmp.(dashboard.HasMenu); ok {
			row["has_dashboard_menu"] = true
		}

		if _, ok := cmp.(dashboard.HasRoutes); ok {
			row["has_dashboard_routes"] = true
		}

		if tpl, ok := cmp.(dashboard.HasTemplates); ok {
			templates := tpl.GetTemplates()

			if templates != nil {
				row["has_dashboard_templates"] = templates.Prefix
			}
		}

		if h.Application.HasComponent(database.ComponentName) {
			if _, ok := cmp.(database.HasMigrations); ok {
				row["has_database_migrations"] = true
			}
		}

		if h.Application.HasComponent(grpc.ComponentName) {
			if _, ok := cmp.(grpc.HasGrpcServer); ok {
				row["has_grpc_server"] = true
			}
		}

		if h.Application.HasComponent(metrics.ComponentName) {
			if _, ok := cmp.(metrics.HasMetrics); ok {
				row["has_metrics"] = true
			}
		}

		contextComponents = append(contextComponents, row)
	}

	h.Render(r.Context(), dashboard.ComponentName, "components", map[string]interface{}{
		"components": contextComponents,
	})
}
