package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/types"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func CoverDownload(c *gin.Context) {
	coverIdStr := c.Param("coverId")

	coverId, err := strconv.Atoi(coverIdStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid coverId format",
		})
		return
	}

	cover, err := repository.GetCoverById(coverId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid coverId format",
		})
		return
	}

	dir, err := repository.GetDirById(cover.DirId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Failed to get dir",
		})
		return
	}

	absolutePath := dir.Path + cover.Path
	file, err := os.Open(absolutePath)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to open cover file",
		})
		return
	}
	defer file.Close()

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(cover.Path))

	c.File(absolutePath)
}
