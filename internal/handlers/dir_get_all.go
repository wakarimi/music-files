package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/types"
	"net/http"
	"time"
)

type DirGetAllResponseOne struct {
	DirId       int        `json:"dirId"`
	Path        string     `json:"path"`
	DateAdded   time.Time  `json:"dateAdded"`
	LastScanned *time.Time `json:"lastScanned,omitempty"`
}

type DirGetAllResponse struct {
	Dirs []DirGetAllResponseOne `json:"dirs" binding:"required"`
}

func DirGetAll(c *gin.Context) {
	dirs, err := repository.GetAllDirs()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get all dirs",
		})
		return
	}

	dirsResponse := make([]DirGetAllResponseOne, 0)
	for _, dir := range dirs {
		dirResponse := DirGetAllResponseOne{
			DirId:       dir.DirId,
			Path:        dir.Path,
			DateAdded:   dir.DateAdded,
			LastScanned: dir.LastScanned,
		}
		dirsResponse = append(dirsResponse, dirResponse)
	}

	c.JSON(http.StatusOK, DirGetAllResponse{
		Dirs: dirsResponse,
	})
}
