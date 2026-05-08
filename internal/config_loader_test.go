package internal

import (
	"os"
	"strings"
	"testing"
)

func TestLoadConfig_OpenError(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent config file")
	}
}

func TestLoadConfig_DecodeError(t *testing.T) {
	f, err := os.CreateTemp("", "bad_yaml_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(f.Name()) }()
	if _, err := f.WriteString("key: [\nunclosed"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(f.Name())
	if err == nil || !strings.Contains(err.Error(), "failed to parse config") {
		t.Fatalf("expected parse error, got: %v", err)
	}
}

func TestLoadConfig_Success(t *testing.T) {
	// Create a temp YAML config for the test
	tmpFile, err := os.CreateTemp("", "photo_dedupe_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	sample := "supported_extensions:\n  - jpg\n  - jpeg\n  - png\n"
	if _, err := tmpFile.Write([]byte(sample)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	expected := []string{"jpg", "jpeg", "png"}
	if len(cfg.SupportedExtensions) != len(expected) {
		t.Fatalf("Expected %d extensions, got %d", len(expected), len(cfg.SupportedExtensions))
	}
	for i, ext := range cfg.SupportedExtensions {
		if ext != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], ext)
		}
	}
}
