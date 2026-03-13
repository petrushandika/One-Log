# 🛠️ Tech Stack

## Unified Log & Activity Monitor (ULAM)

---

## Stack Overview

```text
┌───────────────────────────────────────────────────┐
│  ULAM Technology Stack (Latest Stable — Mar 2026)  │
│                                                    │
│  Frontend     → React 19 + Vite 7                 │
│  Styling      → Tailwind CSS v4                   │
│  Backend      → Golang 1.26                       │
│  HTTP Router  → Gin v1.10                         │
│  ORM          → GORM v2 (v1.25)                   │
│  Database     → PostgreSQL 17                     │
│  Auth         → JWT (golang-jwt v5) + API Key     │
│  Email        → SMTP (net/smtp)                   │
│  AI Engine    → Groq API (Llama 3.3)              │
│  Container    → Docker + Docker Compose v2        │
│  Runtime      → Node.ts 22 LTS (frontend build)  │
│  Version Ctrl → Git                               │
└───────────────────────────────────────────────────┘
```

---

## Frontend

### React 19 + Vite 7

| Spec                | Detail                                                                  |
| ------------------- | ----------------------------------------------------------------------- |
| **Library**         | React 19 (Dec 2024)                                                     |
| **Build Tool**      | Vite 7.x                                                                |
| **Language**        | TypeScript 5.7+                                                         |
| **Why TypeScript?** | Type safety, IntelliSense, auto-complete, catch errors at compile time  |
| **Why React 19?**   | Server Components support, improved Actions API, stable Concurrent Mode |
| **Why Vite 7?**     | Fastest dev server, native ESM, HMR instant                             |

**Key dependencies:**

```json
{
  "@types/react": "^19.0.0",
  "@types/react-dom": "^19.0.0",
  "typescript": "^5.7.0",
  "react": "^19.0.0",
  "react-dom": "^19.0.0",
  "react-router-dom": "^7.0.0",
  "@tanstack/react-query": "^5.0.0",
  "axios": "^1.7.0",
  "recharts": "^2.15.0",
  "date-fns": "^4.1.0",
  "lucide-react": "^0.475.0"
}
```

**Project structure (Feature-Based):**

```text
frontend/
├── src/
│   ├── features/              # Fitur bisnis — unit mandiri
│   │   ├── auth/              # Login, JWT store
│   │   ├── logs/              # Log list, filters, detail
│   │   ├── sources/           # Source registration, rotate key
│   │   ├── activity/          # Auth events, audit trail
│   │   └── stats/             # Overview charts & cards
│   ├── shared/                # Komponen & utils lintas fitur
│   │   ├── components/ui/     # Button, Card, Badge, Modal...
│   │   ├── lib/               # Axios instance, React Query config
│   │   └── utils/             # format-date.ts, level-color.ts
│   ├── pages/                 # Thin layer, assembles features
│   │   ├── LoginPage.tsx
│   │   ├── OverviewPage.tsx
│   │   ├── LogsPage.tsx
│   │   ├── ActivityPage.tsx
│   │   └── SourcesPage.tsx
│   ├── router/
│   │   └── index.tsx          # Routes + auth guard
│   ├── App.tsx
│   └── main.tsx
├── index.html
├── tsconfig.json
└── vite.config.ts
```

---

### Tailwind CSS v4

| Spec        | Detail                                                                                     |
| ----------- | ------------------------------------------------------------------------------------------ |
| **Version** | Tailwind CSS v4.x (Jan 2025)                                                               |
| **Why v4?** | Hingga 5x lebih cepat, konfigurasi berbasis CSS (bukan JS config), first-party Vite plugin |
| **Plugin**  | `@tailwindcss/vite` (tidak perlu PostCSS manual)                                           |

> ⚠️ **Breaking change dari v3**: Tailwind v4 menggunakan CSS-based config, bukan `tailwind.config.ts`.

**Setup di `vite.config.ts`:**

```typescript
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
});
```

**Custom theme di `src/index.css`:**

```css
@import "tailwindcss";

@theme {
  --color-primary: oklch(0.6 0.2 250);
  --color-danger: oklch(0.55 0.22 25);
  --color-warning: oklch(0.72 0.18 70);
  --font-sans: "Inter", sans-serif;
}
```

