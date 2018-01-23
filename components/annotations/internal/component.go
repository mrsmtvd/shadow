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
	"github.com/kihamo/shadow/components/logger"
)

type Component struct {
	mutex sync.RWMutex

	application shadow.Application
	config      config.Component
	logger      logger.Logger

	storages map[string]annotations.Storage
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
	c.storages = make(map[string]annotations.Storage, 0)

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.GetName(), c.application)

	c.initStorageGrafana()
	c.initStorageTelegram()

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
	return nil
}

func (c *Component) RemoveStorage(id string) {
	c.mutex.Lock()
	delete(c.storages, id)
	c.mutex.Unlock()
}

func (c *Component) initStorageGrafana() {
	c.RemoveStorage(annotations.StorageGrafana)

	if !c.config.GetBool(annotations.ConfigStorageGrafanaEnabled) {
		return
	}

	var dashboards []int64

	for _, id := range strings.Split(c.config.GetString(annotations.ConfigStorageGrafanaDashboards), ",") {
		if value, err := strconv.ParseInt(id, 10, 0); err == nil {
			dashboards = append(dashboards, value)
		}
	}

	s := storage.NewGrafana(
		c.config.GetString(annotations.ConfigStorageGrafanaAddress),
		c.config.GetString(annotations.ConfigStorageGrafanaApiKey),
		c.config.GetString(annotations.ConfigStorageGrafanaUsername),
		c.config.GetString(annotations.ConfigStorageGrafanaPassword),
		dashboards,
		c.logger)

	c.AddStorage(annotations.StorageGrafana, s)
}

func (c *Component) initStorageTelegram() {
	c.RemoveStorage(annotations.StorageTelegram)

	if !c.config.GetBool(annotations.ConfigStorageTelegramEnabled) {
		return
	}

	var chats []int64

	for _, id := range strings.Split(c.config.GetString(annotations.ConfigStorageTelegramChats), ",") {
		if value, err := strconv.ParseInt(id, 10, 0); err == nil {
			chats = append(chats, value)
		}
	}

	s, err := storage.NewTelegram(
		c.config.GetString(annotations.ConfigStorageTelegramToken),
		chats,
		c.config.GetBool(config.ConfigDebug))

	if err != nil {
		c.logger.Errorf("Telegram storage failed: %s", err.Error())
		return
	}

	c.AddStorage(annotations.StorageTelegram, s)
}
