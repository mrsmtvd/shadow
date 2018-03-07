package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	pb "github.com/golang/protobuf/ptypes/struct"
	"github.com/kihamo/gotypes"
	"github.com/kihamo/shadow/components/grpc"
)

type Struct map[string]interface{}

func (t Struct) ToMap() map[string]interface{} {
	return t
}

func (t Struct) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (t *Struct) Scan(value interface{}) error {
	input := map[string]interface{}{}

	switch v := value.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &input); err != nil {
			return err
		}

	default:
		converter := gotypes.NewConverter(value, &input)

		if !converter.Valid() {
			return errors.New("Scan failed")
		}
	}

	*t = input
	return nil
}

func (t *Struct) Proto() *pb.Struct {
	return grpc.ConvertMapStringInterfaceToStructProto(t.ToMap())
}
