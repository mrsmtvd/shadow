package internal

// http://docs.translatehouse.org/projects/localization-guide/en/latest/l10n/pluralforms.html

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/chai2010/gettext-go/gettext/mo"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/i18n/internationalization"
	"github.com/kihamo/shadow/components/logger"
	"golang.org/x/text/language"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	manager     *internationalization.Manager
}

func (c *Component) Name() string {
	return i18n.ComponentName
}

func (c *Component) Version() string {
	return i18n.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: logger.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.manager = internationalization.NewManager()

	return nil
}

func (c *Component) Run() error {
	c.logger = logger.NewOrNop(c.Name(), c.application)

	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	for _, cmp := range components {
		if cmpI18n, ok := cmp.(i18n.HasI18n); ok {
			locales := cmpI18n.I18n()
			if len(locales) == 0 {
				continue
			}

			for localeName, reader := range locales {
				b, err := ioutil.ReadAll(reader)
				if err != nil {
					c.logger.Warn("Failed read from file", map[string]interface{}{
						"component": cmp.Name(),
						"locale":    localeName,
						"error":     err.Error(),
					})
					continue
				}

				file, err := mo.LoadData(b)
				if err != nil {
					c.logger.Warn("Failed parse MO file", map[string]interface{}{
						"component": cmp.Name(),
						"locale":    localeName,
						"error":     err.Error(),
					})
					continue
				}

				locale, ok := c.manager.Locale(localeName)
				if !ok {
					locale = internationalization.NewLocale(localeName)
					c.manager.AddLocale(locale)
				}

				pluralRule := internationalization.NewPluralRule(file.MimeHeader.PluralForms)

				messages := make([]*internationalization.Message, 0, len(file.Messages))
				for _, m := range file.Messages {
					messages = append(messages, internationalization.NewMessage(m.MsgId, m.MsgStr, m.MsgIdPlural, m.MsgStrPlural, m.MsgContext))
				}

				locale.AddDomain(internationalization.NewDomain(cmp.Name(), messages, pluralRule))

				c.logger.Debugf("Load %d translations", len(file.Messages), map[string]interface{}{
					"component": cmp.Name(),
					"locale":    localeName,
				})
			}
		}
	}

	return nil
}

func (c *Component) Manager() *internationalization.Manager {
	return c.manager
}

func (c *Component) LocaleFromAcceptLanguage(acceptLanguage string) (*internationalization.Locale, error) {
	tags, _, err := language.ParseAcceptLanguage(acceptLanguage)
	if err != nil {
		return nil, err
	}

	for _, t := range tags {
		locale, ok := c.Manager().Locale(strings.Replace(t.String(), "-", "_", -1))
		if ok {
			return locale, nil
		}
	}

	return nil, errors.New("Locale not found")
}
