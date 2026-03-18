# 🚀 MVP — Minimum Viable Product

## Unified Log & Activity Monitor (ULAM)

| Field        | Detail                   |
| ------------ | ------------------------ |
| **Version**  | MVP v1.0                 |
| **Timeline** | Sprint 1–3 (6 minggu)    |
| **Goal**     | Core logging + dashboard |

---

## MVP Philosophy

> _"Ship fast, iterate faster."_

MVP berfokus pada **4 kapabilitas inti** yang memberikan nilai langsung tanpa over-engineering:

1. **Terima log** dari source eksternal
2. **Simpan** ke database dengan aman
3. **Tampilkan** di dashboard yang bisa diakses
4. **Notifikasi** saat terjadi error kritis

---

## ✅ MVP Feature List

### 1. Ingestion API

| Sub-Feature              | Status  | Notes                 |
| ------------------------ | ------- | --------------------- |
| `POST /v1/logs` endpoint | [x] | Accept JSON payload   |
| API Key authentication   | [x] | Per-source token      |
| Request validation       | [x] | Required fields check |
| Async processing         | [x] | Respond < 100ms       |
| Background DB write      | [x] | Via goroutine         |
| PII Data Masking         | [x] | Automatic sensor sensitive keys |

**Minimal Payload yang Diterima:**

```json
{
  "source_id": "string (required)",
  "category": "SYSTEM_ERROR | USER_ACTIVITY | AUTH_EVENT",
  "level": "CRITICAL | ERROR | WARN | INFO",
  "message": "string (required)",
  "stack_trace": "string (optional)",
  "context": "object (optional, free-form JSON)"
}
```

---

### 2. Data Storage

| Sub-Feature                        | Status  | Notes                     |
| ---------------------------------- | ------- | ------------------------- |
| PostgreSQL setup                   | [x] | Docker atau cloud         |
| `log_entries` table via GORM       | [x] | Auto-migrate              |
| `sources` table                    | [x] | Simpan API key per source |
| JSONB context field                | [x] | Flexible metadata         |
| Indexing on source_id + created_at | [x] | Query performance         |

---

### 3. Admin Dashboard (React)

| Sub-Feature                  | Status  | Notes                           |
| ---------------------------- | ------- | ------------------------------- |
| Login page (hardcoded admin) | 🔲 Todo | Single admin user               |
| Overview stats card          | 🔲 Todo | Total errors, warnings          |
| Log table with pagination    | 🔲 Todo | 20 rows per page                |
| Filter by: Source            | 🔲 Todo | Dropdown                        |
| Filter by: Level             | 🔲 Todo | Dropdown                        |
| Filter by: Date range        | 🔲 Todo | Date picker                     |
| Log detail modal/page        | 🔲 Todo | Show full context + stack trace |
| Search by message            | 🔲 Todo | Simple text search              |

---

### 4. Email Notification

| Sub-Feature                     | Status  | Notes                   |
| ------------------------------- | ------- | ----------------------- |
| SMTP config (Gmail)             | [x] | Via env variables       |
| Trigger on ERROR/CRITICAL       | [x] | Inside goroutine        |
| Email HTML template             | [x] | Source + message + link |
| Throttling (5 menit/error type) | [x] | In-memory map           |
| Log Retention Worker            | [x] | Auto-cleanup > 30 days  |
| AI Insight Engine               | [x] | Manual & Auto analysis via Groq API |

---

## ❌ Explicitly NOT in MVP

| Feature                  | Reason Deferred                                |
| ------------------------ | ---------------------------------------------- |
| WebSocket real-time      | Kompleksitas tambahan, polling cukup untuk MVP |
| CSV Export               | Tidak urgent                                   |
| Slack/Telegram           | SMTP cukup untuk MVP                           |
| Multi-admin / RBAC       | Single admin cukup                             |
| Custom alert rules       | Hardcoded trigger cukup                        |

---

## API Endpoints (MVP Scope)

| Method | Path                               | Description                              | Auth    |
| ------ | ---------------------------------- | ---------------------------------------- | ------- |
| `POST` | `/api/ingest`                   | Kirim log baru                           | API Key |
| `GET`  | `/api/logs`                     | List log dengan filter                   | JWT     |
| `GET`  | `/api/logs/:id`                 | Detail log by ID                         | JWT     |
| `GET`  | `/api/sources`                  | List semua source terdaftar              | JWT     |
| `POST` | `/api/sources`                  | Daftarkan source baru + generate API key | JWT     |
| `GET`  | `/api/sources/:id`              | Detail source                            | JWT     |
| `POST` | `/api/sources/:id/rotate-key`   | Rotate API key                           | JWT     |
| `POST` | `/api/auth/login`               | Admin login                              | Public  |

---

## Sprint Plan

### Sprint 1 (Minggu 1-2): Backend Foundation

**Goal**: API bisa menerima dan menyimpan log

- [x] Setup project Golang (Gin/Fiber)
- [x] Connect PostgreSQL dengan GORM
- [x] Auto-migrate schema `log_entries` & `sources`
- [x] Implementasi `POST /v1/logs` dengan token auth
- [x] Goroutine untuk async DB write
- [x] Unit test untuk ingestion endpoint

**Deliverable**: `curl -X POST /v1/logs` berhasil menyimpan ke DB

---

### Sprint 2 (Minggu 3-4): Dashboard & Auth

**Goal**: Dashboard bisa dilihat dan dipakai

- [x] Setup React (Vite) + React Router
- [x] Admin login page + JWT session
- [x] API endpoint `GET /v1/logs` dengan pagination & filter
- [x] Log table component dengan filter UI
- [x] Log detail modal dengan JSON viewer
- [x] Overview stats (total per level per source)

**Deliverable**: Dashboard live dan bisa filter log

---

### Sprint 3 (Minggu 5-6): Notification & Polish

**Goal**: Sistem berjalan end-to-end dengan notifikasi

- [x] SMTP email integration
- [x] Email template HTML
- [x] Throttling logic (in-memory)
- [x] `POST /api/sources` untuk manage source API keys
- [x] `POST /api/sources/:id/rotate-key` untuk rotate API keys
- [x] Error handling & logging di backend sendiri
- [ ] Deployment: Docker Compose (API + DB)
- [ ] README dokumentasi

**Deliverable**: System berjalan end-to-end di production-like environment

---

## Definition of Done (MVP)

MVP dianggap **selesai** ketika:

1. ✅ Endpoint `/api/ingest` bisa menerima log dari source eksternal
2. ✅ Log tersimpan di PostgreSQL dan muncul di dashboard
3. ✅ Email terkirim dalam < 30 detik untuk level ERROR/CRITICAL
4. ✅ Dashboard bisa filter dan search log
5. ✅ API response time < 100ms (diukur dengan load test sederhana)

---

## Risk & Mitigation

| Risk                                | Probability | Impact | Mitigation                                                                |
| ----------------------------------- | ----------- | ------ | ------------------------------------------------------------------------- |
| SMTP rate limit (Gmail)             | Medium      | Medium | Gunakan App Password, switch ke Resend/SendGrid jika perlu                |
| DB performance dengan volume tinggi | Low         | High   | Index yang tepat + koneksi pooling                                        |
| Token bocor dari source client      | Medium      | High   | Dokumentasikan cara rotate API key via `/api/sources/:slug/rotate-key` |
