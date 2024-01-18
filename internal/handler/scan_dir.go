package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/response"
	"music-files/internal/internal_error"
	"music-files/internal/use_case"
	"net/http"
	"strconv"
)

func (h Handler) ScanDir(c *gin.Context) {
	log.Debug().Msg("Scanning directory")

	lang := c.MustGet("lang").(string)
	localizer := i18n.NewLocalizer(&h.bundle, lang)

	dirIDStr := c.Param("dirId")
	dirID, err := strconv.Atoi(dirIDStr)
	if err != nil {
		log.Error().Err(err).Str("dirIdStr", dirIDStr).Msg("Invalid dirId format")
		messageID := "InvalidInputFormat"
		message, errLoc := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
		if errLoc != nil {
			message = h.engLocalizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
		}
		c.JSON(http.StatusNotFound, response.Error{
			Message: message,
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("dirId", dirID).Msg("Url parameter read successfully")

	scanDirInput := use_case.ScanDirInput{
		DirID: dirID,
	}
	_, err = h.useCase.ScanDir(scanDirInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete root")
		if _, ok := err.(internal_error.NotFound); ok {
			messageID := "DirNotFound"
			message, errLoc := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
			if errLoc != nil {
				message = h.engLocalizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
			}
			c.JSON(http.StatusNotFound, response.Error{
				Message: message,
				Reason:  err.Error(),
			})
			return
		} else {
			messageID := "FailedToScanDir"
			message, errLoc := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
			if errLoc != nil {
				message = h.engLocalizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
			}
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: message,
				Reason:  err.Error(),
			})
			return
		}
	}

	log.Debug().Msg("Scanning started")
	c.Status(http.StatusAccepted)
}
