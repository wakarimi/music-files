package cover_service

type coverRepo interface {
}

type Service struct {
	coverRepo coverRepo
}

func New(coverRepo coverRepo) *Service {
	return &Service{
		coverRepo: coverRepo,
	}
}
