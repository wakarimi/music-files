package file_processor_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

func (s *Service) AbsolutePathToAudioFile(tx *sqlx.Tx, audioFileId int) (absolutePath string, err error) {
	log.Debug().Int("audioFileId", audioFileId).Msg("Calculating absolute path to audio file")

	audioFile, err := s.AudioFileService.GetAudioFile(tx, audioFileId)
	if err != nil {
		log.Debug().Int("audioFileId", audioFileId).Msg("Failed to get audio file")
		return "", err
	}

	absolutePathToDir, err := s.DirService.AbsolutePath(tx, audioFile.DirId)
	if err != nil {
		log.Debug().Int("dirId", audioFile.DirId).Msg("Failed to calculate absolute path to directory")
		return "", err
	}

	absolutePath = filepath.Join(absolutePathToDir, audioFile.Filename)

	log.Debug().Int("audioFileId", audioFileId).Str("absolutePath", absolutePath).Msg("Calculating absolute path to audio file")
	return absolutePath, nil
}
