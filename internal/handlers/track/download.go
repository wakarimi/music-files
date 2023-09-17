package track

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"net/http"
	"path/filepath"
	"strconv"
)

// Download godoc
// @Summary Download a track by its ID
// @Tags Tracks
// @Produce application/octet-stream
// @Param trackId path integer true "Track Identifier"
// @Success 200 {file} byte "Successfully downloaded track file"
// @Failure 400 {object} types.ErrorResponse "Invalid trackId format"
// @Failure 500 {object} types.ErrorResponse "Failed to fetch track"
// @Router /tracks/{trackId}/download [get]
func (h *Handler) Download(c *gin.Context) {
	log.Debug().Msg("Downloading track")

	trackIdStr := c.Param("trackId")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid trackId format")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error: "Invalid trackId format",
		})
		return
	}
	log.Debug().Int("trackId", trackId).Msg("Url parameter read successfully")

	var absolutePath string

	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		absolutePath, err = h.TrackService.Download(tx, trackId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch track")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to fetch track",
		})
		return
	}

	log.Debug().Msg("Cover sent successfully")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(absolutePath))
	c.File(absolutePath)
}
