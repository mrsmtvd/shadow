package alerts

import (
	"html/template"
	"time"
)

type Alert interface {
	Title() string
	Message() string
	MessageAsHTML() template.HTML
	Icon() string
	Date() time.Time
	DateAsMessage() string
}
