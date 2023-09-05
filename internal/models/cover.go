package models

import "time"

type Cover struct {
	CoverId   int       `db:"cover_id"`
	DirId     int       `db:"dir_id"`
	Path      string    `db:"path"`
	Name      string    `db:"name"`
	Size      int64     `db:"size"`
	Format    string    `db:"format"`
	DateAdded time.Time `db:"date_added"`
}
