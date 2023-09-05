package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/types"
	"net/http"
)

func DirScanAll(c *gin.Context) {
	dirs, err := repository.GetAllDirs()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Failed to get directory",
		})
		return
	}

	for _, dir := range dirs {
		dirScanOne(c, dir)
	}
}
