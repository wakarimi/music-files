package file_processor_service

import (
	"music-files/internal/service/audio_file_service"
	"music-files/internal/service/cover_service"
	"music-files/internal/service/dir_service"
)

type Service struct {
	DirService       dir_service.Service
	CoverService     cover_service.Service
	AudioFileService audio_file_service.Service
}

func NewService(dirService dir_service.Service,
	coverService cover_service.Service,
	audioFileService audio_file_service.Service) (s *Service) {

	s = &Service{
		DirService:       dirService,
		CoverService:     coverService,
		AudioFileService: audioFileService,
	}

	return s
}
