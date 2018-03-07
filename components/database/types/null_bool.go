package types

import (
	"database/sql"
)

type NullBool struct {
	sql.NullBool
}

func (t *NullBool) Proto() bool {
	if !t.Valid {
		return false
	}

	return t.Bool
}
