package handlers

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/dashboard"
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
			"name":         cmp.GetName(),
			"version":      cmp.GetVersion(),
			"dependencies": []string{},
		}

		if deps, ok := cmp.(shadow.ComponentDependency); ok {
			row["dependencies"] = deps.GetDependencies()
		}

		contextComponents = append(contextComponents, row)
	}

	h.Render(r.Context(), dashboard.ComponentName, "components", map[string]interface{}{
		"components": contextComponents,
	})
}
