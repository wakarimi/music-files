package cover_service

import "music-files/internal/database/repository"

type Service struct {
	CoverRepo repository.CoverRepositoryInterface
	DirRepo   repository.DirRepositoryInterface
}

func NewService(coverRepo repository.CoverRepositoryInterface,
	dirRepo repository.DirRepositoryInterface) (s *Service) {

	s = &Service{
		CoverRepo: coverRepo,
		DirRepo:   dirRepo,
	}

	return s
}
