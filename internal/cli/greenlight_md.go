package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kinleyrabgay/greenlight/internal/profile"
)

// GreenlightGuideFile is the per-repo audit guide greenlight reads during the
// review and document steps. It is owned by the repo: each project edits it to
// control what greenlight checks.
const GreenlightGuideFile = "GREENLIGHT.md"

// greenlightGuideContent builds the seed GREENLIGHT.md for a project. When a
// framework profile is resolved, its guide body is used as the seed; otherwise
// a generic template with empty sections is written.
func greenlightGuideContent(framework string, p *profile.Profile) string {
	var b strings.Builder
	b.WriteString("# GREENLIGHT.md\n\n")
	b.WriteString("Project guide for [greenlight](https://github.com/kinleyrabgay/greenlight). ")
	b.WriteString("greenlight reads this file during its review/audit and applies the rules below to your changes. ")
	b.WriteString("Edit it freely — it is yours to control.\n\n")

	if p != nil && strings.TrimSpace(p.Rules) != "" {
		fmt.Fprintf(&b, "<!-- seeded from the %q framework profile; tailor it to this repo -->\n\n", framework)
		b.WriteString(strings.TrimSpace(p.Rules))
		b.WriteString("\n")
		return b.String()
	}

	b.WriteString(strings.TrimSpace(`
## Stack

<!-- Languages, frameworks, and key libraries this repo uses. -->

## Architecture

<!-- How the codebase is organized and the conventions to follow. -->

## Do

<!-- Practices greenlight should enforce in review. -->
-

## Don't

<!-- Anti-patterns greenlight should flag. -->
-
`))
	b.WriteString("\n")
	return b.String()
}

// writeGreenlightGuide creates GREENLIGHT.md in repoDir when absent and returns
// whether it wrote the file. An existing GREENLIGHT.md is never overwritten —
// the repo owns its content.
func writeGreenlightGuide(repoDir, framework string, p *profile.Profile) (wrote bool, err error) {
	path := filepath.Join(repoDir, GreenlightGuideFile)
	if _, statErr := os.Stat(path); statErr == nil {
		return false, nil
	} else if !os.IsNotExist(statErr) {
		return false, statErr
	}
	if err := os.WriteFile(path, []byte(greenlightGuideContent(framework, p)), 0o644); err != nil {
		return false, err
	}
	return true, nil
}
