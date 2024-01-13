package audio_service

import (
	"github.com/rs/zerolog/log"
	"github.com/wtolson/go-taglib"
	"music-files/internal/model/audio"
	"os"
	"path/filepath"
	"time"
)

func (s Service) ConstructByPath(path string) (audioFile audio.Audio, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Error().Str("absolutePath", path).Msg("Failed to get file info")
		return audio.Audio{}, err
	}

	fileDetails, err := taglib.Read(path)
	if err != nil {
		log.Error().Str("absolutePath", path).Msg("Failed to read file details")
		return audio.Audio{}, err
	}

	durationMs := int64(fileDetails.Length() / time.Millisecond)

	audioFile = audio.Audio{
		Filename:     fileInfo.Name(),
		Extension:    filepath.Ext(path),
		SizeByte:     fileInfo.Size(),
		DurationMs:   durationMs,
		BitrateKbps:  fileDetails.Bitrate(),
		SampleRateHz: fileDetails.Samplerate(),
		ChannelsN:    fileDetails.Channels(),
	}

	return audioFile, nil
}
