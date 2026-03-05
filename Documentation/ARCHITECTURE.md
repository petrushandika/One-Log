# 🏛️ System Architecture

## Unified Log & Activity Monitor (ULAM)

---

## High-Level Architecture

ULAM memisahkan tanggung jawab ke dalam layer yang jelas. Sisi kiri adalah **producer** (semua aplikasi yang mengirim log — jumlahnya tidak terbatas), sisi kanan adalah **consumer** (admin yang memonitor via dashboard).

```text
╔══════════════════════════════════════════════════════════════════════════╗
║                         PRODUCER LAYER                                   ║
║  (Semua source yang didaftarkan — tidak ada batasan jumlah)              ║
║                                                                           ║
║  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐   ║
║  │    App A     │ │    App B     │ │    App C     │ │    App N     │   ║
║  │ (any lang)   │ │ (any lang)   │ │ (any lang)   │ │ (any lang)   │   ║
║  └──────┬───────┘ └──────┬───────┘ └──────┬───────┘ └──────┬───────┘   ║
╚═════════╪════════════════╪════════════════╪════════════════╪════════════╝
          │                │                │                │
          │  POST /api/ingest               │                │
          │  Header: X-API-Key: ulam_xxx    │                │
          └─────────────────┬───────────────┘                │
                            │◄───────────────────────────────┘
                            │
╔═══════════════════════════▼══════════════════════════════════════════════╗
║                        BACKEND LAYER (Golang)                            ║
║                                                                           ║
║  ┌───────────────────────────────────────────────────────────────────┐   ║
║  │                    HTTP Router (Gin)                               │   ║
║  │  /api/ingest   →  IngestHandler                                   │   ║
║  │  /api/logs     →  LogHandler                                      │   ║
║  │  /api/sources  →  SourceHandler                                   │   ║
║  │  /api/stats    →  StatsHandler                                    │   ║
║  │  /api/auth     →  AuthHandler                                     │   ║
║  └─────────────────────────────┬─────────────────────────────────────┘   ║
║                                │                                          ║
║  ┌─────────────────────────────▼─────────────────────────────────────┐   ║
║  │                     Service Layer                                  │   ║
║  │  LogService · SourceService · NotificationService · AuthService   │   ║
║  └─────────────────────────────┬─────────────────────────────────────┘   ║
║                                │                                          ║
║  ┌─────────────────────────────▼─────────────────────────────────────┐   ║
║  │                   Repository Layer                                 │   ║
║  │  LogRepository · SourceRepository · AdminRepository               │   ║
║  └─────────────────────────────┬─────────────────────────────────────┘   ║
╚═══════════════════════════════╪══════════════════════════════════════════╝
          ┌──────────────────────┤
          │                      │
          ▼                      ▼ (goroutine)
╔════════════════╗     ╔══════════════════════╗
║  PostgreSQL    ║     ║  Notification Worker  ║
║  (GORM)        ║     ║  - Throttle Check     ║
║  log_entries   ║     ║  - SMTP Send          ║
║  sources       ║     ╚══════════════════════╝
╚════════════════╝
          ▲
╔═════════╪════════════════════════════════════════════════════════════════╗
║         │                FRONTEND LAYER (React + TypeScript)             ║
║  ┌──────┴────────────────────────────────────────────────────────────┐   ║
║  │          React Dashboard (Vite 7 + Tailwind v4 + TypeScript)      │   ║
║  │                                                                   │   ║
║  │  Feature: Auth    │ Feature: Logs  │ Feature: Sources │ Feature:  │   ║
║  │  - Login page     │ - Log list     │ - Source list    │ Stats     │   ║
║  │  - Cookie store   │ - Filters      │ - Register       │ - Charts  │   ║
║  │  - Token refresh  │ - Detail view  │ - Rotate key     │ - Overview│   ║
║  └───────────────────────────────────────────────────────────────────┘   ║
╚══════════════════════════════════════════════════════════════════════════╝
```

---

## Backend Architecture — Clean Architecture

Backend menggunakan **Clean Architecture** (Dependency Rule: outer → inner, never reverse):

```text
cmd/
└── api/
    └── main.go          ← Entry point

internal/
├── handler/             ← Layer 1: HTTP (paling luar)
│   ├── ingest.go
│   ├── log.go
│   ├── source.go
│   ├── stats.go
│   └── auth.go
│
├── service/             ← Layer 2: Business Logic
│   ├── log-service.go
│   ├── app-service.go
│   ├── auth-service.go
│   └── notification-service.go
│
├── repository/          ← Layer 3: Data Access
│   ├── log-repo.go
│   ├── app-repo.go
│   └── admin-repo.go
│
├── domain/              ← Layer 4: Core (paling dalam, no dependencies)
│   ├── log.go           ← LogEntry struct + interfaces
│   ├── app.go           ← Source struct + interfaces
│   └── errors.go        ← Domain error types
│
├── middleware/          ← Cross-cutting concerns
│   ├── api-key.go       ← Validasi X-API-Key
│   ├── jwt-auth.go      ← Validasi JWT
│   └── rate-limiter.go  ← Rate limiting per key/IP
│
├── worker/              ← Background processes
│   └── email-worker.go  ← Goroutine email dispatcher
│
├── infra/               ← Infrastructure concerns
│   ├── db/
│   │   └── postgres.go  ← GORM connection & migration
│   ├── smtp/
│   │   └── mailer.go    ← SMTP client wrapper
│   └── config/
│       └── config.go    ← Env loader (godotenv)
│
└── router/
    └── router.go        ← Route registration
```

