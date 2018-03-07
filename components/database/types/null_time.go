package types

import (
	"github.com/go-gorp/gorp"
	pb "github.com/golang/protobuf/ptypes/timestamp"
)

type NullTime struct {
	gorp.NullTime
}

func (t *NullTime) Proto() *pb.Timestamp {
	if !t.Valid {
		return nil
	}

	return Time{
		Time: t.Time,
	}.Proto()
}
