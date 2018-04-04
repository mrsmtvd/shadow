package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/kihamo/shadow/components/messengers/platforms/telegram"
)

// TODO: listen debug mode

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(messengers.ConfigTelegramEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Telegram").
			WithEditable(true),
		config.NewVariable(messengers.ConfigTelegramToken, config.ValueTypeString).
			WithUsage("Token").
			WithGroup("Telegram").
			WithEditable(true),
		config.NewVariable(messengers.ConfigTelegramUpdatesEnabled, config.ValueTypeBool).
			WithUsage("Enabled updates").
			WithGroup("Telegram").
			WithEditable(true),
		config.NewVariable(messengers.ConfigTelegramWebHookEnabled, config.ValueTypeBool).
			WithUsage("Enabled web hooks").
			WithGroup("Telegram").
			WithEditable(true),
		config.NewVariable(messengers.ConfigTelegramAnnotationsStorageEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Telegram annotations storage").
			WithEditable(true),
		config.NewVariable(messengers.ConfigTelegramAnnotationsStorageChats, config.ValueTypeString).
			WithUsage("Chats ID").
			WithGroup("Telegram annotations storage").
			WithEditable(true).
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a ID"}),
		config.NewVariable(messengers.ConfigBaseURL, config.ValueTypeString).
			WithUsage("Base URL for web hooks").
			WithEditable(true),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{
			config.ConfigDebug,
			messengers.ConfigTelegramEnabled,
			messengers.ConfigTelegramToken,
		}, c.watchTelegramWebEnabled),
		config.NewWatcher([]string{
			messengers.ConfigTelegramAnnotationsStorageEnabled,
			messengers.ConfigTelegramAnnotationsStorageChats,
		}, c.watchTelegramAnnotationsStorage),
		config.NewWatcher([]string{
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
