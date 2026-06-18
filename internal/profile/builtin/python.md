---
commands:
  lint: "ruff check ."
  test: "pytest"
  format: "ruff format ."
ignore_patterns:
  - "**/*_pb2.py"
  - ".venv/**"
  - "build/**"
  - "dist/**"
  - "**/__pycache__/**"
---
# Python project guide

## Stack

- Python 3.10+; `ruff` for lint/format, `pytest` for tests.
- Packaging via `pyproject.toml`; dependencies via uv/poetry/pip.

## Architecture

- Clear module boundaries; keep side effects out of import time.
- Type-annotate public functions; validate external input at the boundary.
- Separate I/O from pure logic so logic is easy to test.

## Do

- Add type hints to new public functions and dataclasses/models.
- Use context managers (`with`) for files, connections, and locks.
- Raise specific exceptions; let them propagate with context.
- Keep functions small and pure where practical; cover them with `pytest`.
- Use logging, not `print`, in library/production code.

## Don't

- Don't use mutable default arguments (`def f(x=[])`).
- Don't catch bare `except:` or swallow exceptions silently.
- Don't do heavy work at import time.
- Don't hand-edit generated files (`*_pb2.py`).
