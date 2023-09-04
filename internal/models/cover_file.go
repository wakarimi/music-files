package models

import "time"

type CoverFile struct {
	CoverFileId int       `db:"cover_file_id"`
	DirId       int       `db:"dir_id"`
	Path        string    `db:"path"`
	Size        int64     `db:"size"`
	Format      string    `db:"format"`
	DateAdded   time.Time `db:"date_added"`
}
