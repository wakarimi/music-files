package song_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/handler/responses"
	"music-files/internal/models"
	"net/http"
	"strconv"
	"time"
)

type getCoverResponse struct {
	CoverId           int       `json:"coverId"`
	Extension         string    `json:"extension"`
	SizeByte          int64     `json:"sizeByte"`
	WidthPx           int       `json:"widthPx"`
	HeightPx          int       `json:"heightPx"`
	Sha256            string    `json:"sha256"`
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

func (h *Handler) GetCover(c *gin.Context) {
	log.Debug().Msg("Getting cover for song")

	songIdStr := c.Param("songId")
	songId, err := strconv.Atoi(songIdStr)
	if err != nil {
		log.Error().Err(err).Str("songIdStr", songIdStr).Msg("Invalid songId format")
		c.JSON(http.StatusBadRequest, responses.Error{
			Message: "Invalid songId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("songId", songId).Msg("Url parameter read successfully")

	var cover models.Cover
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		cover, err = h.FileProcessorService.GetCoverForSong(tx, songId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get cover")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, responses.Error{
				Message: "Cover not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Message: "Failed to get cover",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Cover got successfully")
	c.JSON(http.StatusOK, getCoverResponse{
		CoverId:           cover.CoverId,
		Extension:         cover.Extension,
		SizeByte:          cover.SizeByte,
		WidthPx:           cover.WidthPx,
		HeightPx:          cover.HeightPx,
		Sha256:            cover.Sha256,
		LastContentUpdate: cover.LastContentUpdate,
	})
}
