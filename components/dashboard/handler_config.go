package dashboard

import (
	"fmt"
	"net/http"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

type ConfigHandler struct {
	TemplateHandler

	config *config.Component
	logger logger.Logger
}

func (h *ConfigHandler) saveNewValue() error {
	if err := h.Request().ParseForm(); err != nil {
		return err
	}

	key := h.Request().PostForm.Get("key")
	if key == "" {
		return nil
	}

	if !h.config.Has(key) {
		return fmt.Errorf("Variable %s not found", key)
	}

	if !h.config.IsEditable(key) {
		return fmt.Errorf("Variable %s isn't editable", key)
	}

	currentValue := h.config.Get(key)

	newValue := h.Request().PostForm.Get("value")
	err := h.config.Set(key, newValue)

	h.logger.Infof("Change value for %s with %v to %v", key, currentValue, newValue)

	return err
}

func (h *ConfigHandler) Handle() {
	var err error

	if h.IsPost() {
		err = h.saveNewValue()

		if err == nil {
			h.Redirect(h.Request().RequestURI, http.StatusFound)
			return
		}
	}

	h.SetView("dashboard", "config")
	h.SetVar("variables", h.config.GetAllVariables())
	h.SetVar("error", err)
}
