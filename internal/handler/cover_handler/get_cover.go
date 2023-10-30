package cover_handler

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

// getCoverResponse is the response model for GetCover API
type getCoverResponse struct {
	// Unique identifier for the cover
	CoverId int `json:"coverId"`
	// File extension of the cover
	Extension string `json:"extension"`
	// File size of the cover in bytes
	SizeByte int64 `json:"sizeByte"`
	// Width of the cover in pixels
	WidthPx int `json:"widthPx"`
	// Height of the cover in pixels
	HeightPx int `json:"heightPx"`
	// SHA-256 hash of the cover
	Sha256 string `json:"sha256"`
	// Timestamp of the last content update
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

// GetCover
// @Summary Retrieve a cover by ID
// @Description Retrieves detailed information about a cover by its ID
// @Tags Covers
// @Accept  json
// @Produce  json
// @Param   coverId     path    int     true        "Cover ID"
// @Success 200 {object} getCoverResponse
// @Failure 400 {object} responses.Error "Invalid coverId format"
// @Failure 404 {object} responses.Error "Cover not found"
// @Failure 500 {object} responses.Error "Internal Server Error"
// @Router /covers/{coverId} [get]
func (h *Handler) GetCover(c *gin.Context) {
	log.Debug().Msg("Getting cover")

	coverIdStr := c.Param("coverId")
	coverId, err := strconv.Atoi(coverIdStr)
	if err != nil {
		log.Error().Err(err).Str("coverIdStr", coverIdStr).Msg("Invalid coverId format")
		c.JSON(http.StatusBadRequest, responses.Error{
			Message: "Invalid coverId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("coverId", coverId).Msg("Url parameter read successfully")

	var cover models.Cover
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		cover, err = h.CoverService.GetCover(tx, coverId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get cover")
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
