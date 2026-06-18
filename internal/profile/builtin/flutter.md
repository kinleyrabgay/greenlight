---
commands:
  lint: "flutter analyze"
  test: "flutter test"
  format: "dart format ."
ignore_patterns:
  - "**/*.g.dart"
  - "**/*.freezed.dart"
  - ".dart_tool/**"
  - "build/**"
---
# Flutter (Dart) project guide

## Stack

- Flutter + Dart; widget-based UI; `flutter analyze` for lint, `flutter test` for tests.
- State management: one of Riverpod / Bloc / Provider / GetX (follow what the repo uses).

## Architecture

- Compose small widgets; prefer `const` constructors; split big `build` methods.
- Keep business logic out of widgets — in notifiers/blocs/services.
- Feature-first folders; immutable models (often `freezed`).

## Do

- Use `const` widgets wherever possible to cut rebuilds.
- Dispose controllers, animation controllers, and stream subscriptions.
- Handle loading/error/empty states for async data (`FutureBuilder`/`AsyncValue`).
- Keep `setState` scoped narrowly; lift shared state into the chosen state solution.
- Use keys correctly when reordering or preserving widget state.

## Don't

- Don't do expensive work or I/O inside `build()`.
- Don't call `setState` after dispose or during build.
- Don't hardcode sizes/colors — use the theme and `MediaQuery`/`LayoutBuilder`.
- Don't hand-edit generated files (`*.g.dart`, `*.freezed.dart`).
