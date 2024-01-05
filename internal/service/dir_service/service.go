package dir_service

type dirRepo interface {
}

type Service struct {
	dirRepo dirRepo
}

func New(dirRepo dirRepo) *Service {
	return &Service{
		dirRepo: dirRepo,
	}
}
