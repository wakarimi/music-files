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

// getSongResponse represents song in the getSong response
type getSongResponse struct {
	// Unique identifier for the song
	SongId int `json:"songId"`
	// Directory identifier where the song resides
	DirId int `json:"dirId"`
	// Filename of the song
	Filename string `json:"filename"`
	// File extension of the song
	Extension string `json:"extension"`
	// File size in bytes
	SizeByte int64 `json:"sizeByte"`
	// Duration of the song in milliseconds
	DurationMs int64 `json:"durationMs"`
	// Bitrate in kilobits per second
	BitrateKbps int `json:"bitrateKbps"`
	// Sample rate in hertz
	SampleRateHz int `json:"sampleRateHz"`
	// Number of audio channels
	ChannelsN int `json:"channelsN"`
	// SHA-256 hash of the file
	Sha256 string `json:"sha256"`
	// Time of the last update to the song's content
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

// GetSong retrieves a song by its identifier
// @Summary Retrieve a song by its ID
// @Description Retrieves a single song by its ID
// @Tags Songs
// @Accept  json
// @Produce  json
// @Param   songId path int true "Song ID"
// @Success 200 {object} getSongResponse
// @Failure 400,404,500 {object} responses.Error
// @Router /songs/{songId} [get]
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
	c.JSON(http.StatusOK, getSongResponse{
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
