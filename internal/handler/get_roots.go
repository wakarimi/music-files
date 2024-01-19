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

type getRootsResponseItem struct {
	DirID       int        `json:"dirId"`
	Path        string     `json:"absolutePath"`
	LastScanned *time.Time `json:"lastScanned,omitempty"`
}

type getRootsResponse struct {
	Dirs []getRootsResponseItem `json:"dirs"`
}

func (h Handler) GetRoots(c *gin.Context) {
	log.Debug().Msg("Getting root dirs")

	lang := c.MustGet("lang").(string)
	localizer := i18n.NewLocalizer(&h.bundle, lang)

	getRootsInput := use_case.GetRootsInput{}
	getRootsOutput, err := h.useCase.GetRoots(getRootsInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get roots")
		messageID := "FailedToGetRoots"
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

	responseDirItems := make([]getRootsResponseItem, len(getRootsOutput.Dirs))
	for i, root := range getRootsOutput.Dirs {
		responseDirItems[i] = getRootsResponseItem{
			DirID:       root.DirID,
			Path:        root.Path,
			LastScanned: root.LastScanned,
		}
	}

	log.Debug().Msg("Roots directories got")
	c.JSON(http.StatusOK, getRootsResponse{
		Dirs: responseDirItems,
	})
}
