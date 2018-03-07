package types

import (
	"database/sql"
)

type NullString struct {
	sql.NullString
}

func (t *NullString) Proto() string {
	if !t.Valid {
		return ""
	}

	return t.String
}
