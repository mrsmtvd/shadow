package internationalization

import (
	"sync"
)

type Manager struct {
	mutex sync.RWMutex

	locales       map[string]*Locale
	defaultLocale string
}

func NewManager() *Manager {
	return &Manager{
		locales: make(map[string]*Locale),
	}
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
	if locale == "" {
		locale = m.DefaultLocale()
	}

	m.mutex.RLock()
	l, ok := m.locales[locale]
	m.mutex.RUnlock()

	return l, ok
}

func (m *Manager) Translate(locale, domain, ID, context string, format ...interface{}) string {
	return m.TranslatePlural(locale, domain, ID, "", 1, context, format...)
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
