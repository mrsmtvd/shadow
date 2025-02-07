package handlers

import (
	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/database"
	"github.com/mrsmtvd/shadow/components/grpc"
	"github.com/mrsmtvd/shadow/components/metrics"
)

type ComponentsHandler struct {
	dashboard.Handler

	application shadow.Application
}

func NewComponentsHandler(application shadow.Application) *ComponentsHandler {
	return &ComponentsHandler{
		application: application,
	}
}

func (h *ComponentsHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
	components, _ := h.application.GetComponents()
	contextComponents := make([]map[string]interface{}, 0, len(components))

	for _, cmp := range components {
		row := map[string]interface{}{
			"name":                    cmp.Name(),
			"version":                 cmp.Version(),
			"shutdown":                false,
			"dependencies":            []string{},
			"ready":                   false,
			"status":                  h.application.StatusComponent(cmp.Name()).String(),
			"has_assetfs":             false,
			"has_config_variables":    false,
			"has_config_watchers":     false,
			"has_dashboard_menu":      false,
			"has_dashboard_routes":    false,
			"has_dashboard_templates": false,
			"has_database_migrations": false,
			"has_grpc_server":         false,
			"has_metrics":             false,
		}

		switch h.application.StatusComponent(cmp.Name()) {
		case shadow.ComponentStatusReady, shadow.ComponentStatusFinished:
			row["ready"] = true
		}

		if _, ok := cmp.(shadow.ComponentShutdown); ok {
			row["shutdown"] = true
		}

		if deps, ok := cmp.(shadow.ComponentDependency); ok {
			row["dependencies"] = deps.Dependencies()
		}

		if _, ok := cmp.(dashboard.HasAssetFS); ok {
			row["has_assetfs"] = true
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
			templates := tpl.DashboardTemplates()

			if templates != nil {
				row["has_dashboard_templates"] = templates.Prefix
			}
		}

		if h.application.HasComponent(database.ComponentName) {
			if _, ok := cmp.(database.HasMigrations); ok {
				row["has_database_migrations"] = true
			}
		}

		if h.application.HasComponent(grpc.ComponentName) {
			if _, ok := cmp.(grpc.HasGrpcServer); ok {
				row["has_grpc_server"] = true
			}
		}

		if h.application.HasComponent(metrics.ComponentName) {
			if _, ok := cmp.(metrics.HasMetrics); ok {
				row["has_metrics"] = true
			}
		}

		contextComponents = append(contextComponents, row)
	}

	h.Render(r.Context(), "components", map[string]interface{}{
		"components": contextComponents,
	})
}
