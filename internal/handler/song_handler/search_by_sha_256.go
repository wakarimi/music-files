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

// searchBySha256ResponseItem represents a single song item in the search by SHA256 response.
type searchBySha256ResponseItem struct {
	// Unique identifier for the song.
	SongId int `json:"songId"`
	// Directory ID where the song is located.
	DirId int `json:"dirId"`
	// Filename of the song.
	Filename string `json:"filename"`
	// File extension of the song.
	Extension string `json:"extension"`
	// File size of the song in bytes.
	SizeByte int64 `json:"sizeByte"`
	// Duration of the song in milliseconds.
	DurationMs int64 `json:"durationMs"`
	// Bitrate of the song in Kbps.
	BitrateKbps int `json:"bitrateKbps"`
	// Sample rate of the song in Hz.
	SampleRateHz int `json:"sampleRateHz"`
	// Number of channels in the song.
	ChannelsN int `json:"channelsN"`
	// SHA-256 hash of the song.
	Sha256 string `json:"sha256"`
	// Timestamp of the last content update.
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

// searchBySha256Response represents the search by SHA256 API response.
type searchBySha256Response struct {
	// Array of songs that match the search query.
	Songs []searchBySha256ResponseItem `json:"songs"`
}

// SearchBySha256 retrieves a list of songs based on SHA256 hash.
// @Summary Search songs by SHA256 hash
// @Description Retrieves a list of songs that have the specified SHA256 hash.
// @Tags Songs
// @Accept  json
// @Produce  json
// @Param   sha256     path    string  true        "SHA256 hash"
// @Success 200 {object} searchBySha256Response
// @Failure 500 {object} responses.Error "Internal Server Error"
// @Router /songs/sha256/{sha256} [get]
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
