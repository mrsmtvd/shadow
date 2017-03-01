package dashboard

import (
	"net/http"

	"github.com/kihamo/shadow"
)

type ComponentsHandler struct {
	Handler

	application shadow.Application
}

func (h *ComponentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contextComponents := []map[string]interface{}{}

	components, _ := h.application.GetComponents()
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

	h.Render(r.Context(), ComponentName, "components", map[string]interface{}{
		"components": contextComponents,
	})
}
