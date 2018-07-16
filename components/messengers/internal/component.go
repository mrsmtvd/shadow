package internal

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/messengers"
	storages "github.com/kihamo/shadow/components/messengers/internal/annotations"
	"github.com/kihamo/shadow/components/messengers/platforms/telegram"
)

type Component struct {
	mutex sync.RWMutex

	application shadow.Application
	annotations annotations.Component
	config      config.Component
	logger      logger.Logger
	messengers  map[string]messengers.Messenger
}

func (c *Component) Name() string {
	return messengers.ComponentName
}

func (c *Component) Version() string {
	return messengers.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name: annotations.ComponentName,
		},
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
	c.messengers = make(map[string]messengers.Messenger, 0)

	if a.HasComponent(annotations.ComponentName) {
		c.annotations = a.GetComponent(annotations.ComponentName).(annotations.Component)
	}

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.Name(), c.application)

	c.initTelegram()
	c.initAnnotationsStorageTelegram()

	return nil
}

func (c *Component) RegisterMessenger(id string, messenger messengers.Messenger) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.messengers[id]; ok {
		return fmt.Errorf("Storage %s already exists", id)
	}

	c.messengers[id] = messenger
	c.logger.Debugf("Registered messenger %s", id)
	return nil
}

func (c *Component) UnregisterMessenger(id string) {
	c.mutex.Lock()
	delete(c.messengers, id)
	c.logger.Debugf("Unregistered messenger %s", id)
	c.mutex.Unlock()
}

func (c *Component) Messenger(id string) messengers.Messenger {
	c.mutex.RLock()
	m := c.messengers[id]
	c.mutex.RUnlock()

	return m
}

func (c *Component) initTelegram() {
	c.UnregisterMessenger(messengers.MessengerTelegram)

	if !c.config.Bool(messengers.ConfigTelegramEnabled) {
		return
	}

	messenger, err := telegram.New(
		c.config.String(messengers.ConfigTelegramToken),
		c.config.Bool(config.ConfigDebug))
	if err != nil {
		c.logger.Error("Failed init telegram messenger", map[string]interface{}{
			"error":     err.Error(),
			"messenger": messengers.MessengerTelegram,
		})
		return
	}

	c.initTelegramWebHook(messenger, c.config.Bool(messengers.ConfigTelegramWebHookEnabled))

	c.RegisterMessenger(messengers.MessengerTelegram, messenger)
}

func (c *Component) initTelegramWebHook(messenger *telegram.Telegram, enabled bool) {
	if enabled {
		u, err := url.Parse(c.config.String(messengers.ConfigBaseURL))
		if err != nil {
			c.logger.Error("Failed register webhook for telegram messenger", map[string]interface{}{
				"error":     err.Error(),
				"messenger": messengers.MessengerTelegram,
			})
			return
		}

		messenger.RegisterWebHook(u, "")
	} else {
		err := messenger.UnregisterWebHook()
		if err != nil {
			c.logger.Error("Failed unregister webhook for telegram messenger", map[string]interface{}{
				"error":     err.Error(),
				"messenger": messengers.MessengerTelegram,
			})
			return
		}
	}
}

func (c *Component) initAnnotationsStorageTelegram() {
	if c.annotations == nil {
		return
	}

	c.annotations.RemoveStorage(messengers.MessengerTelegram)

	if !c.config.Bool(messengers.ConfigTelegramAnnotationsStorageEnabled) {
		return
	}

	messenger := c.Messenger(messengers.MessengerTelegram)
	if messenger == nil {
		return
	}

	t, ok := messenger.(*telegram.Telegram)
	if !ok {
		return
	}

	var chats []int64

	for _, id := range strings.Split(c.config.String(messengers.ConfigTelegramAnnotationsStorageChats), ",") {
		if value, err := strconv.ParseInt(id, 10, 0); err == nil {
			chats = append(chats, value)
		}
	}

	storage := storages.NewTelegram(t, chats)
	c.annotations.AddStorage(messengers.MessengerTelegram, storage)
}
