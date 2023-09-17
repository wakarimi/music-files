package directory

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"net/http"
	"strconv"
)

// Delete godoc
// @Summary Delete a directory
// @Description Deletes a directory using its unique identifier
// @Tags Directories
// @Accept  json
// @Produce  json
// @Param dirId path int true "Directory ID"
// @Success 204 {string} none "Successfully deleted directory"
// @Failure 400 {object} types.ErrorResponse "Invalid dirId format"
// @Failure 500 {object} types.ErrorResponse "Failed to delete directory"
// @Router /dirs/{dirId} [delete]
func (h *Handler) Delete(c *gin.Context) {
	log.Debug().Msg("Creating new directory")

	dirIdStr := c.Param("dirId")
	dirId, err := strconv.Atoi(dirIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error: "Invalid dirId format",
		})
		return
	}
	log.Debug().Int("dirId", dirId).Msg("Url parameter read successfully")

	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		err = h.DirService.Delete(tx, dirId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete directory")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to delete directory",
		})
		return
	}

	log.Debug().Msg("Directory deleted successfully")
	c.Status(http.StatusNoContent)
}
