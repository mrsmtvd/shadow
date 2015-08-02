package shadow

import (
	"strconv"
)

func ToBool(in interface{}) bool {
	var castIn bool

	switch v := in.(type) {
	case bool:
		castIn = v

	case string:
		castIn = v != "" && v != "0" && v != "false"

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		castIn = v != 0
	}

	return castIn
}

func ToString(in interface{}) string {
	var castIn string

	switch v := in.(type) {
	case string:
		castIn = v

	case bool:
		castIn = strconv.FormatBool(v)

	case int:
		castIn = strconv.Itoa(v)

	case int8:
		castIn = strconv.FormatInt(int64(v), 10)

	case int16:
		castIn = strconv.FormatInt(int64(v), 10)

	case int32:
		castIn = strconv.FormatInt(int64(v), 10)

	case int64:
		castIn = strconv.FormatInt(v, 10)

	case uint:
		castIn = strconv.FormatUint(uint64(v), 10)

	case uint8:
		castIn = strconv.FormatUint(uint64(v), 10)

	case uint16:
		castIn = strconv.FormatUint(uint64(v), 10)

	case uint32:
		castIn = strconv.FormatUint(uint64(v), 10)

	case uint64:
		castIn = strconv.FormatUint(v, 10)

	case float32:
		castIn = strconv.FormatFloat(float64(v), 'f', 6, 32)

	case float64:
		castIn = strconv.FormatFloat(v, 'f', 6, 64)

	default:
		if cast, ok := in.(string); ok {
			castIn = cast
		}
	}

	return castIn
}

func ToUint(in interface{}) uint {
	return uint(ToInt64(in))
}

func ToUint8(in interface{}) uint8 {
	return uint8(ToUint64(in))
}

func ToUint16(in interface{}) uint16 {
	return uint16(ToUint64(in))
}

func ToUint32(in interface{}) uint32 {
	return uint32(ToUint64(in))
}

func ToUint64(in interface{}) uint64 {
	return uint64(ToInt64(in))
}

func ToInt(in interface{}) int {
	return int(ToInt64(in))
}

func ToInt8(in interface{}) int8 {
	return int8(ToInt64(in))
}

func ToInt16(in interface{}) int16 {
	return int16(ToInt64(in))
}

func ToInt32(in interface{}) int32 {
	return int32(ToInt64(in))
}

func ToInt64(in interface{}) int64 {
	var castIn int64

	switch v := in.(type) {
	case string:
		t, _ := strconv.Atoi(v)
		castIn = int64(t)

	case int64:
		castIn = v

	case int:
		castIn = int64(v)

	case int8:
		castIn = int64(v)

	case int16:
		castIn = int64(v)

	case int32:
		castIn = int64(v)

	case uint:
		castIn = int64(v)

	case uint8:
		castIn = int64(v)

	case uint16:
		castIn = int64(v)

	case uint32:
		castIn = int64(v)

	case uint64:
		castIn = int64(v)

	case float32:
		castIn = int64(v)

	case float64:
		castIn = int64(v)
	}

	return castIn
}

func ToFloat32(in interface{}) float32 {
	return float32(ToFloat64(in))
}

func ToFloat64(in interface{}) float64 {
	var castIn float64

	switch v := in.(type) {
	case string:
		castIn, _ = strconv.ParseFloat(v, 64)

	case float64:
		castIn = v

	case int:
		castIn = float64(v)

	case int8:
		castIn = float64(v)

	case int16:
		castIn = float64(v)

	case int32:
		castIn = float64(v)

	case int64:
		castIn = float64(v)

	case uint:
		castIn = float64(v)

	case uint8:
		castIn = float64(v)

	case uint16:
		castIn = float64(v)

	case uint32:
		castIn = float64(v)

	case uint64:
		castIn = float64(v)

	case float32:
		castIn = float64(v)
	}

	return castIn
}
