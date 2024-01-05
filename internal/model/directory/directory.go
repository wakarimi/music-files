package directory

import "time"

type Directory struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	ParentDirID *int       `db:"parent_dir_id"`
	LastScanned *time.Time `db:"last_scanned"`
}
