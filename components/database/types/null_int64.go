package types

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func (t *NullInt64) ToInt64() *int64 {
	if !t.Valid {
		return nil
	}

	return &t.Int64
}

func (t *NullInt64) Proto() int64 {
	if !t.Valid {
		return 0
	}

	return t.Int64
}

func (t *NullInt64) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.Int64)
	}

	return json.Marshal(nil)
}

func (t *NullInt64) UnmarshalJSON(data []byte) error {
	var j *int64

	err := json.Unmarshal(data, &j)
	if err == nil && j != nil {
		t.Int64, t.Valid = *j, true
	} else {
		t.Int64, t.Valid = 0, false
	}

	return err
}
