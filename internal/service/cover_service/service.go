package cover_service

import "music-files/internal/database/repository/cover_repo"

type Service struct {
	CoverRepo cover_repo.Repo
}

func NewService(coverRepo cover_repo.Repo) (s *Service) {

	s = &Service{
		CoverRepo: coverRepo,
	}

	return s
}
