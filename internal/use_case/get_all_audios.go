package use_case

import "time"

type GetAllAudiosInput struct {
}

type GetAllAudiosOutputAudios struct {
	ID                int
	DirID             int
	DurationMs        int64
	SHA256            string
	LastContentUpdate time.Time
}

type GetAllAudiosOutput struct {
	Audios []GetAllAudiosOutputAudios
}

func (u UseCase) GetAllAudios(input GetAllAudiosInput) (output GetAllAudiosOutput, err error) {
	//TODO implement me
	panic("implement me")
}