---

## Backend

### Golang 1.26

| Spec          | Detail                                                              |
| ------------- | ------------------------------------------------------------------- |
| **Version**   | Go 1.26.0 (February 2026)                                           |
| **Why Go?**   | Performa tinggi, goroutines native, single compiled binary          |
| **Strengths** | Concurrency native untuk async email, memory footprint sangat kecil |

**Key packages:**

```go
// go.mod
require (
    github.com/gin-gonic/gin          v1.10.0
    gorm.io/gorm                      v1.25.12
    gorm.io/driver/postgres           v1.5.11
    github.com/golang-jwt/jwt/v5      v5.2.1
    github.com/joho/godotenv          v1.5.1
    golang.org/x/crypto               v0.31.0
    gorm.io/datatypes                 v1.2.5
)
```

---

### Gin v1.10

| Spec            | Detail                                                           |
| --------------- | ---------------------------------------------------------------- |
| **Package**     | `github.com/gin-gonic/gin`                                       |
| **Version**     | v1.10.0 (2024, latest stable)                                    |
| **Why Gin?**    | Fastest Go HTTP router, battle-tested, rich middleware ecosystem |
| **Alternative** | Fiber v3 (if prefer Express-like API)                            |

**Route structure:**

```go
// router/router.go
r := gin.New()
r.Use(gin.Recovery(), middleware.RequestLogger())

api := r.Group("/api")

// Public
api.POST("/auth/login",   authHandler.Login)
api.POST("/auth/refresh", authHandler.Refresh)

// Ingestion — API Key auth (semua source yang terdaftar)
ingestion := api.Group("/")
ingestion.Use(middleware.APIKeyAuth())
{
    ingestion.POST("/ingest", ingestHandler.Handle)
}

// Dashboard — JWT auth (for admin)
dashboard := api.Group("/")
dashboard.Use(middleware.JWTAuth())
{
    dashboard.GET("/logs",               logHandler.List)
    dashboard.GET("/logs/:id",           logHandler.GetByID)
    dashboard.GET("/sources",            sourceHandler.List)
    dashboard.POST("/sources",           sourceHandler.Create)
    dashboard.GET("/sources/:slug",      sourceHandler.GetBySlug)
    dashboard.PATCH("/sources/:slug",    sourceHandler.Update)
    dashboard.POST("/sources/:slug/rotate-key", sourceHandler.RotateKey)
    dashboard.GET("/activity",           activityHandler.List)
    dashboard.GET("/activity/summary",   activityHandler.Summary)
    dashboard.GET("/activity/users/:id", activityHandler.ByUser)
    dashboard.GET("/stats/overview",     statsHandler.Overview)
    dashboard.GET("/stats/sources/:slug", statsHandler.BySource)
}
```

---

### GORM v2 (v1.25.12)

| Spec          | Detail                                                                        |
| ------------- | ----------------------------------------------------------------------------- |
| **Package**   | `gorm.io/gorm`                                                                |
| **Version**   | v1.25.12 (latest stable, 2025)                                                |
| **Driver**    | `gorm.io/driver/postgres` v1.5.11                                             |
| **Why GORM?** | Developer-friendly, auto-migrate, schema as code, JSONB support via datatypes |

**Model definitions:**

