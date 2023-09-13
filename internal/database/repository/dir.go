package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

type DirRepositoryInterface interface {
	Create(dir models.Directory) (dirId int, err error)
	CreateTx(tx *sqlx.Tx, dir models.Directory) (dirId int, err error)
	Read(dirId int) (cover models.Directory, err error)
	ReadTx(tx *sqlx.Tx, dirId int) (cover models.Directory, err error)
	ReadAll() (dirs []models.Directory, err error)
	ReadAllTx(tx *sqlx.Tx) (dirs []models.Directory, err error)
	UpdateLastScanned(dirId int) (err error)
	UpdateLastScannedTx(tx *sqlx.Tx, dirId int) (err error)
	Delete(dirId int) (err error)
	DeleteTx(tx *sqlx.Tx, dirId int) (err error)
	IsExistsByPath(path string) (exists bool, err error)
	IsExistsByPathTx(tx *sqlx.Tx, path string) (exists bool, err error)
}

type DirRepository struct {
	Db *sqlx.DB
}

func NewDirRepository(db *sqlx.DB) DirRepositoryInterface {
	return &DirRepository{Db: db}
}

func (r *DirRepository) Create(dir models.Directory) (dirId int, err error) {
	log.Debug().Str("path", dir.Path).Msg("Creating new directory")
	return r.create(r.Db, dir)
}

func (r *DirRepository) CreateTx(tx *sqlx.Tx, dir models.Directory) (dirId int, err error) {
	log.Debug().Str("path", dir.Path).Msg("Creating new directory transactional")
	return r.create(tx, dir)
}

func (r *DirRepository) create(queryer Queryer, dir models.Directory) (dirId int, err error) {
	query := `
		INSERT INTO directories(path)
		VALUES (:path)
		RETURNING dir_id
	`
	rows, err := queryer.NamedQuery(query, dir)
	if err != nil {
		log.Error().Err(err).Str("path", dir.Path).Msg("Failed to create directory")
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&dirId); err != nil {
			log.Error().Err(err).Msg("Error scanning dirId from result set")
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("no id returned after directory insert")
	}

	log.Debug().Int("dirId", dirId).Msg("Directory created successfully")
	return dirId, nil
}

func (r *DirRepository) Read(dirId int) (cover models.Directory, err error) {
	log.Debug().Int("dirId", dirId).Msg("Fetching directory by id")
	return r.read(r.Db, dirId)
}

func (r *DirRepository) ReadTx(tx *sqlx.Tx, dirId int) (cover models.Directory, err error) {
	log.Debug().Int("dirId", dirId).Msg("Fetching directory by id transactional")
	return r.read(tx, dirId)
}

func (r *DirRepository) read(queryer Queryer, dirId int) (dir models.Directory, err error) {
	query := `
		SELECT *
		FROM directories
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to fetch directory")
		return models.Directory{}, err
	}
	defer rows.Close()
	if rows.Next() {
		if err = rows.StructScan(&dir); err != nil {
			log.Error().Err(err).Int("dirId", dirId).Msg("Failed to scan directory data")
			return models.Directory{}, err
		}
	}

	log.Debug().Str("path", dir.Path).Msg("Directory fetched by id successfully")
	return dir, nil
}

func (r *DirRepository) ReadAll() (dirs []models.Directory, err error) {
	log.Debug().Msg("Fetching all directories")
	return r.readAll(r.Db)
}

func (r *DirRepository) ReadAllTx(tx *sqlx.Tx) (dirs []models.Directory, err error) {
	log.Debug().Msg("Fetching all directories transactional")
	return r.readAll(tx)
}

func (r *DirRepository) readAll(queryer Queryer) (dirs []models.Directory, err error) {
	query := `
		SELECT * 
		FROM directories
	`
	err = queryer.Select(&dirs, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch directories")
		return nil, err
	}

	log.Debug().Int("dirsCount", len(dirs)).Msg("All directories fetched successfully")
	return dirs, nil
}

func (r *DirRepository) UpdateLastScanned(dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Updating last scanned for directory")
	return r.updateLastScanned(r.Db, dirId)
}

func (r *DirRepository) UpdateLastScannedTx(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Updating last scanned for directory transactional")
	return r.updateLastScanned(tx, dirId)
}

func (r *DirRepository) updateLastScanned(queryer Queryer, dirId int) (err error) {
	query := `
		UPDATE directories
		SET last_scanned = CURRENT_TIMESTAMP
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	_, err = queryer.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to update last scanned for directory")
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Last scanned for directory updated successfully")
	return nil
}

func (r *DirRepository) Delete(dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting directory by id")
	return r.delete(r.Db, dirId)
}

func (r *DirRepository) DeleteTx(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting directory by id transactional")
	return r.delete(tx, dirId)
}

func (r *DirRepository) delete(queryer Queryer, dirId int) (err error) {
	query := `
		DELETE FROM directories
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	_, err = queryer.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to delete directory")
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Directory deleted successfully")
	return nil
}

func (r *DirRepository) IsExistsByPath(path string) (exists bool, err error) {
	log.Debug().Str("path", path).Msg("Checking if directory exists")
	return r.isExistsByPath(r.Db, path)
}

func (r *DirRepository) IsExistsByPathTx(tx *sqlx.Tx, path string) (exists bool, err error) {
	log.Debug().Str("path", path).Msg("Checking if directory exists transactional")
	return r.isExistsByPath(tx, path)
}

func (r *DirRepository) isExistsByPath(queryer Queryer, path string) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM directories
			WHERE path = :path
		)
	`
	args := map[string]interface{}{
		"path": path,
	}
	row, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to execute query to check directory existence")
		return false, err
	}
	defer row.Close()
	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Str("path", path).Msg("Failed to scan result of directory existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Str("path", path).Msg("Directory exists")
	} else {
		log.Debug().Str("path", path).Msg("No directory found")
	}
	return exists, nil
}
