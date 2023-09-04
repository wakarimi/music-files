package repository

import (
	"music-files/internal/database"
	"music-files/internal/models"
)

func DeleteTracksByDirId(dirId int) (err error) {
	query := `
		DELETE FROM music_files
		WHERE dir_id = $1
	`
	_, err = database.Db.Exec(query, dirId)
	return err
}

func InsertMusicFile(musicFile models.MusicFile) (musicFileId int, err error) {
	query := `
		INSERT INTO music_files(dir_id, path, size, format) 
		VALUES ($1, $2, $3, $4) 
		RETURNING music_file_id
	`
	err = database.Db.QueryRow(query, musicFile.DirId, musicFile.Path, musicFile.Size, musicFile.Format).Scan(&musicFileId)
	if err != nil {
		return 0, err
	}

	return musicFileId, nil
}

func DeleteMusicFile(musicFileId int) (err error) {
	query := `
		DELETE FROM music_files
		WHERE music_file_id = $1
	`
	_, err = database.Db.Exec(query, musicFileId)
	return err
}

func GetAllMusicFiles() (musicFiles []models.MusicFile, err error) {
	query := `
		SELECT music_file_id, dir_id, path, size, format, date_added
		FROM music_files
	`

	rows, err := database.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var musicFile models.MusicFile
		if err := rows.Scan(&musicFile.MusicFileId, &musicFile.DirId, &musicFile.Path, &musicFile.Size,
			&musicFile.Format, &musicFile.DateAdded); err != nil {
			return nil, err
		}
		musicFiles = append(musicFiles, musicFile)
	}

	return musicFiles, rows.Err()
}
