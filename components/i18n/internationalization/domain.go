package internationalization

import (
	"errors"
)

const (
	DomainMessageKeySeparator = "\x04"
)

var (
	DefaultPluralRule *PluralRule = NewPluralRule("nplurals=2; plural=n != 1;")
)

type Domain struct {
	name           string
	messages       []*Message
	messagesByKeys map[string]*Message
	pluralRule     *PluralRule
}

func NewDomain(name string, messages []*Message, pluralRule *PluralRule) *Domain {
	d := &Domain{
		name:           name,
		messages:       messages,
		messagesByKeys: make(map[string]*Message, len(messages)),
		pluralRule:     pluralRule,
	}

	for _, m := range messages {
		if m.SingleID() == "" {
			continue
		}

		if ctx := m.Context(); ctx != "" {
			d.messagesByKeys[domainMessageKeyWIthContext(m.SingleID(), ctx)] = m
		} else {
			d.messagesByKeys[m.SingleID()] = m
		}
	}

	return d
}

func domainMessageKeyWIthContext(ID, ctx string) string {
	return ctx + DomainMessageKeySeparator + ID
}

func (d *Domain) Name() string {
	if d.name == "" {
		return DefaultDomain
	}

	return d.name
}

func (d *Domain) PluralRule() *PluralRule {
	if d.pluralRule == nil {
		return DefaultPluralRule
	}

	return d.pluralRule
}

func (d *Domain) Messages() []*Message {
	list := make([]*Message, len(d.messages))
	copy(list, d.messages)

	return list
}

func (d *Domain) Merge(domain *Domain) (*Domain, error) {
	if d.PluralRule().NumberPlurals() != domain.PluralRule().NumberPlurals() {
		return nil, errors.New("Plural rule of merging domain is not compatible with the current one")
	}

	return NewDomain(d.Name(), append(d.Messages(), domain.Messages()...), d.PluralRule()), nil
}

func (d *Domain) Translate(ID, context string, format ...interface{}) string {
	return d.TranslatePlural(ID, "", 1, context, format...)
}

func (d *Domain) TranslatePlural(singleID, pluralID string, number int, context string, format ...interface{}) string {
	var message *Message

	// by context
	if context != "" {
		if messageByCtx, ok := d.messagesByKeys[domainMessageKeyWIthContext(singleID, context)]; ok {
			message = messageByCtx
		}
	}

	// by ID
	if message == nil {
		if messageByID, ok := d.messagesByKeys[singleID]; ok {
			message = messageByID
		}
	}

	// check plural
	if message != nil {
		index := d.PluralRule().Number(number)
		translates := message.Translated()

		if index < len(translates) {
			return Format(translates[index], format...)
		}
	}

	// not found
	if number > 1 {
		return Format(pluralID, format...)
	}

	return Format(singleID, format...)
}
