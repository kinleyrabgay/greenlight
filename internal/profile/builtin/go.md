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
# Go project guide

## Stack

- Go (modules), standard toolchain; `golangci-lint` for lint, `go test -race` for tests.

## Architecture

- Small packages with clear ownership; `cmd/` for entrypoints, `internal/` for private code.
- Accept interfaces, return concrete types; keep interfaces small.
- Thread `context.Context` through I/O, subprocess, and networked work.

## Do

- Check every returned error; wrap with `fmt.Errorf("...: %w", err)` for context.
- Take `context.Context` as the first parameter for blocking calls; honor cancellation.
- Give every goroutine a clear exit path; guard shared state; run tests with `-race`.
- Release resources with `defer` (files, locks, rows).
- Document exported identifiers; prefer usable zero values.

## Don't

- Don't discard errors with `_` unless clearly justified.
- Don't store a `context.Context` in a struct.
- Don't leak goroutines or use time-based sleeps in tests.
- Don't hand-edit generated files (`*.pb.go`, `*.generated.go`).
