package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindPhotoFiles_NonexistentRoot(t *testing.T) {
	_, err := FindPhotoFiles("/nonexistent/path/that/does/not/exist", []string{"jpg"})
	if err == nil {
		t.Fatal("expected error for nonexistent root")
	}
}

func TestFindPhotoFiles(t *testing.T) {
	dir := t.TempDir()
	// Create some sample files
	if err := os.WriteFile(filepath.Join(dir, "a.JPG"), []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.png"), []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "c.txt"), []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sub", "d.jpeg"), []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}

	ext := []string{"jpg", "jpeg", "png"}
	files, err := FindPhotoFiles(dir, ext)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 3 {
		t.Errorf("expected 3 photo files, got %d", len(files))
	}
}
