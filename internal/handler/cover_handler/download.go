package cover_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/response"
	"net/http"
	"path/filepath"
	"strconv"
)

// Download
// @Summary Download a cover image by ID
// @Description Downloads a cover image file identified by the coverId
// @Tags Covers
// @Accept  json
// @Produce  octet-stream
// @Param   coverId     path    int     true        "Cover ID"
// @Success 200 {file} byte "Cover File"
// @Header 200 {string} Content-Type "application/octet-stream"
// @Header 200 {string} Content-Disposition "attachment; filename=[name of the file]"
// @Failure 400 {object} response.Error "Invalid coverId format"
// @Failure 500 {object} response.Error "Internal Server Error, Failed to calculate absolute path"
// @Router /covers/{coverId}/download [get]
func (h *Handler) Download(c *gin.Context) {
	log.Debug().Msg("Downloading cover")

	coverIdStr := c.Param("coverId")
	coverId, err := strconv.Atoi(coverIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid coverId format")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Invalid coverId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("coverId", coverId).Msg("Url parameter read successfully")

	var absolutePath string
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		absolutePath, err = h.FileProcessorService.AbsolutePathToCover(tx, coverId)
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate absolute path")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to calculate absolute path",
			Reason:  err.Error(),
		})
		return
	}

	log.Debug().Msg("Cover sent successfully")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(absolutePath))
	c.File(absolutePath)
}
