package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/annotations/storage"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logger"
)

type Component struct {
	mutex sync.RWMutex

	application shadow.Application
	config      config.Component
	logger      logger.Logger

	storages map[string]annotations.Storage
}

func (c *Component) Name() string {
	return annotations.ComponentName
}

func (c *Component) Version() string {
	return annotations.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: i18n.ComponentName,
		},
		{
			Name: logger.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.storages = make(map[string]annotations.Storage, 0)

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.Name(), c.application)

	c.initStorageGrafana()

	return nil
}

func (c *Component) Create(annotation annotations.Annotation) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.storages) == 0 {
		return errors.New("Storage not init")
	}

	for name, s := range c.storages {
		if err := s.Create(annotation); err != nil {
			c.logger.Error("Send annotation failed", map[string]interface{}{
				"storage": name,
				"error":   err.Error(),
			})
		}
	}

	return nil
}

func (c *Component) CreateInStorages(annotation annotations.Annotation, names []string) error {
	c.mutex.RLock()
	l := len(c.storages)
	c.mutex.RUnlock()

	if l == 0 {
		return errors.New("Storage not init")
	}

	for _, name := range names {
		c.mutex.RLock()
		s, ok := c.storages[name]
		c.mutex.RUnlock()

		if ok {
			if err := s.Create(annotation); err != nil {
				c.logger.Error("Send annotation failed", map[string]interface{}{
					"storage": name,
					"error":   err.Error(),
				})
			}
		}
	}

	return nil
}

func (c *Component) AddStorage(id string, s annotations.Storage) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.storages[id]; ok {
		return fmt.Errorf("Storage %s already exists", id)
	}

	c.storages[id] = s
	c.logger.Debugf("Added annotations storage %s", id)
	return nil
}

func (c *Component) RemoveStorage(id string) {
	c.mutex.Lock()
	delete(c.storages, id)
	c.logger.Debugf("Removed annotations storage %s", id)
	c.mutex.Unlock()
}

func (c *Component) initStorageGrafana() {
	c.RemoveStorage(annotations.StorageGrafana)

	if !c.config.Bool(annotations.ConfigStorageGrafanaEnabled) {
		return
	}

	var dashboards []int64

	for _, id := range strings.Split(c.config.String(annotations.ConfigStorageGrafanaDashboards), ",") {
		if value, err := strconv.ParseInt(id, 10, 0); err == nil {
			dashboards = append(dashboards, value)
		}
	}

	s := storage.NewGrafana(
		c.config.String(annotations.ConfigStorageGrafanaAddress),
		c.config.String(annotations.ConfigStorageGrafanaApiKey),
		c.config.String(annotations.ConfigStorageGrafanaUsername),
		c.config.String(annotations.ConfigStorageGrafanaPassword),
		dashboards,
		c.logger)

	c.AddStorage(annotations.StorageGrafana, s)
}
