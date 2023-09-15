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

func (h *Handler) Download(c *gin.Context) {
	log.Debug().Msg("Downloading track")

	trackIdStr := c.Param("trackId")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid trackId format")
		c.JSON(http.StatusBadRequest, types.Error{
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
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch track",
		})
		return
	}

	log.Debug().Msg("Cover sent successfully")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(absolutePath))
	c.File(absolutePath)
}
