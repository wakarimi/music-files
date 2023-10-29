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

// getSongsResponseItem represents each song item in the getSongs response
type getSongsResponseItem struct {
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

// getSongsResponse is the response model for the GetAll API
type getSongsResponse struct {
	// Array containing song items
	Songs []getSongsResponseItem `json:"songs"`
}

// GetAll retrieves all songs
// @Summary Retrieve all songs
// @Description Retrieves a list of all songs in the system
// @Tags Songs
// @Accept  json
// @Produce  json
// @Success 200 {object} getSongsResponse
// @Failure 500 {object} responses.Error "Internal Server Error"
// @Router /songs [get]
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
	c.JSON(http.StatusOK, getSongsResponse{
		Songs: songsResponseItems,
	})
}
