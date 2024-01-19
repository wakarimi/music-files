package use_case

import "time"

type GetAudiosInput struct {
	AudioIDs []int
}

type GetAudiosOutputAudios struct {
	ID                int
	DirID             int
	Filename          string
	Extension         string
	SizeByte          int64
	DurationMs        int64
	BitrateKbps       int
	SampleRateHz      int
	ChannelsN         int
	SHA256            string
	LastContentUpdate time.Time
}

type GetAudiosOutput struct {
	Audios []GetAudiosOutputAudios
}

func (u UseCase) GetAudios(input GetAudiosInput) (output GetAudiosOutput, err error) {
	//TODO implement me
	panic("implement me")
}
