package system

import (
	"fmt"
	"net/http"

	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
	"github.com/kihamo/shadow/service/frontend"
)

type ConfigHandler struct {
	frontend.AbstractFrontendHandler

	config *config.Resource
	logger logger.Logger
}

func (h *ConfigHandler) saveNewValue() error {
	if err := h.Input.ParseForm(); err != nil {
		return err
	}

	key := h.Input.PostForm.Get("key")
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

	newValue := h.Input.PostForm.Get("value")
	err := h.config.Set(key, newValue)

	h.logger.Infof("Change value for %s with %v to %v", key, currentValue, newValue)

	return err
}

func (h *ConfigHandler) Handle() {
	var err error

	if h.IsPost() {
		err = h.saveNewValue()

		if err == nil {
			h.Redirect(h.Input.RequestURI, http.StatusFound)
			return
		}
	}

	h.SetTemplate("config.tpl.html")
	h.SetPageTitle("Configuration")
	h.SetPageHeader("Configuration")
	h.SetVar("Variables", h.config.GetAllVariables())
	h.SetVar("Error", err)
}
