package internationalization

import (
	"sync"
)

const (
	DefaultLocale = "en"
)

type Manager struct {
	mutex sync.RWMutex

	locales       map[string]*Locale
	defaultLocale string
}

func NewManager(defaultLocale string) *Manager {
	if defaultLocale == "" {
		defaultLocale = DefaultLocale
	}

	m := &Manager{
		locales:       make(map[string]*Locale),
		defaultLocale: defaultLocale,
	}
	m.AddLocale(NewLocale(defaultLocale))

	return m
}

func (m *Manager) DefaultLocale() string {
	return m.defaultLocale
}

func (m *Manager) AddLocale(locale *Locale) {
	m.mutex.Lock()
	m.locales[locale.Locale()] = locale
	m.mutex.Unlock()
}

func (m *Manager) Locale(locale string) (*Locale, bool) {
	m.mutex.RLock()
	l, ok := m.locales[locale]
	m.mutex.RUnlock()

	return l, ok
}

func (m *Manager) Locales() []*Locale {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	locales := make([]*Locale, 0, len(m.locales))
	for _, locale := range m.locales {
		locales = append(locales, locale)
	}

	return locales
}

func (m *Manager) Translate(locale, domain, singleID, context string, format ...interface{}) string {
	return m.TranslatePlural(locale, domain, singleID, "", 1, context, format...)
}

func (m *Manager) TranslatePlural(locale, domain, singleID, pluralID string, number int, context string, format ...interface{}) string {
	if locale == "" {
		locale = m.DefaultLocale()
	}

	l, ok := m.Locale(locale)
	if !ok {
		if number > 1 {
			return Format(pluralID, format...)
		}

		return Format(singleID, format...)
	}

	return l.TranslatePlural(domain, singleID, pluralID, number, context, format...)
}
