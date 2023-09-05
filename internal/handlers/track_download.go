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

func TrackDownload(c *gin.Context) {
	trackIdStr := c.Param("trackId")

	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid trackId format",
		})
		return
	}

	track, err := repository.GetTrackById(trackId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Invalid to get track",
		})
		return
	}

	dir, err := repository.GetDirById(track.DirId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get dir",
		})
		return
	}

	absolutePath := dir.Path + track.Path
	file, err := os.Open(absolutePath)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to open track file",
		})
		return
	}
	defer file.Close()

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(track.Path))

	c.File(absolutePath)
}
