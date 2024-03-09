package model

import "time"

type File struct {
	ID         int
	DirID      int
	Name       string
	Extension  Extension
	SizeByte   int64
	LastUpdate time.Time
}
