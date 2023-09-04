package models

import "time"

type MusicFile struct {
	MusicFileId int       `db:"music_file_id"`
	DirId       int       `db:"dir_id"`
	Path        string    `db:"path"`
	Size        int64     `db:"size"`
	Format      string    `db:"format"`
	DateAdded   time.Time `db:"date_added"`
}
