package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadBuiltins(t *testing.T) {
	for _, name := range []string{"angular", "react", "node", "go", "rails"} {
		p, err := Load(name, "")
		if err != nil {
			t.Fatalf("Load(%q) error: %v", name, err)
		}
		if p.Name != name {
			t.Errorf("Load(%q).Name = %q", name, p.Name)
		}
		if p.Lint == "" || p.Test == "" {
			t.Errorf("Load(%q): expected lint+test commands, got lint=%q test=%q", name, p.Lint, p.Test)
		}
		if p.Rules == "" {
			t.Errorf("Load(%q): expected rules body", name)
		}
	}
}

func TestLoadUnknown(t *testing.T) {
	if _, err := Load("nope", ""); err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

func TestUserOverride(t *testing.T) {
	dir := t.TempDir()
	custom := "---\ncommands:\n  lint: \"custom-lint\"\n  test: \"custom-test\"\n---\n# Custom rules\n- be nice\n"
	if err := os.WriteFile(filepath.Join(dir, "angular.md"), []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := Load("angular", dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.Lint != "custom-lint" {
		t.Errorf("user override not applied: lint=%q", p.Lint)
	}
	if p.Rules != "# Custom rules\n- be nice" {
		t.Errorf("unexpected rules body: %q", p.Rules)
	}
}

func TestDetect(t *testing.T) {
	cases := map[string]func(dir string){
		"angular": func(d string) { os.WriteFile(filepath.Join(d, "angular.json"), []byte("{}"), 0o644) },
		"go":      func(d string) { os.WriteFile(filepath.Join(d, "go.mod"), []byte("module x"), 0o644) },
		"react": func(d string) {
			os.WriteFile(filepath.Join(d, "package.json"), []byte(`{"dependencies":{"react":"18"}}`), 0o644)
		},
		"node": func(d string) {
			os.WriteFile(filepath.Join(d, "package.json"), []byte(`{"dependencies":{"express":"4"}}`), 0o644)
		},
		"rails": func(d string) { os.WriteFile(filepath.Join(d, "Gemfile"), []byte("gem 'rails'"), 0o644) },
	}
	for want, setup := range cases {
		dir := t.TempDir()
		setup(dir)
		if got := Detect(dir); got != want {
			t.Errorf("Detect() = %q, want %q", got, want)
		}
	}
	if got := Detect(t.TempDir()); got != "" {
		t.Errorf("Detect(empty) = %q, want \"\"", got)
	}
}
