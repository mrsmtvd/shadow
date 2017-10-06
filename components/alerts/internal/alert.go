package internal

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

func (a *Alert) Title() string {
	return a.title
}

func (a *Alert) Message() string {
	return a.message
}

func (a *Alert) MessageAsHTML() template.HTML {
	return template.HTML(a.message)
}

func (a *Alert) Icon() string {
	return a.icon
}

func (a *Alert) Date() time.Time {
	return a.date
}

func (a *Alert) DateAsMessage() string {
	return shadow.DateSinceAsMessage(a.date)
}
