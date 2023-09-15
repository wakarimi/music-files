package directory

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"net/http"
	"strconv"
)

func (h *Handler) Scan(c *gin.Context) {
	log.Debug().Msg("Creating new directory")

	dirIdStr := c.Param("dirId")
	dirId, err := strconv.Atoi(dirIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid dirId format",
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
		log.Error().Err(err).Msg("Failed to scan directory")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to scan directory",
		})
		return
	}

	log.Debug().Int("dirId", dirId).Msg("Directory scanned successfully")
	c.Status(http.StatusOK)
}
