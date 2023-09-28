package dir_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/responses"
	"music-files/internal/models"
	"net/http"
	"strconv"
)

type readContentItem struct {
	DirId int    `json:"dirId"`
	Name  string `json:"name"`
}

type readContentResponse struct {
	Dirs []readContentItem `json:"subdirs"`
}

func (h *Handler) ReadContent(c *gin.Context) {
	log.Debug().Msg("Fetching directory content")

	dirIdStr := c.Param("dirId")
	dirId, err := strconv.Atoi(dirIdStr)
	if err != nil {
		log.Error().Err(err)
		c.JSON(http.StatusBadRequest, responses.Error{
			Error: "Invalid dirId format",
		})
		return
	}
	log.Debug().Int("dirId", dirId).
		Msg("Url parameter read successfully")

	var dirs []models.Directory

	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		dirs, err = h.DirService.ReadContent(tx, dirId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err)
		c.JSON(http.StatusInternalServerError, responses.Error{
			Error: "Failed to fetch directory content",
		})
		return
	}

	dirsResponse := make([]readContentItem, len(dirs))
	for i, dir := range dirs {
		dirsResponse[i] = readContentItem{
			DirId: dir.DirId,
			Name:  dir.Name,
		}
	}

	response := readContentResponse{
		Dirs: dirsResponse,
	}

	log.Debug().Int("dirId", dirId).
		Msg("Directory content fetched successfully")
	c.JSON(http.StatusOK, response)
}
