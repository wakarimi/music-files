package repository

import (
	"music-files/internal/database"
	"music-files/internal/models"
)

func DeleteTracksByDirId(dirId int) (err error) {
	query := `
		DELETE FROM tracks
		WHERE dir_id = $1
	`
	_, err = database.Db.Exec(query, dirId)
	return err
}

func InsertMusicFile(musicFile models.Track) (musicFileId int, err error) {
	query := `
		INSERT INTO tracks(dir_id, path, size, format) 
		VALUES ($1, $2, $3, $4) 
		RETURNING track_id
	`
	err = database.Db.QueryRow(query, musicFile.DirId, musicFile.Path, musicFile.Size, musicFile.Format).Scan(&musicFileId)
	if err != nil {
		return 0, err
	}

	return musicFileId, nil
}

func DeleteMusicFile(musicFileId int) (err error) {
	query := `
		DELETE FROM tracks
		WHERE track_id = $1
	`
	_, err = database.Db.Exec(query, musicFileId)
	return err
}

func GetAllMusicFilesByDirId(dirId int) (musicFiles []models.Track, err error) {
	query := `
		SELECT track_id, dir_id, path, size, format, date_added
		FROM tracks
		WHERE dir_id = $1
	`

	rows, err := database.Db.Query(query, dirId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var musicFile models.Track
		if err := rows.Scan(&musicFile.TrackId, &musicFile.DirId, &musicFile.Path, &musicFile.Size,
			&musicFile.Format, &musicFile.DateAdded); err != nil {
			return nil, err
		}
		musicFiles = append(musicFiles, musicFile)
	}

	return musicFiles, rows.Err()
}
