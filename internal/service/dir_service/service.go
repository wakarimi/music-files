package dir_service

import (
	"music-files/internal/database/repository/dir_repo"
	"music-files/internal/service/audio_file_service"
	"music-files/internal/service/cover_service"
)

type Service struct {
	DirRepo dir_repo.Repo

	CoverService     cover_service.Service
	AudioFileService audio_file_service.Service
}

func NewService(dirRepo dir_repo.Repo,
	coverService cover_service.Service,
	audioFileService audio_file_service.Service) (s *Service) {

	s = &Service{
		DirRepo:          dirRepo,
		CoverService:     coverService,
		AudioFileService: audioFileService,
	}

	return s
}
