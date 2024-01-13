package cover_service

import (
	"github.com/rs/zerolog/log"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"music-files/internal/model/cover"
	"os"
	"path/filepath"
)

func (s Service) ConstructByPath(path string) (coverFile cover.Cover, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to get file info")
		return cover.Cover{}, err
	}

	f, err := os.Open(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to open file")
		return cover.Cover{}, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close file")
		}
	}(f)

	img, _, err := image.DecodeConfig(f)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to decode config")
		return cover.Cover{}, err
	}

	coverFile = cover.Cover{
		Filename:  fileInfo.Name(),
		Extension: filepath.Ext(path),
		SizeByte:  fileInfo.Size(),
		WidthPx:   img.Width,
		HeightPx:  img.Height,
	}

	return coverFile, nil
}
