package i18n

import (
	"io"
)

type HasI18n interface {
	I18n() map[string]io.ReadSeeker
}