### Dependency Rule

```text
handler → service → repository → domain
   ↑           ↑          ↑
(can use)  (can use)  (can use)

domain ← TIDAK boleh bergantung ke layer manapun
```

### Kenapa Clean Architecture?

| Benefit              | Detail                                                     |
| -------------------- | ---------------------------------------------------------- |
| **Testable**         | Service layer bisa di-unit-test tanpa DB (mock repository) |
| **Maintainable**     | Perubahan di handler tidak mempengaruhi business logic     |
| **Swappable**        | Ganti PostgreSQL ke MySQL? Cukup ubah repository layer     |
| **Clear boundaries** | Setiap file punya satu tanggung jawab yang jelas           |

---

### Contoh: Flow Data dari Request ke Response

```go
// handler/ingest.go — hanya tahu HTTP
func (h *IngestHandler) Handle(c *gin.Context) {
    var payload domain.LogPayload
    c.ShouldBindJSON(&payload)

    // Delegasi ke service, tidak tahu detail DB
    h.logService.Ingest(c.Request.Context(), payload, appID)

    c.JSON(202, gin.H{"status": "accepted"})
}

// service/log-service.go — hanya tahu business rules
func (s *LogService) Ingest(ctx context.Context, p domain.LogPayload, appID string) {
    entry := domain.NewLogEntry(p, appID)

    go func() {
        s.repo.Save(ctx, entry)           // Simpan ke DB
        if entry.NeedsAlert() {           // Business rule: level ERROR/CRITICAL
            s.notifier.Dispatch(entry)    // Kirim email
        }
    }()
}

// repository/log-repo.go — hanya tahu SQL / GORM
func (r *LogRepository) Save(ctx context.Context, e *domain.LogEntry) error {
    return r.db.WithContext(ctx).Create(e).Error
}
```

---

## Frontend Architecture — Feature-Based Structure

Frontend menggunakan **Feature-Based Folder Structure** (bukan atomic design), di mana setiap fitur adalah unit mandiri:

```text
src/
├── features/                    ← Inti: tiap folder = 1 fitur bisnis
│   │
│   ├── auth/                    ← Fitur: Autentikasi
│   │   ├── components/
│   │   │   └── LoginForm.tsx
│   │   ├── hooks/
│   │   │   └── use-auth.ts       ← useLogin, useLogout, useToken
│   │   ├── services/
│   │   │   └── auth-api.ts       ← POST /auth/login, POST /auth/refresh
│   │   └── store/
│   │       └── auth-store.ts     ← Token, user state (Zustand / Context)
│   │
│   ├── logs/                    ← Fitur: Log Management
│   │   ├── components/
│   │   │   ├── LogTable.tsx     ← Tabel utama dengan pagination
│   │   │   ├── LogFilters.tsx   ← Filter bar (app, level, date, search)
│   │   │   ├── LogDetail.tsx    ← Modal/drawer detail log
│   │   │   ├── LogBadge.tsx     ← Level badge (CRITICAL/ERROR/...)
│   │   │   └── StackTrace.tsx   ← Formatted stack trace viewer
│   │   ├── hooks/
│   │   │   ├── use-logs.ts       ← Fetch + filter logic (React Query)
│   │   │   └── use-log-detail.ts
│   │   └── services/
│   │       └── logs-api.ts       ← GET /logs, GET /logs/:id
│   │
│   ├── sources/                 ← Fitur: Registered Source Management
│   │   ├── components/
│   │   │   ├── SourceList.tsx      ← Daftar semua source terdaftar
│   │   │   ├── SourceCard.tsx      ← Card per source + stats ringkasan
│   │   │   ├── RegisterSourceForm.tsx
│   │   │   └── ApiKeyDisplay.tsx   ← One-time display dengan copy button
│   │   ├── hooks/
│   │   │   └── use-sources.ts
│   │   └── services/
│   │       └── sources-api.ts      ← GET/POST/PATCH /sources, POST /sources/:slug/rotate-key
│   │
│   └── stats/                   ← Fitur: Overview & Statistics
│       ├── components/
│       │   ├── OverviewCards.tsx ← Total logs, errors, etc.
│       │   ├── LogTrendChart.tsx ← Line chart (Recharts)
│       │   └── AppBreakdown.tsx  ← Bar chart per app
│       ├── hooks/
│       │   └── use-stats.ts
│       └── services/
│           └── stats-api.ts      ← GET /stats/overview, GET /stats/sources/:slug
│
├── shared/                      ← Komponen & utils yang dipakai banyak fitur
│   ├── components/
│   │   ├── ui/
│   │   │   ├── Button.tsx
│   │   │   ├── Card.tsx
│   │   │   ├── Badge.tsx
│   │   │   ├── Modal.tsx
│   │   │   ├── Table.tsx
│   │   │   ├── Pagination.tsx
│   │   │   └── Spinner.tsx
│   │   ├── Layout.tsx           ← Sidebar + topbar wrapper
│   │   └── ErrorBoundary.tsx
│   ├── hooks/
│   │   └── use-debounce.ts
│   ├── lib/
│   │   ├── axios.ts             ← Axios instance + interceptors (JWT attach)
│   │   └── query-client.ts       ← React Query global config
│   └── utils/
│       ├── format-date.ts
│       ├── truncate.ts
│       └── level-color.ts        ← Map level → Tailwind color class
│
├── pages/                       ← Halaman (thin layer, assembles features)
│   ├── LoginPage.tsx
│   ├── OverviewPage.tsx
│   ├── LogsPage.tsx
│   ├── LogDetailPage.tsx
│   └── SourcesPage.tsx
│
├── router/
│   └── index.tsx                ← React Router v7 routes + auth guard
│
├── App.tsx
└── main.tsx
```

