package types

import (
	"time"

	"github.com/go-gorp/gorp"
	"github.com/golang/protobuf/ptypes"
	pb "github.com/golang/protobuf/ptypes/timestamp"
)

type NullTime struct {
	gorp.NullTime
}

func (t *NullTime) Scan(value interface{}) error {
	if value == nil {
		t.Valid, t.Time = false, time.Time{}
		return nil
	}

	return t.NullTime.Scan(value)
}

func (t *NullTime) Proto() *pb.Timestamp {
	if !t.Valid {
		return nil
	}

	p, _ := ptypes.TimestampProto(t.Time)
	return p
}
