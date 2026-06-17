---
commands:
  lint: "yarn lint"
  test: "yarn nx run-many -t test"
  format: "yarn nx format:write"
ignore_patterns:
  - "**/*.generated.ts"
  - "dist/**"
  - ".nx/**"
  - "node_modules/**"
---
# Angular (v17+/v21) review rules

Apply these when reviewing or documenting Angular changes:

- **Standalone components** — no NgModules for new code.
- **Signals** — `signal()`, `computed()`, `effect()` for state. Prefer `computed()` over manual subscriptions. Flag `effect()` that merely copies one signal into another.
- **Signal forms** — use `form()` from `@angular/forms/signals`, not `FormBuilder`/`new FormGroup()`, for new forms and modals.
- **Inputs/outputs** — `input()`/`output()`/`model()`, not `@Input()`/`@Output()` decorators.
- **DI** — `inject()`, not constructor injection, in new code.
- **Control flow** — `@if`/`@for`/`@switch`, not `*ngIf`/`*ngFor`/`*ngSwitch`.
- **Cleanup** — `DestroyRef` + `takeUntilDestroyed()` for subscriptions.
- **Zoneless** — no `zone.js` assumptions; change detection is zoneless.
- **i18n** — every new translation key must exist in ALL configured languages. Missing keys render as raw strings in production. Flag keys added to only one language file.
- **Server state** — TanStack Query (`injectQuery`/`injectMutation`) via facade services; do not call `HttpClient` directly in components.
- **Shared state** — NgRx Signal Store (`signalStore`, `patchState`), not classic NgRx actions/reducers.

Do not flag generated files (e.g. GraphQL codegen output) for style.
