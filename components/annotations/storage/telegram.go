package storage

import (
	"fmt"
	"time"

	"github.com/kihamo/shadow/components/annotations"
	"gopkg.in/telegram-bot-api.v4"
)

type Telegram struct {
	bot   *tgbotapi.BotAPI
	chats []int64
}

func NewTelegram(token string, chats []int64, debug bool) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = debug

	return &Telegram{
		bot:   bot,
		chats: chats,
	}, nil
}

func (s *Telegram) Create(annotation annotations.Annotation) (err error) {
	msg := tgbotapi.NewMessage(-1, fmt.Sprintf("*%s*\n%s\nStart at %s", annotation.Title(), annotation.Text(), annotation.Time().Format(time.RFC822Z)))
	msg.ParseMode = tgbotapi.ModeMarkdown

	for _, chatId := range s.chats {
		msg.ChatID = chatId

		if _, err := s.bot.Send(msg); err != nil {
			break
		}
	}

	return err
}
