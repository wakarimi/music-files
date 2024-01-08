package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/response"
	"music-files/internal/internal_error"
	"net/http"
)

type addRootRequest struct {
	Path string `json:"path" bind:"required"`
}

type addRootResponse struct {
	DirId int    `json:"dirId"`
	Path  string `json:"path"`
}

func (h Handler) AddRoot(c *gin.Context) {
	log.Debug().Msg("Adding music root directory")

	lang := c.MustGet("lang").(string)
	localizer := i18n.NewLocalizer(&h.bundle, lang)

	var request addRootRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("Failed to encode request")
		messageID := "FailedToEncodeRequest"
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

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		log.Error().Err(err).Msg("Validation failed for request")
		messageID := "ValidationFailedForRequest"
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

	addRootInput := AddRootInput{
		Path: request.Path,
	}
	addRootOutput, err := h.useCase.AddRoot(addRootInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add music root dir")
		if _, ok := err.(internal_error.NotFound); ok {
			messageID := "DirectoryNotFoundOnDisk"
			message, errLoc := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
			if errLoc != nil {
				message = h.engLocalizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
			}
			c.JSON(http.StatusNotFound, response.Error{
				Message: message,
				Reason:  err.Error(),
			})
			return
		} else if _, ok := err.(internal_error.Conflict); ok {
			messageID := "TheDirectoryIsAlreadyBeingTracked"
			message, errLoc := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
			if errLoc != nil {
				message = h.engLocalizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
			}
			c.JSON(http.StatusConflict, response.Error{
				Message: message,
				Reason:  err.Error(),
			})
			return
		} else {
			messageID := "FailedToAddMusicRoot"
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

	log.Debug().Msg("Music root added")
	c.JSON(http.StatusOK, addRootResponse{
		DirId: addRootOutput.DirID,
		Path:  addRootOutput.Path,
	})
}
