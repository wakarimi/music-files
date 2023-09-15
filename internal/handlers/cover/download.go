package cover

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
	log.Debug().Msg("Downloading cover")

	coverIdStr := c.Param("coverId")
	coverId, err := strconv.Atoi(coverIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid coverId format")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid coverId format",
		})
		return
	}
	log.Debug().Int("coverId", coverId).Msg("Url parameter read successfully")

	var absolutePath string

	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		absolutePath, err = h.CoverService.Download(tx, coverId)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch cover")
			return err
		}
		return err
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch cover")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch cover",
		})
		return
	}

	log.Debug().Msg("Cover sent successfully")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(absolutePath))
	c.File(absolutePath)
}
