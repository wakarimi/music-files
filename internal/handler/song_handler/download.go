package song_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/responses"
	"net/http"
	"path/filepath"
	"strconv"
)

func (h *Handler) Download(c *gin.Context) {
	log.Debug().Msg("Downloading song")

	songIdStr := c.Param("songId")
	songId, err := strconv.Atoi(songIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid songId format")
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: "Invalid songId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("songId", songId).Msg("Url parameter read successfully")

	var absolutePath string
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		absolutePath, err = h.FileProcessorService.AbsolutePathToSong(tx, songId)
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate absolute path")
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: "Failed to calculate absolute path",
			Reason:  err.Error(),
		})
		return
	}

	log.Debug().Msg("Song sent successfully")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(absolutePath))
	c.File(absolutePath)
}
