---
commands:
  lint: "npm run lint"
  test: "npm test -- --run"
  format: "npm run format"
ignore_patterns:
  - "dist/**"
  - "build/**"
  - ".next/**"
  - "node_modules/**"
  - "**/*.generated.ts"
---
# React (+ TypeScript) review rules

Apply these when reviewing or documenting React changes:

- **Rules of Hooks** — hooks only at the top level of components/custom hooks, never in conditionals or loops. Flag any violation.
- **Dependency arrays** — `useEffect`/`useMemo`/`useCallback` deps must be complete and correct. Flag missing deps and effects that should be derived state instead.
- **Derived state** — compute during render, don't mirror props into state with an effect.
- **Keys** — stable, unique `key` on list items; never the array index when the list reorders.
- **Effects** — an effect that only transforms data is usually wrong; prefer computing inline or `useMemo`.
- **Server Components / Next.js** — keep `"use client"` boundaries minimal; don't put server-only code (secrets, fs) in client components. Data fetching belongs in Server Components or route handlers.
- **Accessibility** — interactive elements are real buttons/links with labels; images have `alt`.
- **TypeScript** — no `any` in new code; prefer discriminated unions over boolean flags for state.
- **Performance** — memoize expensive children only when profiling justifies it; avoid premature `memo`/`useCallback` noise.
