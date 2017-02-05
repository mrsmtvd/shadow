package alerts

import (
	"html/template"
	"time"

	"github.com/kihamo/shadow"
)

type Alert struct {
	title   string
	message string
	icon    string
	date    time.Time
}

func NewAlert(title, message, icon string, date time.Time) *Alert {
	return &Alert{
		icon:    icon,
		title:   title,
		message: message,
		date:    date,
	}
}

func (a *Alert) GetTitle() string {
	return a.title
}

func (a *Alert) GetMessage() string {
	return a.message
}

func (a *Alert) GetMessageAsHTML() template.HTML {
	return template.HTML(a.message)
}

func (a *Alert) GetIcon() string {
	return a.icon
}

func (a *Alert) GetDate() time.Time {
	return a.date
}

func (a *Alert) GetDateAsMessage() string {
	return shadow.DateSinceAsMessage(a.date)
}
