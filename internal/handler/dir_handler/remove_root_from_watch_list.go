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

// RemoveRootFromWatchList
// @Summary Remove a tracked root directory
// @Description Stops tracking a directory with a specified root directory ID
// @Tags Directories
// @Accept  json
// @Produce  json
// @Param   dirId path int true "Root Directory ID"
// @Success 204
// @Failure 400 {object} response.Error "Invalid dirId format, The directory is not root"
// @Failure 404 {object} response.Error "Directory not found"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router  /roots/{dirId} [delete]
func (h *Handler) RemoveRootFromWatchList(c *gin.Context) {
	log.Debug().Msg("Removing directories from watch list")

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
		err = h.DirService.RemoveRootFromWatchList(tx, dirId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to remove directory from watch list")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Directory not found in database",
				Reason:  err.Error(),
			})
		} else if _, ok = err.(errors.BadRequest); ok {
			c.JSON(http.StatusBadRequest, response.Error{
				Message: "The directory is not root",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to remove directory from watch list",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Directory removed from tracked")
	c.Status(http.StatusNoContent)
}
