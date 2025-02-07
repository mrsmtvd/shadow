package internal

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/annotations"
	"github.com/mrsmtvd/shadow/components/annotations/storage"
	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/i18n"
	"github.com/mrsmtvd/shadow/components/logging"
)

type Component struct {
	mutex sync.RWMutex

	config config.Component
	logger logging.Logger

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
			Name: logging.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.storages = make(map[string]annotations.Storage)

	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	c.logger = logging.DefaultLazyLogger(c.Name())

	<-a.ReadyComponent(config.ComponentName)

	ready <- struct{}{}

	c.initStorageGrafana()

	return nil
}

func (c *Component) Create(annotation annotations.Annotation) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.storages) == 0 {
		return errors.New("storage not init")
	}

	for name, s := range c.storages {
		if err := s.Create(annotation); err != nil {
			c.logger.Error("Send annotation failed",
				"storage", name,
				"error", err.Error(),
			)
		}
	}

	return nil
}

func (c *Component) CreateInStorages(annotation annotations.Annotation, names []string) error {
	c.mutex.RLock()
	l := len(c.storages)
	c.mutex.RUnlock()

	if l == 0 {
		return errors.New("storage not init")
	}

	for _, name := range names {
		c.mutex.RLock()
		s, ok := c.storages[name]
		c.mutex.RUnlock()

		if ok {
			if err := s.Create(annotation); err != nil {
				c.logger.Error("Send annotation failed",
					"storage", name,
					"error", err.Error(),
				)
			}
		}
	}

	return nil
}

func (c *Component) AddStorage(id string, s annotations.Storage) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.storages[id]; ok {
		return errors.New("storage " + id + " already exists")
	}

	c.storages[id] = s

	return nil
}

func (c *Component) RemoveStorage(id string) {
	c.mutex.Lock()
	delete(c.storages, id)
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
		c.config.String(annotations.ConfigStorageGrafanaAPIKey),
		c.config.String(annotations.ConfigStorageGrafanaUsername),
		c.config.String(annotations.ConfigStorageGrafanaPassword),
		dashboards,
		&logger{c.logger})

	_ = c.AddStorage(annotations.StorageGrafana, s)
}
