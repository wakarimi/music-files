package directory

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"net/http"
)

func (h *Handler) ScanAll(c *gin.Context) {
	log.Info().Msg("Scanning all directories")

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		err = h.DirService.ScanAll(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to scan directory")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to scan directories",
		})
		return
	}

	log.Info().Msg("Directories scanned")
	c.Status(http.StatusOK)
}
