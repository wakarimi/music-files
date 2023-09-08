package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

type DirRepositoryInterface interface {
	Create(dir models.Directory) (dirId int, err error)
	Read(dirId int) (cover models.Directory, err error)
	ReadAll() (dirs []models.Directory, err error)
	UpdateLastScanned(dirId int) (err error)
	Delete(dirId int) (err error)
	IsExistsByPath(path string) (exists bool, err error)
}

type DirRepository struct {
	Db *sqlx.DB
}

func NewDirRepository(db *sqlx.DB) DirRepositoryInterface {
	return &DirRepository{Db: db}
}

func (r *DirRepository) Create(dir models.Directory) (dirId int, err error) {
	log.Debug().Str("path", dir.Path).Msg("Creating new directory")

	query := `
		INSERT INTO directories(path)
		VALUES (:path)
		RETURNING dir_id
	`
	rows, err := r.Db.NamedQuery(query, dir)
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

func (r *DirRepository) Read(dirId int) (dir models.Directory, err error) {
	log.Debug().Int("dirId", dirId).Msg("Fetching directory by ID")

	query := `
		SELECT *
		FROM directories
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	rows, err := r.Db.NamedQuery(query, args)
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

	log.Debug().Str("path", dir.Path).Msg("Directory fetched by ID successfully")
	return dir, nil
}

func (r *DirRepository) ReadAll() (dirs []models.Directory, err error) {
	log.Debug().Msg("Fetching all directories")

	query := `
		SELECT * 
		FROM directories
	`
	err = r.Db.Select(&dirs, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch directories")
		return nil, err
	}

	log.Debug().Int("dirsCount", len(dirs)).Msg("All directories fetched successfully")
	return dirs, nil
}

func (r *DirRepository) UpdateLastScanned(dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Updating last scanned for directory")

	query := `
		UPDATE directories
		SET last_scanned = CURRENT_TIMESTAMP
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	_, err = r.Db.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to update last scanned for directory")
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Last scanned for directory updated successfully")
	return nil
}

func (r *DirRepository) Delete(dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting directory by ID")

	query := `
		DELETE FROM directories
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	_, err = r.Db.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to delete directory")
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Directory deleted successfully")
	return nil
}

func (r *DirRepository) IsExistsByPath(path string) (exists bool, err error) {
	log.Debug().Str("path", path).Msg("Checking if directory exists")

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
	row, err := r.Db.NamedQuery(query, args)
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
