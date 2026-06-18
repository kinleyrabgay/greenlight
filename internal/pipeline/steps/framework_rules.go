package steps

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kinleyrabgay/greenlight/internal/pipeline"
)

// greenlightGuideFile is the per-repo audit guide greenlight reads during
// review/document. It lives at the worktree root and is owned by the repo.
const greenlightGuideFile = "GREENLIGHT.md"

// frameworkRulesPromptSection returns a prompt fragment carrying the project's
// audit rules. It prefers the repo's own GREENLIGHT.md (so each project
// controls what greenlight checks) and falls back to the embedded framework
// profile when the repo has no guide. Empty when neither is available, so steps
// can append it unconditionally.
//
// These rules are trusted, repo-owner-curated guidance, so they are embedded as
// authoritative instructions rather than untrusted data.
func frameworkRulesPromptSection(sctx *pipeline.StepContext) string {
	if sctx == nil {
		return ""
	}

	if sctx.WorkDir != "" {
		if data, err := os.ReadFile(filepath.Join(sctx.WorkDir, greenlightGuideFile)); err == nil {
			if rules := strings.TrimSpace(string(data)); rules != "" {
				return "\n\nProject audit guide (from GREENLIGHT.md) — apply these project rules when judging the change:\n" +
					rules + "\n"
			}
		}
	}

	if sctx.Config != nil {
		if rules := strings.TrimSpace(sctx.Config.FrameworkRules); rules != "" {
			header := "\n\nFramework conventions"
			if name := sctx.Config.Framework; name != "" {
				header += " (" + name + ")"
			}
			header += " — apply these project rules when judging the change:\n"
			return header + rules + "\n"
		}
	}

	return ""
}
