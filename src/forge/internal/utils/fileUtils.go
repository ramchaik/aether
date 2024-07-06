package utils

import (
	"os"
	"path/filepath"
)

func ReadDirRecursive(dir string) ([][]byte, error) {
	var files [][]byte
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			files = append(files, content)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
