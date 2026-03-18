# AGENTS.md

## Learned User Preferences

- Backend API paths must not use a `/v1` prefix — routes are `/api/...` (e.g., `/api/auth/login`, `/api/logs`), never `/api/v1/...`.
- All `<select>` elements must use the reusable `SelectField` component with `appearance-none`, `pr-8`, and a custom `ChevronDown` icon; never use bare `<select>` with a tight or missing arrow icon.
- Every main page header (`h1`) must include a consistent, relevant Lucide icon (e.g., `LayoutDashboard` for Overview, `ScrollText` for Logs, `Gauge` for APM, `Signal` for Status, `Server` for Sources).
- Save/submit buttons must never include an icon — use plain text label only.
- All tables that have pagination must expose a rows-per-page selector with options 10, 20, 50, 100, and All.
- Source API keys must be shown in full exactly once (on creation or rotation) with an eye-toggle; after hiding or on any subsequent page load the key is masked as `ulam_live_****`. No "save this key" warning banner is needed.
- The floating chat widget must be named "One Log AI" (not "AI Copilot"), must appear directly above its trigger button (not float to the top of the page), and must have a fixed panel size so it never shifts layout.
- Chat widget AI responses must wrap within the widget width; code blocks must be copyable and must not overflow.
- The AI chatbot system prompt must be broad enough to cover debugging, SQL, DevOps, security, performance, and any project-related topic — not limited to only log queries.
- Tailwind v4 shorthand syntax is required: use `bg-white/3` not `bg-white/[0.03]`, `bg-linear-to-r` not `bg-gradient-to-r`, etc.
- The pre-commit hook must gate on: `go fmt`, `go vet`, `go build` (backend) and `npm run lint`, `tsc --noEmit`, `npm run build` (frontend); all must pass with zero errors and zero warnings before a commit is allowed.
- `seed.go` must seed only the admin user (credentials from `ADMIN_EMAIL`/`ADMIN_PASSWORD` env vars, falling back to `admin@onelog.com` / `123456`); no dummy sources, logs, or other fixture data.

## Learned Workspace Facts

- **Project**: ULAM — Unified Log & Activity Monitor ("One-Log"). A centralized observability platform.
- **Tech stack**: Go 1.26 + Gin + GORM (backend), React 19 + Vite 7 + Tailwind v4 + TypeScript + TanStack Query (frontend), PostgreSQL 17.
- **Architecture**: Clean Architecture — `handler → service → repository → domain`; no business logic in handlers, no DB calls in services.
- **Auth**: JWT stored in httpOnly cookies (`ulam_access`, `ulam_refresh`); API Key authentication for log ingestion. Frontend sends `withCredentials: true`.
- **CORS**: `CORS_ALLOWED_ORIGIN` must be set to `http://localhost:5173` in `Backend/.env`; wildcard `*` is incompatible with `withCredentials: true`. `PATCH` must be in `Access-Control-Allow-Methods`.
- **AI**: Groq API used for log analysis and the "One Log AI" chatbot; model configurable via `GROQ_MODEL` env var; client is `pkg/ai/groq.go`.
- **Migrations**: Goose controls schema; SQL files live in `Backend/migrations/`; GORM `AutoMigrate` is used only for new additive columns.
- **API key format**: `ulam_live_<random>`; keys are hashed with `utils.HashAPIKey()` before storage; raw key is returned once on creation/rotation.
- **GitGuardian**: Historical false-positive findings are suppressed with `ignore-known-secrets: true` in `.gitguardian.yaml`; specific `ignore-matches` entries added for legacy `ulam_live_*` and `dev_*` keys.
- **Frontend ports**: Backend on `localhost:8080`, Frontend on `localhost:5173`.
- **JSONB queries**: Use `jsonb_exists(column, 'key')` instead of `column ? 'key'` in raw SQL passed through GORM to avoid placeholder conflicts.
- **Bundle size**: `vite.config.ts` sets `build.chunkSizeWarningLimit: 1500` to suppress the cosmetic oversized-bundle warning from `react-markdown` / `recharts`.
