package types

import (
	"time"

	"encoding/json"

	"github.com/go-gorp/gorp"
	"github.com/golang/protobuf/ptypes"
	pb "github.com/golang/protobuf/ptypes/timestamp"
)

type NullTime struct {
	gorp.NullTime
}

func (t *NullTime) Scan(value interface{}) error {
	if value == nil {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

	return t.NullTime.Scan(value)
}

func (t *NullTime) ToTime() *time.Time {
	if !t.Valid {
		return nil
	}

	return &t.Time
}

func (t *NullTime) Proto() *pb.Timestamp {
	if !t.Valid {
		return nil
	}

	p, _ := ptypes.TimestampProto(t.Time)

	return p
}

func (t *NullTime) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.Time)
	}

	return json.Marshal(nil)
}

func (t *NullTime) UnmarshalJSON(data []byte) error {
	var j *time.Time

	err := json.Unmarshal(data, &j)
	if err == nil && j != nil {
		t.Time, t.Valid = *j, true
	} else {
		t.Time, t.Valid = time.Time{}, false
	}

	return err
}
