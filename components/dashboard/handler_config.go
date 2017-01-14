package dashboard

import (
	"fmt"
	"net/http"
)

type ConfigHandler struct {
	Handler
}

func (h *ConfigHandler) saveNewValue(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	key := r.PostForm.Get("key")
	if key == "" {
		return nil
	}

	config := ConfigFromContext(r.Context())

	if !config.Has(key) {
		return fmt.Errorf("Variable %s not found", key)
	}

	if !config.IsEditable(key) {
		return fmt.Errorf("Variable %s isn't editable", key)
	}

	currentValue := config.Get(key)

	newValue := r.PostForm.Get("value")
	err := config.Set(key, newValue)

	LoggerFromContext(r.Context()).Infof("Change value for %s with '%v' to '%v'", key, currentValue, newValue)

	return err
}

func (h *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	if h.IsPost(r) {
		err = h.saveNewValue(r)

		if err == nil {
			h.Redirect(r.RequestURI, http.StatusFound, w, r)
			return
		}
	}

	h.Render(r.Context(), "dashboard", "config", map[string]interface{}{
		"variables": ConfigFromContext(r.Context()).GetAllVariables(),
		"error":     err,
	})
}
