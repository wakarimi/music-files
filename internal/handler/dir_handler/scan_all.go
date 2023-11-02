package dir_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/response"
	"net/http"
)

// ScanAll scans all directories for new or updated files.
// @Summary Scan all directories
// @Description Initiates a scan in all directories to identify new or updated files.
// @Tags Directories
// @Accept  json
// @Produce  json
// @Success 200 "All directories scanned successfully"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /dirs/scan [post]
func (h *Handler) ScanAll(c *gin.Context) {
	log.Debug().Msg("Scanning directory")

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		err = h.DirService.ScanAll(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to scan directories")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to scan directories",
			Reason:  err.Error(),
		})
		return
	}

	log.Debug().Msg("Directories scanned successfully")
	c.Status(http.StatusOK)
}
