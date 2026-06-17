<h1 align="center"><code>git push greenlight</code></h1>
<p align="center">
  <a href="https://github.com/kinleyrabgay/greenlight/actions/workflows/release.yml"
    ><img
      alt="Release"
      src="https://img.shields.io/github/actions/workflow/status/kinleyrabgay/greenlight/release.yml?style=flat-square&label=release"
  /></a>
  <a href="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=flat-square"
    ><img
      alt="Platform"
      src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=flat-square"
  /></a>
  <a href="https://x.com/kinleyrabgay"
    ><img
      alt="X"
      src="https://img.shields.io/badge/X-@kinleyrabgay-black?style=flat-square"
  /></a>
  <a href="https://discord.gg/Wsy2NpnZDu"
    ><img
      alt="Discord"
      src="https://img.shields.io/discord/1439901831038763092?style=flat-square&label=discord"
  /></a>
</p>

<h3 align="center">Kill all the slop. Raise clean PR.</h3>

<p align="center"><strong>English</strong> · <a href="README.zh-CN.md">简体中文</a></p>

<p align="center">
  <img src="https://raw.githubusercontent.com/kinleyrabgay/greenlight/main/demo.gif" alt="greenlight demo" width="800" />
</p>

`greenlight` puts a local git proxy in front of your real remote. Push to `greenlight` instead of `origin`, and it spins up a disposable worktree, runs an AI-driven validation pipeline, forwards upstream only after every check passes, and opens a clean PR automatically.

- **Non-blocking** - the pipeline runs in an isolated worktree without disrupting your work.
- **Agent-agnostic** - `claude`, `codex`, `rovodev`, `opencode`, `pi`, or `acp:<target>` via `acpx`.
- **Agent-native** - `/greenlight` lets your coding agent do a task and gate it, or gate existing committed work: it runs the pipeline, has the pipeline apply safe fixes, and escalates the rest to you.
- **Human stays in charge** - auto-fix or review findings, your call.
- **Clean PRs by default** - push, open PR, watch CI, and auto-fix failures in one shot.

Full documentation: <https://kinleyrabgay.github.io/greenlight/>

## How it works

```
        your branch
            │  git push greenlight
            ▼
   ┌──────────────────────────────────────────────┐
   │  disposable worktree — your work stays put     │
   │  review → test → docs → lint → push → PR → CI  │
   └──────────────────────────────────────────────┘
            │  every check green
            ▼
        clean PR, opened for you
```

Each step either passes on its own or stops with a **finding** for you to act on. Safe, mechanical fixes are applied automatically; anything that touches your intent is escalated for you to **approve**, **fix**, or **skip**. Nothing reaches your real remote until every check is green.

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/kinleyrabgay/greenlight/main/docs/install.sh | sh
```

Windows, Go install, and build-from-source instructions are in the [installation guide](https://kinleyrabgay.github.io/greenlight/start-here/installation/).

## Quick Start

```sh
$ greenlight init
  ✓ Gate initialized

    repo  /Users/you/src/my-repo
    gate  greenlight → /Users/you/.greenlight/repos/abc123def456.git
  remote  git@github.com:you/my-repo.git
   skill  /greenlight installed for agents at user level

  Push through the gate with:
  git push greenlight <branch>

$ git checkout my-branch

# do some work in the branch...

$ git push greenlight
  * Pipeline started

  Run greenlight to review.

$ greenlight
# opens the TUI for the active run
```

From the TUI you act on each **finding**: **auto-fix** ones are applied for you (or approve to let them), **ask-user** ones are a judgement call you approve, fix, or skip. Once every check is green, the gate forwards your branch upstream and opens the PR for you — no manual `git push origin`, no hand-written PR body. Prefer to let your coding agent drive the same flow headlessly? Use `/greenlight` (see below).

## Three ways to trigger the gate

Every change runs through the same pipeline. Pick the entry point that fits how you're working when the change is ready:

- **`git push greenlight`** - the explicit Git path. Push a committed branch to the gate remote instead of `origin`.
- **`greenlight`** - the TUI. Run it after making changes (no commit needed) and a wizard walks you through creating a branch, committing, and pushing through the gate, then attaches to the run. `greenlight -y` does all of that automatically.
- **`/greenlight`** - the agent skill. Tell the coding agent to do a task and gate it with `/greenlight <task>`, or use bare `/greenlight` to gate existing committed work. It runs the pipeline, has the pipeline apply safe fixes, and stops to ask you about anything that needs a human call.

`greenlight init` installs the `/greenlight` skill for Claude Code and other agents. Under the hood the skill drives `greenlight axi`, a non-interactive TOON interface to the same approval flow.

See the [quick start](https://kinleyrabgay.github.io/greenlight/start-here/quick-start/) for the full first-run walkthrough.

## Development

```sh
make build   # Build bin/greenlight with version info
make test    # Run go test -race ./... (excludes the e2e suite)
make e2e     # Run the tagged end-to-end agent journey suite
make e2e-record # Re-record e2e fixtures when agent wire formats change
make lint    # Check generated skill drift and run go vet ./...
make skill   # Regenerate committed greenlight skill files
make fmt     # Run gofmt -w .
make demo    # Regenerate demo.gif and demo.mp4 (needs vhs and ffmpeg)
make docs    # Build the Astro docs site in docs/dist
```

See `Makefile` for the full target list.

`make e2e-record` overwrites `internal/e2e/fixtures/` from the real `claude`, `codex`, and `opencode` CLIs, spends real API quota, and should be reviewed before committing.

## Star History

<a href="https://www.star-history.com/?repos=kinleyrabgay%2Fgreenlight&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=kinleyrabgay/greenlight&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=kinleyrabgay/greenlight&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://api.star-history.com/chart?repos=kinleyrabgay/greenlight&type=date&legend=top-left" />
 </picture>
</a>
