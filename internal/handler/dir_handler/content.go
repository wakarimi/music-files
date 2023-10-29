package dir_handler

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

// contentResponseDirItem represents a directory item in the content response
type contentResponseDirItem struct {
	// Unique identifier for the directory
	DirId int `json:"dirId"`
	// Name of the directory
	Name string `json:"name"`
}

// contentResponseDirItem represents a song item in the content response
type contentResponseSongItem struct {
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

// contentResponse is the response model for the Content API
type contentResponse struct {
	// Array containing directory items
	Dirs []contentResponseDirItem `json:"dirs"`
	// Array containing song items
	Songs []contentResponseSongItem `json:"songs"`
}

// Content
// @Summary Retrieve content of a directory by ID
// @Description Retrieves a list of subdirectories for a given directory ID
// @Tags Directories
// @Accept  json
// @Produce  json
// @Param   dirId     path    int     true        "Directory ID"
// @Success 200 {object} contentResponse
// @Failure 400 {object} responses.Error "Invalid dirId format"
// @Failure 404 {object} responses.Error "Directory not found"
// @Failure 500 {object} responses.Error "Internal Server Error"
// @Router /dirs/{dirId}/content [get]
func (h *Handler) Content(c *gin.Context) {
	log.Debug().Msg("Getting directory content")

	dirIdStr := c.Param("dirId")
	dirId, err := strconv.Atoi(dirIdStr)
	if err != nil {
		log.Error().Err(err).Str("dirIdStr", dirIdStr).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, responses.Error{
			Message: "Invalid dirId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("dirId", dirId).Msg("Url parameter read successfully")

	var subDirs []models.Directory
	var songs []models.Song
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		subDirs, err = h.DirService.SubDirs(tx, dirId)
		if err != nil {
			return err
		}
		songs, err = h.DirService.GetSongs(tx, dirId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get directory content")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, responses.Error{
				Message: "Directory not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Message: "Failed to get directory content",
				Reason:  err.Error(),
			})
		}
		return
	}

	subDirsResponse := make([]contentResponseDirItem, len(subDirs))
	for i := range subDirs {
		subDirsResponse[i].DirId = subDirs[i].DirId
		subDirsResponse[i].Name = subDirs[i].Name
	}

	songsResponse := make([]contentResponseSongItem, len(songs))
	for i, song := range songs {
		songsResponse[i] = contentResponseSongItem{
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

	log.Debug().Msg("Directory content got successfully")
	c.JSON(http.StatusOK, contentResponse{
		Dirs:  subDirsResponse,
		Songs: songsResponse,
	})
}
