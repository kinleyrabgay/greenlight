---
commands:
  lint: "bundle exec rubocop"
  test: "bundle exec rspec"
  format: "bundle exec rubocop -A"
ignore_patterns:
  - "db/schema.rb"
  - "tmp/**"
  - "log/**"
  - "vendor/bundle/**"
---
# Ruby on Rails review rules

Apply these when reviewing or documenting Rails changes:

- **Strong parameters** — controllers permit params explicitly; never `params.permit!` or mass-assign raw params.
- **N+1 queries** — eager-load associations (`includes`) where iteration hits the DB; flag obvious N+1 patterns.
- **Migrations** — reversible (`change` or paired `up`/`down`); avoid data + schema changes in one irreversible migration; add indexes for foreign keys and lookups.
- **Security** — no SQL string interpolation (use parameterized queries / ActiveRecord); escape output; guard against mass-assignment; verify authorization on every action.
- **Fat model / skinny controller** — business logic in models/services, not controllers or views.
- **Callbacks** — avoid heavy or side-effecting `before_save`/`after_*` callbacks that hide control flow; prefer explicit service objects.
- **Background work** — long or external work goes to a job (ActiveJob/Sidekiq), not the request cycle.
- **Tests** — cover model validations, request specs for controllers; no reliance on `Time.now` without freezing.
- **Style** — follow RuboCop; keep methods short; prefer Ruby idioms (`&.`, `presence`, guard clauses).
