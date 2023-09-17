package directory

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"net/http"
)

// ScanAll godoc
// @Summary Scan all directories
// @Tags Directories
// @Accept json
// @Produce json
// @Success 200 {string} none "Directories scanned successfully"
// @Failure 500 {object} types.ErrorResponse "Failed to scan directories"
// @Router /dirs/scan-all [post]
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
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to scan directories",
		})
		return
	}

	log.Info().Msg("Directories scanned")
	c.Status(http.StatusOK)
}
