package dir_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/responses"
	"music-files/internal/models"
	"net/http"
	"time"
)

// getRootsResponseItem is the model for each item in the getRoots response
type getRootsResponseItem struct {
	// Unique identifier for the directory
	DirId int `json:"dirId"`
	// Name of the directory
	Name string `json:"name"`
	// Last time the directory was scanned
	LastScanned *time.Time `json:"lastScanned,omitempty"`
}

// getRootsResponse is the response model for GetRoots API
type getRootsResponse struct {
	// Array containing root directories
	Dirs []getRootsResponseItem `json:"dirs"`
}

// GetRoots
// @Summary Retrieve root directories
// @Description Retrieves a list of all root directories that are tracked
// @Tags Directories
// @Accept  json
// @Produce  json
// @Success 200 {object} getRootsResponse
// @Failure 500 {object} responses.Error "Internal Server Error"
// @Router /roots [get]
func (h *Handler) GetRoots(c *gin.Context) {
	log.Debug().Msg("Getting root directories")

	var roots []models.Directory
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		roots, err = h.DirService.GetRoots(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get root directories")
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: "Failed to get root directories",
			Reason:  err.Error(),
		})
		return
	}

	responseRootItems := make([]getRootsResponseItem, len(roots))
	for i, root := range responseRootItems {
		responseRootItems[i] = getRootsResponseItem{
			DirId:       root.DirId,
			Name:        root.Name,
			LastScanned: root.LastScanned,
		}
	}

	log.Debug().Msg("Root directories got successfully")
	c.JSON(http.StatusOK, getRootsResponse{
		Dirs: responseRootItems,
	})
}
