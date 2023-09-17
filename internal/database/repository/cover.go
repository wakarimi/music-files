package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

type CoverRepositoryInterface interface {
	Create(cover models.Cover) (coverId int, err error)
	CreateTx(tx *sqlx.Tx, cover models.Cover) (coverId int, err error)
	Read(coverId int) (cover models.Cover, err error)
	ReadTx(tx *sqlx.Tx, coverId int) (cover models.Cover, err error)
	ReadByDirIdAndRelativePath(dirId int, relativePath string) (cover models.Cover, err error)
	ReadByDirIdAndRelativePathTx(tx *sqlx.Tx, dirId int, relativePath string) (cover models.Cover, err error)
	ReadAllByDirId(dirId int) (covers []models.Cover, err error)
	ReadAllByDirIdTx(tx *sqlx.Tx, dirId int) (covers []models.Cover, err error)
	Update(coverId int, cover models.Cover) (err error)
	UpdateTx(tx *sqlx.Tx, coverId int, cover models.Cover) (err error)
	Delete(coverId int) (err error)
	DeleteTx(tx *sqlx.Tx, coverId int) (err error)
	DeleteByDirId(dirId int) (err error)
	DeleteByDirIdTx(tx *sqlx.Tx, dirId int) (err error)
	IsExists(coverId int) (exists bool, err error)
	IsExistsTx(tx *sqlx.Tx, coverId int) (exists bool, err error)
}

type CoverRepository struct {
	Db *sqlx.DB
}

func NewCoverRepository(db *sqlx.DB) CoverRepositoryInterface {
	return &CoverRepository{Db: db}
}

func (r *CoverRepository) Create(cover models.Cover) (coverId int, err error) {
	log.Debug().Str("filename", cover.Filename).Msg("Creating new cover")
	return r.create(r.Db, cover)
}

func (r *CoverRepository) CreateTx(tx *sqlx.Tx, cover models.Cover) (coverId int, err error) {
	log.Debug().Str("filename", cover.Filename).Msg("Creating new cover transactional")
	return r.create(tx, cover)
}

func (r *CoverRepository) create(queryer Queryer, cover models.Cover) (coverId int, err error) {
	query := `
		INSERT INTO covers(dir_id, relative_path, filename, format, width_px, height_px, size, hash_sha_256)
		VALUES (:dir_id, :relative_path, :filename, :format, :width_px, :height_px, :size, :hash_sha_256)
		RETURNING cover_id
	`
	rows, err := queryer.NamedQuery(query, cover)
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
	log.Debug().Int("coverId", coverId).Msg("Fetching cover by id")
	return r.read(r.Db, coverId)
}

func (r *CoverRepository) ReadTx(tx *sqlx.Tx, coverId int) (cover models.Cover, err error) {
	log.Debug().Int("coverId", coverId).Msg("Fetching cover by id transactional")
	return r.read(tx, coverId)
}

func (r *CoverRepository) read(queryer Queryer, coverId int) (cover models.Cover, err error) {
	query := `
		SELECT *
		FROM covers
		WHERE cover_id = :cover_id
	`
	args := map[string]interface{}{
		"cover_id": coverId,
	}
	rows, err := queryer.NamedQuery(query, args)
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

	log.Debug().Int("dirId", cover.DirId).Str("relativePath", cover.RelativePath).Msg("Cover fetched by id successfully")
	return cover, nil
}

func (r *CoverRepository) ReadByDirIdAndRelativePath(dirId int, relativePath string) (cover models.Cover, err error) {
	log.Debug().Int("dirId", dirId).Str("relativePath", relativePath).Msg("Fetching cover by id")
	return r.readByDirIdAndRelativePath(r.Db, dirId, relativePath)
}

func (r *CoverRepository) ReadByDirIdAndRelativePathTx(tx *sqlx.Tx, dirId int, relativePath string) (cover models.Cover, err error) {
	log.Debug().Int("dirId", dirId).Str("relativePath", relativePath).Msg("Fetching cover by id transactional")
	return r.readByDirIdAndRelativePath(tx, dirId, relativePath)
}

