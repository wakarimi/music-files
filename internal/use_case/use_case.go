package use_case

import "github.com/jmoiron/sqlx"

type transactor interface {
	WithTransaction(do func(tx *sqlx.Tx) (err error)) (err error)
}

type audioService interface {
}

type coverService interface {
}

type dirService interface {
}

type UseCase struct {
	transactor   transactor
	audioService audioService
	coverService coverService
	dirService   dirService
}

func New(transactor transactor,
	audioService audioService,
	coverService coverService,
	dirService dirService) *UseCase {
	return &UseCase{
		transactor:   transactor,
		audioService: audioService,
		coverService: coverService,
		dirService:   dirService,
	}
}
