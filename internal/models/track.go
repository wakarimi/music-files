package models

import "time"

type Track struct {
	TrackId      int       `db:"track_id"`
	DirId        int       `db:"dir_id"`
	CoverId      *int      `db:"cover_id"`
	RelativePath string    `db:"relative_path"`
	Filename     string    `db:"filename"`
	Extension    string    `db:"extension"`
	Size         int64     `db:"size"`
	Hash         string    `db:"hash"`
	DateAdded    time.Time `db:"date_added"`
}
