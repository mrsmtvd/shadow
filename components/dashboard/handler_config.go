package dashboard

import (
	"net/http"
	"strings"

	"github.com/kihamo/shadow"
)

type ConfigHandler struct {
	Handler

	application shadow.Application
}

func (h *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	config := ConfigFromContext(r.Context())
	vars := config.GetAllVariables()

	if h.IsPost(r) {
		err = r.ParseForm()
		if err == nil {
			for key, values := range r.PostForm {
				if !config.Has(key) || !config.IsEditable(key) || len(values) == 0 {
					continue
				}

				config.Set(key, values[0])
			}

			h.Redirect(r.URL.String(), http.StatusFound, w, r)
			return
		}
	}

	variables := map[string]map[string]interface{}{}
	for k, v := range vars {
		parts := strings.SplitN(k, ".", 2)

		cmpName := parts[0]
		if !h.application.HasComponent(cmpName) {
			cmpName = "main"
		}

		cmp, ok := variables[cmpName]
		if !ok {
			variables[cmpName] = map[string]interface{}{}
			cmp = variables[cmpName]
		}

		cmp[k] = v
		variables[cmpName] = cmp
	}

	h.Render(r.Context(), ComponentName, "config", map[string]interface{}{
		"variables": variables,
		"error":     err,
	})
}
