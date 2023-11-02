package model

import "time"

type Directory struct {
	DirId       int        `db:"dir_id"`
	Name        string     `db:"name"`
	ParentDirId *int       `db:"parent_dir_id"`
	LastScanned *time.Time `db:"last_scanned"`
}
