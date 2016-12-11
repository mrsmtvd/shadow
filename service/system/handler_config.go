package system

import (
	"fmt"

	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
	"github.com/kihamo/shadow/service/frontend"
)

type ConfigHandler struct {
	frontend.AbstractFrontendHandler

	config *config.Resource
	logger logger.Logger
}

func (h *ConfigHandler) saveNewValue(c *config.Resource) error {
	if err := h.Input.ParseForm(); err != nil {
		return err
	}

	key := h.Input.PostForm.Get("key")
	if key == "" {
		return nil
	}

	if !c.Has(key) {
		return fmt.Errorf("Variable %s not found", key)
	}

	if !c.IsEditable(key) {
		return fmt.Errorf("Variable %s isn't editable", key)
	}

	currentValue := c.Get(key)
	newValue := h.Input.PostForm.Get("value")
	err := c.Set(key, newValue)

	h.logger.Infof("Change value for %s with %v to %v", key, currentValue, newValue)

	return err
}

func (h *ConfigHandler) Handle() {
	resourceConfig, _ := h.Application.GetResource("config")
	config := resourceConfig.(*config.Resource)

	var err error

	if h.IsPost() {
		err = h.saveNewValue(config)
	}

	h.SetTemplate("config.tpl.html")
	h.SetPageTitle("Configuration")
	h.SetPageHeader("Configuration")
	h.SetVar("Variables", config.GetAllVariables())
	h.SetVar("Error", err)
}
