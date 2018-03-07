package types

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/golang/protobuf/ptypes/timestamp"
)

const (
	timeFormat = "2006-01-02 15:04:05.999999"
)

type Time struct {
	time.Time
}

func (t Time) ToTime() time.Time {
	return t.Time
}

func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t *Time) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
	case []byte:
		t.Time, err = parseDateTime(string(v), time.UTC)
	case string:
		t.Time, err = parseDateTime(v, time.UTC)
	default:
		err = fmt.Errorf("Can't convert %T to time.Time", value)
	}

	return err
}

func (t *Time) Proto() *pb.Timestamp {
	p, _ := ptypes.TimestampProto(t.Time)
	return p
}

func parseDateTime(str string, loc *time.Location) (t time.Time, err error) {
	base := "0000-00-00 00:00:00.0000000"
	switch len(str) {
	case 10, 19, 21, 22, 23, 24, 25, 26: // up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
		if str == base[:len(str)] {
			return
		}
		t, err = time.Parse(timeFormat[:len(str)], str)
	default:
		err = fmt.Errorf("invalid time string: %s", str)
		return
	}

	// Adjust location
	if err == nil && loc != time.UTC {
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		t, err = time.Date(y, mo, d, h, mi, s, t.Nanosecond(), loc), nil
	}

	return
}
