package dir_service

import "music-files/internal/database/repository"

type Service struct {
	DirRepo   repository.DirRepositoryInterface
	CoverRepo repository.CoverRepositoryInterface
	TrackRepo repository.TrackRepositoryInterface
}

func NewService(dirRepo repository.DirRepositoryInterface,
	coverRepo repository.CoverRepositoryInterface,
	trackRepo repository.TrackRepositoryInterface) (s *Service) {

	s = &Service{
		DirRepo: dirRepo,
	}

	return s
}
