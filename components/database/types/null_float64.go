package types

import (
	"database/sql"
)

type NullFloat64 struct {
	sql.NullFloat64
}

func (t *NullFloat64) Proto() float64 {
	if !t.Valid {
		return 0
	}

	return t.Float64
}
