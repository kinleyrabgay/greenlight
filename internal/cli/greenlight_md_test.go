package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kinleyrabgay/greenlight/internal/profile"
)

func TestWriteGreenlightGuide_CreatesFromProfile(t *testing.T) {
	dir := t.TempDir()
	prof, err := profile.Load("go", "")
	if err != nil {
		t.Fatal(err)
	}
	wrote, err := writeGreenlightGuide(dir, "go", prof)
	if err != nil || !wrote {
		t.Fatalf("expected create, got wrote=%v err=%v", wrote, err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "GREENLIGHT.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, "# GREENLIGHT.md") || !strings.Contains(text, "Go project guide") {
		t.Fatalf("unexpected guide content:\n%s", text)
	}
}

func TestWriteGreenlightGuide_GenericWhenNoProfile(t *testing.T) {
	dir := t.TempDir()
	wrote, err := writeGreenlightGuide(dir, "", nil)
	if err != nil || !wrote {
		t.Fatalf("expected create, got wrote=%v err=%v", wrote, err)
	}
	data, _ := os.ReadFile(filepath.Join(dir, "GREENLIGHT.md"))
	for _, want := range []string{"## Stack", "## Architecture", "## Do", "## Don't"} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("generic template missing %q", want)
		}
	}
}

func TestWriteGreenlightGuide_NeverOverwrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "GREENLIGHT.md")
	if err := os.WriteFile(path, []byte("MY CONTENT"), 0o644); err != nil {
		t.Fatal(err)
	}
	wrote, err := writeGreenlightGuide(dir, "go", nil)
	if err != nil {
		t.Fatal(err)
	}
	if wrote {
		t.Fatal("should not overwrite an existing GREENLIGHT.md")
	}
	data, _ := os.ReadFile(path)
	if string(data) != "MY CONTENT" {
		t.Fatalf("existing content was modified: %q", string(data))
	}
}
