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
)

type contentResponseDirItem struct {
	DirId int    `json:"dirId"`
	Name  string `json:"name"`
}

type contentResponse struct {
	Dirs []contentResponseDirItem `json:"dirs"`
}

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
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		subDirs, err = h.DirService.SubDirs(tx, dirId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get directory content")
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

	log.Debug().Msg("Directory content got successfully")
	c.JSON(http.StatusOK, contentResponse{
		Dirs: subDirsResponse,
	})
	c.Status(http.StatusOK)
}
