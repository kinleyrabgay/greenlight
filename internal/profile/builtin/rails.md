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
# Ruby on Rails project guide

## Stack

- Ruby on Rails (MVC), Ruby; RSpec/Minitest for tests; RuboCop for lint/format.
- ActiveRecord ORM; background jobs via ActiveJob/Sidekiq.

## Architecture

- Skinny controllers, fat models/service objects; views render, controllers orchestrate.
- Business logic in models/POROs/services, not controllers or views.
- Long or external work goes to a background job, not the request cycle.

## Do

- Permit params explicitly with strong parameters.
- Eager-load associations (`includes`) to avoid N+1 queries.
- Write reversible migrations; index foreign keys and lookup columns.
- Use parameterized queries / ActiveRecord; verify authorization on every action.
- Cover model validations and add request specs for controllers.

## Don't

- Don't `params.permit!` or mass-assign raw params.
- Don't interpolate user input into SQL.
- Don't hide side effects in heavy `before_save`/`after_*` callbacks.
- Don't hand-edit `db/schema.rb`.
