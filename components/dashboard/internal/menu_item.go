package internal

import (
	"github.com/mrsmtvd/shadow/components/dashboard"
)

type MenuItem struct {
	dashboard.Menu

	menu   dashboard.Menu
	childs []dashboard.Menu
	source string
}

func NewMenuItem(menu dashboard.Menu, source string) *MenuItem {
	if source == "" {
		source = "unknown"
	}

	m := &MenuItem{
		menu:   menu,
		childs: make([]dashboard.Menu, 0, len(menu.Childs())),
		source: source,
	}

	for _, child := range menu.Childs() {
		m.childs = append(m.childs, NewMenuItem(child, source))
	}

	return m
}

func (m *MenuItem) Source() string {
	return m.source
}

func (m *MenuItem) Title() string {
	return m.menu.Title()
}

func (m *MenuItem) URL() string {
	return m.menu.URL()
}

func (m *MenuItem) Route() dashboard.Route {
	return m.menu.Route()
}

func (m *MenuItem) Icon() string {
	return m.menu.Icon()
}

func (m *MenuItem) Childs() []dashboard.Menu {
	return m.childs
}

func (m *MenuItem) IsShow(request *dashboard.Request) bool {
	return m.menu.IsShow(request)
}
