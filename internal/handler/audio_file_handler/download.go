package audio_file_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/responses"
	"net/http"
	"path/filepath"
	"strconv"
)

// Download
// @Summary Download a audio file by ID
// @Description Downloads a audio file identified by the audioFileId
// @Tags AudioFiles
// @Accept  json
// @Produce  octet-stream
// @Param   audioFileId      path    int     true        "Audio File ID"
// @Success 200 {file} byte "Audio File"
// @Header 200 {string} Content-Type "application/octet-stream"
// @Header 200 {string} Content-Disposition "attachment; filename=[name of the file]"
// @Failure 400 {object} responses.Error "Invalid audioFileId format"
// @Failure 500 {object} responses.Error "Internal Server Error, Failed to calculate absolute path"
// @Router /audio-files/{audioFileId}/download [get]
func (h *Handler) Download(c *gin.Context) {
	log.Debug().Msg("Downloading audio file")

	audioFileIdStr := c.Param("audioFileId")
	audioFileId, err := strconv.Atoi(audioFileIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid audioFileId format")
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: "Invalid audioFileId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("audioFileId", audioFileId).Msg("Url parameter read successfully")

	var absolutePath string
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		absolutePath, err = h.FileProcessorService.AbsolutePathToAudioFile(tx, audioFileId)
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate absolute path")
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: "Failed to calculate absolute path",
			Reason:  err.Error(),
		})
		return
	}

	log.Debug().Msg("Audio file sent successfully")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(absolutePath))
	c.File(absolutePath)
}
