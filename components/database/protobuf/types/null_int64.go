package types

import (
	"database/sql"
)

type NullInt64 struct {
	sql.NullInt64
}

func (n *NullInt64) Proto() int64 {
	if !n.Valid {
		return 0
	}

	return n.Int64
}
