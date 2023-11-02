package audio_file_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/handler/response"
	"music-files/internal/model"
	"net/http"
	"strconv"
	"time"
)

// getCoverResponse represents the response model for GetCoverForAudioFile API.
type getCoverResponse struct {
	// Unique identifier for the cover.
	CoverId int `json:"coverId"`
	// File extension of the cover.
	Extension string `json:"extension"`
	// File size of the cover in bytes.
	SizeByte int64 `json:"sizeByte"`
	// Width of the cover in pixels.
	WidthPx int `json:"widthPx"`
	// Height of the cover in pixels.
	HeightPx int `json:"heightPx"`
	// SHA-256 hash of the cover.
	Sha256 string `json:"sha256"`
	// Timestamp of the last content update.
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

// GetCover retrieves a cover for a specific audioFile.
// @Summary Retrieve a cover for a audioFile by ID
// @Description Retrieves detailed information about a cover for a audioFile by its ID.
// @Tags Covers
// @Accept  json
// @Produce  json
// @Param   audioFileId     path    int     true        "AudioFile ID"
// @Success 200 {object} getCoverResponse
// @Failure 400 {object} response.Error "Invalid audioFileId format"
// @Failure 404 {object} response.Error "Cover not found"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /audio-files/{audioFileId}/cover [get]
func (h *Handler) GetCover(c *gin.Context) {
	log.Debug().Msg("Getting cover for audioFile")

	audioFileIdStr := c.Param("audioFileId")
	audioFileId, err := strconv.Atoi(audioFileIdStr)
	if err != nil {
		log.Error().Err(err).Str("audioFileIdStr", audioFileIdStr).Msg("Invalid audioFileId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid audioFileId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("audioFileId", audioFileId).Msg("Url parameter read successfully")

	var cover model.Cover
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		cover, err = h.FileProcessorService.GetCoverForAudioFile(tx, audioFileId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get cover")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Cover not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
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
