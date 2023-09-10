package models

type Cover struct {
	CoverId      int    `db:"cover_id"`
	DirId        int    `db:"dir_id"`
	RelativePath string `db:"relative_path"`
	Filename     string `db:"filename"`
	Extension    string `db:"extension"`
	Size         int64  `db:"size"`
	Hash         string `db:"hash"`
}
