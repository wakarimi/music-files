package use_case

import "time"

type GetCoversInput struct {
	CoverIDs []int
}

type GetCoversOutputCovers struct {
	ID                int
	DirID             int
	Filename          string
	Extension         string
	SizeByte          int64
	WidthPx           int
	HeightPx          int
	SHA256            string
	LastContentUpdate time.Time
}

type GetCoversOutput struct {
	Covers []GetCoversOutputCovers
}

func (u UseCase) GetCovers(input GetCoversInput) (output GetCoversOutput, err error) {
	//TODO implement me
	panic("implement me")
}
