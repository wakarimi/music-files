package directory

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"net/http"
)

func (h *DirHandler) ScanAll(c *gin.Context) {
	log.Info().Msg("Scanning all directories")

	dirs, err := h.DirRepo.ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get directories")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get directories",
		})
		return
	}

	for _, dir := range dirs {
		log.Debug().Int("dirId", dir.DirId).Msg("Scanning directory")
		h.dirScan(c, dir)
	}

	log.Info().Msg("Directories scanned")
	c.Status(http.StatusOK)
}
