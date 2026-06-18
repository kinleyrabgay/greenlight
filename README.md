<div align="center">

# greenlight

**A local git proxy that validates your code before it reaches the remote.**

Push to `greenlight` instead of `origin`. It runs an AI-driven pipeline — review, tests, docs, lint — in a disposable worktree, then pushes upstream and opens a clean PR only after every check is green.

[![CI](https://github.com/kinleyrabgay/greenlight/actions/workflows/ci.yml/badge.svg)](https://github.com/kinleyrabgay/greenlight/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)

</div>

---

## Contents

- [Why](#why)
- [Install](#install)
- [Quick start](#quick-start)
- [How it works](#how-it-works)
- [GREENLIGHT.md — per-repo audit guide](#greenlightmd--per-repo-audit-guide)
- [Framework profiles](#framework-profiles)
- [Configuration](#configuration)
- [Commands](#commands)
- [Troubleshooting](#troubleshooting)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## Why

- **Non-blocking** — the pipeline runs in an isolated worktree; your working tree stays put.
- **Agent-agnostic** — `claude`, `codex`, `rovodev`, `opencode`, `pi`, or `acp:<target>`.
- **Framework-aware** — ships review guides for Angular, React, Next.js, Node, Go, Rails, Flutter, Svelte, Python.
- **You own the rules** — each repo gets a `GREENLIGHT.md` it fully controls.
- **Human in charge** — auto-fix the mechanical stuff, escalate judgment calls to you.

## Install

Homebrew (builds from source, needs the `go` toolchain):

```sh
brew install kinleyrabgay/greenlight/greenlight
```

<details>
<summary>Other methods</summary>

```sh
# install script
curl -fsSL https://raw.githubusercontent.com/kinleyrabgay/greenlight/main/docs/install.sh | sh

# from source
go build -o ~/bin/greenlight ./cmd/greenlight
```
</details>

Check your environment any time with `greenlight doctor`.

## Quick start

```sh
# 0. prerequisites: git repo with an `origin` remote, `gh` authenticated,
#    and a coding agent on PATH (claude / codex / opencode / acli / pi).

# 1. set up the gate (once per repo). seeds a GREENLIGHT.md for your stack.
greenlight init                      # auto-detects the framework
greenlight init --framework angular  # ...or pick one explicitly

# 2. work on a feature branch (the gate validates committed, non-default branches)
git checkout -b my-feature
git add -A && git commit -m "my change"

# 3. push to the gate instead of origin
git push greenlight

# 4. act on findings
greenlight   # opens the TUI: approve / fix / skip each finding
```

Once every step is green, greenlight pushes to `origin` and opens the PR for you — no manual `git push origin`, no hand-written PR body.

**Prefer your agent to drive it?** `init` installs a `/greenlight` skill:

```
/greenlight                       # gate the work you already committed
/greenlight add a --json flag …   # do the task, then gate it
```

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

Each step passes or stops with a **finding** (severity + an action: `auto-fix`, `ask-user`, `no-op`). Safe, mechanical fixes are applied automatically; anything touching intent is escalated for you to **approve**, **fix**, or **skip**. Nothing reaches your real remote until every check is green. The pipeline order is fixed and opinionated so "passed the gate" means the same thing across repos — see [docs/pipeline.md](docs/pipeline.md).

## GREENLIGHT.md — per-repo audit guide

`greenlight init` drops a **`GREENLIGHT.md`** at your repo root describing the stack, architecture, and do/don't rules for your project. greenlight reads it during review, so **each repo controls exactly what gets audited** — edit it freely.

- Seeded from the resolved [framework profile](#framework-profiles), or a generic template if none matches.
- Never overwritten on re-init — it's yours.
- If a repo has no `GREENLIGHT.md`, greenlight falls back to the built-in profile for its detected framework.

See ready-made examples in [`example/`](example/) — [Angular](example/angular/GREENLIGHT.md), [React](example/react/GREENLIGHT.md).

## Framework profiles

A profile bundles default `lint`/`test`/`format` commands, ignore patterns, and a review guide. It picks the `GREENLIGHT.md` seed and supplies commands when your `.greenlight.yaml` doesn't.

| Profile | Detected by |
|---|---|
| `angular` | `angular.json`, `nx.json`, or `@angular/core` |
| `nextjs` | `next.config.*` or `next` dependency |
| `react` | `react` dependency |
| `svelte` | `svelte.config.*` or `svelte` dependency |
| `flutter` | `pubspec.yaml` |
| `rails` | `Gemfile` with Rails |
| `go` | `go.mod` |
| `python` | `pyproject.toml` / `requirements.txt` / … |
| `node` | any other `package.json` |

Selection precedence: `greenlight init --framework <name>` > the `framework:` field in `.greenlight.yaml` > auto-detection. List them and see what's detected here with `greenlight profiles`. Add your own at `~/.greenlight/profiles/<name>.md`.

## Configuration

Per-repo `.greenlight.yaml` (all fields optional):

```yaml
framework: angular            # profile to use; omit to auto-detect
agent: claude                 # claude | codex | rovodev | opencode | pi | acp:<target>
commands:                     # override the profile's commands
  lint: "yarn lint"
  test: "yarn nx run-many -t test"
  format: "yarn nx format:write"
ignore_patterns:
  - "**/*.generated.ts"
auto_fix:                     # max auto-fix attempts per step (0 = always ask)
  review: 3
  test: 3
  lint: 5
```

Full reference: [docs/configuration.md](docs/configuration.md).

## Commands

| Command | What it does |
|---|---|
| `greenlight init [--framework X]` | Set up the gate, daemon, `/greenlight` skill, and seed `GREENLIGHT.md`. |
| `greenlight` | Attach to the active run (TUI). |
| `greenlight profiles` | List framework profiles and the one detected here. |
| `greenlight doctor` | Check git, gh, agents, daemon, and data dir. |
| `greenlight status` / `runs` / `stats` | Inspect repo, runs, and history. |
| `greenlight rerun` / `eject` / `update` | Re-run the pipeline / remove the gate / self-update. |

Full reference: [docs/cli.md](docs/cli.md).

## Troubleshooting

- **`greenlight doctor` shows a missing agent** — install one of `claude`, `codex`, `opencode`, `acli`, `pi` and authenticate it.
- **PR / CI steps fail** — install and authenticate the [GitHub CLI](https://cli.github.com/): `gh auth login`.
- **"must be on a non-default branch"** — the gate validates committed history on a feature branch; `git checkout -b`.
- **Heavy/slow runs** — trim the `test` command in `.greenlight.yaml` (some profiles run a full build).
- **Skip a step for one push** — `git push -o greenlight.skip=test,lint greenlight <branch>` or `greenlight --skip test,lint`.

## Documentation

- [Configuration](docs/configuration.md) — `.greenlight.yaml`, framework profiles, `GREENLIGHT.md`.
- [Pipeline](docs/pipeline.md) — the nine steps and what each does.
- [CLI](docs/cli.md) — every command and flag.

## Contributing

This repo *is* greenlight — contributions go through the tool. See [CONTRIBUTING.md](CONTRIBUTING.md) and [AGENTS.md](AGENTS.md).

## License

MIT. Forked from [no-mistakes](https://github.com/kunchenguid/no-mistakes) (© 2026 Kun Chen).
