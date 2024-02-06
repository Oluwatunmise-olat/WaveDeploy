package files

import (
	"os"
	"path/filepath"
)

func WriteToFileWithOverride(content string, path string) error {
	directory := GetDirectoryFromPath(path)

	if !PathExists(directory) {
		if err := os.MkdirAll(directory, 0777); err != nil {
			return err
		}
	}

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return err
	}
	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetDirectoryFromPath(dir string) string {
	return filepath.Dir(dir)
}

func GetFileContent(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
