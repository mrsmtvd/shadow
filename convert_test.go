package shadow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BoolTrue_ToBoolTrue(t *testing.T) {
	val := true

	result := ToBool(val)

	assert.True(t, result)
}

func Test_BoolFalse_ToBoolFalse(t *testing.T) {
	val := false

	result := ToBool(val)

	assert.False(t, result)
}

func Test_StringNonEmpty_ToBoolTrue(t *testing.T) {
	val := "123"

	result := ToBool(val)

	assert.True(t, result)
}

func Test_StringEmpty_ToBoolFalse(t *testing.T) {
	val := ""

	result := ToBool(val)

	assert.False(t, result)
}

func Test_StringZero_ToBoolFalse(t *testing.T) {
	val := "0"

	result := ToBool(val)

	assert.False(t, result)
}

func Test_StringFalse_ToBoolFalse(t *testing.T) {
	val := "false"

	result := ToBool(val)

	assert.False(t, result)
}

func Test_IntOne_ToBoolTrue(t *testing.T) {
	val := 1

	result := ToBool(val)

	assert.True(t, result)
}

func Test_IntZero_ToBoolFalse(t *testing.T) {
	val := 0

	result := ToBool(val)

	assert.False(t, result)
}

func Test_FloatOne_ToBoolTrue(t *testing.T) {
	val := 1.0

	result := ToBool(val)

	assert.True(t, result)
}

func Test_StringNonEmpty_ToString(t *testing.T) {
	val := "1"

	result := ToString(val)

	assert.Equal(t, "1", result)
}

func Test_BoolTrue_ToString(t *testing.T) {
	val := true

	result := ToString(val)

	assert.Equal(t, "true", result)
}

func Test_BoolFalse_ToString(t *testing.T) {
	val := false

	result := ToString(val)

	assert.Equal(t, "false", result)
}

func Test_IntOne_ToString(t *testing.T) {
	val := 1

	result := ToString(val)

	assert.Equal(t, "1", result)
}

func Test_IntZero_ToString(t *testing.T) {
	val := 0

	result := ToString(val)

	assert.Equal(t, "0", result)
}

func Test_Int8One_ToString(t *testing.T) {
	val := int8(1)

	result := ToString(val)

	assert.Equal(t, "1", result)
}

func Test_Int16One_ToString(t *testing.T) {
	val := int16(1)

	result := ToString(val)

	assert.Equal(t, "1", result)
}

func Test_Int32One_ToString(t *testing.T) {
	val := int32(1)

	result := ToString(val)

	assert.Equal(t, "1", result)
}

func Test_Int64One_ToString(t *testing.T) {
	val := int64(1)

	result := ToString(val)

	assert.Equal(t, "1", result)
}

func Test_UintOne_ToString(t *testing.T) {
	val := uint(1)

	result := ToString(val)

	assert.Equal(t, "1", result)
}

func Test_Uint8One_ToString(t *testing.T) {
	val := uint8(1)

	result := ToString(val)

	assert.Equal(t, "1", result)
}

func Test_Uint16One_ToString(t *testing.T) {
	val := uint16(1)

	result := ToString(val)

	assert.Equal(t, "1", result)
}

func Test_Uint32One_ToString(t *testing.T) {
	val := uint32(1)

	result := ToString(val)

	assert.Equal(t, "1", result)
}

func Test_Uint64One_ToString(t *testing.T) {
	val := uint64(1)

	result := ToString(val)

	assert.Equal(t, "1", result)
}

func Test_Float32One_ToString(t *testing.T) {
	val := float32(1)

	result := ToString(val)

	assert.Equal(t, "1.000000", result)
}

func Test_Float64One_ToString(t *testing.T) {
	val := float64(1)

	result := ToString(val)

	assert.Equal(t, "1.000000", result)
}
