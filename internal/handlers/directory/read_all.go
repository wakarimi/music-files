package directory

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"music-files/internal/models"
	"net/http"
	"time"
)

type readAllResponseItem struct {
	DirId       int        `json:"dirId"`
	Path        string     `json:"path"`
	LastScanned *time.Time `json:"lastScanned,omitempty"`
}

type readAllResponse struct {
	Dirs []readAllResponseItem `json:"directories"`
}

func (h *Handler) ReadAll(c *gin.Context) {
	log.Info().Msg("Fetching all directories")

	var dirs []models.Directory

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		dirs, err = h.DirService.ReadAll(tx)
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

	dirsResponse := make([]readAllResponseItem, 0)
	for _, dir := range dirs {
		dirResponse := readAllResponseItem{
			DirId:       dir.DirId,
			Path:        dir.Path,
			LastScanned: dir.LastScanned,
		}
		dirsResponse = append(dirsResponse, dirResponse)
	}

	log.Debug().Msg("All directories fetched successfully")
	c.JSON(http.StatusOK, readAllResponse{
		Dirs: dirsResponse,
	})
}
