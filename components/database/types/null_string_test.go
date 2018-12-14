package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullString_NewInstance_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}

	a.False(typ.Valid)
}

func TestNullString_NewInstance_ValueEmptyString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}

	a.Equal(typ.String, "")
}

func TestNullString_NewInstance_ProtoReturnsEmptyString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}

	a.Equal(typ.Proto(), "")
}

func TestNullString_NewInstance_MarshalJSONReturnsNullAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`null`))
	a.Nil(e)
}

func TestNullString_NewInstanceAndScanOne_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	e := typ.Scan(1)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullString_NewInstanceAndScanOne_ValueOne(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	e := typ.Scan("one")

	a.Equal(typ.String, "one")
	a.Nil(e)
}

func TestNullString_NewInstanceAndScanOne_ProtoReturnsOne(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	e := typ.Scan("one")

	a.Equal(typ.Proto(), "one")
	a.Nil(e)
}

func TestNullString_NewInstanceAndScanOne_MarshalJSONReturnsOneAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	es := typ.Scan("one")
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`"one"`))
	a.Nil(e)
	a.Nil(es)
}

func TestNullString_NewInstanceAndScanEmptyString_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	e := typ.Scan("")

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullString_NewInstanceAndScanEmptyString_ValueEmptyString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	e := typ.Scan("")

	a.Equal(typ.String, "")
	a.Nil(e)
}

func TestNullString_NewInstanceAndScanEmptyString_ProtoReturnsEmptyString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	e := typ.Scan("")

	a.Equal(typ.Proto(), "")
	a.Nil(e)
}

func TestNullString_NewInstanceAndScanEmptyString_MarshalJSONReturnsEmptyStringAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	es := typ.Scan("")
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`""`))
	a.Nil(e)
	a.Nil(es)
}

func TestNullString_NewInstanceAndScanNil_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	e := typ.Scan(nil)

	a.False(typ.Valid)
	a.Nil(e)
}

func TestNullString_NewInstanceAndScanNil_ValueEmptyString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	e := typ.Scan(nil)

	a.Equal(typ.String, "")
	a.Nil(e)
}

func TestNullString_NewInstanceAndScanNil_ProtoReturnsEmptyString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	e := typ.Scan(nil)

	a.Equal(typ.String, "")
	a.Nil(e)
}

func TestNullString_NewInstanceAndScanNil_MarshalJSONReturnsNilAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullString{}
	es := typ.Scan(nil)
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`null`))
	a.Nil(e)
	a.Nil(es)
}

func TestNullString_NewInstanceAndUnmarshalJSONOne_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullString
	e := json.Unmarshal([]byte(`"one"`), &typ)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullString_NewInstanceAndUnmarshalJSONOne_ValueIsOne(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullString
	e := json.Unmarshal([]byte(`"one"`), &typ)

	a.Equal(typ.String, "one")
	a.Nil(e)
}

func TestNullString_NewInstanceAndUnmarshalJSONOne_ProtoReturnsOne(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullString
	e := json.Unmarshal([]byte(`"one"`), &typ)

	a.Equal(typ.String, "one")
	a.Nil(e)
}

func TestNullString_NewInstanceAndUnmarshalJSONEmptyString_ValidationIsTrue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullString
	e := json.Unmarshal([]byte(`""`), &typ)

	a.True(typ.Valid)
	a.Nil(e)
}

func TestNullString_NewInstanceAndUnmarshalJSONEmptyString_ValueIsEmptyString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullString
	e := json.Unmarshal([]byte(`""`), &typ)

	a.Equal(typ.String, "")
	a.Nil(e)
}

func TestNullString_NewInstanceAndUnmarshalJSONEmptyString_ProtoReturnsEmptyString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullString
	e := json.Unmarshal([]byte(`""`), &typ)

	a.Equal(typ.Proto(), "")
	a.Nil(e)
}

func TestNullString_NewInstanceAndUnmarshalJSONNil_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullString
	e := json.Unmarshal([]byte(`null`), &typ)

	a.False(typ.Valid)
	a.Nil(e)
}

func TestNullString_NewInstanceAndUnmarshalJSONNil_ValueIsEmptyString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullString
	e := json.Unmarshal([]byte(`null`), &typ)

	a.Equal(typ.String, "")
	a.Nil(e)
}

func TestNullString_NewInstanceAndUnmarshalJSONNil_ProtoReturnsEmptyString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullString
	e := json.Unmarshal([]byte(`null`), &typ)

	a.Equal(typ.Proto(), "")
	a.Nil(e)
}
