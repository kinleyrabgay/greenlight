# GREENLIGHT.md

Project guide for [greenlight](https://github.com/kinleyrabgay/greenlight). greenlight reads this file during its review/audit and applies the rules below to your changes. Edit it freely — it is yours to control.

<!-- Example for an Angular project. Copy to your repo root (greenlight init seeds one for you) and tailor it. -->

## Stack

- Angular v21, TypeScript, standalone components, Signals + Signal Forms.
- Nx monorepo, PrimeNG, Tailwind, Transloco i18n (en/de/fr/it), NgRx Signal Store, TanStack Query, Apollo GraphQL.

## Architecture

- Three-layer feature pattern: `*-api.ts` (GraphQL/REST) → `*-bl.ts` (pure logic) → `*-facade.ts` (public API for components). Never bypass layers.
- Server state lives in TanStack Query behind facades; shared client state in Signal Store; components never call `HttpClient` directly.
- Feature libraries under `libs/features/{domain}`; shared primitives in `@sbh/ui`.

## Do

- Use `signal()` / `computed()` / `effect()`; prefer `computed()` over manual subscriptions.
- Use `input()` / `output()` / `model()` and `inject()` (not decorators / constructor injection).
- Use `@if` / `@for` / `@switch` control flow.
- Use Signal Forms (`form()` from `@angular/forms/signals`) for new forms.
- Add every new i18n key to ALL four languages (en/de/fr/it).

## Don't

- Don't use `FormBuilder` / `new FormGroup()` for new forms.
- Don't use an `effect()` just to copy one signal into another.
- Don't duplicate server data into a store.
- Don't hand-edit generated GraphQL types (`libs/shared/references/src/lib/generated/graphql.ts`).
