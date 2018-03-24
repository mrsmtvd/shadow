package internationalization

import (
	"errors"
	"fmt"
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
		for _, key := range domainMessageKeys(m) {
			d.messagesByKeys[key] = m
		}
	}

	return d
}

func domainMessageKeys(message *Message) []string {
	id := message.SingleID()
	if id == "" {
		return nil
	}

	keys := []string{
		message.SingleID(),
	}

	if ctx := message.Context(); ctx != "" {
		keys = append(keys, domainMessageKeyWIthContext(id, ctx))
	}

	return keys
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
	list := make([]*Message, 0, len(d.messages))
	copy(list, d.messages)

	return list
}

func (d *Domain) Merge(domain *Domain) (*Domain, error) {
	if d.PluralRule().NumberPlurals() != domain.PluralRule().NumberPlurals() {
		return nil, errors.New("Plural rule of merging domain is not compatible with the current one")
	}

	return NewDomain(d.Name(), append(d.Messages(), domain.Messages()...), d.PluralRule()), nil
}

func (d *Domain) Translate(ID string, context string, format ...interface{}) string {
	return d.TranslatePlural(ID, "", 1, context, format...)
}

func (d *Domain) TranslatePlural(singleID string, pluralID string, number int, context string, format ...interface{}) string {
	var message *Message

	if messageByCtx, ok := d.messagesByKeys[domainMessageKeyWIthContext(singleID, context)]; ok { // by context
		message = messageByCtx
	} else if messageByID, ok := d.messagesByKeys[singleID]; ok { // by ID
		message = messageByID
	}

	// check plural
	if message != nil {
		index := d.PluralRule().Number(number)
		translates := message.Translated()

		if index < len(translates) {
			if len(format) > 0 {
				return fmt.Sprintf(translates[index], format...)
			}

			return translates[index]
		}
	}

	// not found
	if number > 1 {
		return pluralID
	}

	return singleID
}
