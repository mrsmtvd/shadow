package api

import (
	"errors"
	"fmt"
	"strings"
)

type procedureArgs []interface{}
type procedureKwargs map[string]interface{}

func (a *procedureArgs) String() string {
	return fmt.Sprintf("%v", *a)
}

func (a *procedureArgs) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func (k *procedureKwargs) String() string {
	return fmt.Sprintf("%v", *k)
}

func (k *procedureKwargs) Set(value string) error {
	v := strings.Split(value, "=")

	if len(v) != 2 {
		return errors.New("Not valid format. Use \"key=value\"")
	}

	(*k)[v[0]] = v[1]
	return nil
}
