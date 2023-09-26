package dir_service

import (
	"music-files/internal/database/repository/dir_repo"
	"music-files/internal/service/cover_service"
	"music-files/internal/service/track_service"
)

type Service struct {
	DirRepo dir_repo.Repo

	CoverService cover_service.Service
	TrackService track_service.Service
}

func NewService(dirRepo dir_repo.Repo,
	coverService cover_service.Service,
	trackService track_service.Service) (s *Service) {

	s = &Service{
		DirRepo:      dirRepo,
		CoverService: coverService,
		TrackService: trackService,
	}

	return s
}
