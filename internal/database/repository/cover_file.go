package repository

import (
	"music-files/internal/database"
	"music-files/internal/models"
)

func InsertCoverFile(coverFile models.CoverFile) (coverFileId int, err error) {
	query := `
		INSERT INTO cover_files(dir_id, path, size, format) 
		VALUES ($1, $2, $3, $4) 
		RETURNING cover_file_id
	`
	err = database.Db.QueryRow(query, coverFile.DirId, coverFile.Path, coverFile.Size, coverFile.Format).Scan(&coverFileId)
	if err != nil {
		return 0, err
	}

	return coverFileId, nil
}

func DeleteCoverFile(coverFileId int) (err error) {
	query := `
		DELETE FROM cover_files
		WHERE cover_file_id = $1
	`
	_, err = database.Db.Exec(query, coverFileId)
	return err
}

func GetAllCoverFiles() (coverFiles []models.CoverFile, err error) {
	query := `
		SELECT cover_file_id, dir_id, path, size, format, date_added
		FROM cover_files
	`

	rows, err := database.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var coverFile models.CoverFile
		if err := rows.Scan(&coverFile.CoverFileId, &coverFile.DirId, &coverFile.Path, &coverFile.Size,
			&coverFile.Format, &coverFile.DateAdded); err != nil {
			return nil, err
		}
		coverFiles = append(coverFiles, coverFile)
	}

	return coverFiles, rows.Err()
}
