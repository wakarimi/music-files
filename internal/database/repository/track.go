package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

type TrackRepositoryInterface interface {
	Create(track models.Track) (trackId int, err error)
	CreateTx(tx *sqlx.Tx, track models.Track) (trackId int, err error)
	Read(trackId int) (track models.Track, err error)
	ReadTx(tx *sqlx.Tx, trackId int) (track models.Track, err error)
	ReadAll() (tracks []models.Track, err error)
	ReadAllTx(tx *sqlx.Tx) (tracks []models.Track, err error)
	ReadAllByDirId(dirId int) (tracks []models.Track, err error)
	ReadAllByDirIdTx(tx *sqlx.Tx, dirId int) (tracks []models.Track, err error)
	Update(trackId int, track models.Track) (err error)
	UpdateTx(tx *sqlx.Tx, trackId int, track models.Track) (err error)
	ResetCover(coverId int) (err error)
	ResetCoverTx(tx *sqlx.Tx, coverId int) (err error)
	Delete(trackId int) (err error)
	DeleteTx(tx *sqlx.Tx, trackId int) (err error)
	DeleteByDirId(dirId int) (err error)
	DeleteByDirIdTx(tx *sqlx.Tx, dirId int) (err error)
	IsExists(trackId int) (exists bool, err error)
	IsExistsTx(tx *sqlx.Tx, trackId int) (exists bool, err error)
}

type TrackRepository struct {
	Db *sqlx.DB
}

func NewTrackRepository(db *sqlx.DB) TrackRepositoryInterface {
	return &TrackRepository{Db: db}
}

func (r *TrackRepository) Create(track models.Track) (trackId int, err error) {
	log.Debug().Str("filename", track.Filename).Msg("Creating new track")
	return r.create(r.Db, track)
}

func (r *TrackRepository) CreateTx(tx *sqlx.Tx, track models.Track) (trackId int, err error) {
	log.Debug().Str("filename", track.Filename).Msg("Creating new track transactional")
	return r.create(tx, track)
}

func (r *TrackRepository) create(queryer Queryer, track models.Track) (trackId int, err error) {
	query := `
		INSERT INTO tracks(dir_id, cover_id, relative_path, filename, duration_ms, size_byte, audio_codec, bitrate_kbps, sample_rate_hz, channels, hash_sha_256)
		VALUES (:dir_id, :cover_id, :relative_path, :filename, :duration_ms, :size_byte, :audio_codec, :bitrate_kbps, :sample_rate_hz, :channels, :hash_sha_256)
		RETURNING track_id
	`
	rows, err := queryer.NamedQuery(query, track)
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
	log.Debug().Int("trackId", trackId).Msg("Fetching track by id")
	return r.read(r.Db, trackId)
}

func (r *TrackRepository) ReadTx(tx *sqlx.Tx, trackId int) (track models.Track, err error) {
	log.Debug().Int("trackId", trackId).Msg("Fetching track by id transactional")
	return r.read(tx, trackId)
}

func (r *TrackRepository) read(queryer Queryer, trackId int) (track models.Track, err error) {
	query := `
		SELECT *
		FROM tracks
		WHERE track_id = :track_id
	`
	args := map[string]interface{}{
		"track_id": trackId,
	}
	rows, err := queryer.NamedQuery(query, args)
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

	log.Debug().Int("dirId", track.DirId).Str("relativePath", track.RelativePath).Msg("Track fetched by id successfully")
	return track, nil
}

func (r *TrackRepository) ReadAll() (tracks []models.Track, err error) {
	log.Debug().Msg("Fetching all tracks")
	return r.readAll(r.Db)
}

func (r *TrackRepository) ReadAllTx(tx *sqlx.Tx) (tracks []models.Track, err error) {
	log.Debug().Msg("Fetching all tracks transactional")
	return r.readAll(tx)
}

func (r *TrackRepository) readAll(queryer Queryer) (tracks []models.Track, err error) {
	query := `
		SELECT *
		FROM tracks
	`
	rows, err := queryer.Queryx(query)
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
	log.Debug().Int("dirId", dirId).Msg("Fetching all tracks in directory")
	return r.readAllByDirId(r.Db, dirId)
}

func (r *TrackRepository) ReadAllByDirIdTx(tx *sqlx.Tx, dirId int) (tracks []models.Track, err error) {
	log.Debug().Int("dirId", dirId).Msg("Fetching all tracks in directory transactional")
	return r.readAllByDirId(tx, dirId)
}

