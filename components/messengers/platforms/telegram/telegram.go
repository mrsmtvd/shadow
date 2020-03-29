package telegram

import (
	"fmt"
	"io"
	"strconv"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Telegram struct {
	bot *tgbotapi.BotAPI
}

func New(token string, debug bool) (platform *Telegram, err error) {
	platform = &Telegram{}

	platform.bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	platform.bot.Debug = debug

	return platform, nil
}

func (p *Telegram) chatID(to string) (int64, error) {
	chatID, err := strconv.Atoi(to)
	if err != nil {
		return -1, err
	}

	return int64(chatID), err
}

func (p *Telegram) Me() (tgbotapi.User, error) {
	return p.bot.GetMe()
}

func (p *Telegram) SendMessage(to, message string) error {
	chatID, err := p.chatID(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, message)
	_, err = p.bot.Send(msg)

	return err
}

func (p *Telegram) SendPhoto(to, name string, file io.Reader) error {
	chatID, err := p.chatID(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileReader{
		Name:   name,
		Reader: file,
		Size:   -1,
	})
	msg.Caption = name

	_, err = p.bot.Send(msg)

	return err
}

func (p *Telegram) SendMessageRaw(msg tgbotapi.Chattable) (tgbotapi.Message, error) {
	return p.bot.Send(msg)
}

func (p *Telegram) RegisterWebHook(link fmt.Stringer, cert string) error {
	var config tgbotapi.WebhookConfig

	if cert != "" {
		config = tgbotapi.NewWebhookWithCert(link.String(), cert)
	} else {
		config = tgbotapi.NewWebhook(link.String())
	}

	_, err := p.bot.SetWebhook(config)

	return err
}

func (p *Telegram) UnregisterWebHook() error {
	_, err := p.bot.RemoveWebhook()
	return err
}

func (p *Telegram) WebHook() (tgbotapi.WebhookInfo, error) {
	return p.bot.GetWebhookInfo()
}
