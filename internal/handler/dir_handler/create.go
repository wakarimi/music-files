package dir_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handler/responses"
	"music-files/internal/models"
	"net/http"
)

type createRequest struct {
	Path string `json:"path" bind:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	log.Debug().
		Msg("Adding new root directory")

	var request createRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err)
		c.JSON(http.StatusBadRequest, responses.Error{
			Error: "Failed to encode request",
		})
		return
	}

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		dir := models.Directory{
			ParentDirId: nil,
			Name:        request.Path,
		}
		err = h.DirService.Create(tx, dir)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err)
		c.JSON(http.StatusInternalServerError, responses.Error{
			Error: "Failed to create directory",
		})
		return
	}

	log.Debug().
		Msg("Root directory added successfully")
	c.Status(http.StatusCreated)
}
