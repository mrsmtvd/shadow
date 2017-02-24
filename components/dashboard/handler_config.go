package dashboard

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/kihamo/shadow"
)

type ConfigHandler struct {
	Handler

	application shadow.Application
}

func (h *ConfigHandler) saveNewValue(r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", err
	}

	key := r.PostForm.Get("key")
	if key == "" {
		return "", nil
	}

	config := ConfigFromContext(r.Context())

	if !config.Has(key) {
		return "", fmt.Errorf("Variable %s not found", key)
	}

	if !config.IsEditable(key) {
		return "", fmt.Errorf("Variable %s isn't editable", key)
	}

	err := config.Set(key, r.PostForm.Get("value"))

	return key, err
}

func (h *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	if h.IsPost(r) {
		key, err := h.saveNewValue(r)

		if err == nil {
			parts := strings.SplitN(key, ".", 2)

			redirectUrl := &url.URL{}
			*redirectUrl = *r.URL
			redirectUrl.RawQuery = "tab=" + parts[0]

			h.Redirect(redirectUrl.String(), http.StatusFound, w, r)
			return
		}
	}

	variables := map[string]map[string]interface{}{}
	for k, v := range ConfigFromContext(r.Context()).GetAllVariables() {
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

	h.Render(r.Context(), "dashboard", "config", map[string]interface{}{
		"variables": variables,
		"error":     err,
	})
}