```go
// domain/log.go
type LogEntry struct {
    ID         uint           `gorm:"primaryKey;autoIncrement"          json:"id"`
    SourceID   string         `gorm:"index;type:varchar(50);not null"   json:"source_id"`
    Category   string         `gorm:"index;type:varchar(30);not null"   json:"category"`
    Level      string         `gorm:"type:varchar(20);not null"         json:"level"`
    Message    string         `gorm:"type:text;not null"                json:"message"`
    StackTrace string         `gorm:"type:text"                         json:"stack_trace,omitempty"`
    Context    datatypes.JSON `gorm:"type:jsonb"                        json:"context,omitempty"`
    IPAddress  string         `gorm:"type:varchar(45)"                  json:"ip_address,omitempty"`
    CreatedAt  time.Time      `gorm:"autoCreateTime"                    json:"created_at"`
}

// domain/source.go
type Source struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"                json:"id"`
    Name      string    `gorm:"type:varchar(100);not null"              json:"name"`
    Slug      string    `gorm:"type:varchar(50);uniqueIndex;not null"   json:"slug"`
    APIKey    string    `gorm:"type:varchar(64);uniqueIndex;not null"   json:"-"`
    IsActive  bool      `gorm:"default:true"                            json:"is_active"`
    CreatedAt time.Time `gorm:"autoCreateTime"                          json:"created_at"`
}
```

---

## Database

### PostgreSQL 17

| Spec                   | Detail                                                      |
| ---------------------- | ----------------------------------------------------------- |
| **Version**            | PostgreSQL 17 (Sep 2024, latest stable)                     |
| **Why PostgreSQL 17?** | JSONB improvements, `MERGE` command GA, faster vacuum, ACID |
| **Features used**      | JSONB untuk context field, GIN indexes, partial indexes     |

**Key features di PostgreSQL 17:**

```sql
-- JSONB query (field di dalam context)
SELECT * FROM log_entries
WHERE context->>'auth_method' = 'google_oauth';

-- Composite index untuk dashboard query yang paling sering
CREATE INDEX idx_log_source_level_date
ON log_entries (source_id, level, created_at DESC);

-- GIN index untuk query JSONB context
CREATE INDEX idx_log_context_gin
ON log_entries USING GIN (context);

-- Partial index untuk query error saja (lebih kecil, lebih cepat)
CREATE INDEX idx_log_errors_only
ON log_entries (source_id, created_at DESC)
WHERE level IN ('ERROR', 'CRITICAL');
```

---

## Authentication & Security

### JWT via httpOnly Cookies

Auth menggunakan **dua token** yang disimpan sebagai **httpOnly cookie** — bukan localStorage — untuk mencegah serangan XSS.

| Spec          | Detail                                            |
| ------------- | ------------------------------------------------- |
| **Package**   | `github.com/golang-jwt/jwt/v5`                    |
| **Version**   | v5.2.1 (latest stable)                            |
| **Algorithm** | HS256                                             |
| **Storage**   | `httpOnly` cookie (tidak bisa diakses JavaScript) |
| **SameSite**  | `Strict` (CSRF protection)                        |
| **Secure**    | `true` (HTTPS only)                               |

**Dua token yang digunakan:**

| Token             | Cookie Name    | Expiry     | Tujuan                                       |
| ----------------- | -------------- | ---------- | -------------------------------------------- |
| **Access Token**  | `ulam_access`  | **24 jam** | Otentikasi setiap request API                |
| **Refresh Token** | `ulam_refresh` | **7 hari** | Generate access token baru tanpa login ulang |

```go
// domain/auth.go

