package annotations

import (
	"fmt"
	"time"

	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/messengers/platforms/telegram"
	"gopkg.in/telegram-bot-api.v4"
)

type Telegram struct {
	messenger *telegram.Telegram
	chats     []int64
}

func NewTelegram(messenger *telegram.Telegram, chats []int64) *Telegram {
	return &Telegram{
		messenger: messenger,
		chats:     chats,
	}
}

func (s *Telegram) Create(annotation annotations.Annotation) (err error) {
	msg := tgbotapi.NewMessage(-1, fmt.Sprintf("*%s*\n%s\nStart at %s", annotation.Title(), annotation.Text(), annotation.Time().Format(time.RFC822Z)))
	msg.ParseMode = tgbotapi.ModeMarkdown

	for _, chatId := range s.chats {
		msg.ChatID = chatId

		if _, err := s.messenger.SendMessageRaw(msg); err != nil {
			break
		}
	}

	return err
}
