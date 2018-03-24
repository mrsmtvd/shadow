package internal

// http://docs.translatehouse.org/projects/localization-guide/en/latest/l10n/pluralforms.html

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/chai2010/gettext-go/gettext/mo"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
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

			for localeName, readers := range locales {
				var domain *internationalization.Domain

				locale, ok := c.manager.Locale(localeName)
				if !ok {
					locale = internationalization.NewLocale(localeName)
					c.manager.AddLocale(locale)
				}

				for _, reader := range readers {
					b, err := ioutil.ReadAll(reader)
					if err != nil {
						c.logger.Warn("Failed read from file", map[string]interface{}{
							"locale": localeName,
							"error":  err.Error(),
							"domain": cmp.Name(),
						})
						continue
					}

					file, err := mo.LoadData(b)
					if err != nil {
						c.logger.Warn("Failed parse MO file", map[string]interface{}{
							"locale": localeName,
							"error":  err.Error(),
							"domain": cmp.Name(),
						})
						continue
					}

					pluralRule := internationalization.NewPluralRule(file.MimeHeader.PluralForms)

					messages := make([]*internationalization.Message, 0, len(file.Messages))
					for _, m := range file.Messages {
						messages = append(messages, internationalization.NewMessage(m.MsgId, m.MsgStr, m.MsgIdPlural, m.MsgStrPlural, m.MsgContext))
					}

					domainTmp := internationalization.NewDomain(cmp.Name(), messages, pluralRule)

					if domain == nil {
						domain = domainTmp
					} else {
						mergeDomain, err := domain.Merge(domainTmp)
						if err != nil {
							c.logger.Warn("Failed merge domains", map[string]interface{}{
								"locale": localeName,
								"domain": cmp.Name(),
								"error":  err.Error(),
							})
						} else {
							domain = mergeDomain
						}
					}
				}

				if domain != nil {
					locale.AddDomain(domain)

					c.logger.Debugf("Load %d translations", len(domain.Messages()), map[string]interface{}{
						"locale": localeName,
						"domain": cmp.Name(),
					})
				}
			}
		}
	}

	return nil
}

func (c *Component) Manager() *internationalization.Manager {
	return c.manager
}

func (c *Component) LocaleFromRequest(request *dashboard.Request) (*internationalization.Locale, error) {
	// in session
	localeSession, err := c.LocaleFromSession(request.Session())
	if err == nil {
		return localeSession, err
	}

	// in request
	return c.LocaleFromAcceptLanguage(request.Original().Header.Get("Accept-Language"))
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

func (c *Component) LocaleFromSession(session dashboard.Session) (*internationalization.Locale, error) {
	localeName, err := session.GetString(i18n.SessionLocale)
	if err != nil {
		return nil, err
	}

	locale, ok := c.Manager().Locale(localeName)
	if ok {
		return locale, nil
	}

	return nil, errors.New("Locale not found")
}

func (c *Component) SessionSave(session dashboard.Session, locale *internationalization.Locale) error {
	return session.PutString(i18n.SessionLocale, locale.Locale())
}
