package dashboard

import (
	"sort"
	"strings"
)

type Menu struct {
	Name    string
	Url     string
	Direct  bool
	Icon    string
	SubMenu []*Menu
}

type hasMenu interface {
	GetDashboardMenu() *Menu
}

type orderedMenus []*Menu

func (m orderedMenus) Len() int {
	return len(m)
}
func (m orderedMenus) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m orderedMenus) Less(i, j int) bool {
	return strings.Compare(m[i].Name, m[j].Name) < 0
}

func (c *Component) loadMenu() error {
	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	menus := make([]*Menu, 0)

	for _, component := range components {
		if componentMenu, ok := component.(hasMenu); ok {
			menu := componentMenu.GetDashboardMenu()
			if menu != nil {
				c.changeUrlMenu(menu, component.GetName())

				if component == c {
					menus = append([]*Menu{menu}, menus...)
				} else {
					menus = append(menus, menu)
				}
			}
		}
	}

	contextMenus := orderedMenus(menus)
	sort.Sort(contextMenus)

	c.renderer.AddGlobalVar("Menu", contextMenus)
	return nil
}

func (c *Component) changeUrlMenu(m *Menu, p string) {
	if !m.Direct {
		m.Url = "/" + p + m.Url
	}

	if len(m.SubMenu) > 0 {
		for i := range m.SubMenu {
			c.changeUrlMenu(m.SubMenu[i], p)
		}
	}
}
