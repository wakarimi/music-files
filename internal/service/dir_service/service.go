package dir_service

import "music-files/internal/database/repository/dir_repo"

type Service struct {
	DirRepo dir_repo.Repo
}

func NewService(dirRepo dir_repo.Repo) (s *Service) {

	s = &Service{
		DirRepo: dirRepo,
	}

	return s
}
