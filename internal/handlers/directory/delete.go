package directory

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-files/internal/handlers/types"
	"net/http"
	"strconv"
)

func (h *DirHandler) Delete(c *gin.Context) {
	log.Debug().Msg("Creating new directory")

	dirIdStr := c.Param("dirId")
	dirId, err := strconv.Atoi(dirIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid dirId format",
		})
		return
	}
	log.Debug().Int("dirId", dirId).Msg("Url parameter read successfully")

	err = h.TrackRepo.DeleteByDirId(dirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete tracks associated with the directory")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to delete tracks associated with the directory",
		})
		return
	}

	err = h.CoverRepo.DeleteByDirId(dirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete covers associated with the directory")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to delete covers associated with the directory",
		})
		return
	}

	err = h.DirRepo.Delete(dirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete directory")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to delete directory",
		})
		return
	}

	log.Debug().Msg("Directory deleted successfully")
	c.Status(http.StatusNoContent)
}
