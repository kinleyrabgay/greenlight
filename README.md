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

## Quick start

```sh
greenlight init                 # set up the gate in this repo
greenlight init --framework go  # ...or scaffold a framework profile too

git checkout -b my-branch
# ...do some work...

git push greenlight             # run the pipeline
greenlight                      # open the TUI to act on findings
```

`init` installs the `/greenlight` agent skill (Claude Code, Codex, OpenCode, Rovo Dev, Pi). Under the hood it drives `greenlight axi`, a non-interactive interface to the same flow.

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
