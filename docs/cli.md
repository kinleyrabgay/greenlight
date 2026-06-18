# CLI

## `greenlight`

Attach to the active run on the current branch. If none exists in a TTY, starts the setup wizard (branch → commit → push → attach). `-y` accepts wizard defaults; without a TTY it falls back to listing recent runs.

```sh
greenlight
greenlight --skip test,lint
```

| Flag | Default | Description |
|---|---|---|
| `-y`, `--yes` | `false` | Run the wizard and accept defaults |
| `--skip` | — | Comma-separated steps to skip for a new run |

Step names: `intent`, `rebase`, `review`, `test`, `document`, `lint`, `push`, `pr`, `ci`.

## `greenlight init`

Initialize or refresh the gate for the current repo: create/refresh the bare gate repo, install the post-receive hook, add the `greenlight` remote, record the repo, install the `/greenlight` skill at user level, and ensure the daemon runs. Re-running is safe and reports `Gate already initialized (refreshed)`.

```sh
greenlight init
greenlight init --framework angular
```

| Flag | Default | Description |
|---|---|---|
| `--framework` | — | Scaffold `.greenlight.yaml` selecting a framework profile (validated against `greenlight profiles`). Won't overwrite an existing config. |

## `greenlight profiles`

List built-in and user framework profiles, marking user overrides, and show the profile auto-detected for the current directory.

```sh
greenlight profiles
```

User profiles live in `~/.greenlight/profiles/<name>.md` and override the built-in of the same name. See [Configuration](configuration.md#framework-profiles).

## `greenlight axi`

Non-interactive interface for agents (TOON on stdout, progress on stderr, exit `1` on operational failure, `2` on bad usage). The `/greenlight` skill drives this.

```sh
greenlight axi                                   # status + next steps
greenlight axi run --intent "the user's goal"    # start/reattach, block at first gate
greenlight axi run --intent "..." --skip test,lint
greenlight axi respond --action approve
greenlight axi respond --action fix --findings F1,F2 --instructions "..."
greenlight axi respond --action skip
greenlight axi status [--run <id>]
greenlight axi logs --step review [--full] [--run <id>]
greenlight axi abort
```

`--intent` is the user's goal (not a diff description); greenlight uses it verbatim. `--yes` treats `auto-fix` and `ask-user` findings as standing consent and resolves gates until a decision point or outcome. When CI is green but the PR still needs a human merge, `axi run` exits with `outcome: checks-passed`.

## Other commands

| Command | Description |
|---|---|
| `greenlight attach [--run <id>]` | Open the TUI for the active run (or a specific run). |
| `greenlight rerun` | Start a new run from the last head SHA on the current branch. |
| `greenlight eject` | Remove the gate (remote, bare repo, worktrees, DB record). Leaves skill files. |
| `greenlight status` | Repo, daemon, and active-run status. |
| `greenlight runs [--limit <n>]` | List recent runs (default 10). |
| `greenlight stats` | Historical usage: changes, rescue rate, fixes by step, top repos. |
| `greenlight doctor` | Check `git`, `gh`, data dir, SQLite, daemon, and agent binaries. |
| `greenlight update [--beta] [-y]` | Download + verify + atomically replace the binary, reset the daemon. |
| `greenlight daemon start\|stop\|restart\|status` | Manage the background daemon. |

## Skipping steps for one push

```sh
greenlight --skip test,lint
git push -o greenlight.skip=test,lint greenlight <branch>
greenlight axi run --intent "..." --skip test,lint
```

Per-run skips don't change the pipeline — it always has all nine steps.
