package directory

import (
	"github.com/gin-gonic/gin"
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

	exists, err := h.DirRepo.IsExistsByPath(request.Path)
	if err != nil {
		log.Debug().Msg("Failed to directory existence")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to directory existence",
		})
		return
	}
	if exists {
		log.Info().Msg("Directory already added")
		c.JSON(http.StatusConflict, types.Error{
			Error: "Directory already added",
		})
		return
	}

	dir := models.Directory{
		Path: request.Path,
	}
	_, err = h.DirRepo.Create(dir)
	if err != nil {
		log.Info().Msg("Failed to add directory")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to add directory",
		})
		return
	}

	log.Debug().Msg("Directory added successfully")
	c.Status(http.StatusCreated)
}