### Kenapa Feature-Based?

| Benefit              | Detail                                                            |
| -------------------- | ----------------------------------------------------------------- |
| **Koherensi**        | Semua yang berhubungan dengan `logs` ada di satu folder           |
| **Scalable**         | Tambah fitur baru? Buat folder fitur baru, tidak ganggu yang lain |
| **Onboarding mudah** | Developer baru mudah menemukan kode fitur yang ingin diubah       |
| **Lazy loading**     | Setiap fitur bisa di-code-split dengan mudah (React.lazy)         |

### Contoh: Menambah Fitur Baru

Jika di masa depan ingin tambah fitur **Notifications** (history email yang terkirim):

```text
src/features/notifications/
├── components/
│   ├── NotificationList.tsx
│   └── NotificationItem.tsx
├── hooks/
│   └── use-notifications.ts
└── services/
    └── notifications-api.ts
```

Tidak ada perubahan di fitur lain. ✅

---

## Request Lifecycle (End-to-End)

```text
[Browser / App Client]
         │
         │ POST /api/ingest
         ▼
[Nginx / Reverse Proxy]   ← TLS termination, rate limit per IP
         │
         ▼
[Golang API — Gin Router]
         │
         ├── Middleware: Rate Limiter (per API key)
         ├── Middleware: API Key Auth → lookup DB → attach app to context
         ├── Middleware: Request Logger
         │
         ▼
[IngestHandler]
         │ validate payload
         │ extract appID from context
         │
         ▼
[LogService.Ingest()]
         │
         ├─── respond 202 immediately ──────────────────► [Client]
         │
         └─── go func() {
                  LogRepository.Save()   ─────────────► [PostgreSQL]
                  if needsAlert {
                      EmailWorker.Dispatch()  ────────► [SMTP]
                  }
              }()
```

---

## Deployment Architecture

```yaml
# docker-compose.yml

services:
  nginx:
    image: nginx:alpine
    ports: ["80:80", "443:443"]
    depends_on: [api, frontend]

  api:
    build:
      context: ./backend
      dockerfile: Dockerfile
    environment:
      - DATABASE_URL
      - JWT_SECRET
      - SMTP_HOST
      - SMTP_USER
      - SMTP_PASS
      - ALERT_EMAIL
    depends_on: [postgres]
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    environment:
      - VITE_API_URL

  postgres:
    image: postgres:16-alpine
    volumes:
      - pgdata:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB
      - POSTGRES_USER
      - POSTGRES_PASSWORD
    restart: unless-stopped

volumes:
  pgdata:
```

---

## Environment Variables

| Variable       | Layer    | Description                              |
| -------------- | -------- | ---------------------------------------- |
| `DATABASE_URL` | Backend  | `postgres://user:pass@host:5432/dbname`  |
| `JWT_SECRET`   | Backend  | Random 32+ char string untuk signing JWT |
| `SMTP_HOST`    | Backend  | `smtp.gmail.com`                         |
| `SMTP_PORT`    | Backend  | `587`                                    |
| `SMTP_USER`    | Backend  | Email pengirim                           |
| `SMTP_PASS`    | Backend  | App Password Gmail                       |
| `ALERT_EMAIL`  | Backend  | Email penerima notifikasi (admin)        |
| `SERVER_PORT`  | Backend  | Default `8080`                           |
| `VITE_API_URL` | Frontend | `https://api.ulam.your-domain.com`       |
