package types

import (
	"github.com/go-gorp/gorp"
	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
)

type NullTime struct {
	gorp.NullTime
}

func (n *NullTime) Proto() *tspb.Timestamp {
	if !n.Valid {
		return nil
	}

	p, _ := ptypes.TimestampProto(n.Time)
	return p
}
