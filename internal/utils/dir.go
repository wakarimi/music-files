package utils

import "os"

func DirectoryExists(dirPath string) (bool, error) {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
