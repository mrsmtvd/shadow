package internal

import (
	"sort"
	"strings"

	"github.com/kihamo/shadow/components/dashboard"
)

type orderedMenus []dashboard.Menu

func (m orderedMenus) Len() int {
	return len(m)
}
func (m orderedMenus) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m orderedMenus) Less(i, j int) bool {
	return strings.Compare(m[i].Title(), m[j].Title()) < 0
}

func (c *Component) loadMenu() error {
	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	menus := make([]dashboard.Menu, 0)

	for _, component := range components {
		if componentMenu, ok := component.(dashboard.HasMenu); ok {
			menu := componentMenu.DashboardMenu()
			if menu != nil {
				menu = NewMenuItem(menu, component.Name())

				if component == c {
					menus = append([]dashboard.Menu{menu}, menus...)
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
