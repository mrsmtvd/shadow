package annotations

import (
	"time"
)

type Annotation interface {
	Time() time.Time
	TimeEnd() *time.Time
	Title() string
	Text() string
	Tags() []string
}

type AnnotationBase struct {
	timeStart time.Time
	timeEnd   *time.Time
	title     string
	text      string
	tags      []string
}

func NewAnnotation(title, text string, tags []string, timeStart *time.Time, timeEnd *time.Time) *AnnotationBase {
	a := &AnnotationBase{
		timeEnd: timeEnd,
		title:   title,
		text:    text,
		tags:    tags,
	}

	if timeStart == nil {
		a.timeStart = time.Now()
	} else {
		a.timeStart = *timeStart
	}

	return a
}

func (a *AnnotationBase) Time() time.Time {
	return a.timeStart
}

func (a *AnnotationBase) TimeEnd() *time.Time {
	return a.timeEnd
}

func (a *AnnotationBase) Title() string {
	return a.title
}

func (a *AnnotationBase) Text() string {
	return a.text
}

func (a *AnnotationBase) Tags() []string {
	return a.tags
}
