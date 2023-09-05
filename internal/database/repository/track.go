package repository

import (
	"music-files/internal/database"
	"music-files/internal/models"
)

func GetTrackById(trackId int) (track models.Track, err error) {
	query := `
		SELECT track_id, dir_id, cover_id, path, name, size, format, date_added
		FROM tracks
		WHERE track_id = $1
	`
	err = database.Db.QueryRow(query, trackId).Scan(&track.TrackId, &track.DirId, &track.CoverId, &track.Path,
		&track.Name, &track.Size, &track.Format, &track.DateAdded)
	if err != nil {
		return models.Track{}, err
	}

	return track, nil
}

func GetTracks() (tracks []models.Track, err error) {
	query := `
		SELECT track_id, dir_id, cover_id, path, name, size, format, date_added
		FROM tracks
	`

	rows, err := database.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var track models.Track
		if err := rows.Scan(&track.TrackId, &track.DirId, &track.CoverId, &track.Path, &track.Name, &track.Size,
			&track.Format, &track.DateAdded); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, rows.Err()
}

func DeleteTracksByDirId(dirId int) (err error) {
	query := `
		DELETE FROM tracks
		WHERE dir_id = $1
	`
	_, err = database.Db.Exec(query, dirId)
	return err
}

func InsertTrack(track models.Track) (trackId int, err error) {
	query := `
		INSERT INTO tracks(dir_id, cover_id, path, name, size, format) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING track_id
	`
	err = database.Db.QueryRow(query, track.DirId, track.CoverId, track.Path, track.Name, track.Size, track.Format).
		Scan(&trackId)
	if err != nil {
		return 0, err
	}

	return trackId, nil
}

func DeleteTrackById(trackId int) (err error) {
	query := `
		DELETE FROM tracks
		WHERE track_id = $1
	`
	_, err = database.Db.Exec(query, trackId)
	return err
}

func GetTracksByDirId(dirId int) (tracks []models.Track, err error) {
	query := `
		SELECT track_id, dir_id, cover_id, path, name, size, format, date_added
		FROM tracks
		WHERE dir_id = $1
	`

	rows, err := database.Db.Query(query, dirId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var track models.Track
		if err := rows.Scan(&track.TrackId, &track.DirId, &track.CoverId, &track.Path, &track.Name, &track.Size,
			&track.Format, &track.DateAdded); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, rows.Err()
}
