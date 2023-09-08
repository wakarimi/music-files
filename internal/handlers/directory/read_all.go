package directory

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"net/http"
	"time"
)

type readAllResponseItem struct {
	DirId       int        `json:"dirId"`
	Path        string     `json:"path"`
	DateAdded   time.Time  `json:"dateAdded"`
	LastScanned *time.Time `json:"lastScanned,omitempty"`
}

type readAllResponse struct {
	Dirs []readAllResponseItem `json:"directories"`
}

func (h *DirHandler) ReadAll(c *gin.Context) {
	log.Info().Msg("Fetching all directories")

	dirs, err := h.DirRepo.ReadAll()
	if err != nil {
		log.Info().Msg("Failed to get all directories")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get all dirs",
		})
		return
	}

	dirsResponse := make([]readAllResponseItem, 0)
	for _, dir := range dirs {
		dirResponse := readAllResponseItem{
			DirId:       dir.DirId,
			Path:        dir.Path,
			DateAdded:   dir.DateAdded,
			LastScanned: dir.LastScanned,
		}
		dirsResponse = append(dirsResponse, dirResponse)
	}

	log.Debug().Msg("All directories fetched successfully")
	c.JSON(http.StatusOK, readAllResponse{
		Dirs: dirsResponse,
	})
}