func (r *TrackRepository) readAllByDirId(queryer Queryer, dirId int) (tracks []models.Track, err error) {
	query := `
		SELECT *
		FROM tracks
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	rows, err := queryer.NamedQuery(query, args)
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
	return r.update(r.Db, trackId, track)
}

func (r *TrackRepository) UpdateTx(tx *sqlx.Tx, trackId int, track models.Track) (err error) {
	log.Debug().Int("trackId", trackId).Msg("Updating track transactional")
	return r.update(tx, trackId, track)
}

func (r *TrackRepository) update(queryer Queryer, trackId int, track models.Track) (err error) {
	query := `
		UPDATE tracks 
		SET dir_id = :dir_id, cover_id = :cover_id, relative_path = :relative_path, filename = :filename,
		    duration_ms = :duration_ms, size_byte = :size_byte, audio_codec = :audio_codec, bitrate_kbps = :bitrate_kbps,
		    sample_rate_hz = :sample_rate_hz, channels = :channels, hash_sha_256 = :hash_sha_256
		WHERE track_id = :track_id
	`
	track.TrackId = trackId
	_, err = queryer.NamedExec(query, track)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to update track")
		return err
	}

	log.Debug().Int("trackId", trackId).Msg("Track updated successfully")
	return nil
}

func (r *TrackRepository) ResetCover(coverId int) (err error) {
	log.Debug().Int("coverId", coverId).Msg("Resetting cover_id for tracks")
	return r.resetCover(r.Db, coverId)
}

func (r *TrackRepository) ResetCoverTx(tx *sqlx.Tx, coverId int) (err error) {
	log.Debug().Int("coverId", coverId).Msg("Resetting cover_id for tracks transactional")
	return r.resetCover(tx, coverId)
}

func (r *TrackRepository) resetCover(queryer Queryer, coverId int) (err error) {
	query := `
		UPDATE tracks
		SET cover_id = NULL
		WHERE cover_id = :cover_id
	`
	args := map[string]interface{}{
		"cover_id": coverId,
	}
	_, err = queryer.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverId).Msg("Failed to reset cover for tracks")
		return err
	}

	log.Debug().Int("coverId", coverId).Msg("cover for tracks reset successfully")
	return nil
}

func (r *TrackRepository) Delete(trackId int) (err error) {
	log.Debug().Int("trackId", trackId).Msg("Deleting track")
	return r.delete(r.Db, trackId)
}

func (r *TrackRepository) DeleteTx(tx *sqlx.Tx, trackId int) (err error) {
	log.Debug().Int("trackId", trackId).Msg("Deleting track transactional")
	return r.delete(tx, trackId)
}

func (r *TrackRepository) delete(queryer Queryer, trackId int) (err error) {
	query := `
		DELETE FROM tracks
		WHERE track_id = :track_id
	`
	args := map[string]interface{}{
		"track_id": trackId,
	}
	_, err = queryer.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to delete track")
		return err
	}

	log.Debug().Int("trackId", trackId).Msg("Track deleted successfully")
	return nil
}

func (r *TrackRepository) DeleteByDirId(dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting tracks by directory id")
	return r.deleteByDirId(r.Db, dirId)
}

func (r *TrackRepository) DeleteByDirIdTx(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting tracks by directory id transactional")
	return r.deleteByDirId(tx, dirId)
}

func (r *TrackRepository) deleteByDirId(queryer Queryer, dirId int) (err error) {
	query := `
		DELETE FROM tracks
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	_, err = queryer.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to delete tracks by directory id")
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Tracks deleted by directory id successfully")
	return nil
}

func (r *TrackRepository) IsExists(trackId int) (exists bool, err error) {
	log.Debug().Int("trackId", trackId).Msg("Checking if track exists")
	return r.isExists(r.Db, trackId)
}

func (r *TrackRepository) IsExistsTx(tx *sqlx.Tx, trackId int) (exists bool, err error) {
	log.Debug().Int("trackId", trackId).Msg("Checking if track exists transactional")
	return r.isExists(tx, trackId)
}

func (r *TrackRepository) isExists(queryer Queryer, trackId int) (exists bool, err error) {
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
	row, err := queryer.NamedQuery(query, args)
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
