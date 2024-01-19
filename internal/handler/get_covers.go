package handler

import (
	"github.com/gin-gonic/gin"
	"time"
)

type getCoversResponseCoverItem struct {
	ID                int       `json:"id"`
	DirID             int       `json:"dirId"`
	Filename          string    `json:"filename"`
	Extension         string    `json:"extension"`
	SizeByte          int64     `json:"sizeByte"`
	WidthPx           int       `json:"widthPx"`
	HeightPx          int       `json:"heightPx"`
	SHA256            string    `json:"sha256"`
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

type getCoversResponse struct {
	Covers []getCoversResponseCoverItem `json:"covers"`
}

func (h Handler) GetCovers(c *gin.Context) {

}
