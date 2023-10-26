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

type getDirResponse struct {
	DirId        int        `json:"dirId"`
	Name         string     `json:"name"`
	AbsolutePath string     `json:"absolutePath"`
	LastScanned  *time.Time `json:"lastScanned,omitempty"`
}

func (h *Handler) GetDir(c *gin.Context) {
	log.Debug().Msg("Getting directory")

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

	var dir models.Directory
	var absolutePath string
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		dir, err = h.DirService.GetDir(tx, dirId)
		if err != nil {
			return err
		}
		absolutePath, err = h.DirService.AbsolutePath(tx, dirId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get directory")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, responses.Error{
				Message: "Directory not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Message: "Failed to get content",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Directory content got successfully")
	c.JSON(http.StatusOK, getDirResponse{
		DirId:        dir.DirId,
		Name:         dir.Name,
		AbsolutePath: absolutePath,
		LastScanned:  dir.LastScanned,
	})
}
