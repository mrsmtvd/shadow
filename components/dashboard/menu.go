package dashboard

type Menu struct {
	Name    string
	Url     string
	Icon    string
	SubMenu []*Menu
}

type hasMenu interface {
	GetDashboardMenu() *Menu
}

func (c *Component) loadMenu() {
	menus := make([]*Menu, 0)

	for _, component := range c.application.GetComponents() {
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

	c.renderer.AddGlobalVar("Menu", menus)
}

func (c *Component) changeUrlMenu(m *Menu, p string) {
	m.Url = "/" + p + m.Url

	if len(m.SubMenu) > 0 {
		for i := range m.SubMenu {
			c.changeUrlMenu(m.SubMenu[i], p)
		}
	}
}
