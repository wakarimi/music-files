package use_case

import (
	"time"
)

type GetDirsInput struct {
	DirIDs []int
}

type GetDirsOutputDirItem struct {
	ID          int
	Name        string
	ParentDirID *int
	LastScanned *time.Time
}

type GetDirsOutput struct {
	Dirs []GetDirsOutputDirItem
}

func (u UseCase) GetDirs(input GetDirsInput) (output GetDirsOutput, err error) {
	//TODO implement me
	panic("implement me")
}
