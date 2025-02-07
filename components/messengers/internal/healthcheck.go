package internal

import (
	"errors"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/messengers"
	"github.com/mrsmtvd/shadow/components/messengers/platforms/telegram"
)

const (
	requestTimeout  = time.Second * 3
	requestInterval = time.Minute
)

func (c *Component) LivenessCheck() map[string]dashboard.HealthCheck {
	return map[string]dashboard.HealthCheck{
		"telegram_webhook": c.TelegramWebHookCheck(),
	}
}

func (c *Component) ReadinessCheck() map[string]dashboard.HealthCheck {
	return map[string]dashboard.HealthCheck{
		"telegram": c.TelegramCheck(),
	}
}

func (c *Component) TelegramCheck() dashboard.HealthCheck {
	return healthcheck.Async(healthcheck.Timeout(func() error {
		if !c.config.Bool(messengers.ConfigTelegramEnabled) && !c.config.Bool(messengers.ConfigTelegramWebHookEnabled) && !c.config.Bool(messengers.ConfigTelegramAnnotationsStorageEnabled) {
			return nil
		}

		messenger := c.Messenger(messengers.MessengerTelegram)
		if messenger == nil {
			return errors.New("telegram messenger isn't initialization")
		}

		tg, ok := messenger.(*telegram.Telegram)
		if !ok {
			return errors.New("telegram messenger isn't initialization")
		}

		_, err := tg.Me()
		return err
	}, requestTimeout), requestInterval)
}

func (c *Component) TelegramWebHookCheck() dashboard.HealthCheck {
	return healthcheck.Async(healthcheck.Timeout(func() error {
		if !c.config.Bool(messengers.ConfigTelegramWebHookEnabled) {
			return nil
		}

		messenger := c.Messenger(messengers.MessengerTelegram)
		if messenger == nil {
			return errors.New("telegram messenger isn't initialization")
		}

		tg, ok := messenger.(*telegram.Telegram)
		if !ok {
			return errors.New("telegram messenger isn't initialization")
		}

		_, err := tg.WebHook()
		return err
	}, requestTimeout), requestInterval)
}
