package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/kihamo/shadow/components/messengers/platforms/telegram"
)

// TODO: listen debug mode

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			messengers.ConfigTelegramEnabled,
			config.ValueTypeBool,
			false,
			"Enabled",
			true,
			"Telegram",
			nil,
			nil),
		config.NewVariable(
			messengers.ConfigTelegramToken,
			config.ValueTypeString,
			nil,
			"Token",
			true,
			"Telegram",
			nil,
			nil),
		config.NewVariable(
			messengers.ConfigTelegramUpdatesEnabled,
			config.ValueTypeBool,
			false,
			"Enabled updates",
			true,
			"Telegram",
			nil,
			nil),
		config.NewVariable(
			messengers.ConfigTelegramWebHookEnabled,
			config.ValueTypeBool,
			false,
			"Enabled web hooks",
			true,
			"Telegram",
			nil,
			nil),
		config.NewVariable(
			messengers.ConfigTelegramAnnotationsStorageEnabled,
			config.ValueTypeBool,
			false,
			"Enabled",
			true,
			"Telegram annotations storage",
			nil,
			nil),
		config.NewVariable(
			messengers.ConfigTelegramAnnotationsStorageChats,
			config.ValueTypeString,
			nil,
			"Chats ID",
			true,
			"Telegram annotations storage",
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a ID",
			}),
		config.NewVariable(
			messengers.ConfigBaseURL,
			config.ValueTypeString,
			nil,
			"Base URL for web hooks",
			true,
			"Others",
			nil,
			nil),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher(messengers.ComponentName, []string{
			config.ConfigDebug,
			messengers.ConfigTelegramEnabled,
			messengers.ConfigTelegramToken,
		}, c.watchTelegramWebEnabled),
		config.NewWatcher(messengers.ComponentName, []string{
			messengers.ConfigTelegramAnnotationsStorageEnabled,
			messengers.ConfigTelegramAnnotationsStorageChats,
		}, c.watchTelegramAnnotationsStorage),
		config.NewWatcher(messengers.ComponentName, []string{
			messengers.ConfigTelegramWebHookEnabled,
		}, c.watchTelegramWebHookEnabled),
	}
}

func (c *Component) watchTelegramWebEnabled(_ string, _ interface{}, _ interface{}) {
	c.initTelegram()
	c.initAnnotationsStorageTelegram()
}

func (c *Component) watchTelegramAnnotationsStorage(_ string, _ interface{}, _ interface{}) {
	c.initAnnotationsStorageTelegram()
}

func (c *Component) watchTelegramWebHookEnabled(_ string, newValue interface{}, _ interface{}) {
	messenger := c.Messenger(messengers.MessengerTelegram)
	if messenger == nil {
		return
	}

	if t, ok := messenger.(*telegram.Telegram); ok {
		c.initTelegramWebHook(t, newValue.(bool))
	}
}
