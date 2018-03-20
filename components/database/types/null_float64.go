package types

import (
	"database/sql"
	"encoding/json"
)

type NullFloat64 struct {
	sql.NullFloat64
}

func (t *NullFloat64) ToFloat64() *float64 {
	if !t.Valid {
		return nil
	}

	return &t.Float64
}

func (t *NullFloat64) Proto() float64 {
	if !t.Valid {
		return 0
	}

	return t.Float64
}

func (t *NullFloat64) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.Float64)
	}

	return json.Marshal(nil)
}

func (t *NullFloat64) UnmarshalJSON(data []byte) error {
	var j *float64

	err := json.Unmarshal(data, &j)
	if err == nil && j != nil {
		t.Float64, t.Valid = *j, true
	} else {
		t.Float64, t.Valid = 0, false
	}

	return err
}
