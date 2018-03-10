package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VariableSuite struct {
	suite.Suite
}

func TestVariableSuite(t *testing.T) {
	suite.Run(t, new(VariableSuite))
}

func (s *VariableSuite) Test_NewInstance_ValueReturnsNil() {
	t := NewVariable(
		"test",
		ValueTypeBool,
		nil,
		"Test",
		true,
		"group",
		nil,
		nil)

	assert.Nil(s.T(), t.Value())
}

func (s *VariableSuite) Test_NewInstance_ValueReturnsDefaultValue() {
	t := NewVariable(
		"test",
		ValueTypeBool,
		"12345",
		"Test",
		true,
		"group",
		nil,
		nil)

	assert.Equal(s.T(), t.Value(), "12345")
}

func (s *VariableSuite) Test_NewInstanceAndSetNilValue_ValueReturnsNil() {
	t := NewVariable(
		"test",
		ValueTypeBool,
		"12345",
		"Test",
		true,
		"group",
		nil,
		nil)

	t.Change(nil)

	assert.Nil(s.T(), t.Value())
}

func (s *VariableSuite) Test_NewInstanceAndSetStringValue_ValueReturnsString() {
	t := NewVariable(
		"test",
		ValueTypeBool,
		"12345",
		"Test",
		true,
		"group",
		nil,
		nil)

	t.Change("54321")

	assert.Equal(s.T(), t.Value(), "54321")
}

func BenchmarkVariableChange(b *testing.B) {
	t := NewVariable(
		"test",
		ValueTypeBool,
		"22",
		"Test",
		true,
		"group",
		nil,
		nil)

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			t.Change(nil)
			t.Value()
		}
	})
}
