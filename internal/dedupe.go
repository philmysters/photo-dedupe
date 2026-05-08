package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DuplicatePair holds a matched pair of files with identical content.
type DuplicatePair struct {
	First  PhotoFile
	Second PhotoFile
}

// DedupeResult categorises files from two input sets by their relationship.
type DedupeResult struct {
	UniqueToFirst  []PhotoFile     // present only in the first set
	UniqueToSecond []PhotoFile     // present only in the second set
	Duplicates     []DuplicatePair // identical content in both sets
}

// Deduplicate compares two slices of PhotoFile by SHA-256 hash and returns
// which files are unique to each set and which are duplicates across both.
func Deduplicate(files1, files2 []PhotoFile) (DedupeResult, error) {
	hashToFirst := make(map[string]PhotoFile, len(files1))
	for _, pf := range files1 {
		h, err := HashFile(pf.Path)
		if err != nil {
			return DedupeResult{}, fmt.Errorf("hashing %s: %w", pf.Path, err)
		}
		hashToFirst[h] = pf
	}

	var result DedupeResult
	matched := make(map[string]struct{})

	for _, pf := range files2 {
		h, err := HashFile(pf.Path)
		if err != nil {
			return DedupeResult{}, fmt.Errorf("hashing %s: %w", pf.Path, err)
		}
		if f1, ok := hashToFirst[h]; ok {
			result.Duplicates = append(result.Duplicates, DuplicatePair{First: f1, Second: pf})
			matched[h] = struct{}{}
		} else {
			result.UniqueToSecond = append(result.UniqueToSecond, pf)
		}
	}

	for h, pf := range hashToFirst {
		if _, ok := matched[h]; !ok {
			result.UniqueToFirst = append(result.UniqueToFirst, pf)
		}
	}

	return result, nil
}

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
