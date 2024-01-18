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
	"time"
)

type getDirContentResponseDirItem struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	LastScanned *time.Time `json:"lastScanned"`
}

type getDirContentResponseAudioItem struct {
	ID                int       `json:"id"`
	DirID             int       `json:"dirId"`
	DurationMs        int64     `json:"durationMs"`
	SHA256            string    `json:"sha256"`
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

type getDirContentResponse struct {
	Dirs   []getDirContentResponseDirItem   `json:"dirs"`
	Audios []getDirContentResponseAudioItem `json:"audios"`
}

func (h Handler) GetDirContent(c *gin.Context) {
	log.Debug().Msg("Getting dir's content")

	lang := c.MustGet("lang").(string)
	localizer := i18n.NewLocalizer(&h.bundle, lang)

	dirIDStr := c.Param("dirId")
	dirID, err := strconv.Atoi(dirIDStr)
	if err != nil {
		log.Error().Err(err).Str("dirIDStr", dirIDStr).Msg("Invalid audioId format")
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
	log.Debug().Int("dirId", dirID).Msg("Url parameter read successfully")

	getDirContentInput := use_case.GetDirContentInput{
		DirID: dirID,
	}
	getDirContentOutput, err := h.useCase.GetDirContent(getDirContentInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get static link to an audio file")
		if _, ok := err.(internal_error.NotFound); ok {
			messageID := "DirNotFound"
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
			messageID := "FailedToGetDirContent"
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

	responseDirItems := make([]getDirContentResponseDirItem, len(getDirContentOutput.Dirs))
	for i, dir := range getDirContentOutput.Dirs {
		responseDirItems[i] = getDirContentResponseDirItem{
			ID:          dir.ID,
			Name:        dir.Name,
			LastScanned: dir.LastScanned,
		}
	}

	responseAudioItems := make([]getDirContentResponseAudioItem, len(getDirContentOutput.Audios))
	for i, audio := range getDirContentOutput.Audios {
		responseAudioItems[i] = getDirContentResponseAudioItem{
			ID:                audio.ID,
			DirID:             audio.DirID,
			DurationMs:        audio.DurationMs,
			SHA256:            audio.SHA256,
			LastContentUpdate: audio.LastContentUpdate,
		}
	}

	log.Debug().Msg("Dir's content got")
	c.JSON(http.StatusOK, getDirContentResponse{
		Dirs:   responseDirItems,
		Audios: responseAudioItems,
	})
}
