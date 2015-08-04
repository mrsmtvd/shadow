package api

import (
	"encoding/json"
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
	var target interface{}
	if err := json.Unmarshal([]byte(value), &target); err != nil {
		*a = append(*a, value)
	} else {
		*a = append(*a, target)
	}

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

	var target interface{}
	if err := json.Unmarshal([]byte(v[1]), &target); err != nil {
		(*k)[v[0]] = v[1]
	} else {
		(*k)[v[0]] = target
	}

	return nil
}
