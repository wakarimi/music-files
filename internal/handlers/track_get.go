package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/types"
	"net/http"
	"strconv"
	"time"
)

type TrackGetResponse struct {
	TrackId   int       `json:"trackId"`
	DirId     int       `json:"dirId"`
	CoverId   *int      `json:"coverId,omitempty"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	Format    string    `json:"format"`
	DateAdded time.Time `json:"dateAdded"`
}

func TrackGet(c *gin.Context) {
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
			Error: "Failed to get track",
		})
		return
	}

	c.JSON(http.StatusOK, TrackGetResponse{
		TrackId:   track.TrackId,
		DirId:     track.DirId,
		CoverId:   track.CoverId,
		Path:      track.Path,
		Size:      track.Size,
		Format:    track.Format,
		DateAdded: track.DateAdded,
	})

}
