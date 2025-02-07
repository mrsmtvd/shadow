package i18n

import (
	"context"

	i18n "github.com/mrsmtvd/shadow/components/i18n/internationalization"
)

type contextKey string

var (
	localeContextKey = contextKey("locale")
)

func ContextWithLocale(ctx context.Context, locale *i18n.Locale) context.Context {
	return context.WithValue(ctx, localeContextKey, locale)
}

func LocaleFromContext(ctx context.Context) *i18n.Locale {
	v := ctx.Value(localeContextKey)
	if v != nil {
		return v.(*i18n.Locale)
	}

	return i18n.NewLocale(i18n.DefaultLocale)
}

func Locale(ctx context.Context) *i18n.Locale {
	return LocaleFromContext(ctx)
}
