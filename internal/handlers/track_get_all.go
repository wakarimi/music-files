package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/types"
	"net/http"
	"time"
)

type TrackGetAllResponseOne struct {
	TrackId   int       `json:"trackId"`
	DirId     int       `json:"dirId"`
	CoverId   *int      `json:"coverId,omitempty"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Format    string    `json:"format"`
	DateAdded time.Time `json:"dateAdded"`
}

type TrackGetAllResponse struct {
	Dirs []TrackGetAllResponseOne `json:"tracks" binding:"required"`
}

func TrackGetAll(c *gin.Context) {
	tracks, err := repository.GetTracks()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get all tracks",
		})
		return
	}

	tracksResponse := make([]TrackGetAllResponseOne, 0)
	for _, track := range tracks {
		trackResponse := TrackGetAllResponseOne{
			TrackId:   track.TrackId,
			DirId:     track.TrackId,
			CoverId:   track.CoverId,
			Path:      track.Path,
			Name:      track.Name,
			Size:      track.Size,
			Format:    track.Format,
			DateAdded: track.DateAdded,
		}
		tracksResponse = append(tracksResponse, trackResponse)
	}

	c.JSON(http.StatusOK, TrackGetAllResponse{
		Dirs: tracksResponse,
	})
}
