package track_service

import (
	"music-files/internal/database/repository/track_repo"
)

type Service struct {
	TrackRepo track_repo.Repo
}

func NewService(trackRepo track_repo.Repo) (s *Service) {

	s = &Service{
		TrackRepo: trackRepo,
	}

	return s
}
