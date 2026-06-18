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

## GREENLIGHT.md (per-repo audit guide)

`greenlight init` seeds a **`GREENLIGHT.md`** at the repo root: stack summary, architecture, and do/don't rules. greenlight reads it during the review and document steps, so **each repo controls what gets audited**.

- Seeded from the resolved framework profile, or a generic template when none matches.
- **Never overwritten** on re-init — the repo owns it.
- When a repo has no `GREENLIGHT.md`, greenlight falls back to the built-in profile for its detected framework.

Edit it freely; commit it so the gate (which validates committed history) picks it up. Examples live in [`example/`](../example/).

## Framework profiles

A **profile** is a bundle of defaults for a stack: lint/test/format commands, ignore patterns, and a review guide. It picks the `GREENLIGHT.md` seed and supplies commands when `.greenlight.yaml` doesn't.

**Resolution precedence:**

1. `greenlight init --framework <name>` — scaffolds the field into a new `.greenlight.yaml`.
2. The `framework:` field in `.greenlight.yaml`.
3. Auto-detection from repo files.

Repo `commands` / `ignore_patterns` always win over a profile — the profile only fills what you left empty.

**Built-in profiles:** `angular`, `react`, `nextjs`, `node`, `go`, `rails`, `flutter`, `svelte`, `python`. List everything available (and what's detected here) with:

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
# Guide body — seeds GREENLIGHT.md and feeds the review/document agent
## Stack
## Architecture
## Do
## Don't
```

Save it at `~/.greenlight/profiles/<name>.md`. A user file overrides the built-in of the same name, so you can tune `angular` or add a new stack without rebuilding.
