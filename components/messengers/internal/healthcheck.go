package internal

import (
	"errors"
	"time"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/kihamo/shadow/components/messengers/platforms/telegram"
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
	return func() error {
		if !c.config.Bool(messengers.ConfigTelegramEnabled) && !c.config.Bool(messengers.ConfigTelegramWebHookEnabled) && !!c.config.Bool(messengers.ConfigTelegramAnnotationsStorageEnabled) {
			return nil
		}

		messenger := c.Messenger(messengers.MessengerTelegram)
		if messenger == nil {
			return errors.New("Telegram messenger isn't initialization")
		}

		tg, ok := messenger.(*telegram.Telegram)
		if !ok {
			return errors.New("Telegram messenger isn't initialization")
		}

		result := make(chan error, 1)
		go func() {
			_, err := tg.Me()
			result <- err
		}()

		select {
		case err := <-result:
			return err
		case <-time.After(time.Second * 3):
			return errors.New("Request timeout")
		}
	}
}

func (c *Component) TelegramWebHookCheck() dashboard.HealthCheck {
	return func() error {
		if !c.config.Bool(messengers.ConfigTelegramWebHookEnabled) {
			return nil
		}

		messenger := c.Messenger(messengers.MessengerTelegram)
		if messenger == nil {
			return errors.New("Telegram messenger isn't initialization")
		}

		tg, ok := messenger.(*telegram.Telegram)
		if !ok {
			return errors.New("Telegram messenger isn't initialization")
		}

		result := make(chan error, 1)
		go func() {
			_, err := tg.WebHook()
			result <- err
		}()

		select {
		case err := <-result:
			return err
		case <-time.After(time.Second * 3):
			return errors.New("Request timeout")
		}
	}
}
