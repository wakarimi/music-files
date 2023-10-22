package dir_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/handler/responses"
	"net/http"
	"strconv"
)

// UntrackRoot
// @Summary Remove a tracked root directory
// @Description Stops tracking a directory with a specified root directory ID
// @Tags Directories
// @Accept  json
// @Produce  json
// @Param   dirId path int true "Root Directory ID"
// @Success 204
// @Failure 400 {object} responses.Error "Invalid dirId format, The directory is not root"
// @Failure 404 {object} responses.Error "Directory not found"
// @Failure 500 {object} responses.Error "Internal Server Error"
// @Router  /roots/{dirId} [delete]
func (h *Handler) UntrackRoot(c *gin.Context) {
	log.Debug().Msg("Removing directories from tracked")

	dirIdStr := c.Param("dirId")
	dirId, err := strconv.Atoi(dirIdStr)
	if err != nil {
		log.Error().Err(err).Str("dirIdStr", dirIdStr).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, responses.Error{
			Message: "Invalid dirId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("dirId", dirId).Msg("Url parameter read successfully")

	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		err = h.DirService.UntrackRoot(tx, dirId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to stop tracking directory")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, responses.Error{
				Message: "Directory not found",
				Reason:  err.Error(),
			})
		} else if _, ok = err.(errors.BadRequest); ok {
			c.JSON(http.StatusBadRequest, responses.Error{
				Message: "The directory is not root",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Message: "Failed to stop tracking directory",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Directory removed from tracked")
	c.Status(http.StatusNoContent)
}
