package file_processor_service

import (
	"music-files/internal/service/cover_service"
	"music-files/internal/service/dir_service"
	"music-files/internal/service/song_service"
)

type Service struct {
	DirService   dir_service.Service
	CoverService cover_service.Service
	SongService  song_service.Service
}

func NewService(dirService dir_service.Service,
	coverService cover_service.Service,
	songService song_service.Service) (s *Service) {

	s = &Service{
		DirService:   dirService,
		CoverService: coverService,
		SongService:  songService,
	}

	return s
}
