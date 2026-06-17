---
commands:
  lint: "golangci-lint run ./..."
  test: "go test -race ./..."
  format: "gofmt -w ."
ignore_patterns:
  - "*.generated.go"
  - "**/*.pb.go"
  - "vendor/**"
---
# Go review rules

Apply these when reviewing or documenting Go changes:

- **Error handling** — check every returned error; wrap with `fmt.Errorf("...: %w", err)` for context. Never discard errors with `_` unless clearly justified.
- **Context** — accept `context.Context` as the first parameter for I/O and blocking calls; honor cancellation; never store a context in a struct.
- **Goroutines** — every goroutine has a clear exit path; no leaks. Guard shared state with mutexes or channels; run with `-race`.
- **defer** — release resources with `defer` (files, locks, rows); beware `defer` in loops.
- **Interfaces** — accept interfaces, return concrete types; keep interfaces small.
- **Nil & zero values** — guard nil maps/slices/pointers; prefer usable zero values.
- **Naming & idioms** — exported identifiers documented; errors are `ErrXxx` or wrapped; no stutter (`pkg.PkgThing`).
- **Tests** — table-driven where it fits; `t.Parallel()` when safe; no time-based flakiness.
