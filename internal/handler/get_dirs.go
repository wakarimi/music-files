package handler

import (
	"github.com/gin-gonic/gin"
	"time"
)

type getDirsResponseCoverItem struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	ParentDirID *int       `json:"parentDirId,omitempty"`
	LastScanned *time.Time `json:"lastScanned,omitempty"`
}

type getDirsResponse struct {
	Dirs []getDirsResponseCoverItem `json:"dirs"`
}

func (h Handler) GetDirs(context *gin.Context) {}
