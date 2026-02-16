package directory

import (
	"io"
	"os"
)

func FolderExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func RemoveFolderIfExists(folderPath string) error {
	if FolderExists(folderPath) {
		return os.RemoveAll(folderPath)
	}

	return nil
}

func IsDirEmpty(path string) (bool, error) {
	if !FolderExists(path) {
		return true, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == nil {
		return false, nil
	}
	if err == io.EOF {
		return true, nil
	}

	return false, err
}
