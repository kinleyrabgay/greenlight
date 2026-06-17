---
title: Installation
description: All install options, prerequisites, update, and uninstall.
---

## macOS / Linux

```sh
curl -fsSL https://raw.githubusercontent.com/kinleyrabgay/greenlight/main/docs/install.sh | sh
```

The installer keeps the real binary in `~/.greenlight/bin` and exposes `greenlight` through a symlink in `~/.local/bin` or `/usr/local/bin`. That keeps future `greenlight update` runs in a user-owned location instead of rewriting a system binary in place.

It also installs or refreshes the background daemon for you by running `greenlight daemon restart`, preferring a managed service (launchd on macOS, systemd user service on Linux) and falling back to a detached daemon if that path is unavailable. If the restart fails, the install command fails.

Official release binaries installed this way include the default self-hosted telemetry host and website ID. Disable telemetry with `GREENLIGHT_TELEMETRY=0`, or override the host and website ID with `GREENLIGHT_UMAMI_HOST` and `GREENLIGHT_UMAMI_WEBSITE_ID`.

## Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/kinleyrabgay/greenlight/main/docs/install.ps1 | iex
```

Installs the binary and restarts the background daemon automatically with `greenlight.exe daemon restart`, preferring a managed Task Scheduler task and falling back to a detached daemon if needed. If the restart fails, the install command fails.

Official release binaries installed this way include the default self-hosted telemetry host and website ID. Disable telemetry with `GREENLIGHT_TELEMETRY=0`, or override the host and website ID with `GREENLIGHT_UMAMI_HOST` and `GREENLIGHT_UMAMI_WEBSITE_ID`.

## Go install

```sh
go install github.com/kinleyrabgay/greenlight/cmd/greenlight@latest
```

`go install` builds the CLI without an embedded telemetry website ID, so telemetry stays off by default unless you later set `GREENLIGHT_UMAMI_WEBSITE_ID` at runtime.

## From source

```sh
git clone git@github.com:kinleyrabgay/greenlight.git
cd greenlight
make build
make install
```

`make build` embeds the telemetry host from `GREENLIGHT_UMAMI_HOST` in a repo-local `.env` first, then `UMAMI_HOST` from the shell, then the default self-hosted host. It embeds the telemetry website ID from `GREENLIGHT_UMAMI_WEBSITE_ID` in `.env` first, then `UMAMI_WEBSITE_ID` from the shell, then the default website ID.

## Prerequisites

- **git** - required
- **One supported agent binary** - `claude`, `codex`, `acli` (Rovo Dev), `opencode`, or `pi`, or a separately installed `acpx` binary for `agent: acp:<target>`
- **Optional, for PRs and CI:**
  - `gh` CLI (GitHub)
  - `glab` CLI (GitLab)
  - `GREENLIGHT_BITBUCKET_EMAIL` and `GREENLIGHT_BITBUCKET_API_TOKEN` (Bitbucket Cloud)

Run `greenlight doctor` to check native agents and provider tools.
For ACP agents, verify `acpx` or `acpx_path` separately because `doctor` does not validate ACP targets.

See [Provider Integration](/greenlight/guides/provider-integration/) for PR and CI setup per host.

## Update

```sh
greenlight update
greenlight update --beta
greenlight update -y
```

This downloads the latest release from GitHub, verifies the SHA-256 checksum, atomically replaces the binary, and resets the daemon so it picks up the new executable. It prefers the managed service path and falls back to a detached daemon if service startup is unavailable or fails.

`greenlight update` installs the latest stable release.
Use `greenlight update --beta` to opt into prereleases and install the latest beta when one is newer than the current stable release.
Use `greenlight update -y` to answer yes to update safety prompts.

Because `update` installs the latest official release binary, it installs a binary with the default self-hosted telemetry host and website ID. Disable telemetry with `GREENLIGHT_TELEMETRY=0`, or override the host and website ID with `GREENLIGHT_UMAMI_HOST` and `GREENLIGHT_UMAMI_WEBSITE_ID`.

If pending or running pipeline runs exist, the update warns that restarting the daemon can cause those runs to fail, prints each active run's ID, status, branch, and short head SHA, and prompts before continuing.
If the running daemon was started from a different binary, the update prompts before replacing it.
Pass `-y` or `--yes` to continue through these prompts while still printing warnings.
If the daemon executable path cannot be determined, the update aborts before replacing the binary.
If the daemon does not come back cleanly after a successful replacement, the new binary stays installed but the command reports the daemon reset failure.

Background update checks run automatically on each CLI invocation (except `update` itself). Suppress with `GREENLIGHT_NO_UPDATE_CHECK=1`.

## Remove from a repo

```sh
greenlight eject
```

Removes the `greenlight` remote, deletes the bare repo, cleans up worktrees, and removes the database record.
It does not remove repo-local agent skill files created by `greenlight init`.

## Uninstall

Stop the daemon, delete the binary, and clear state:

```sh
greenlight daemon stop
rm -f ~/.local/bin/greenlight /usr/local/bin/greenlight
rm -rf ~/.greenlight
```

On macOS, also remove `~/Library/LaunchAgents/com.kinleyrabgay.greenlight.daemon.*.plist`. On Linux, also remove `~/.config/systemd/user/greenlight-daemon-*.service`. On Windows, remove the `greenlight-daemon-*` Task Scheduler task.
