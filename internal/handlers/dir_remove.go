package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/types"
	"net/http"
	"strconv"
)

func DirRemove(c *gin.Context) {
	dirIdStr := c.Param("dirId")

	dirId, err := strconv.Atoi(dirIdStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid dirId format",
		})
		return
	}

	err = repository.DeleteTracksByDirId(dirId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to delete tracks associated with the directory",
		})
		return
	}

	err = repository.DeleteCoversByDirId(dirId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to delete covers associated with the directory",
		})
		return
	}

	err = repository.DeleteDir(dirId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to remove directory",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
