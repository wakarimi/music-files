package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

type CoverRepositoryInterface interface {
	Create(cover models.Cover) (coverId int, err error)
	Read(coverId int) (cover models.Cover, err error)
	ReadByDirIdAndRelativePath(dirId int, relativePath string) (cover models.Cover, err error)
	ReadAllByDirId(dirId int) (covers []models.Cover, err error)
	Update(coverId int, cover models.Cover) (err error)
	Delete(coverId int) (err error)
	DeleteByDirId(dirId int) (err error)
	IsExists(coverId int) (exists bool, err error)
}

type CoverRepository struct {
	Db *sqlx.DB
}

func NewCoverRepository(db *sqlx.DB) CoverRepositoryInterface {
	return &CoverRepository{Db: db}
}

func (r *CoverRepository) Create(cover models.Cover) (coverId int, err error) {
	log.Debug().Str("filename", cover.Filename).Msg("Creating new cover")

	query := `
		INSERT INTO covers(dir_id, relative_path, filename, extension, size, hash)
		VALUES (:dir_id, :relative_path, :filename, :extension, :size, :hash)
		RETURNING cover_id
	`
	rows, err := r.Db.NamedQuery(query, cover)
	if err != nil {
		log.Error().Err(err).Str("filename", cover.Filename).Msg("Failed to create cover")
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&coverId); err != nil {
			log.Error().Err(err).Msg("Error scanning coverId from result set")
		}
	} else {
		return 0, fmt.Errorf("no id returned after cover insert")
	}

	log.Debug().Int("coverId", coverId).Msg("Cover created successfully")
	return coverId, nil
}

func (r *CoverRepository) Read(coverId int) (cover models.Cover, err error) {
	log.Debug().Int("coverId", coverId).Msg("Fetching cover by ID")

	query := `
		SELECT *
		FROM covers
		WHERE cover_id = :cover_id
	`
	args := map[string]interface{}{
		"cover_id": coverId,
	}
	rows, err := r.Db.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverId).Msg("Failed to fetch cover")
		return models.Cover{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&cover); err != nil {
			log.Error().Err(err).Msg("Error scanning row into struct")
			return models.Cover{}, err
		}
	} else {
		err := fmt.Errorf("no cover found with id: %d", coverId)
		log.Info().Err(err).Int("coverId", coverId).Msg("No cover found")
		return models.Cover{}, err
	}

	log.Debug().Int("dirId", cover.DirId).Str("relativePath", cover.RelativePath).Msg("Cover fetched by ID successfully")
	return cover, nil
}

func (r *CoverRepository) ReadByDirIdAndRelativePath(dirId int, relativePath string) (cover models.Cover, err error) {
	log.Debug().Int("dirId", dirId).Str("relativePath", relativePath).Msg("Fetching cover by ID")

	query := `
		SELECT *
		FROM covers
		WHERE dir_id = :dir_id 
		  AND relative_path = :relative_path
	`
	args := map[string]interface{}{
		"dir_id":        dirId,
		"relative_path": relativePath,
	}
	rows, err := r.Db.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Str("relativePath", relativePath).Msg("Failed to fetch cover")
		return models.Cover{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&cover); err != nil {
			log.Error().Err(err).Msg("Error scanning row into struct")
			return models.Cover{}, err
		}
	} else {
		err := fmt.Errorf("no cover found with dir_id: %d, relativePath: %s", dirId, relativePath)
		log.Info().Err(err).Int("dirId", dirId).Str("relativePath", relativePath).Msg("No cover found")
		return models.Cover{}, err
	}

	log.Debug().Int("coverId", cover.CoverId).Msg("Cover fetched by ID successfully")
	return cover, nil
}

func (r *CoverRepository) ReadAllByDirId(dirId int) (covers []models.Cover, err error) {
	log.Debug().Msg("Fetching all covers")

	query := `
		SELECT *
		FROM covers
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	rows, err := r.Db.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch covers")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cover models.Cover
		if err = rows.StructScan(&cover); err != nil {
			log.Error().Err(err).Msg("Failed to scan cover data")
			return nil, err
		}
		covers = append(covers, cover)
	}

	log.Debug().Int("coversCount", len(covers)).Msg("All covers fetched successfully")
	return covers, nil
}

func (r *CoverRepository) Update(coverId int, cover models.Cover) (err error) {
	log.Debug().Int("coverId", coverId).Msg("Updating cover")

	query := `
		UPDATE covers 
		SET dir_id = :dir_id, relative_path = :relative_path, filename = :filename, extension = :extension, size = :size, hash = :hash
		WHERE cover_id = :cover_id
	`
	cover.CoverId = coverId
	_, err = r.Db.NamedExec(query, cover)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverId).Msg("Failed to update cover")
		return err
	}

	log.Debug().Int("coverId", coverId).Msg("Cover updated successfully")
	return nil
}

func (r *CoverRepository) Delete(coverId int) (err error) {
	log.Debug().Int("coverId", coverId).Msg("Deleting cover")

	query := `
		DELETE FROM covers
		WHERE cover_id = :cover_id
	`
	args := map[string]interface{}{
		"cover_id": coverId,
	}
	_, err = r.Db.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverId).Msg("Failed to delete cover")
		return err
	}

	log.Debug().Int("coverId", coverId).Msg("Cover deleted successfully")
	return nil
}

func (r *CoverRepository) DeleteByDirId(dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting covers by directory ID")

	query := `
		DELETE FROM covers
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	_, err = r.Db.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to delete covers by directory ID")
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Covers deleted by directory ID successfully")
	return nil
}

func (r *CoverRepository) IsExists(coverId int) (exists bool, err error) {
	log.Debug().Int("coverId", coverId).Msg("Checking if cover exists")

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM covers
			WHERE cover_id = :cover_id
		)
	`
	args := map[string]interface{}{
		"cover_id": coverId,
	}
	row, err := r.Db.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverId).Msg("Failed to execute query to check cover existence")
		return false, err
	}
	defer row.Close()
	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Int("coverId", coverId).Msg("Failed to scan result of cover existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Int("coverId", coverId).Msg("Cover exists")
	} else {
		log.Debug().Int("coverId", coverId).Msg("No cover found")
	}
	return exists, nil
}
