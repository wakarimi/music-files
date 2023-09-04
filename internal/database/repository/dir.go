package repository

import (
	"music-files/internal/database"
	"music-files/internal/models"
)

func InsertDir(path string) (dirId int, err error) {
	query := `
		INSERT INTO directories(path)
		VALUES ($1)
		RETURNING dir_id
	`
	err = database.Db.QueryRow(query, path).Scan(&dirId)
	if err != nil {
		return 0, err
	}

	return dirId, nil
}

func DirExist(path string) (exists bool, err error) {
	var count int

	query := `
		SELECT COUNT(dir_id)
		FROM directories
		WHERE path = $1
	`
	err = database.Db.Get(&count, query, path)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DeleteDir(path string) (err error) {
	query := `
		DELETE FROM directories
		WHERE path = $1
	`
	_, err = database.Db.Exec(query, path)
	return err
}

func GetAllDirs() (dirs []models.Directory, err error) {
	query := `
		SELECT dir_id, path, date_added, last_scanned
		FROM directories
	`

	rows, err := database.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dir models.Directory
		if err := rows.Scan(&dir.DirId, &dir.Path, &dir.DateAdded, &dir.LastScanned); err != nil {
			return nil, err
		}
		dirs = append(dirs, dir)
	}

	return dirs, rows.Err()
}
