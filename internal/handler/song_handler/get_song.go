package song_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/handler/responses"
	"music-files/internal/models"
	"net/http"
	"strconv"
	"time"
)

type getDirResponse struct {
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

func (h *Handler) GetSong(c *gin.Context) {
	log.Debug().Msg("Getting song")

	songIdStr := c.Param("songId")
	songId, err := strconv.Atoi(songIdStr)
	if err != nil {
		log.Error().Err(err).Str("songIdStr", songIdStr).Msg("Invalid songId format")
		c.JSON(http.StatusBadRequest, responses.Error{
			Message: "Invalid songId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("songId", songId).Msg("Url parameter read successfully")

	var song models.Song
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		song, err = h.SongService.GetSong(tx, songId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get song")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, responses.Error{
				Message: "Song not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Message: "Failed to get song",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Song got successfully")
	c.JSON(http.StatusOK, getDirResponse{
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
	})
}
