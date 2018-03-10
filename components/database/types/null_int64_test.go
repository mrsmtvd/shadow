package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullInt64_NewInstance_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}

	a.False(typ.Valid)
}

func TestNullInt64_NewInstance_ValueZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}

	a.Equal(typ.Int64, int64(0))
}

func TestNullInt64_NewInstance_ProtoReturnsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}

	a.Equal(typ.Proto(), int64(0))
}

func TestNullInt64_NewInstance_MarshalJSONReturnsNullAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`null`))
	a.Nil(e)
}

func TestNullInt64_NewInstanceAndScanOne_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(1)

	a.True(typ.Valid)
}

func TestNullInt64_NewInstanceAndScanOne_ValueOne(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(1)

	a.Equal(typ.Int64, int64(1))
}

func TestNullInt64_NewInstanceAndScanOne_ProtoReturnsOne(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(1)

	a.Equal(typ.Proto(), int64(1))
}

func TestNullInt64_NewInstanceAndScanOne_MarshalJSONReturnsOneAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(1)
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`1`))
	a.Nil(e)
}

func TestNullInt64_NewInstanceAndScanZero_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(0)

	a.True(typ.Valid)
}

func TestNullInt64_NewInstanceAndScanZero_ValueZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(0)

	a.Equal(typ.Int64, int64(0))
}

func TestNullInt64_NewInstanceAndScanZero_ProtoReturnsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(0)

	a.Equal(typ.Proto(), int64(0))
}

func TestNullInt64_NewInstanceAndScanZero_MarshalJSONReturnsZeroAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(0)
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`0`))
	a.Nil(e)
}

func TestNullInt64_NewInstanceAndScanNil_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(nil)

	a.False(typ.Valid)
}

func TestNullInt64_NewInstanceAndScanNil_ValueZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(nil)

	a.Equal(typ.Int64, int64(0))
}

func TestNullInt64_NewInstanceAndScanNil_ProtoReturnsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(nil)

	a.Equal(typ.Int64, int64(0))
}

func TestNullInt64_NewInstanceAndScanNil_MarshalJSONReturnsNilAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullInt64{}
	typ.Scan(nil)
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`null`))
	a.Nil(e)
}

func TestNullInt64_NewInstanceAndUnmarshalJSONOne_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullInt64
	e := json.Unmarshal([]byte(`1`), &typ)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullInt64_NewInstanceAndUnmarshalJSONOne_ValueIsOne(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullInt64
	e := json.Unmarshal([]byte(`1`), &typ)

	a.Equal(typ.Int64, int64(1))
	a.Nil(e)
}

func TestNullInt64_NewInstanceAndUnmarshalJSONOne_ProtoReturnsOne(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullInt64
	json.Unmarshal([]byte(`1`), &typ)

	a.Equal(typ.Int64, int64(1))
}

func TestNullInt64_NewInstanceAndUnmarshalJSONZero_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullInt64
	e := json.Unmarshal([]byte(`0`), &typ)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullInt64_NewInstanceAndUnmarshalJSONZero_ValueIsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullInt64
	e := json.Unmarshal([]byte(`0`), &typ)

	a.Equal(typ.Int64, int64(0))
	a.Nil(e)
}

func TestNullInt64_NewInstanceAndUnmarshalJSONZero_ProtoReturnsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullInt64
	json.Unmarshal([]byte(`0`), &typ)

	a.Equal(typ.Proto(), int64(0))
}

func TestNullInt64_NewInstanceAndUnmarshalJSONNil_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullInt64
	e := json.Unmarshal([]byte(`null`), &typ)

	a.False(typ.Valid)
	a.Nil(e)
}

func TestNullInt64_NewInstanceAndUnmarshalJSONNil_ValueIsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullInt64
	e := json.Unmarshal([]byte(`null`), &typ)

	a.Equal(typ.Proto(), int64(0))
	a.Nil(e)
}

func TestNullInt64_NewInstanceAndUnmarshalJSONNil_ProtoReturnsZero(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullInt64
	json.Unmarshal([]byte(`null`), &typ)

	a.Equal(typ.Proto(), int64(0))
}