func (r *CoverRepository) readByDirIdAndRelativePath(queryer Queryer, dirId int, relativePath string) (cover models.Cover, err error) {
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
	rows, err := queryer.NamedQuery(query, args)
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
		log.Debug().Int("dirId", dirId).Str("relativePath", relativePath).Msg("No cover found")
		return models.Cover{}, err
	}

	log.Debug().Int("dirId", dirId).Str("relativePath", relativePath).Msg("Cover fetched successfully")
	return cover, nil
}

func (r *CoverRepository) ReadAllByDirId(dirId int) (covers []models.Cover, err error) {
	log.Debug().Int("dirId", dirId).Msg("Fetching covers")
	return r.readAllByDirId(r.Db, dirId)
}

func (r *CoverRepository) ReadAllByDirIdTx(tx *sqlx.Tx, dirId int) (covers []models.Cover, err error) {
	log.Debug().Int("dirId", dirId).Msg("Fetching covers transactional")
	return r.readAllByDirId(tx, dirId)
}

func (r *CoverRepository) readAllByDirId(queryer Queryer, dirId int) (covers []models.Cover, err error) {
	query := `
		SELECT *
		FROM covers
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	rows, err := queryer.NamedQuery(query, args)
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

	log.Debug().Int("coversCount", len(covers)).Msg("Covers fetched successfully")
	return covers, nil
}

func (r *CoverRepository) Update(coverId int, cover models.Cover) (err error) {
	log.Debug().Int("coverId", coverId).Msg("Updating cover")
	return r.update(r.Db, coverId, cover)
}

func (r *CoverRepository) UpdateTx(tx *sqlx.Tx, coverId int, cover models.Cover) (err error) {
	log.Debug().Int("coverId", coverId).Msg("Updating cover transactional")
	return r.update(tx, coverId, cover)
}

func (r *CoverRepository) update(queryer Queryer, coverId int, cover models.Cover) (err error) {
	query := `
		UPDATE covers 
		SET
		    dir_id = :dir_id, relative_path = :relative_path, filename = :filename, format = :format,
		    width_px = :width_px, height_px = :height_px, size = :size, hash_sha_256 = :hash_sha_256
		WHERE cover_id = :cover_id
	`
	cover.CoverId = coverId
	_, err = queryer.NamedExec(query, cover)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverId).Msg("Failed to update cover")
		return err
	}

	log.Debug().Int("coverId", coverId).Msg("Cover updated successfully")
	return nil
}

func (r *CoverRepository) Delete(coverId int) (err error) {
	log.Debug().Int("coverId", coverId).Msg("Deleting cover")
	return r.delete(r.Db, coverId)
}

func (r *CoverRepository) DeleteTx(tx *sqlx.Tx, coverId int) (err error) {
	log.Debug().Int("coverId", coverId).Msg("Deleting cover transactional")
	return r.delete(tx, coverId)
}

func (r *CoverRepository) delete(queryer Queryer, coverId int) (err error) {
	query := `
		DELETE FROM covers
		WHERE cover_id = :cover_id
	`
	args := map[string]interface{}{
		"cover_id": coverId,
	}
	_, err = queryer.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverId).Msg("Failed to delete cover")
		return err
	}

	log.Debug().Int("coverId", coverId).Msg("Cover deleted successfully")
	return nil
}

func (r *CoverRepository) DeleteByDirId(dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting covers by directory id")
	return r.deleteByDirId(r.Db, dirId)
}

func (r *CoverRepository) DeleteByDirIdTx(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting covers by directory id")
	return r.deleteByDirId(tx, dirId)
}

func (r *CoverRepository) deleteByDirId(queryer Queryer, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting covers by directory id")

	query := `
		DELETE FROM covers
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	_, err = queryer.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to delete covers by directory id")
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Covers deleted by directory id successfully")
	return nil
}

func (r *CoverRepository) IsExists(coverId int) (exists bool, err error) {
	log.Debug().Int("coverId", coverId).Msg("Checking if cover exists")
	return r.isExists(r.Db, coverId)
}

func (r *CoverRepository) IsExistsTx(tx *sqlx.Tx, coverId int) (exists bool, err error) {
	log.Debug().Int("coverId", coverId).Msg("Checking if cover exists")
	return r.isExists(tx, coverId)
}

func (r *CoverRepository) isExists(queryer Queryer, coverId int) (exists bool, err error) {
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
	row, err := queryer.NamedQuery(query, args)
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
