package cover

import "time"

type Cover struct {
	ID                int       `db:"id"`
	DirID             int       `db:"dir_id"`
	Filename          string    `db:"filename"`
	Extension         string    `db:"extension"`
	SizeByte          int64     `db:"size_byte"`
	WidthPx           int       `db:"width_px"`
	HeightPx          int       `db:"height_px"`
	SHA256            string    `db:"sha_256"`
	LastContentUpdate time.Time `db:"last_content_update"`
}
