package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

type TrackRepositoryInterface interface {
	Create(track models.Track) (trackId int, err error)
	Read(trackId int) (track models.Track, err error)
	ReadAll() (tracks []models.Track, err error)
	ReadAllByDirId(dirId int) (tracks []models.Track, err error)
	Update(trackId int, track models.Track) (err error)
	Delete(trackId int) (err error)
	DeleteByDirId(dirId int) (err error)
	IsExists(trackId int) (exists bool, err error)
}

type TrackRepository struct {
	Db *sqlx.DB
}

func NewTrackRepository(db *sqlx.DB) TrackRepositoryInterface {
	return &TrackRepository{Db: db}
}

func (r *TrackRepository) Create(track models.Track) (trackId int, err error) {
	log.Debug().Str("filename", track.Filename).Msg("Creating new track")

	query := `
		INSERT INTO tracks(dir_id, cover_id, relative_path, filename, extension, size, hash)
		VALUES (:dir_id, :cover_id, :relative_path, :filename, :extension, :size, :hash)
		RETURNING track_id
	`
	rows, err := r.Db.NamedQuery(query, track)
	if err != nil {
		log.Error().Err(err).Str("filename", track.Filename).Msg("Failed to create track")
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&trackId); err != nil {
			log.Error().Err(err).Msg("Error scanning trackId from result set")
		}
	} else {
		return 0, fmt.Errorf("no id returned after track insert")
	}

	log.Debug().Int("trackId", trackId).Msg("Track created successfully")
	return trackId, nil
}

func (r *TrackRepository) Read(trackId int) (track models.Track, err error) {
	log.Debug().Int("trackId", trackId).Msg("Fetching track by ID")

	query := `
		SELECT *
		FROM tracks
		WHERE track_id = :track_id
	`
	args := map[string]interface{}{
		"track_id": trackId,
	}
	rows, err := r.Db.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to fetch track")
		return models.Track{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&track); err != nil {
			log.Error().Err(err).Msg("Error scanning row into struct")
			return models.Track{}, err
		}
	} else {
		err := fmt.Errorf("no track found with id: %d", trackId)
		log.Info().Err(err).Int("trackId", trackId).Msg("No track found")
		return models.Track{}, err
	}

	log.Debug().Int("dirId", track.DirId).Str("relativePath", track.RelativePath).Msg("Track fetched by ID successfully")
	return track, nil
}

func (r *TrackRepository) ReadAll() (tracks []models.Track, err error) {
	log.Debug().Msg("Fetching all tracks")

	query := `
		SELECT *
		FROM tracks
	`
	rows, err := r.Db.Queryx(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch tracks")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var track models.Track
		if err = rows.StructScan(&track); err != nil {
			log.Error().Err(err).Msg("Failed to scan track data")
			return nil, err
		}
		tracks = append(tracks, track)
	}

	log.Debug().Int("tracksCount", len(tracks)).Msg("All tracks fetched successfully")
	return tracks, nil
}

func (r *TrackRepository) ReadAllByDirId(dirId int) (tracks []models.Track, err error) {
	log.Debug().Msg("Fetching all tracks")

	query := `
		SELECT *
		FROM tracks
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	rows, err := r.Db.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch tracks")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var track models.Track
		if err = rows.StructScan(&track); err != nil {
			log.Error().Err(err).Msg("Failed to scan track data")
			return nil, err
		}
		tracks = append(tracks, track)
	}

	log.Debug().Int("tracksCount", len(tracks)).Msg("All tracks fetched successfully")
	return tracks, nil
}

func (r *TrackRepository) Update(trackId int, track models.Track) (err error) {
	log.Debug().Int("trackId", trackId).Msg("Updating track")

	query := `
		UPDATE tracks 
		SET dir_id = :dir_id, cover_id = :cover_id, relative_path = :relative_path, filename = :filename, extension = :extension, size = :size, hash = :hash
		WHERE track_id = :track_id
	`
	track.TrackId = trackId
	_, err = r.Db.NamedExec(query, track)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to update track")
		return err
	}

	log.Debug().Int("trackId", trackId).Msg("Track updated successfully")
	return nil
}

func (r *TrackRepository) Delete(trackId int) (err error) {
	log.Debug().Int("trackId", trackId).Msg("Deleting track")

	query := `
		DELETE FROM tracks
		WHERE track_id = :track_id
	`
	args := map[string]interface{}{
		"track_id": trackId,
	}
	_, err = r.Db.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to delete track")
		return err
	}

	log.Debug().Int("trackId", trackId).Msg("Track deleted successfully")
	return nil
}

func (r *TrackRepository) DeleteByDirId(dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting tracks by directory ID")

	query := `
		DELETE FROM tracks
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	_, err = r.Db.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to delete tracks by directory ID")
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Tracks deleted by directory ID successfully")
	return nil
}

func (r *TrackRepository) IsExists(trackId int) (exists bool, err error) {
	log.Debug().Int("trackId", trackId).Msg("Checking if track exists")

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM tracks
			WHERE track_id = :track_id
		)
	`
	args := map[string]interface{}{
		"track_id": trackId,
	}
	row, err := r.Db.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to execute query to track cover existence")
		return false, err
	}
	defer row.Close()
	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Int("trackId", trackId).Msg("Failed to scan result of track existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Int("trackId", trackId).Msg("Track exists")
	} else {
		log.Debug().Int("trackId", trackId).Msg("No track found")
	}
	return exists, nil
}
