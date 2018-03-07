package types

import (
	"database/sql"
)

type NullInt64 struct {
	sql.NullInt64
}

func (t *NullInt64) Proto() int64 {
	if !t.Valid {
		return 0
	}

	return t.Int64
}
