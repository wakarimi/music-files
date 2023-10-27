package song_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/responses"
	"music-files/internal/models"
	"net/http"
	"time"
)

type searchBySha256ResponseItem struct {
	SongId            int       `json:"songId"`
	DirId             int       `json:"dirId"`
	Filename          string    `json:"filename"`
	Extension         string    `json:"extension"`
	SizeByte          int64     `json:"sizeByte"`
	DurationMs        int64     `json:"durationMs"`
	BitrateKbps       int       `json:"bitrateKbps"`
	SampleRateHz      int       `json:"sampleRateHz"`
	ChannelsN         int       `json:"channelsN"`
	Sha256            string    `json:"sha256"`
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

type searchBySha256Response struct {
	Songs []searchBySha256ResponseItem `json:"songs"`
}

func (h *Handler) SearchBySha256(c *gin.Context) {
	sha256 := c.Param("sha256")
	log.Debug().Str("sha256", sha256).Msg("Url parameter read successfully")

	var songs []models.Song
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		songs, err = h.SongService.SearchBySha256(tx, sha256)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: "Failed to get songs",
			Reason:  err.Error(),
		})
		return
	}

	songsResponseItems := make([]searchBySha256ResponseItem, len(songs))
	for i, song := range songs {
		songsResponseItems[i] = searchBySha256ResponseItem{
			SongId:            song.SongId,
			DirId:             song.DirId,
			Filename:          song.Filename,
			Extension:         song.Extension,
			SizeByte:          song.SizeByte,
			DurationMs:        song.DurationMs,
			BitrateKbps:       song.BitrateKbps,
			SampleRateHz:      song.SampleRateHz,
			ChannelsN:         song.ChannelsN,
			Sha256:            song.Sha256,
			LastContentUpdate: song.LastContentUpdate,
		}
	}

	c.JSON(http.StatusOK, searchBySha256Response{
		Songs: songsResponseItems,
	})
}
