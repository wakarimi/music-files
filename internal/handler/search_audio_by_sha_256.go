package handler

import (
	"github.com/gin-gonic/gin"
	"time"
)

type SearchAudioBySHA256ResponseAudioItem struct {
	ID                int       `json:"id"`
	DirID             int       `json:"dirId"`
	SHA256            string    `json:"sha256"`
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

type SearchAudioBySHA256Response struct {
	Audios SearchAudioBySHA256ResponseAudioItem `json:"audios"`
}

func (h Handler) SearchAudioBySHA256(c *gin.Context) {

}
