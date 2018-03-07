package types

import (
	"database/sql"
)

type NullFloat64 struct {
	sql.NullFloat64
}

func (n *NullFloat64) Proto() float64 {
	if !n.Valid {
		return 0
	}

	return n.Float64
}
