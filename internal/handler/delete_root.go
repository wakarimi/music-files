package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/response"
	"music-files/internal/internal_error"
	"net/http"
	"strconv"
)

func (h Handler) DeleteRoot(c *gin.Context) {
	log.Debug().Msg("Deleting music root directory")

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

	deleteRootInput := DeleteRootInput{
		DirID: dirID,
	}
	_, err = h.useCase.DeleteRoot(deleteRootInput)
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
		} else if _, ok := err.(internal_error.BadRequest); ok {
			messageID := "DirectoryIsNotRoot"
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
			messageID := "FailedToDeleteRoot"
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

	log.Debug().Msg("Music root deleted")
	c.Status(http.StatusNoContent)
}
