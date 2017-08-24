package config

const (
	ViewPassword = "password"
	ViewTags     = "tags"

	ViewOptionTagsDefaultText = "default-text"
)

type Variable struct {
	Key         string
	Default     interface{}
	Value       interface{}
	Type        string
	Usage       string
	Editable    bool
	View        []string
	ViewOptions map[string]interface{}
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

func (v Variable) GetViewOption(o string) interface{} {
	if len(v.ViewOptions) > 0 {
		if opt, ok := v.ViewOptions[o]; ok {
			return opt
		}
	}

	return nil
}
