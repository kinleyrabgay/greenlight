# Contributing

This repo _is_ greenlight, so contributions go through the tool itself.

## Workflow

1. Fork and clone.
2. Create a branch and make your changes.
3. Set up the gate: `greenlight init`.
4. Commit, then push through the gate instead of `origin`:

   ```sh
   git push greenlight
   ```

5. Run `greenlight` to attach to the pipeline, act on findings, and let it open the PR once green.

## Repo conventions

- Go 1.25+, standard toolchain. See `AGENTS.md` for the full agent/dev guide.
- Run `make fmt`, `make lint`, and `make test` before pushing. Add `make e2e` when you touch agent integrations, the e2e harness, or recorded fixtures.
- Run `make skill` when you change the canonical skill content under `internal/skill`; `make lint` fails on skill drift.
- Use `make e2e-record` only when an agent wire format changes or you add a fixture — it overwrites `internal/e2e/fixtures/`, spends real API quota, and should be reviewed before committing.
- Keep `README.md` high-level; deep reference belongs in `docs/`.
