package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func captureExit(t *testing.T) *int {
	t.Helper()
	code := -1
	orig := osExit
	osExit = func(c int) { code = c }
	t.Cleanup(func() { osExit = orig })
	return &code
}

func makeTestDirs(t *testing.T) (in1, in2, out string) {
	t.Helper()
	base := t.TempDir()
	in1 = filepath.Join(base, "in1")
	in2 = filepath.Join(base, "in2")
	out = filepath.Join(base, "out")
	if err := os.MkdirAll(in1, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(in2, 0755); err != nil {
		t.Fatal(err)
	}
	return
}

func makeConfig(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "cfg_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("supported_extensions:\n  - jpg\n"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(f.Name()) })
	return f.Name()
}

func TestRun_TooFewArgs(t *testing.T) {
	err := run([]string{"-in1", "a"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for too few args")
	}
}

func TestRun_MissingRequiredArgs(t *testing.T) {
	err := run([]string{"-in1", "a", "-in2", "b"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for missing -out")
	}
}

func TestRun_BadConfig(t *testing.T) {
	in1, in2, out := makeTestDirs(t)
	err := run([]string{"-in1", in1, "-in2", in2, "-out", out, "-config", "nonexistent.yaml"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestRun_Success(t *testing.T) {
	in1, in2, out := makeTestDirs(t)
	if err := os.WriteFile(filepath.Join(in1, "a.jpg"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := makeConfig(t)

	var buf bytes.Buffer
	err := run([]string{"-in1", in1, "-in2", in2, "-out", out, "-config", cfg, "--dryrun"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("output dir not created: %v", err)
	}
}

func TestRun_BadInputDir(t *testing.T) {
	_, in2, out := makeTestDirs(t)
	cfg := makeConfig(t)

	err := run([]string{"-in1", "/nonexistent/dir", "-in2", in2, "-out", out, "-config", cfg}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for bad input dir")
	}
}

func TestRun_BadInput2Dir(t *testing.T) {
	in1, _, out := makeTestDirs(t)
	cfg := makeConfig(t)

	err := run([]string{"-in1", in1, "-in2", "/nonexistent/dir", "-out", out, "-config", cfg}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for bad input2 dir")
	}
}

func TestRun_MkdirAllFails(t *testing.T) {
	in1, in2, _ := makeTestDirs(t)
	cfg := makeConfig(t)

	// Use a file path as the output parent so MkdirAll cannot create it.
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	err := run([]string{"-in1", in1, "-in2", in2, "-out", filepath.Join(blocker, "sub"), "-config", cfg}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error when output dir cannot be created")
	}
}

func TestRun_BadYAML(t *testing.T) {
	in1, in2, out := makeTestDirs(t)

	f, err := os.CreateTemp("", "bad_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(": :\n"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(f.Name()) })

	err = run([]string{"-in1", in1, "-in2", in2, "-out", out, "-config", f.Name()}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestMain_Error(t *testing.T) {
	exitCode := captureExit(t)
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"photo-dedupe"} // too few args → run returns error

	main()

	if *exitCode != 1 {
		t.Fatalf("expected exit 1, got %d", *exitCode)
	}
}

func TestMain_Success(t *testing.T) {
	exitCode := captureExit(t)
	in1, in2, out := makeTestDirs(t)
	cfg := makeConfig(t)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"photo-dedupe", "-in1", in1, "-in2", in2, "-out", out, "-config", cfg}

	main()

	if *exitCode != -1 {
		t.Fatalf("expected no exit, got exit code %d", *exitCode)
	}
}
