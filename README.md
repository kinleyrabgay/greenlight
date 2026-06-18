# `git push greenlight`

A local git proxy that validates code before it reaches your real remote.

Push to `greenlight` instead of `origin`. It spins up a disposable worktree, runs an AI-driven pipeline (review → test → docs → lint), pushes upstream only after every check passes, and opens a clean PR for you.

- **Non-blocking** — runs in an isolated worktree; your working tree stays put.
- **Agent-agnostic** — `claude`, `codex`, `rovodev`, `opencode`, `pi`, or `acp:<target>`.
- **Framework-aware** — profiles supply per-stack commands + review rules (`angular`, `react`, `node`, `go`, `rails`, or your own).
- **Human in charge** — auto-fix the mechanical stuff, escalate judgment calls to you.

## How it works

```
        your branch
            │  git push greenlight
            ▼
   ┌────────────────────────────────────────────────┐
   │  disposable worktree — your work stays put       │
   │  review → test → docs → lint → push → PR → CI    │
   └────────────────────────────────────────────────┘
            │  every check green
            ▼
        clean PR, opened for you
```

Each step passes or stops with a **finding**. Safe, mechanical fixes are applied automatically; anything touching intent is escalated for you to **approve**, **fix**, or **skip**. Nothing reaches your real remote until every check is green.

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/kinleyrabgay/greenlight/main/docs/install.sh | sh
```

Or from source:

```sh
go build -o ~/bin/greenlight ./cmd/greenlight
```

## Step-by-step

### 0. Prerequisites

- **git** and a repo with an `origin` remote.
- **[gh](https://cli.github.com/)**, authenticated (`gh auth login`) — needed for the PR and CI steps.
- A **coding agent** on `PATH`: `claude`, `codex`, `opencode`, `acli` (Rovo Dev), or `pi`.

Check everything at once:

```sh
greenlight doctor
```

### 1. Install greenlight

```sh
curl -fsSL https://raw.githubusercontent.com/kinleyrabgay/greenlight/main/docs/install.sh | sh
# or from source:  go build -o ~/bin/greenlight ./cmd/greenlight
```

### 2. Initialize the gate (once per repo)

From inside the repo:

```sh
greenlight init                      # sets up the gate + daemon + /greenlight skill
greenlight init --framework angular  # ...and scaffold a framework profile (optional)
```

This creates a `greenlight` git remote, starts the background daemon, and installs the `/greenlight` agent skill. Re-running is safe.

### 3. Work on a feature branch

The gate validates committed history on a **non-default** branch:

```sh
git checkout -b my-feature
# ...make changes...
git add -A && git commit -m "my change"
```

### 4. Run the pipeline

Push to `greenlight` instead of `origin`:

```sh
git push greenlight
```

### 5. Act on findings

```sh
greenlight        # open the TUI for the active run
```

Each step that needs a decision shows **findings**. For each: **approve** (accept as-is), **fix** (let the pipeline fix it), or **skip**. Auto-fixable findings can be applied for you. Once every step is green, greenlight pushes to `origin` and opens the PR — no manual `git push origin`, no hand-written PR body.

### Optional: let your agent drive it

```
/greenlight                         # gate the work you already committed
/greenlight add a --json flag ...   # do the task, then gate it
```

`init` installs `/greenlight` for Claude Code, Codex, OpenCode, Rovo Dev, and Pi. Under the hood it drives `greenlight axi`, a non-interactive interface to the same flow.

## Three ways to trigger the gate

- **`git push greenlight`** — the explicit Git path; push a committed branch to the gate remote.
- **`greenlight`** — the TUI; run after making changes and a wizard branches, commits, and pushes for you (`greenlight -y` does it automatically).
- **`/greenlight`** — the agent skill; `/greenlight <task>` does a task and gates it, bare `/greenlight` gates existing committed work.

## Framework profiles

A profile bundles default `lint`/`test`/`format` commands, ignore patterns, and review rules for a stack. Select one per repo:

```yaml
# .greenlight.yaml
framework: angular
```

Precedence: `greenlight init --framework <name>` > the `framework:` field > auto-detection from repo files. Repo values always override profile defaults. List them with `greenlight profiles`; add your own at `~/.greenlight/profiles/<name>.md`.

## Docs

- [Configuration](docs/configuration.md) — `.greenlight.yaml`, framework profiles, auto-fix, evidence.
- [Pipeline](docs/pipeline.md) — the nine steps and what each does.
- [CLI](docs/cli.md) — every command and flag.

## Development

```sh
make build   # build bin/greenlight with version info
make test    # go test -race ./... (excludes e2e)
make e2e     # tagged end-to-end agent suite
make lint    # skill-drift check + go vet ./...
make skill   # regenerate the committed /greenlight skill files
make fmt     # gofmt -w .
```

`make e2e-record` overwrites `internal/e2e/fixtures/` from the real `claude`, `codex`, and `opencode` CLIs and spends real API quota — review before committing.

## License

MIT. Forked from [no-mistakes](https://github.com/kunchenguid/no-mistakes) (© 2026 Kun Chen).
