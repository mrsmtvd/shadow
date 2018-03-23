package dashboard

import (
	"html/template"
	"sync"
)

var DefaultTemplateFunctions = &TemplateFunctions{
	funcsMap: template.FuncMap{},
}

type TemplateFunctions struct {
	mutex    sync.RWMutex
	funcsMap template.FuncMap
}

func (t *TemplateFunctions) AddFunction(key string, fn interface{}) {
	t.mutex.Lock()
	t.funcsMap[key] = fn
	t.mutex.Unlock()
}

func (t *TemplateFunctions) FuncMap() template.FuncMap {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.funcsMap
}
