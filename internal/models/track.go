package models

import "time"

type Track struct {
	TrackId   int       `db:"track_id"`
	DirId     int       `db:"dir_id"`
	CoverId   *int      `db:"cover_id"`
	Path      string    `db:"path"`
	Name      string    `db:"name"`
	Size      int64     `db:"size"`
	Format    string    `db:"format"`
	DateAdded time.Time `db:"date_added"`
}
