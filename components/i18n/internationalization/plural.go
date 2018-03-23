package internationalization

import (
	"strconv"
	"strings"

	"github.com/mattn/kinako/ast"
	"github.com/mattn/kinako/parser"
	"github.com/mattn/kinako/vm"
)

type PluralRule struct {
	numberPlurals int
	plural        []ast.Stmt
}

func NewPluralRule(rule string) *PluralRule {
	r := &PluralRule{}
	r.parse(rule)

	return r
}

func (r *PluralRule) parse(rule string) {
	for _, sub := range strings.Split(rule, ";") {
		kv := strings.SplitN(sub, "=", 2)

		if len(kv) != 2 {
			continue
		}

		switch strings.TrimSpace(kv[0]) {
		case "nplurals":
			r.numberPlurals, _ = strconv.Atoi(kv[1])

		case "plural":
			stmts, err := parser.ParseSrc(kv[1])
			if err == nil {
				r.plural = stmts
			}
		}
	}
}

func (r *PluralRule) Number(number int) int {
	env := vm.NewEnv()
	env.Define("n", number)

	plural, err := vm.Run(r.plural, env)
	if err != nil {
		return 0
	}

	if plural.Type().Name() == "bool" {
		if plural.Bool() {
			return 1
		}

		return 0
	}

	if int(plural.Int()) > r.numberPlurals {
		return 0
	}

	return int(plural.Int())
}

func (r *PluralRule) NumberPlurals() int {
	return r.numberPlurals
}
