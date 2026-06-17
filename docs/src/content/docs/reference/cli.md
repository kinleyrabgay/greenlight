---
title: CLI Commands
description: Complete reference for all greenlight commands and flags.
---

## greenlight

Attach to the active pipeline run for the current branch when one exists. If none exists, bare `greenlight` can start the setup wizard to create a branch, commit changes, push through the gate, wait for the daemon to register the new run, and then attach. If the push succeeds but no run is registered, that wizard path now exits with an explicit error instead of silently falling through. By default this wizard path is interactive and only runs in a TTY session. In non-interactive contexts, bare `greenlight` falls back to showing the last 5 runs inline unless you pass `-y` or `--yes` to run the wizard and accept defaults automatically. When a TTY is available, `-y` keeps the wizard visible, shows a brief `waiting for run…` state after push, and auto-advances the default path; without a TTY it falls back to the headless path.

```sh
greenlight
greenlight --skip test,lint
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `-y`, `--yes` | `bool` | `false` | Run setup wizard and accept defaults automatically |
| `--skip` | `string` | (none) | Comma-separated pipeline steps to skip for a new run |

Unlike `greenlight attach`, bare `greenlight` only auto-attaches to an active run on the current branch.
`--skip` only applies when bare `greenlight` starts a new pipeline run through the wizard; it does not skip a step on an already-active run.
Valid step names are `intent`, `rebase`, `review`, `test`, `document`, `lint`, `push`, `pr`, and `ci`.

## greenlight init

Initialize or refresh the gate for the current repository.

```sh
greenlight init
```

Creates or refreshes a local bare repo, installs the post-receive hook, best-effort isolates the gate repo's hook path from shared git config changes when Git supports `config --worktree`, adds or repairs the `greenlight` git remote, detects the default branch, records or updates the repo in SQLite, installs the `/greenlight` agent skill at user level into `~/.claude/skills/greenlight/SKILL.md` and `~/.agents/skills/greenlight/SKILL.md`, and ensures the daemon is running, installing the managed service when available and falling back to a detached daemon otherwise.
`init` writes no skill files into the repo; the user-level copies cover every supported agent (`~/.claude/skills` for Claude Code, `~/.agents/skills` for Codex, OpenCode, Rovo Dev, and Pi) across all repos.
If the home `.claude` links to `.agents`, `.claude/skills` links to `.agents/skills`, or the reverse, `init` follows that layout and still makes the skill readable from both logical paths.
If the repo still contains a vendored skill copy written by an older greenlight version, `init` leaves it untouched and prints a notice that it is no longer needed and can be removed.
The gate advertises Git push-option support, so you can skip steps for one push with `git push -o greenlight.skip=test,lint greenlight <branch>`.

Re-running `init` on an already-initialized repo succeeds and reports `Gate already initialized (refreshed)`.
It refreshes managed gate wiring, origin/default-branch metadata, hook-path isolation, and the installed agent skill, overwriting any stale `SKILL.md` content from an older binary.
If you rename or move an initialized working directory and the old path no longer exists, re-running `init` from the new path reattaches the existing gate, preserves the repo ID and run history, and updates the stored working path.
If you copy an initialized working directory while the original still exists, the copy is treated as a separate repo and gets a fresh gate.
Fresh init rolls back gate setup when a required gate or daemon step fails; refresh does not eject a pre-existing gate if daemon startup fails.
Skill installation is best-effort: if the skill write fails, init reports it and leaves the working gate in place.

## greenlight axi

Agent eXperience Interface for non-interactive agents.
Most agent workflows use the installed `/greenlight` skill, which drives this command surface underneath.
It prints TOON to stdout, prints progress to stderr, and uses structured stdout errors with exit code `1` for operational failures and `2` for bad usage.

```sh
greenlight axi
```

With no subcommand, shows the executable path, description, repo, current branch, daemon state, recent runs, and next-step help.
When the current branch has an active run, that run appears as `active_run` with any approval gate and help for `axi respond` or `axi abort`.
When only another branch has an active run, that run appears as `other_branch_active_run`; the help tells agents to leave it alone and start validation for the current branch.

## greenlight axi run

Start or reattach to validation for the current branch, blocking until the first approval gate, CI-ready decision point, or final outcome.
An active run on another branch does not block starting validation for the current branch.

```sh
greenlight axi run --intent "the user's goal"
greenlight axi run --intent "the user's goal" --skip test,lint
greenlight axi run --intent "the user's goal" --yes
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--intent` | `string` | (none) | What the user set out to accomplish; required to start a new run |
| `-y`, `--yes` | `bool` | `false` | Auto-resolve every gate until a decision point or outcome |
| `--skip` | `string` | (none) | Comma-separated pipeline steps to skip |

`--intent` is not a description of the diff.
It is the user's goal or request, and greenlight uses it verbatim instead of transcript inference.
Err on the side of completeness: include the goal, important decisions and tradeoffs, constraints or approaches ruled in or out, and explicit requests that might otherwise look surprising in the diff.
When starting a new run, `axi run` refuses the default branch and uncommitted working trees with actionable errors instead of auto-branching or auto-committing.
Reattaching to an in-flight run does not require `--intent`.
With `--yes`, `axi run` treats both `action: auto-fix` and `action: ask-user` findings as standing consent for the pipeline to fix them by selecting every finding, then accepts the resulting fix review.
Gates with no findings or only `action: no-op` findings are approved as-is, and each step is fixed at most once so unresolved findings do not loop forever.
Without `--yes`, an agent driving `axi run` should stop when a gate contains `action: ask-user` findings and relay each finding's ID, file, and full description to the user before responding.
When the CI step is still monitoring an open PR and checks are green, `axi run` exits successfully with `outcome: checks-passed` instead of waiting for a human merge.
Treat that as the agent stopping point: ask the user to review and merge the PR from the `help` line.
Successful outcomes (`checks-passed` and `passed`) also carry `help` instructions telling the agent to summarize the run.
When the pipeline applied fixes, they include a `fixes` table and a `help` instruction to acknowledge the misses and list those fixes for the user's review.

## greenlight axi respond

Answer the current approval gate and continue until the next gate, CI-ready decision point, or final outcome.

```sh
greenlight axi respond --action approve
greenlight axi respond --action fix --findings F1,F2 --instructions "optional guidance"
greenlight axi respond --action fix --add-finding '{"description":"...","action":"auto-fix"}'
greenlight axi respond --action skip
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--action` | `string` | (none) | `approve`, `fix`, or `skip`; required |
| `--step` | `string` | awaiting step | Step to respond to |
| `--findings` | `string` | (none) | Comma-separated finding IDs for `--action fix` |
| `--instructions` | `string` | (none) | Guidance applied to selected findings |
| `--add-finding` | `string` | (none) | JSON finding object to add and fix |
| `-y`, `--yes` | `bool` | `false` | Auto-resolve every subsequent gate until a decision point or outcome |

After the explicit response, `--yes` uses the same auto-resolution behavior as `axi run --yes`: have the pipeline fix `auto-fix` and `ask-user` findings once, approve the fix review, approve gates that only contain non-actionable `no-op` findings, and stop at `outcome: checks-passed` when CI is green but the PR still needs a human merge.
The same successful-output reporting instructions apply to `axi respond` results.

## greenlight axi status

Show a run, preferring the current branch's active or most recent run before falling back to repo-wide active or recent runs.

```sh
greenlight axi status
greenlight axi status --run <id>
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--run` | `string` | resolved run | Inspect a specific run ID |

## greenlight axi logs

Show the log output of one pipeline step.

```sh
greenlight axi logs --step review
greenlight axi logs --step review --full
greenlight axi logs --step review --run <id>
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--step` | `string` | (none) | Step name; required |
| `--run` | `string` | resolved run | Run ID to inspect |
| `--full` | `bool` | `false` | Show the entire log instead of the tail |

Without `--full`, long logs show the last 40 lines and a help hint for the full log.

## greenlight axi abort

Cancel the active run for the current branch.
Active runs on other branches are left alone.

```sh
greenlight axi abort
```

If there is no active run, this succeeds as a no-op.

## greenlight eject

Remove the gate from the current repository.

```sh
greenlight eject
```

Removes the `greenlight` remote, deletes the bare repo directory, cleans up worktrees, and deletes the database record (cascades to runs and steps).
It does not remove repo-local agent skill files created by `init`.

## greenlight attach

Attach to the active pipeline run.

```sh
greenlight attach [--run <id>]
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--run` | `string` | (none) | Attach to a specific run ID instead of the active run |

Opens the TUI for the active run anywhere in the current repo. If `--run` is specified, attaches to that specific run regardless of branch. Unlike bare `greenlight`, this does not stay branch-scoped before falling back.

## greenlight rerun

Rerun the pipeline for the current branch.

```sh
greenlight rerun
```

Starts a new pipeline run using the last-known head SHA on the current branch. Useful for retrying after a fix or configuration change.

## greenlight status

Show repo, daemon, and active run status.

```sh
greenlight status
```

Displays:
- Repo path and upstream URL
- Gate path
- Daemon status (running/stopped, PID)
- Active run details: ID, branch, status, head SHA, start time

## greenlight runs

List recorded pipeline runs for the current repo.

```sh
greenlight runs [--limit <n>]
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--limit` | `int` | `10` | Maximum number of runs to display |

Shows runs newest-first with branch, status (styled), short SHA, timestamp, and PR URL if set.

## greenlight stats

Show historical usage stats across all repos.

```sh
greenlight stats
```

Displays total changes, rescued changes, rescue rate, reported and fixed mistakes, fixes by pipeline step, and the top repos by rescue activity.

## greenlight doctor

Check system health and dependencies.

```sh
greenlight doctor
```

Checks:
- `git` binary
- `gh` CLI (optional, needed for GitHub PR and CI steps)
- Data directory (`~/.greenlight/`)
- SQLite database
- Daemon status
- Native agent binaries: `claude`, `codex`, `acli`, `opencode`, `pi`

Uses indicators: `✓` (available), `–` (not found, optional), `✗` (problem detected).

`doctor` does not validate `acpx` or ACP targets. For `agent: acp:<target>`, verify `acpx_path` yourself.

`doctor` currently checks `gh` availability only. For GitLab PR and CI steps, install and authenticate `glab`. For Bitbucket Cloud PR and CI steps, set `GREENLIGHT_BITBUCKET_EMAIL` and `GREENLIGHT_BITBUCKET_API_TOKEN`.

## greenlight update

Update the installed binary and reset the daemon.

```sh
greenlight update
greenlight update --beta
greenlight update -y
```

Downloads the latest release, verifies the SHA-256 checksum, atomically replaces the running binary, and resets the daemon when it is running or stale daemon artifacts exist so the new executable is picked up, preferring the managed service path and falling back to a detached daemon if service startup is unavailable or fails.
By default this installs the latest stable release.
Pass `--beta` to include prereleases and install the latest beta when one is newer than the current stable release.
If pending or running pipeline runs exist, update warns that restarting the daemon can cause those runs to fail, prints each active run's ID, status, branch, and short head SHA, and prompts before continuing.
If the daemon is running from a different executable path, update prompts before replacing it.
Pass `-y` or `--yes` to continue through update safety prompts while still printing warnings.
If the daemon executable path cannot be determined, the update aborts before replacement.
If the daemon does not come back cleanly after a successful replacement, the command reports that failure.
On macOS, removes the quarantine extended attribute.

Because `update` installs the latest official release binary, the replacement binary includes the default self-hosted telemetry host and website ID. Disable telemetry with `GREENLIGHT_TELEMETRY=0`, or override the host and website ID with `GREENLIGHT_UMAMI_HOST` and `GREENLIGHT_UMAMI_WEBSITE_ID`.

Background update checks run automatically on each CLI invocation (except `update` itself). If a newer version is available, a notification is printed to stderr. Suppressed for dev builds or when `GREENLIGHT_NO_UPDATE_CHECK=1` is set.

## greenlight daemon start

Start the daemon, installing or refreshing the managed service when possible.

```sh
greenlight daemon start
```

Prefers the managed service path and falls back to a detached daemon if service install or startup is unavailable or fails. If the daemon is already running, the command refreshes a stale macOS `launchd` or Linux `systemd` service definition and restarts through the managed service; if the definition is unchanged, it reports that the daemon is already running.

## greenlight daemon stop

Stop the running daemon process.

```sh
greenlight daemon stop
```

This does not remove the managed service. A later `greenlight`, `greenlight daemon start`, `init`, `attach`, `rerun`, or `update` can start the daemon again through the same service manager when available, or as a detached daemon otherwise.

## greenlight daemon restart

Restart the daemon.

```sh
greenlight daemon restart
```

Stops the current daemon and starts it again. This works whether the daemon is currently running or not.

## greenlight daemon status

Check whether the daemon is running.

```sh
greenlight daemon status
```

Shows the PID if the daemon is running.
