package dir_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/handler/response"
	"net/http"
	"strconv"
)

// Scan scans a directory for new or updated files.
// @Summary Scan a directory by ID
// @Description Initiates a scan in the specified directory to identify new or updated files.
// @Tags Directories
// @Accept  json
// @Produce  json
// @Param   dirId     path    int     true        "Directory ID"
// @Success 200 "Directory scanned successfully"
// @Failure 400 {object} response.Error "Invalid dirId format"
// @Failure 404 {object} response.Error "Directory not found"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /dirs/{dirId}/scan [post]
func (h *Handler) Scan(c *gin.Context) {
	log.Debug().Msg("Scanning directory")

	dirIdStr := c.Param("dirId")
	dirId, err := strconv.Atoi(dirIdStr)
	if err != nil {
		log.Error().Err(err).Str("dirIdStr", dirIdStr).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid dirId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("dirId", dirId).Msg("Url parameter read successfully")

	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		err = h.DirService.Scan(tx, dirId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get scan directory")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Directory not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get scan directory",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Directory scanned successfully")
	c.Status(http.StatusOK)
}
