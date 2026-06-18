# GREENLIGHT.md

Project guide for [greenlight](https://github.com/kinleyrabgay/greenlight). greenlight reads this file during its review/audit and applies the rules below to your changes. Edit it freely — it is yours to control.

<!-- Example for a React project. Copy to your repo root (greenlight init seeds one for you) and tailor it. -->

## Stack

- React 18+, TypeScript, Vite. Vitest + Testing Library, ESLint + Prettier.
- Server state via TanStack Query; shared client state via Zustand.

## Architecture

- Feature folders: component + hook + test co-located. Shared logic in custom hooks under `src/hooks`.
- Data fetching only through the query layer (`src/api`), never ad-hoc `fetch` in components.
- Presentational vs. container split for non-trivial screens.

## Do

- Follow the Rules of Hooks; keep dependency arrays complete.
- Derive state during render / with `useMemo`; lift shared state into hooks or the store.
- Give list items stable, unique keys.
- Keep components typed; prefer discriminated unions over boolean-flag state.

## Don't

- Don't use an effect that only transforms data — compute it inline.
- Don't introduce `any` in new code.
- Don't fetch the same data in multiple components — share via the query cache.
- Don't add `memo`/`useCallback` without a profiling reason.
