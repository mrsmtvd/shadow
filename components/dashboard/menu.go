package dashboard

type Menu interface {
	Title() string
	Url() string
	Route() Route
	Icon() string
	Childs() []Menu
	IsShow(request *Request) bool
}

type HasMenu interface {
	DashboardMenu() Menu
}

type MenuSimple struct {
	title  string
	url    string
	route  Route
	icon   string
	childs []Menu
	show   func(*Request) bool
}

func NewMenu(title string) *MenuSimple {
	return &MenuSimple{
		title:  title,
		childs: make([]Menu, 0),
	}
}

func (m *MenuSimple) Title() string {
	return m.title
}

func (m *MenuSimple) WithTitle(title string) *MenuSimple {
	m.title = title
	return m
}

func (m *MenuSimple) Url() string {
	if m.url == "" && m.route != nil {
		return m.route.Path()
	}

	return m.url
}

func (m *MenuSimple) WithUrl(url string) *MenuSimple {
	m.url = url
	return m
}

func (m *MenuSimple) Route() Route {
	return m.route
}

func (m *MenuSimple) WithRoute(route Route) *MenuSimple {
	m.route = route
	return m
}

func (m *MenuSimple) Icon() string {
	return m.icon
}

func (m *MenuSimple) WithIcon(icon string) *MenuSimple {
	m.icon = icon
	return m
}

func (m *MenuSimple) Childs() []Menu {
	return m.childs
}

func (m *MenuSimple) WithChild(child Menu) *MenuSimple {
	m.childs = append(m.childs, child)
	return m
}

func (m *MenuSimple) WithChilds(childs []Menu) *MenuSimple {
	m.childs = append(m.childs, childs...)
	return m
}

func (m *MenuSimple) IsShow(request *Request) bool {
	if m.show != nil {
		return m.show(request)
	}

	return true
}

func (m *MenuSimple) WithShow(show func(*Request) bool) *MenuSimple {
	m.show = show
	return m
}
