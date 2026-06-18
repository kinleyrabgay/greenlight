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
# Angular project guide

## Stack

- Angular v17+ (v20/v21 idioms), TypeScript, standalone components.
- Signals for reactivity; Signal Forms (`@angular/forms/signals`) for new forms.
- Often Nx monorepo, PrimeNG/Material, Transloco i18n, NgRx Signal Store, TanStack Query.

## Architecture

- Standalone components — no NgModules for new code.
- Feature-first layout: components, services (api/bl/facade), state, models per feature.
- Server state via TanStack Query behind facade services; shared client state in Signal Store; never call `HttpClient` directly from components.

## Do

- Use `signal()` / `computed()` / `effect()`; prefer `computed()` over manual subscriptions.
- Use `input()` / `output()` / `model()`, not `@Input()` / `@Output()` decorators.
- Use `inject()`, not constructor injection, in new code.
- Use `@if` / `@for` / `@switch`, not `*ngIf` / `*ngFor` / `*ngSwitch`.
- Clean up subscriptions with `DestroyRef` + `takeUntilDestroyed()`.
- Add every new i18n key to ALL configured languages (missing keys render as raw strings).

## Don't

- Don't use `FormBuilder` / `new FormGroup()` for new forms — use Signal Forms.
- Don't use an `effect()` just to copy one signal into another — derive with `computed()`.
- Don't duplicate server data into a store — keep it in the query cache.
- Don't assume `zone.js` — apps are zoneless.
- Don't hand-edit generated files (e.g. GraphQL codegen output).
