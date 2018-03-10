package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullBool_NewInstance_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}

	a.False(typ.Valid)
}

func TestNullBool_NewInstance_ValueFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}

	a.False(typ.Bool)
}

func TestNullBool_NewInstance_ProtoReturnsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}

	a.False(typ.Proto())
}

func TestNullBool_NewInstance_MarshalJSONReturnsNullAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`null`))
	a.Nil(e)
}

func TestNullBool_NewInstanceAndScanTrue_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(true)

	a.True(typ.Valid)
}

func TestNullBool_NewInstanceAndScanTrue_ValueTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(true)

	a.True(typ.Bool)
}

func TestNullBool_NewInstanceAndScanTrue_ProtoReturnsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(true)

	a.True(typ.Proto())
}

func TestNullBool_NewInstanceAndScanTrue_MarshalJSONReturnsTrueAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(true)
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`true`))
	a.Nil(e)
}

func TestNullBool_NewInstanceAndScanFalse_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(false)

	a.True(typ.Valid)
}

func TestNullBool_NewInstanceAndScanFalse_ValueFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(false)

	a.False(typ.Bool)
}

func TestNullBool_NewInstanceAndScanFalse_ProtoReturnsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(false)

	a.False(typ.Proto())
}

func TestNullBool_NewInstanceAndScanFalse_MarshalJSONReturnsFalseAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(false)
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`false`))
	a.Nil(e)
}

func TestNullBool_NewInstanceAndScanNil_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(nil)

	a.False(typ.Valid)
}

func TestNullBool_NewInstanceAndScanNil_ValueFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(nil)

	a.False(typ.Bool)
}

func TestNullBool_NewInstanceAndScanNil_ProtoReturnsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(nil)

	a.False(typ.Proto())
}

func TestNullBool_NewInstanceAndScanNil_MarshalJSONReturnsNilAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullBool{}
	typ.Scan(nil)
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`null`))
	a.Nil(e)
}

func TestNullBool_NewInstanceAndUnmarshalJSONTrue_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullBool
	e := json.Unmarshal([]byte(`true`), &typ)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullBool_NewInstanceAndUnmarshalJSONTrue_ValueIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullBool
	e := json.Unmarshal([]byte(`true`), &typ)

	a.True(typ.Bool)
	a.Nil(e)
}

func TestNullBool_NewInstanceAndUnmarshalJSONTrue_ProtoReturnsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullBool
	json.Unmarshal([]byte(`true`), &typ)

	a.True(typ.Proto())
}

func TestNullBool_NewInstanceAndUnmarshalJSONFalse_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullBool
	e := json.Unmarshal([]byte(`false`), &typ)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullBool_NewInstanceAndUnmarshalJSONFalse_ValueIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullBool
	e := json.Unmarshal([]byte(`false`), &typ)

	a.False(typ.Bool)
	a.Nil(e)
}

func TestNullBool_NewInstanceAndUnmarshalJSONFalse_ProtoReturnsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullBool
	json.Unmarshal([]byte(`false`), &typ)

	a.False(typ.Proto())
}

func TestNullBool_NewInstanceAndUnmarshalJSONNil_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullBool
	e := json.Unmarshal([]byte(`null`), &typ)

	a.False(typ.Valid)
	a.Nil(e)
}

func TestNullBool_NewInstanceAndUnmarshalJSONNil_ValueIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullBool
	e := json.Unmarshal([]byte(`null`), &typ)

	a.False(typ.Bool)
	a.Nil(e)
}

func TestNullBool_NewInstanceAndUnmarshalJSONNil_ProtoReturnsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullBool
	json.Unmarshal([]byte(`null`), &typ)

	a.False(typ.Proto())
}
