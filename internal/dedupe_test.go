package internal

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTestFile creates a file with the given content and returns its path.
func writeTestFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, content, 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestDeduplicate_WithDuplicates(t *testing.T) {
	dir := t.TempDir()
	same := []byte("shared content")

	f1a := writeTestFile(t, dir, "1a.jpg", same)
	f1b := writeTestFile(t, dir, "1b.jpg", []byte("only in set1"))
	f2a := writeTestFile(t, dir, "2a.jpg", same) // duplicate of f1a
	f2b := writeTestFile(t, dir, "2b.jpg", []byte("only in set2"))

	result, err := Deduplicate(
		[]PhotoFile{{Path: f1a}, {Path: f1b}},
		[]PhotoFile{{Path: f2a}, {Path: f2b}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Duplicates) != 1 {
		t.Errorf("expected 1 duplicate pair, got %d", len(result.Duplicates))
	}
	if len(result.UniqueToFirst) != 1 {
		t.Errorf("expected 1 unique-to-first, got %d", len(result.UniqueToFirst))
	}
	if len(result.UniqueToSecond) != 1 {
		t.Errorf("expected 1 unique-to-second, got %d", len(result.UniqueToSecond))
	}
}

func TestDeduplicate_AllUnique(t *testing.T) {
	dir := t.TempDir()
	f1 := writeTestFile(t, dir, "a.jpg", []byte("aaa"))
	f2 := writeTestFile(t, dir, "b.jpg", []byte("bbb"))

	result, err := Deduplicate([]PhotoFile{{Path: f1}}, []PhotoFile{{Path: f2}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Duplicates) != 0 {
		t.Errorf("expected 0 duplicates, got %d", len(result.Duplicates))
	}
	if len(result.UniqueToFirst) != 1 || len(result.UniqueToSecond) != 1 {
		t.Errorf("expected 1 unique in each set")
	}
}

func TestDeduplicate_EmptySets(t *testing.T) {
	result, err := Deduplicate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Duplicates) != 0 || len(result.UniqueToFirst) != 0 || len(result.UniqueToSecond) != 0 {
		t.Error("expected all empty slices for empty input")
	}
}

func TestDeduplicate_HashErrorFiles1(t *testing.T) {
	_, err := Deduplicate(
		[]PhotoFile{{Path: "/nonexistent/file.jpg"}},
		nil,
	)
	if err == nil {
		t.Fatal("expected error when hashing files1 fails")
	}
}

func TestDeduplicate_HashErrorFiles2(t *testing.T) {
	dir := t.TempDir()
	f1 := writeTestFile(t, dir, "a.jpg", []byte("content"))

	_, err := Deduplicate(
		[]PhotoFile{{Path: f1}},
		[]PhotoFile{{Path: "/nonexistent/file.jpg"}},
	)
	if err == nil {
		t.Fatal("expected error when hashing files2 fails")
	}
}

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
