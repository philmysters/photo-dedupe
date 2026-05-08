package internal

import (
	"os"
	"testing"
)

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
