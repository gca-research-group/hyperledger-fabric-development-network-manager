package directory

import "os"

func FolderExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func RemoveFolderIfExists(folderPath string) error {
	if FolderExists(folderPath) {
		return os.RemoveAll(folderPath)
	}

	return nil
}
