package cover_service

import (
	"github.com/h2non/filetype"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

func (s Service) IsCoverByPath(path string) (isCover bool, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return false, err
	}

	file, err := os.Open(path)
	if err != nil {
		log.Warn().Err(err).Str("path", path).Msg("Failed to open file to check on image")
		return false, err
	}
	defer file.Close()

	if !strings.Contains(file.Name(), "cover") {
		return false, err
	}

	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		log.Warn().Err(err).Str("path", path).Msg("Failed to read file to check on image")
		return false, err
	}

	kind, _ := filetype.Match(head)
	if kind == filetype.Unknown {
		return false, nil
	}

	isCover = kind.MIME.Value == "image/jpeg" ||
		kind.MIME.Value == "image/png" ||
		kind.MIME.Value == "image/gif"

	return isCover, nil
}
