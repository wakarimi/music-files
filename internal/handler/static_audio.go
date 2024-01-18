package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/response"
	"music-files/internal/internal_error"
	"music-files/internal/use_case"
	"net/http"
	"path/filepath"
	"strconv"
)

func (h Handler) StaticAudio(c *gin.Context) {
	log.Debug().Msg("Getting a static url to an audio file")

	lang := c.MustGet("lang").(string)
	localizer := i18n.NewLocalizer(&h.bundle, lang)

	audioIDStr := c.Param("audioId")
	audioID, err := strconv.Atoi(audioIDStr)
	if err != nil {
		log.Error().Err(err).Str("audioIdStr", audioIDStr).Msg("Invalid audioId format")
		messageID := "InvalidInputFormat"
		message, errLoc := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
		if errLoc != nil {
			message = h.engLocalizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
		}
		c.JSON(http.StatusBadRequest, response.Error{
			Message: message,
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("audioId", audioID).Msg("Url parameter read successfully")

	staticAudioInput := use_case.StaticAudioInput{
		AudioID: audioID,
	}
	staticAudioOutput, err := h.useCase.StaticAudio(staticAudioInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get static link to an audio file")
		if _, ok := err.(internal_error.NotFound); ok {
			messageID := "AudioNotFound"
			message, errLoc := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
			if errLoc != nil {
				message = h.engLocalizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
			}
			c.JSON(http.StatusBadRequest, response.Error{
				Message: message,
				Reason:  err.Error(),
			})
			return
		} else {
			messageID := "FailedToGetStaticLinkToAudio"
			message, errLoc := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
			if errLoc != nil {
				message = h.engLocalizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
			}
			c.JSON(http.StatusBadRequest, response.Error{
				Message: message,
				Reason:  err.Error(),
			})
			return
		}
	}

	log.Debug().Msg("Audio file sent")
	c.Header("Content-Disposition", "inline; filename=\""+filepath.Base(staticAudioOutput.AbsolutePath)+"\"")
	c.Header("Content-Type", staticAudioOutput.Mime)
	c.Status(http.StatusOK)
	c.File(staticAudioOutput.AbsolutePath)
}