type AccessClaims struct {
    AdminID  string `json:"admin_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims // exp: 24 jam
}

type RefreshClaims struct {
    AdminID string `json:"admin_id"`
    jwt.RegisteredClaims // exp: 7 hari
}
```

**Set cookie di response (Gin):**

```go
// handler/auth.go — setelah login berhasil
func setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
    c.SetCookie("ulam_access",  accessToken,  60*60*24,   "/", "", true, true) // 24 jam
    c.SetCookie("ulam_refresh", refreshToken, 60*60*24*7, "/", "", true, true) // 7 hari
}
```

**Alur token refresh:**

```text
Browser → GET /api/dashboard
  Cookie: ulam_access=<expired>
Backend → 401 Unauthorized

Browser → POST /api/auth/refresh
  Cookie: ulam_refresh=<valid, 7 hari>
Backend → Set-Cookie: ulam_access=<new, 24 jam>
        → 200 OK

Browser → retry GET /api/dashboard ✓
```

### API Key Auth (Source Ingestion)

| Spec           | Detail                                     |
| -------------- | ------------------------------------------ |
| **Format**     | `X-API-Key: ulam_<32-random-hex>`          |
| **Storage**    | Hashed bcrypt (cost 12) di tabel `sources` |
| **Generation** | `crypto/rand` 32 bytes → hex encode        |

---

## Email — `net/smtp` (Standard Library)

| Spec                 | Detail                                 |
| -------------------- | -------------------------------------- |
| **Package**          | Go standard library `net/smtp`         |
| **Provider (MVP)**   | Gmail SMTP + App Password              |
| **Provider (Scale)** | Resend atau Brevo (jika volume tinggi) |
| **Template**         | `html/template` standard library       |

```go
// Konfigurasi env
SMTP_HOST = "smtp.gmail.com"
SMTP_PORT = "587"
SMTP_USER = "your-email@gmail.com"
SMTP_PASS = "xxxx xxxx xxxx xxxx"   // App Password, bukan password akun
ALERT_EMAIL = "admin@your-domain.com"
GROQ_API_KEY = "gsk_xxxx..."        // API Key dari Groq Console
```

---

## AI Engine

### Groq API (Llama 3.3)

| Spec           | Detail                                                  |
| -------------- | ------------------------------------------------------- |
| **Model**      | `llama-3.3-70b-versatile`                               |
| **Inference**  | LPUs (Language Processing Units) by Groq                |
| **Latency**    | < 500ms per analysis                                    |
| **Why Groq?**  | Kecepatan luar biasa, biaya rendah, model open-weights  |
| **Function**   | Summarization, RCA, and Solution suggestion for logs    |

---

## DevOps & Infrastructure

### Docker + Docker Compose v2

| Spec               | Detail                                     |
| ------------------ | ------------------------------------------ |
| **API Image**      | `golang:1.26-alpine` (multi-stage build)   |
| **Frontend Image** | `node:22-alpine` → `nginx:alpine`          |
| **DB Image**       | `postgres:17-alpine`                       |
| **Compose**        | Docker Compose v2 (built-in Docker plugin) |

**Multi-stage Dockerfile (Backend):**

```dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ulam-api ./cmd/api

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/ulam-api /ulam-api
EXPOSE 8080
ENTRYPOINT ["/ulam-api"]
```

---

## Technology Decision Log

| Decision       | Chosen            | Alternatives                   | Reason                                |
| -------------- | ----------------- | ------------------------------ | ------------------------------------- |
| HTTP Framework | Gin v1.10         | Fiber v3, Echo v5              | Gin paling mature & banyak middleware |
| ORM            | GORM v2           | sqlx, pgx/v5 raw               | Speed of development, auto-migrate    |
| Database       | PostgreSQL 17     | MySQL 9, CockroachDB           | JSONB & partial indexes krusial       |
| CSS            | Tailwind v4       | CSS Modules, styled-components | Utility-first, Vite plugin native     |
| Frontend       | React 19          | Vue 3.5, SvelteKit             | Server Components, ekosistem terbesar |
| Router         | React Router v7   | TanStack Router                | File-based routing, stable            |
| State          | TanStack Query v5 | SWR, Redux RTK Query           | Server state management terbaik       |
| Async          | Goroutines native | Redis Queue, RabbitMQ          | Zero dependency, cukup untuk MVP      |
| Email          | SMTP              | Resend, SendGrid               | Zero cost untuk MVP, mudah switch     |

---

## Dependency Matrix

```text
Frontend Dependencies:
├── react@19                  → Core UI
├── react-dom@19              → DOM renderer
├── react-router-dom@7        → Client-side routing
├── @tanstack/react-query@5   → Server state management
├── axios@1.7                 → HTTP client
├── recharts@2.15             → Charts untuk overview
├── date-fns@4                → Date formatting
└── lucide-react@0.475        → Icon library

Backend Dependencies:
├── gin@v1.10                 → HTTP server
├── gorm@v1.25.12             → ORM
├── gorm/driver/postgres      → PG 17 driver
├── golang-jwt/jwt@v5.2.1     → JWT auth
├── joho/godotenv@v1.5.1      → .env loading
├── gorm/datatypes            → JSONB support
└── golang.org/x/crypto       → bcrypt hashing

DevOps:
├── Docker                    → Containerization
├── Docker Compose v2         → Local orchestration
├── postgres:17-alpine        → Database image
├── golang:1.26-alpine        → Build image
└── nginx:alpine              → Frontend serving
```

---

_Stack ini menggunakan versi paling baru yang sudah stable per Maret 2026. Dipilih untuk memaksimalkan developer velocity pada MVP dengan path upgrade yang jelas._
