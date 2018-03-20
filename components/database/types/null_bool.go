package types

import (
	"database/sql"
	"encoding/json"
)

type NullBool struct {
	sql.NullBool
}

func (t *NullBool) ToBool() *bool {
	if !t.Valid {
		return nil
	}

	return &t.Bool
}

func (t *NullBool) Proto() bool {
	if !t.Valid {
		return false
	}

	return t.Bool
}

func (t *NullBool) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.Bool)
	}

	return json.Marshal(nil)
}

func (t *NullBool) UnmarshalJSON(data []byte) error {
	var j *bool

	err := json.Unmarshal(data, &j)
	if err == nil && j != nil {
		t.Bool, t.Valid = *j, true
	} else {
		t.Bool, t.Valid = false, false
	}

	return err
}
