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

func (h Handler) StaticCover(c *gin.Context) {
	log.Debug().Msg("Getting a static url to an cover file")

	lang := c.MustGet("lang").(string)
	localizer := i18n.NewLocalizer(&h.bundle, lang)

	coverIDStr := c.Param("coverId")
	coverID, err := strconv.Atoi(coverIDStr)
	if err != nil {
		log.Error().Err(err).Str("coverIdStr", coverIDStr).Msg("Invalid coverId format")
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
	log.Debug().Int("coverId", coverID).Msg("Url parameter read successfully")

	staticCoverInput := use_case.StaticCoverInput{
		CoverID: coverID,
	}
	staticCoverOutput, err := h.useCase.StaticCover(staticCoverInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get static link to an cover file")
		if _, ok := err.(internal_error.NotFound); ok {
			messageID := "CoverNotFound"
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
			messageID := "FailedToGetStaticLinkToCover"
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

	log.Debug().Msg("Cover file sent")
	c.Header("Content-Disposition", "inline; filename=\""+filepath.Base(staticCoverOutput.AbsolutePath)+"\"")
	c.Header("Content-Type", staticCoverOutput.Mime)
	c.Status(http.StatusOK)
	c.File(staticCoverOutput.AbsolutePath)
}
