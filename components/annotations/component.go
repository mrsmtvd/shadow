package annotations

import (
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	AddStorage(string, Storage) error
	RemoveStorage(string)
	Create(Annotation) error
	CreateInStorages(Annotation, []string) error
}
