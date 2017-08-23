package config

const (
	ViewPassword = "password"
)

type Variable struct {
	Key      string
	Default  interface{}
	Value    interface{}
	Type     string
	Usage    string
	Editable bool
	View     []string
}

func (v Variable) HasView(n string) bool {
	if len(v.View) == 0 {
		return false
	}

	for _, nv := range v.View {
		if nv == n {
			return true
		}
	}

	return false
}
