package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNullTime_NewInstance_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullTime{}

	a.False(typ.Valid)
}

func TestNullTime_NewInstance_ValueDefaultTime(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullTime{}

	a.Equal(typ.Time, time.Time{})
}

func TestNullTime_NewInstance_ProtoReturnsNil(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullTime{}

	a.Nil(typ.Proto())
}

func TestNullTime_NewInstance_MarshalJSONReturnsNullAsBytes(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	typ := &NullTime{}
	j, e := json.Marshal(typ)

	a.Equal(j, []byte(`null`))
	a.Nil(e)
}

func TestNullTime_NewInstanceAndUnmarshalJSONNil_ValidationIsFalse(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullTime
	e := json.Unmarshal([]byte(`null`), &typ)

	a.False(typ.Valid)
	a.Nil(e)
}

func TestNullTime_NewInstanceAndUnmarshalJSONNil_ValueIsDefaultTime(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullTime
	e := json.Unmarshal([]byte(`null`), &typ)

	a.Equal(typ.Time, time.Time{})
	a.Nil(e)
}

func TestNullTime_NewInstanceAndUnmarshalJSONNil_ProtoReturnsNil(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ NullTime
	e := json.Unmarshal([]byte(`null`), &typ)

	a.Nil(typ.Proto())
	a.Nil(e)
}
