package audio_service

type audioRepo interface {
}

type Service struct {
	audioRepo audioRepo
}

func New(audioRepo audioRepo) *Service {
	return &Service{
		audioRepo: audioRepo,
	}
}
