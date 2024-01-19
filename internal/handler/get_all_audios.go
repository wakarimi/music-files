package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/response"
	"music-files/internal/use_case"
	"net/http"
	"time"
)

type getAllAudiosResponseAudioItem struct {
	ID                int       `json:"id"`
	DirID             int       `json:"dirId"`
	SHA256            string    `json:"sha256"`
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

type getAllAudiosResponse struct {
	Audios []getAllAudiosResponseAudioItem `json:"audios"`
}

func (h Handler) GetAllAudios(c *gin.Context) {
	log.Debug().Msg("Getting all audios")

	lang := c.MustGet("lang").(string)
	localizer := i18n.NewLocalizer(&h.bundle, lang)

	getAllAudiosInput := use_case.GetAllAudiosInput{}
	getAllAudiosOutput, err := h.useCase.GetAllAudios(getAllAudiosInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all audios")
		messageID := "FailedToGetAllAudios"
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

	responseAudioItems := make([]getAllAudiosResponseAudioItem, len(getAllAudiosOutput.Audios))
	for i, audio := range getAllAudiosOutput.Audios {
		responseAudioItems[i] = getAllAudiosResponseAudioItem{
			ID:                audio.ID,
			DirID:             audio.DirID,
			SHA256:            audio.SHA256,
			LastContentUpdate: audio.LastContentUpdate,
		}
	}

	log.Debug().Msg("All audios got")
	c.JSON(http.StatusOK, getAllAudiosResponse{
		Audios: responseAudioItems,
	})
}
