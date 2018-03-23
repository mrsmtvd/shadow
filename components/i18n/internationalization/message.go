package internationalization

type Message struct {
	singleID   string
	pluralID   string
	translated []string
	context    string
}

func NewMessage(singleID, translatedSingle, pluralID string, translatedPlural []string, context string) *Message {
	m := &Message{
		singleID: singleID,
		pluralID: pluralID,
		context:  context,
	}

	if len(translatedPlural) > 0 {
		m.translated = translatedPlural
	} else {
		m.translated = []string{translatedSingle}
	}

	return m
}

func (m *Message) SingleID() string {
	return m.singleID
}

func (m *Message) PluralID() string {
	return m.pluralID
}

func (m *Message) Translated() []string {
	return m.translated
}

func (m *Message) Context() string {
	return m.context
}
