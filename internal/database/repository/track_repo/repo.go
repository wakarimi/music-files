package track_repo

type Repo interface {
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}
