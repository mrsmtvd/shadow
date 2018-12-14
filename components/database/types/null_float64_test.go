package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullFloat64_NewInstance_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}

	a.False(typ.Valid)
}

func TestNullFloat64_NewInstance_ValueZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}

	a.Equal(typ.Float64, float64(0))
}

func TestNullFloat64_NewInstance_ProtoReturnsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}

	a.Equal(typ.Proto(), float64(0))
}

func TestNullFloat64_NewInstance_MarshalJSONReturnsNullAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`null`))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanInteger_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(1)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanInteger_ValueInteger(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(1)

	a.Equal(typ.Float64, float64(1))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanInteger_ProtoReturnsInteger(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(1)

	a.Equal(typ.Proto(), float64(1))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanInteger_MarshalJSONReturnsIntegerAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	es := typ.Scan(1)
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`1`))
	a.Nil(e)
	a.Nil(es)
}

func TestNullFloat64_NewInstanceAndScanZero_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(0)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanZero_ValueZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(0)

	a.Equal(typ.Float64, float64(0))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanZero_ProtoReturnsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(0)

	a.Equal(typ.Proto(), float64(0))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanZero_MarshalJSONReturnsZeroAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	es := typ.Scan(0)
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`0`))
	a.Nil(e)
	a.Nil(es)
}

func TestNullFloat64_NewInstanceAndScanDecimal_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(1.23)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanDecimal_ValueDecimal(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(1.23)

	a.Equal(typ.Float64, float64(1.23))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanDecimal_ProtoReturnsDecimal(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(1.23)

	a.Equal(typ.Proto(), float64(1.23))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanDecimal_MarshalJSONReturnsDecimalAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	es := typ.Scan(1.23)
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`1.23`))
	a.Nil(e)
	a.Nil(es)
}

func TestNullFloat64_NewInstanceAndScanNil_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(nil)

	a.False(typ.Valid)
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanNil_ValueZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(nil)

	a.Equal(typ.Float64, float64(0))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanNil_ProtoReturnsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	e := typ.Scan(nil)

	a.Equal(typ.Float64, float64(0))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndScanNil_MarshalJSONReturnsNilAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullFloat64{}
	es := typ.Scan(nil)
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`null`))
	a.Nil(e)
	a.Nil(es)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONInteger_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`1`), &typ)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONInteger_ValueIsInteger(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`1`), &typ)

	a.Equal(typ.Float64, float64(1))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONInteger_ProtoReturnsInteger(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`1`), &typ)

	a.Equal(typ.Float64, float64(1))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONZero_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`0`), &typ)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONZero_ValueIsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`0`), &typ)

	a.Equal(typ.Float64, float64(0))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONZero_ProtoReturnsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`0`), &typ)

	a.Equal(typ.Proto(), float64(0))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONDecimal_ValidationIsDecimal(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`1.23`), &typ)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONDecimal_ValueIsDecimal(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`1.23`), &typ)

	a.Equal(typ.Float64, float64(1.23))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONDecimal_ProtoReturnsDecimal(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`1.23`), &typ)

	a.Equal(typ.Float64, float64(1.23))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONNil_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`null`), &typ)

	a.False(typ.Valid)
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONNil_ValueIsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`null`), &typ)

	a.Equal(typ.Float64, float64(0))
	a.Nil(e)
}

func TestNullFloat64_NewInstanceAndUnmarshalJSONNil_ProtoReturnsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullFloat64
	e := json.Unmarshal([]byte(`null`), &typ)

	a.Equal(typ.Proto(), float64(0))
	a.Nil(e)
}
