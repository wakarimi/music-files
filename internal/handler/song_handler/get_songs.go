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

type getSongsResponseItem struct {
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

type getRootsResponse struct {
	Songs []getSongsResponseItem `json:"songs"`
}

func (h *Handler) GetAll(c *gin.Context) {
	log.Debug().Msg("Getting songs")

	var songs []models.Song
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		songs, err = h.SongService.GetAll(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get songs")
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: "Failed to get songs",
			Reason:  err.Error(),
		})
		return
	}

	songsResponseItems := make([]getSongsResponseItem, len(songs))
	for i, song := range songs {
		songsResponseItems[i] = getSongsResponseItem{
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

	log.Debug().Msg("Songs got successfully")
	c.JSON(http.StatusOK, getRootsResponse{
		Songs: songsResponseItems,
	})
}
