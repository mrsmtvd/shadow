package shadow

import (
	"testing"
)

func TestConvertToBool(t *testing.T) {
	values := [][]interface{}{
		[]interface{}{true, true},
		[]interface{}{false, false},
		[]interface{}{"", false},
		[]interface{}{"0", false},
		[]interface{}{"00", true},
		[]interface{}{"false", false},
		[]interface{}{"true", true},
		[]interface{}{"1", true},
		[]interface{}{nil, false},
		[]interface{}{1, true},
		[]interface{}{0, false},
		[]interface{}{int8(1), true},
		[]interface{}{int16(1), true},
		[]interface{}{int32(1), true},
		[]interface{}{int64(1), true},
		[]interface{}{1.23, true},
		[]interface{}{float64(1.23), true},
		[]interface{}{float32(1.23), true},
	}

	for _, data := range values {
		r := ToBool(data[0])
		if r != data[1] {
			t.Fatalf("Can not convert to bool. %#v != %#v", data[0], data[1])
		}
	}
}
