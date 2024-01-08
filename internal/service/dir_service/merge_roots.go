package dir_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/directory"
	"strings"
)

func (s Service) MergeRoots(tx *sqlx.Tx, dirID1 int, dirID2 int) error {
	dir1, err := s.dirRepo.Read(tx, dirID1)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID1).Msg("Failed to get directory")
		return err
	}
	dir2, err := s.dirRepo.Read(tx, dirID2)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID1).Msg("Failed to get directory")
		return err
	}

	if dir1.ParentDirID != nil || dir2.ParentDirID != nil {
		err := fmt.Errorf("the directory is not root")
		log.Error().Err(err).Msg("The directory is not root")
		return err
	}

	var parentDir directory.Directory
	var childDir directory.Directory

	if strings.HasPrefix(dir1.Name, dir2.Name) {
		parentDir = dir2
		childDir = dir1
	} else if strings.HasPrefix(dir2.Name, dir1.Name) {
		parentDir = dir1
		childDir = dir2
	} else {
		err := fmt.Errorf("failed to find parent dir")
		log.Error().Err(err).Msg("Failed to find parent dir")
		return err
	}

	lostPath := strings.TrimPrefix(childDir.Name, parentDir.Name)
	if strings.HasPrefix(lostPath, "/") {
		lostPath = strings.TrimPrefix(lostPath, "/")
	}
	lostPathParts := strings.Split(lostPath, "/")
	prevDir := parentDir
	for i, lostPathPart := range lostPathParts {
		if i == len(lostPathParts)-1 {
			break
		}
		alreadyExists, err := s.dirRepo.IsExistsByParentAndName(tx, &prevDir.ID, lostPathPart)
		if err != nil {
			log.Error().Msg("Failed to check dir existence")
			return err
		}
		if alreadyExists {
			prevDir, err = s.dirRepo.ReadByParentAndName(tx, &prevDir.ID, lostPathPart)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get directory")
				return err
			}
		} else {
			dirToCreate := directory.Directory{
				ParentDirID: &prevDir.ID,
				Name:        lostPathPart,
			}
			createdDirId, err := s.dirRepo.Create(tx, dirToCreate)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create directory")
				return err
			}
			prevDir, err = s.dirRepo.Read(tx, createdDirId)
			if err != nil {
				log.Error().Err(err).Msg("Failed to read directory")
			}
		}
	}

	newChildDirectory := directory.Directory{
		ParentDirID: &prevDir.ID,
		Name:        lostPathParts[len(lostPathParts)-1],
	}
	err = s.dirRepo.Update(tx, childDir.ID, newChildDirectory)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update ")
		return err
	}

	return nil
}
