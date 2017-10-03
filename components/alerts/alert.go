package alerts

import (
	"html/template"
	"time"
)

type Alert interface {
	GetTitle() string
	GetMessage() string
	GetMessageAsHTML() template.HTML
	GetIcon() string
	GetDate() time.Time
	GetDateAsMessage() string
}
