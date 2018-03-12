package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStruct_NilEntry_ToMapReturnsEmptyMap(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var typ Struct

	a.Equal(typ.ToMap(), map[string]interface{}{})
}
