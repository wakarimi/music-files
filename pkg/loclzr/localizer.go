package loclzr

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Localizer struct {
	bundle       *i18n.Bundle
	engLocalizer *i18n.Localizer
}

func New(bundle *i18n.Bundle) *Localizer {
	return &Localizer{
		bundle:       bundle,
		engLocalizer: i18n.NewLocalizer(bundle, "en_US"),
	}
}

func (l *Localizer) TryLocalize(lang string, messageID string) *string {
	localizeConfig := i18n.LocalizeConfig{
		MessageID: messageID,
	}

	var langMessage *string
	langLocalizer := i18n.NewLocalizer(l.bundle, lang)
	localizedMessage, err := langLocalizer.Localize(&localizeConfig)
	if err == nil {
		langMessage = &localizedMessage
	}

	return langMessage
}

func (l *Localizer) English(messageID string) string {
	localizeConfig := i18n.LocalizeConfig{
		MessageID: messageID,
	}

	engMessage := l.engLocalizer.MustLocalize(&localizeConfig)

	return engMessage
}
