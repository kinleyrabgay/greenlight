package steps

import (
	"strings"

	"github.com/kinleyrabgay/greenlight/internal/pipeline"
)

// frameworkRulesPromptSection returns a prompt fragment carrying the active
// framework profile's review rules (e.g. Angular/React/Go conventions). The
// fragment is empty when no framework profile is active, so steps can append
// it unconditionally.
//
// Unlike user intent, these rules come from a trusted, repo-owner-curated
// profile file, so they are embedded as authoritative guidance rather than
// untrusted data.
func frameworkRulesPromptSection(sctx *pipeline.StepContext) string {
	if sctx == nil || sctx.Config == nil {
		return ""
	}
	rules := strings.TrimSpace(sctx.Config.FrameworkRules)
	if rules == "" {
		return ""
	}
	name := sctx.Config.Framework
	header := "\n\nFramework conventions"
	if name != "" {
		header += " (" + name + ")"
	}
	header += " — apply these project rules when judging the change:\n"
	return header + rules + "\n"
}
