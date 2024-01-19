package handler

import (
	"github.com/gin-gonic/gin"
	"time"
)

type getAudiosResponseAudioItem struct {
	ID                int       `json:"id"`
	DirID             int       `json:"dirId"`
	Filename          string    `json:"filename"`
	Extension         string    `json:"extension"`
	SizeByte          int64     `json:"sizeByte"`
	DurationMs        int64     `json:"durationMs"`
	BitrateKbps       int       `json:"bitrateKbps"`
	SampleRateHz      int       `json:"sampleRateHz"`
	ChannelsN         int       `json:"channelsN"`
	SHA256            string    `json:"sha256"`
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

type getAudiosResponse struct {
	Audios []getAudiosResponseAudioItem `json:"audios"`
}

func (h Handler) GetAudios(c *gin.Context) {

}
