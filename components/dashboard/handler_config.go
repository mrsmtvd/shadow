package dashboard

import (
	"net/http"
	"strings"

	"github.com/kihamo/shadow"
)

const (
	defaultComponentName = "main"
)

type ConfigHandler struct {
	Handler

	application shadow.Application
}

func (h *ConfigHandler) ServeHTTP(w *Response, r *Request) {
	var err error

	vars := r.Config().GetAllVariables()

	if r.IsPost() {
		err = r.Original().ParseForm()
		if err == nil {
			for key, values := range r.Original().PostForm {
				if !r.Config().Has(key) || !r.Config().IsEditable(key) || len(values) == 0 {
					continue
				}

				r.Config().Set(key, values[0])
			}

			h.Redirect(r.URL().String(), http.StatusFound, w, r)
			return
		}
	}

	variables := map[string]map[string]interface{}{}
	for k, v := range vars {
		parts := strings.SplitN(k, ".", 2)

		cmpName := parts[0]
		if !h.application.HasComponent(cmpName) {
			cmpName = defaultComponentName
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
