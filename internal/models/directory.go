package models

import "time"

type Directory struct {
	DirId       int        `db:"dir_id"`
	Path        string     `db:"path"`
	LastScanned *time.Time `db:"last_scanned"`
}
