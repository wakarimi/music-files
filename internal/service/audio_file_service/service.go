package audio_file_service

import (
	"music-files/internal/database/repository/audio_file_repo"
)

type Service struct {
	AudioFileRepo audio_file_repo.Repo
}

func NewService(audioFileRepo audio_file_repo.Repo) (s *Service) {

	s = &Service{
		AudioFileRepo: audioFileRepo,
	}

	return s
}
