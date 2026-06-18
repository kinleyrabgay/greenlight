---
commands:
  lint: "npm run lint"
  test: "npm test -- --run"
  format: "npm run format"
ignore_patterns:
  - "dist/**"
  - "build/**"
  - "node_modules/**"
  - "**/*.generated.ts"
---
# React project guide

## Stack

- React 18+, TypeScript, function components + hooks.
- Vite or CRA build; Vitest/Jest + Testing Library; ESLint + Prettier.
- State: local hooks, Context, or a store (Redux Toolkit/Zustand) for shared state.

## Architecture

- Components are pure functions of props + state; side effects isolated in hooks.
- Co-locate component, styles, and tests; lift shared logic into custom hooks.
- Data fetching via a cache layer (TanStack Query/SWR), not ad-hoc effects.

## Do

- Follow the Rules of Hooks: hooks only at the top level, never in conditionals/loops.
- Keep `useEffect`/`useMemo`/`useCallback` dependency arrays complete and correct.
- Derive state during render; compute with `useMemo` instead of mirroring props into state.
- Give list items stable, unique `key`s (not the array index when the list reorders).
- Make interactive elements real buttons/links with labels; images get `alt`.

## Don't

- Don't use an effect that only transforms data — compute it inline.
- Don't introduce `any` in new code; prefer discriminated unions over boolean-flag state.
- Don't put server-only code (secrets, fs) in client components.
- Don't sprinkle `memo`/`useCallback` without a profiling reason.
