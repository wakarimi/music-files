package model

import "time"

type Directory struct {
	ID           int
	ParentDirID  *int
	RelativePath string
	LastScanned  *time.Time
}
