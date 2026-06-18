---
commands:
  lint: "next lint"
  test: "npm test -- --run"
  format: "npm run format"
ignore_patterns:
  - ".next/**"
  - "out/**"
  - "node_modules/**"
  - "next-env.d.ts"
  - "**/*.generated.ts"
---
# Next.js project guide

## Stack

- Next.js (App Router), React 18+, TypeScript.
- Server Components by default; Client Components opt in with `"use client"`.
- Server Actions, Route Handlers; deploy target often Vercel/Node.

## Architecture

- App Router: `app/` route segments with `layout`, `page`, `loading`, `error`, `route` files.
- Keep the Server/Client boundary explicit and minimal — push `"use client"` to leaves.
- Fetch data in Server Components / Route Handlers; mutate via Server Actions.
- Use the framework cache (`fetch` cache, `revalidate`, tags) instead of ad-hoc caching.

## Do

- Default to Server Components; only mark a component client when it needs state/effects/browser APIs.
- Keep secrets and `fs`/DB access server-side; expose only `NEXT_PUBLIC_*` to the client.
- Use `<Image>`, `<Link>`, and `next/font` for images, navigation, and fonts.
- Co-locate route files under `app/`; use route groups for organization.

## Don't

- Don't add `"use client"` to a whole page when only one widget is interactive.
- Don't import server-only modules into client components.
- Don't fetch the same data in many client components — fetch on the server and pass down.
- Don't disable caching globally to "fix" staleness — revalidate by tag/path.
