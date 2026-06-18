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
# Node.js (TypeScript) backend guide

## Stack

- Node.js + TypeScript; HTTP via Express/Fastify/Nest or similar.
- Vitest/Jest for tests; ESLint + Prettier; tsc/tsx or esbuild for build.

## Architecture

- Layered: routes/controllers → services (logic) → data access. Keep handlers thin.
- Validate external input at the boundary; never trust HTTP bodies, params, or env.
- Configuration from env; structured, leveled logging.

## Do

- Await or explicitly handle every Promise; no floating promises.
- Throw `Error` (or subclasses) with context; wrap with cause where useful.
- Validate and narrow all external input before use.
- Close DB connections, file handles, and timers; avoid leaks in long-lived processes.
- Keep exported function return types explicit; prefer `unknown` + narrowing over `any`.

## Don't

- Don't swallow errors with empty `catch` blocks.
- Don't log secrets or tokens.
- Don't hardcode config that belongs in env.
- Don't share mutable global state without guarding it.
