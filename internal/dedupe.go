package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PhotoFile represents a discovered photo file.
type PhotoFile struct {
	Path string
}

// FindPhotoFiles walks a directory and returns files matching supported extensions (case-insensitive).
func FindPhotoFiles(root string, exts []string) ([]PhotoFile, error) {
	// Normalize and put extensions in a set (map).
	extSet := make(map[string]struct{})
	for _, ext := range exts {
		extSet[strings.ToLower(ext)] = struct{}{}
	}

	var files []PhotoFile
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(info.Name()), "."))
		if _, ok := extSet[ext]; ok {
			files = append(files, PhotoFile{Path: path})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking %s: %w", root, err)
	}
	return files, nil
}
