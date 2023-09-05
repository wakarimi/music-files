package repository

import (
	"music-files/internal/database"
	"music-files/internal/models"
)

func GetCoverById(trackId int) (cover models.Cover, err error) {
	query := `
		SELECT cover_id, dir_id, path, size, format, date_added
		FROM covers
		WHERE dir_id = $1
	`
	err = database.Db.QueryRow(query, trackId).Scan(&cover.CoverId, &cover.DirId, &cover.Path, &cover.Size, &cover.Format, &cover.DateAdded)
	if err != nil {
		return models.Cover{}, err
	}

	return cover, nil
}

func DeleteCoversByDirId(dirId int) (err error) {
	query := `
		DELETE FROM covers
		WHERE dir_id = $1
	`
	_, err = database.Db.Exec(query, dirId)
	return err
}

func InsertCover(cover models.Cover) (coverId int, err error) {
	query := `
		INSERT INTO covers(dir_id, path, size, format) 
		VALUES ($1, $2, $3, $4) 
		RETURNING cover_id
	`
	err = database.Db.QueryRow(query, cover.DirId, cover.Path, cover.Size, cover.Format).Scan(&coverId)
	if err != nil {
		return 0, err
	}

	return coverId, nil
}

func DeleteCover(coverId int) (err error) {
	query := `
		DELETE FROM covers
		WHERE cover_id = $1
	`
	_, err = database.Db.Exec(query, coverId)
	return err
}

func GetAllCoversByDirId(dirId int) (covers []models.Cover, err error) {
	query := `
		SELECT cover_id, dir_id, path, size, format, date_added
		FROM covers
		WHERE dir_id = $1
	`

	rows, err := database.Db.Query(query, dirId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cover models.Cover
		if err := rows.Scan(&cover.CoverId, &cover.DirId, &cover.Path, &cover.Size,
			&cover.Format, &cover.DateAdded); err != nil {
			return nil, err
		}
		covers = append(covers, cover)
	}

	return covers, rows.Err()
}
