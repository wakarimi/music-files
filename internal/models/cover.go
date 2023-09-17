package models

type Cover struct {
	CoverId      int    `db:"cover_id"`
	DirId        int    `db:"dir_id"`
	RelativePath string `db:"relative_path"`
	Filename     string `db:"filename"`
	Format       string `db:"format"`
	WidthPx      int    `db:"width_px"`
	HeightPx     int    `db:"height_px"`
	Size         int64  `db:"size"`
	HashSha256   string `db:"hash_sha_256"`
}
