# Configuration

greenlight reads two optional config files:

- **`~/.greenlight/config.yaml`** — global defaults (agent, telemetry, auto-fix limits).
- **`.greenlight.yaml`** — per-repo, at the repository root.

Per-repo values override global ones. Neither file is required: with no config, greenlight auto-detects the agent and framework and lets each step's agent figure out commands.

## `.greenlight.yaml`

```yaml
framework: angular        # framework profile (see below); omit to auto-detect

agent: claude             # claude | codex | rovodev | opencode | pi | acp:<target>

commands:                 # explicit step commands; override the profile's
  lint: "yarn lint"
  test: "yarn nx run-many -t test"
  format: "yarn nx format:write"

ignore_patterns:          # paths excluded from review + documentation checks
  - "**/*.generated.ts"
  - "dist/**"

auto_fix:                 # per-step auto-fix attempt limits (0 = always ask)
  rebase: 3
  review: 3
  test: 3
  document: 3
  lint: 5
  ci: 3

intent:
  enabled: true           # infer author intent from local agent transcripts

test:
  evidence:
    store_in_repo: true   # commit captured test evidence into the repo
    dir: .greenlight/evidence
```

### Fields

| Field | Purpose |
|---|---|
| `framework` | Selects a framework profile (commands + ignore patterns + review rules). |
| `agent` | Overrides the global agent for this repo. `auto` picks the first agent found on `PATH`. |
| `commands.lint` / `.test` / `.format` | Exact shell commands run for those steps. Empty = the agent auto-detects. Set values override the framework profile. |
| `ignore_patterns` | Basename glob (`*.generated.ts`), subtree (`vendor/**`), or full-path glob. Overrides the profile when non-empty. |
| `auto_fix.<step>` | Max auto-fix attempts before the step pauses for approval. Unset inherits global. |
| `intent.enabled` | Toggle transcript-based intent inference (used when no `--intent` is supplied). |
| `test.evidence.store_in_repo` | Persist test evidence in the repo instead of only the run dir. |

## Framework profiles

A **profile** is a bundle of defaults for a stack: lint/test/format commands, ignore patterns, and review rules that get injected into the review and document steps. It lets a repo say `framework: angular` instead of hand-writing everything.

**Resolution precedence:**

1. `greenlight init --framework <name>` — scaffolds the field into a new `.greenlight.yaml`.
2. The `framework:` field in `.greenlight.yaml`.
3. Auto-detection from repo files (`angular.json`/`nx.json`, `package.json` deps, `go.mod`, `Gemfile`, …).

Repo `commands` / `ignore_patterns` always win over a profile — the profile only fills what you left empty.

**Built-in profiles:** `angular`, `react`, `node`, `go`, `rails`. List everything available (and what's detected here) with:

```sh
greenlight profiles
```

### Custom & override profiles

A profile is a markdown file with optional YAML frontmatter:

```markdown
---
commands:
  lint: "make lint"
  test: "make test"
  format: "make fmt"
ignore_patterns:
  - "**/*.pb.go"
---
# House rules (this body is injected into the review/document agent)
- Rule one
- Rule two
```

Save it at `~/.greenlight/profiles/<name>.md`. A user file overrides the built-in of the same name, so you can tune `angular` or add `svelte` without rebuilding.
