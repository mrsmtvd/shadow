package shadow

import (
	"errors"
	"sync"

	mapset "github.com/deckarep/golang-set"
)

type components struct {
	sync.Map

	resolved sync.Once
}

func (c *components) Add(n string, cmp Component) {
	c.Store(n, newComponent(cmp))
	c.resolved = sync.Once{}
}

func (c *components) Get(n string) (*component, bool) {
	cmp, ok := c.Load(n)
	if ok {
		return cmp.(*component), ok
	}

	return nil, ok
}

func (c *components) All() ([]*component, error) {
	var err error
	c.resolved.Do(func() {
		err = c.Resolve()
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

func (c *components) Resolve() (err error) {
	dependencies := make(map[string]mapset.Set)
	c.Range(func(_, value interface{}) bool {
		ms := mapset.NewSet()
		cmp := value.(*component)

		if dependency, ok := cmp.instance.(ComponentDependency); ok {
			for _, dep := range dependency.Dependencies() {
				depCmp, exist := c.Get(dep.Name)
				if exist {
					depCmp.AddReverseDep(cmp.Name())
				} else {
					if dep.Required {
						err = errors.New("component \"" + cmp.Name() + "\" has required dependency \"" + dep.Name + "\"")
						return false
					} else {
						dependencies[dep.Name] = mapset.NewSet()
					}
				}

				ms.Add(dep.Name)
			}
		}

		dependencies[cmp.Name()] = ms
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

			if cmp, exist := c.Get(name); exist {
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
