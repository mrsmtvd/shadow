package mail

import (
	"github.com/mrsmtvd/shadow"
	"gopkg.in/gomail.v2"
)

type Component interface {
	shadow.Component

	Send(message *gomail.Message)
	SendAndReturn(message *gomail.Message) error
}
