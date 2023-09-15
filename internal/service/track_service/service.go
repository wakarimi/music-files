package track_service

import "music-files/internal/database/repository"

type Service struct {
	TrackRepo repository.TrackRepositoryInterface
	DirRepo   repository.DirRepositoryInterface
}

func NewService(trackRepo repository.TrackRepositoryInterface,
	dirRepo repository.DirRepositoryInterface) (s *Service) {

	s = &Service{
		TrackRepo: trackRepo,
		DirRepo:   dirRepo,
	}

	return s
}
