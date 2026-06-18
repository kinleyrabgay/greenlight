# Pipeline

Every gated push runs a fixed, opinionated sequence. The order is not configurable — what each step *runs* is.

```
intent → rebase → review → test → document → lint → push → pr → ci
```

## The nine steps

| # | Step | What it does | Default auto-fix |
|---|---|---|---|
| 1 | **Intent** | Use supplied `--intent`, or infer it from recent local agent transcripts | n/a |
| 2 | **Rebase** | Fetch upstream, rebase your branch onto it | `3` |
| 3 | **Review** | AI code review of your diff (framework rules injected here) | `0` (asks) |
| 4 | **Test** | Run the baseline test command and gather evidence | `3` |
| 5 | **Document** | Update docs when needed, report unresolved gaps | initial pass |
| 6 | **Lint** | Run lint / static analysis | `3` |
| 7 | **Push** | Push the validated branch upstream | n/a |
| 8 | **PR** | Create or update the pull request | n/a |
| 9 | **CI** | Watch CI + mergeability, auto-fix failures | `3` |

## Why this order

- **Intent first** so downstream prompts and the PR description carry the author's goal.
- **Rebase next** so everything runs against fresh upstream. No diff after rebase → the rest is skipped.
- **Review before test** so the agent reads fresh code, not code it just touched.
- **Document after test** so docs describe code that's known to work.
- **Lint last** among local checks so it doesn't churn over code that may still change.
- **Push → PR → CI** only after every local check passes. CI is the only step that talks to the outside world.

## What a step can do

- **Complete** and advance.
- **Return findings** with a severity (`error`/`warning`/`info`) and an action (`auto-fix`, `ask-user`, `no-op`).
- **Auto-fix** when the step's `auto_fix` limit is above 0 and a finding is `auto-fix`-eligible.
- **Pause for approval** when blocking findings remain or any finding is `ask-user`.
- **Skip** when there's nothing to do (no diff, unsupported host).
- **Fail** on a fatal error and stop the pipeline.

## What "passed the gate" means

Because the sequence is fixed, a green gate always means the same thing:

- the branch was checked against fresh upstream first;
- review, tests, docs, and lint ran before any upstream push;
- a human stayed in control wherever judgment was needed;
- push, PR, and CI happened only after the local gate was satisfied.

## What you can configure

You can't reorder, add, or permanently remove steps. You *can*: swap the agent, set `commands.test`/`.lint`/`.format` (or inherit them from a [framework profile](configuration.md#framework-profiles)), tune `auto_fix` limits, ignore paths, toggle intent inference, store test evidence in-repo, and skip steps for a single run with `greenlight --skip <steps>` or `git push -o greenlight.skip=<steps>`.
