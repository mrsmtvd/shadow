package internationalization

import (
	"sync"
)

const (
	DefaultDomain = "default"
)

type Locale struct {
	mutex sync.RWMutex

	locale  string
	domains map[string]*Domain
}

func NewLocale(locale string) *Locale {
	return &Locale{
		locale:  locale,
		domains: make(map[string]*Domain, 0),
	}
}

func (l *Locale) Locale() string {
	return l.locale
}

func (l *Locale) AddDomain(domain *Domain) {
	l.mutex.Lock()
	l.domains[domain.Name()] = domain
	l.mutex.Unlock()
}

func (l *Locale) Domain(domain string) (*Domain, bool) {
	l.mutex.RLock()
	d, ok := l.domains[domain]
	l.mutex.RUnlock()

	return d, ok
}

func (l *Locale) Translate(domain, ID, context string, format ...interface{}) string {
	return l.TranslatePlural(domain, ID, "", 1, context, format...)
}

func (l *Locale) TranslatePlural(domain, singleID, pluralID string, number int, context string, format ...interface{}) string {
	if domain == "" {
		domain = DefaultDomain
	}

	d, ok := l.Domain(domain)
	if !ok {
		if number > 1 {
			return Format(pluralID, format...)
		}

		return Format(singleID, format...)
	}

	return d.TranslatePlural(singleID, pluralID, number, context, format...)
}
