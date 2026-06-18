// Package profile provides framework profiles: named bundles of default
// commands, ignore patterns, and agent review rules. A profile lets a repo
// declare `framework: angular` (or rely on auto-detection) instead of
// hand-writing lint/test/format commands and house rules.
//
// Profiles are markdown files with optional YAML frontmatter:
//
//	---
//	commands:
//	  lint: "yarn lint"
//	  test: "yarn nx run-many -t test"
//	  format: "yarn nx format:write"
//	ignore_patterns:
//	  - "**/*.generated.ts"
//	---
//	# Framework rules (this body is injected into the review/document agent)
//	- Rule one
//	- Rule two
//
// Built-in profiles are embedded in the binary. A file of the same name in
// the user profiles dir (~/.greenlight/profiles/<name>.md) overrides the
// built-in, so users can tune or add frameworks without rebuilding.
package profile

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed builtin/*.md
var builtinFS embed.FS

// Profile is a resolved framework profile.
type Profile struct {
	Name           string
	Lint           string
	Test           string
	Format         string
	IgnorePatterns []string
	// Rules is the markdown body, injected verbatim into agent prompts.
	Rules string
}

type frontmatter struct {
	Commands struct {
		Lint   string `yaml:"lint"`
		Test   string `yaml:"test"`
		Format string `yaml:"format"`
	} `yaml:"commands"`
	IgnorePatterns []string `yaml:"ignore_patterns"`
}

// parse splits optional `---`-delimited YAML frontmatter from the markdown
// body and returns a Profile. A file without frontmatter is treated as
// all-rules (empty commands).
func parse(name string, data []byte) (*Profile, error) {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	p := &Profile{Name: name}

	if strings.HasPrefix(text, "---\n") {
		rest := text[len("---\n"):]
		end := strings.Index(rest, "\n---\n")
		if end == -1 && strings.HasSuffix(rest, "\n---") {
			end = len(rest) - len("\n---")
		}
		if end == -1 {
			return nil, fmt.Errorf("profile %q: unterminated frontmatter", name)
		}
		var fm frontmatter
		if err := yaml.Unmarshal([]byte(rest[:end]), &fm); err != nil {
			return nil, fmt.Errorf("profile %q frontmatter: %w", name, err)
		}
		p.Lint = fm.Commands.Lint
		p.Test = fm.Commands.Test
		p.Format = fm.Commands.Format
		p.IgnorePatterns = fm.IgnorePatterns
		body := rest[end:]
		body = strings.TrimPrefix(body, "\n---\n")
		body = strings.TrimPrefix(body, "\n---")
		p.Rules = strings.TrimSpace(body)
		return p, nil
	}

	p.Rules = strings.TrimSpace(text)
	return p, nil
}

func userPath(userDir, name string) string {
	return filepath.Join(userDir, name+".md")
}

// Load resolves a profile by name. A user override at
// userDir/<name>.md takes precedence over the embedded built-in.
// userDir may be empty to consult built-ins only.
func Load(name, userDir string) (*Profile, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return nil, fmt.Errorf("empty profile name")
	}

	if userDir != "" {
		if data, err := os.ReadFile(userPath(userDir, name)); err == nil {
			return parse(name, data)
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read user profile %q: %w", name, err)
		}
	}

	data, err := builtinFS.ReadFile("builtin/" + name + ".md")
	if err != nil {
		return nil, fmt.Errorf("unknown framework profile %q (available: %s)", name, strings.Join(List(userDir), ", "))
	}
	return parse(name, data)
}

// List returns the sorted union of built-in and user profile names.
func List(userDir string) []string {
	seen := map[string]bool{}
	if entries, err := builtinFS.ReadDir("builtin"); err == nil {
		for _, e := range entries {
			seen[strings.TrimSuffix(e.Name(), ".md")] = true
		}
	}
	if userDir != "" {
		if entries, err := os.ReadDir(userDir); err == nil {
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".md") {
					seen[strings.TrimSuffix(e.Name(), ".md")] = true
				}
			}
		}
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Detect inspects repoDir and returns a best-guess framework name, or "" if
// none is recognized. Checks are ordered most-specific first.
func Detect(repoDir string) string {
	exists := func(rel string) bool {
		_, err := os.Stat(filepath.Join(repoDir, rel))
		return err == nil
	}
	readPkg := func() string {
		data, err := os.ReadFile(filepath.Join(repoDir, "package.json"))
		if err != nil {
			return ""
		}
		return string(data)
	}

	pkg := readPkg()
	hasFile := func(globs ...string) bool {
		for _, g := range globs {
			if exists(g) {
				return true
			}
		}
		return false
	}
	switch {
	case exists("angular.json") || exists("nx.json") || strings.Contains(pkg, "@angular/core"):
		return "angular"
	case exists("pubspec.yaml"):
		return "flutter"
	case exists("Gemfile") && (exists("config/application.rb") || gemfileHasRails(repoDir)):
		return "rails"
	case hasFile("next.config.js", "next.config.mjs", "next.config.ts") || strings.Contains(pkg, "\"next\""):
		return "nextjs"
	case hasFile("svelte.config.js", "svelte.config.ts") || strings.Contains(pkg, "\"svelte\""):
		return "svelte"
	case strings.Contains(pkg, "\"react\""):
		return "react"
	case exists("go.mod"):
		return "go"
	case hasFile("pyproject.toml", "setup.py", "setup.cfg", "requirements.txt", "Pipfile"):
		return "python"
	case pkg != "":
		return "node"
	}
	return ""
}

func gemfileHasRails(repoDir string) bool {
	data, err := os.ReadFile(filepath.Join(repoDir, "Gemfile"))
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "rails")
}
