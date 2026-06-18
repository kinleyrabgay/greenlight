---
commands:
  lint: "npm run lint"
  test: "npm test"
  format: "npm run format"
ignore_patterns:
  - ".svelte-kit/**"
  - "build/**"
  - "node_modules/**"
---
# Svelte / SvelteKit project guide

## Stack

- Svelte (5 runes or 4 stores) + SvelteKit, TypeScript, Vite.
- Vitest/Playwright for tests; ESLint + Prettier (with `prettier-plugin-svelte`).

## Architecture

- SvelteKit routing in `src/routes`; `+page`, `+layout`, `+server`, and `load` functions.
- Load data in `load` (server or universal); keep secrets in `+server`/server `load`.
- Reactivity via runes (`$state`, `$derived`, `$effect`) in Svelte 5, or `$:`/stores in v4.

## Do

- Derive values with `$derived` (or `$:`) instead of manually syncing state.
- Fetch route data in `load`, not in component lifecycle where possible.
- Keep server-only logic in `+server.ts`/server `load`; gate secrets behind `$env/dynamic/private`.
- Clean up subscriptions and effects; prefer stores/runes over ad-hoc globals.

## Don't

- Don't put secrets or DB access in client-reachable code.
- Don't overuse `$effect`/`$:` for values that should be derived.
- Don't mutate props directly; use events/bindings or callbacks.
