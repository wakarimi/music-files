package dir_handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"music-files/internal/database/repository"
	"music-files/internal/handlers/types"
	"net/http"
)

type DirRemoveRequest struct {
	Path string `json:"path" binding:"required"`
}

func Remove(c *gin.Context) {
	var request DirRemoveRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.Error{
			Error: err.Error(),
		})
		return
	}

	err := repository.DeleteDir(request.Path)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to remove directory",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
