package types

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (t *NullString) Proto() string {
	if !t.Valid {
		return ""
	}

	return t.String
}

func (t *NullString) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.String)
	}

	return json.Marshal(nil)
}

func (t *NullString) UnmarshalJSON(data []byte) error {
	var j *string

	err := json.Unmarshal(data, &j)
	if err == nil && j != nil {
		t.String, t.Valid = *j, true
	} else {
		t.String, t.Valid = "", false
	}

	return err
}
