package models

import "time"

type Directory struct {
	DirId       int       `db:"dir_id"`
	Path        string    `db:"path"`
	DateAdded   time.Time `db:"date_added"`
	LastScanned time.Time `db:"last_scanned"`
}
