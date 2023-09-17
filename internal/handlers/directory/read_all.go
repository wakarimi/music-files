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

// readAllResponseItem godoc
// @Description Directory details
// @Property DirId (integer) The unique identifier of the directory
// @Property Path (string) Path of the directory
// @Property LastScanned (time, optional) Timestamp of the last scan for the directory
type readAllResponseItem struct {
	DirId       int        `json:"dirId"`
	Path        string     `json:"path"`
	LastScanned *time.Time `json:"lastScanned,omitempty"`
}

// readAllResponse godoc
// @Description List of directories
// @Property Dirs (array) Array containing details of directories
type readAllResponse struct {
	Dirs []readAllResponseItem `json:"directories"`
}

// ReadAll godoc
// @Summary Get all added directories for scanning
// @Tags Directories
// @Accept json
// @Produce json
// @Success 200 {object} readAllResponse "Successfully fetched all directories"
// @Failure 500 {object} types.ErrorResponse "Failed to fetch directories"
// @Router /dirs [get]
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
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
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
