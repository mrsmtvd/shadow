package shadow

import (
	"errors"
	"sync"

	"github.com/deckarep/golang-set"
)

type components struct {
	sync.Map

	resolved sync.Once
}

func (c *components) add(n string, cmp Component) {
	c.Store(n, newComponent(cmp))
	c.resolved = sync.Once{}
}

func (c *components) get(n string) (*component, bool) {
	cmp, ok := c.Load(n)
	if ok {
		return cmp.(*component), ok
	}

	return nil, ok
}

func (c *components) all() ([]*component, error) {
	var err error
	c.resolved.Do(func() {
		err = c.resolve()
	})

	if err != nil {
		return nil, err
	}

	result := make([]*component, 0)

	c.Range(func(_, value interface{}) bool {
		result = append(result, value.(*component))
		return true
	})

	return result, nil
}

func (c *components) resolve() (err error) {
	dependencies := make(map[string]mapset.Set)
	c.Range(func(_, value interface{}) bool {
		ms := mapset.NewSet()
		cmp := value.(*component)

		if dependency, ok := cmp.instance.(ComponentDependency); ok {
			for _, dep := range dependency.Dependencies() {
				if dep.Required {
					if _, exist := c.get(dep.Name); !exist {
						err = errors.New("component \"" + cmp.instance.Name() + "\" has required dependency \"" + dep.Name + "\"")
						return false
					}
				} else if _, exist := c.get(dep.Name); !exist {
					dependencies[dep.Name] = mapset.NewSet()
				}

				ms.Add(dep.Name)
			}
		}

		dependencies[cmp.instance.Name()] = ms
		return true
	})

	if err != nil {
		return err
	}

	var index int64

	for len(dependencies) > 0 {
		readyMs := mapset.NewSet()
		for name, ms := range dependencies {
			if ms.Cardinality() == 0 {
				readyMs.Add(name)
			}
		}

		if readyMs.Cardinality() == 0 {
			return errors.New("circular dependency found")
		}

		for n := range readyMs.Iter() {
			name := n.(string)
			delete(dependencies, name)

			if cmp, exist := c.get(name); exist {
				cmp.SetOrder(index)
				index++
			}
		}

		for name, ms := range dependencies {
			dependencies[name] = ms.Difference(readyMs)
		}
	}

	return nil
}
