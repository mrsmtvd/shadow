package types

import (
	"database/sql"
)

type NullBool struct {
	sql.NullBool
}

func (n *NullBool) Proto() bool {
	if !n.Valid {
		return false
	}

	return n.Bool
}
