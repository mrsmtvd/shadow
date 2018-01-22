package internal

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/annotations/storage"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

type Component struct {
	mutex sync.RWMutex

	application shadow.Application
	config      config.Component
	logger      logger.Logger

	storage annotations.Storage
}

func (c *Component) GetName() string {
	return annotations.ComponentName
}

func (c *Component) GetVersion() string {
	return annotations.ComponentVersion
}

func (c *Component) GetDependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		}, {
			Name: logger.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.GetName(), c.application)
	c.initStorageGrafana()

	return nil
}

func (c *Component) Create(annotation annotations.Annotation) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.storage == nil {
		return errors.New("Storage not init")
	}

	return c.storage.Create(annotation)
}

func (c *Component) initStorageGrafana() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var dashboards []int64

	for _, id := range strings.Split(c.config.GetString(annotations.ConfigStorageGrafanaDashboards), ",") {
		if value, err := strconv.ParseInt(id, 10, 0); err == nil {
			dashboards = append(dashboards, value)
		}
	}

	c.storage = storage.NewGrafana(
		c.config.GetString(annotations.ConfigStorageGrafanaAddress),
		c.config.GetString(annotations.ConfigStorageGrafanaApiKey),
		c.config.GetString(annotations.ConfigStorageGrafanaUsername),
		c.config.GetString(annotations.ConfigStorageGrafanaPassword),
		dashboards,
		c.logger)
}
