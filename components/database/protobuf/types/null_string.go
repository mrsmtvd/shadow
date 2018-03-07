package types

import (
	"database/sql"
)

type NullString struct {
	sql.NullString
}

func (n *NullString) Proto() string {
	if !n.Valid {
		return ""
	}

	return n.String
}
