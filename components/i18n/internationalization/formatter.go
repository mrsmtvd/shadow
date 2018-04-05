package internationalization

import (
	"fmt"
)

func Format(message string, format ...interface{}) string {
	if len(format) > 0 {
		return fmt.Sprintf(message, format...)
	}

	return message
}
