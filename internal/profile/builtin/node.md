---
commands:
  lint: "npm run lint"
  test: "npm test"
  format: "npm run format"
ignore_patterns:
  - "dist/**"
  - "build/**"
  - "coverage/**"
  - "node_modules/**"
  - "**/*.generated.ts"
---
# Node.js (+ TypeScript) backend review rules

Apply these when reviewing or documenting Node/TypeScript backend changes:

- **Async correctness** — every Promise is awaited or explicitly handled; no floating promises. Flag missing `await` in `try/catch`.
- **Errors** — throw `Error` (or subclasses), not strings; don't swallow errors with empty `catch`. Propagate or log with context.
- **Input validation** — validate and narrow all external input (HTTP bodies, query params, env) at the boundary before use.
- **Secrets** — never log secrets/tokens; read config from env, not hardcoded.
- **Resource cleanup** — close DB connections, file handles, timers; avoid leaks in long-lived processes.
- **TypeScript** — no `any` in new code; prefer `unknown` + narrowing at boundaries; keep return types explicit on exported functions.
- **Concurrency** — guard shared mutable state; prefer immutable data flow.
- **Logging** — structured, leveled logging; no stray `console.log` in production paths.
