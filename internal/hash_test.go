package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"testing"
)

type errReader struct{ err error }

func (e errReader) Read([]byte) (int, error) { return 0, e.err }

func TestHashReader_ReadError(t *testing.T) {
	_, err := hashReader(errReader{errors.New("simulated read error")})
	if err == nil {
		t.Fatal("expected error from hashReader")
	}
}

func TestHashFile_Success(t *testing.T) {
	f, err := os.CreateTemp("", "hash_test_*.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(f.Name()) }()

	content := []byte("photo content")
	if _, err := f.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	got, err := HashFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum := sha256.Sum256(content)
	want := hex.EncodeToString(sum[:])
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestHashFile_SameContentSameHash(t *testing.T) {
	write := func(content []byte) string {
		f, err := os.CreateTemp("", "hash_test_*.bin")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.Remove(f.Name()) }()
		if _, err := f.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
		h, err := HashFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		return h
	}

	content := []byte("duplicate photo bytes")
	h1, h2 := write(content), write(content)
	if h1 != h2 {
		t.Error("identical content should produce identical hashes")
	}
}

func TestHashFile_DifferentContentDifferentHash(t *testing.T) {
	hash := func(content []byte) string {
		f, err := os.CreateTemp("", "hash_test_*.bin")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.Remove(f.Name()) }()
		if _, err := f.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
		h, err := HashFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		return h
	}

	if hash([]byte("photo A")) == hash([]byte("photo B")) {
		t.Error("different content should produce different hashes")
	}
}

func TestHashFile_OpenError(t *testing.T) {
	_, err := HashFile("/nonexistent/path/file.jpg")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}
