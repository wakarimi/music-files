package directory

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"music-files/internal/models"
	"net/http"
)

type createRequest struct {
	Path string `json:"path" bind:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	log.Debug().Msg("Creating new directory")

	var request createRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("Failed to encode request")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Failed to encode request",
		})
		return
	}

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		dir := models.Directory{
			Path: request.Path,
		}
		err = h.DirService.Create(tx, dir)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create directory")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to create directory",
		})
		return
	}

	log.Debug().Msg("Directory added successfully")
	c.Status(http.StatusCreated)
}
