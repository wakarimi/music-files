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

// Download godoc
// @Summary Download a cover by its ID
// @Tags Covers
// @Produce application/octet-stream
// @Param coverId path integer true "Cover Identifier"
// @Success 200 {file} byte "Successfully downloaded cover file"
// @Failure 400 {object} types.ErrorResponse "Invalid coverId format"
// @Failure 500 {object} types.ErrorResponse "Failed to fetch cover"
// @Router /covers/{coverId}/download [get]
func (h *Handler) Download(c *gin.Context) {
	log.Debug().Msg("Downloading cover")

	coverIdStr := c.Param("coverId")
	coverId, err := strconv.Atoi(coverIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid coverId format")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
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
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to fetch cover",
		})
		return
	}

	log.Debug().Msg("Cover sent successfully")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(absolutePath))
	c.File(absolutePath)
}
