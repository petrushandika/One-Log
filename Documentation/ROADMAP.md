# 🗺️ Development Roadmap

## Unified Log & Activity Monitor (ULAM)

> Roadmap ini bersifat **feature-driven**, bukan time-driven.
> Setiap fase dikerjakan ketika fase sebelumnya sudah stabil dan production-ready.
> Tidak ada tanggal target yang ketat — prioritas berubah sesuai kebutuhan nyata.

---

## ✅ Fase 1 — Core (MVP)

**Status**: Stable

Fondasi sistem: terima log, simpan, tampilkan, dan notifikasi.

### Backend

- [x] Setup project Golang + Gin
- [x] PostgreSQL + GORM AutoMigrate
- [x] `POST /api/ingest` — ingestion endpoint + API Key auth
- [x] Async goroutine untuk DB write (response < 100ms)
- [x] Validasi payload (category, level enum)
- [x] JWT Auth untuk admin dashboard
- [x] `GET /api/logs` + filter + pagination
- [x] `GET /api/logs/:id`
- [x] `GET/POST /api/sources` — Source management
- [x] `POST /api/sources/:slug/rotate-key`
- [x] `GET /api/stats/overview`

### Notification

- [x] SMTP email integration
- [x] HTML email template
- [x] Throttle in-memory (5 menit per error type)
- [x] **AI Insight Integration (Groq API)** — Auto analysis for CRITICAL logs

### Frontend

- [x] Setup React 19 + Vite 7 + Tailwind v4 (Feature-Based structure)
- [x] Login page + JWT session
- [x] Overview dashboard — stats cards
- [x] Log Table + filter (source, level, category, date range)
- [x] Log Detail modal — JSON viewer + stack trace
- [x] **AI "Analyze" Button** — Manual trigger on log selection
- [x] Sources Management page

### Infra

- [ ] Docker Compose (API + DB + frontend)
- [ ] README + dokumentasi deployment

**Done ketika:**

- Ingestion API bisa menerima log dari source luar ✅
- Dashboard bisa filter & search log ✅
- Email terkirim untuk ERROR/CRITICAL dalam < 30 detik ✅

---

## ✅ Fase 2 — Activity Monitor & Audit Trail

**Status**: Implemented

Mencatat **siapa melakukan apa, dari mana, kapan** — bukan hanya error, tapi seluruh jejak aktivitas pengguna dan riwayat perubahan data yang tidak bisa diubah.

### Auth Event Tracking

- [x] **Standardized AUTH_EVENT payload** — kontrak field untuk semua auth methods
- [ ] **Auth Method Dashboard** — breakdown pie chart: Google vs Manual vs GitHub vs lainnya
- [ ] **Login Timeline** — visual timeline login events per source per hari
- [ ] **Failed Login Heatmap** — visualisasi jam-jam dengan high failed login rate
- [ ] **Recent Sessions Table** — user, auth_method, IP, browser, device, timestamp, source
- [x] **Brute Force Detection** — alert jika `login_failed` > threshold dari satu IP dalam 10 menit

**Auth methods yang didukung:**

| Method            | Notes                      |
| ----------------- | -------------------------- |
| `google_oauth`    | Login via Google OAuth 2.0 |
| `github_oauth`    | Login via GitHub           |
| `facebook_oauth`  | Login via Facebook         |
| `twitter_oauth`   | Login via Twitter/X        |
| `discord_oauth`   | Login via Discord          |
| `system_password` | Login manual internal      |
| `magic_link`      | Passwordless email link    |
| `sso`             | Enterprise SSO / SAML      |

### Audit Trail — Immutable Logs

Catat setiap perubahan data penting di aplikasi client (CMS, Absensi, dll.) sebagai bukti otentik yang **tidak bisa diedit**.

- [x] **User Activity Log** — catat setiap kali admin mengubah/menghapus data (di CMS, absensi, dll.)
- [x] **Immutable Log Flag** — log dengan `category: AUDIT_TRAIL` tidak bisa dihapus via API, hanya via DB migration
- [x] **IP & Device Tracking** — field `ip_address` tersimpan pada setiap log entry
- [ ] **Before/After Diff** — context menyimpan nilai sebelum dan sesudah perubahan (future)
- [x] **Audit Trail Page** — halaman khusus di dashboard untuk audit, filter & pagination aktif
- [ ] **Compliance Export** — export audit trail ke PDF/CSV untuk keperluan audit eksternal

**Contoh payload Audit Trail:**

```json
{
  "category": "AUDIT_TRAIL",
  "level": "INFO",
  "message": "Admin deleted attendance record",
  "context": {
    "action": "delete",
    "actor_id": "admin_001",
    "actor_role": "admin",
    "resource_type": "attendance",
    "resource_id": "att_march_2026_usr_123",
    "before": { "status": "present", "check_in": "08:02" },
    "after": null,
    "ip_address": "192.168.1.5",
    "device_type": "desktop",
    "reason": "Data entry error"
  }
}
```

### User Activity Trail

- [ ] **Activity Feed** — feed aktivitas user (page_view, create, update, delete, export)
- [ ] **User Profile View** — click user_id → lihat semua aktivitas user lintas semua source
- [ ] **Top Active Users** — siapa user yang paling aktif per source
- [ ] **Activity by Resource** — berapa kali resource tertentu diakses/diubah

### API Enhancements

- [x] `GET /api/activity` — filter activity logs (category=AUTH_EVENT,USER_ACTIVITY,AUDIT_TRAIL)
- [x] `GET /api/activity/summary` — agregat login count by auth_method & event_type
- [x] `GET /api/activity/users/:user_id` — semua aktivitas satu user ID lintas semua source
- [x] `GET /api/activity/suspicious` — login mencurigakan & anomali

---

## ✅ Fase 3 — Performance Monitoring (APM)

**Status**: Implemented

Jangan hanya catat **kapan** error terjadi — catat juga **seberapa lambat** sistem berjalan.

### Response Time Tracking

- [x] **Endpoint Latency Log** — catat waktu response setiap API di aplikasi client dengan `category: PERFORMANCE`
- [x] **P50 / P95 / P99 Stats** — agregasi persentil untuk setiap endpoint via `GET /api/apm/endpoints`
- [ ] **Threshold Alert** — kirim email jika rata-rata response time > X ms (konfigurabel per source)
- [ ] **Response Time Timeline** — Line chart di dashboard untuk melihat tren dari waktu ke waktu per endpoint

### Slow Query Detector

- [ ] **Slow Query Log** — aplikasi client kirim log dengan `category: PERFORMANCE` saat query > threshold
- [ ] **Payload standar slow query:**

```json
{
  "category": "PERFORMANCE",
  "level": "WARN",
  "message": "Slow query detected",
  "context": {
    "query_type": "SELECT",
    "table": "log_entries",
    "duration_ms": 3200,
    "threshold_ms": 2000,
    "query_preview": "SELECT * FROM log_entries WHERE...",
    "endpoint": "/api/attendance/report",
    "user_id": "admin_001"
  }
}
```

- [ ] **Slow Query Table** — tabel di dashboard yang menampilkan query terlambat per source
- [ ] **Query Trend Chart** — grafik frekuensi slow query per hari per source

### APM Dashboard

- [x] **APM Overview Page** — halaman khusus performance di frontend
- [x] **Per-Source APM** — pilih source lalu lihat latency trend per endpoint
- [ ] **Apdex Score** (future) — skor kepuasan performa agregat berdasarkan threshold

---

## ✅ Fase 4 — Status Page & Uptime Monitoring

**Status**: Implemented

Mirip [uptime.com](https://uptime.com) / [betterstack.com](https://betterstack.com) — tapi terintegrasi langsung dengan data log ULAM.

### Health Check Worker

- [x] **Health Check Goroutine** — background worker di ULAM yang ping URL terdaftar setiap N menit (default: 5 menit)
- [x] **Per-Source Health Endpoint** — tiap source bisa mendaftarkan URL health check (misal `https://absensi.app/health`)
- [x] **Status Enum**: `ONLINE`, `DEGRADED`, `OFFLINE`, `MAINTENANCE`
- [x] **Downtime Detection** — jika ping gagal 3x berturut-turut → status `OFFLINE` + kirim email alert
- [x] **Auto Log on Down** — sistem otomatis kirim log `level: CRITICAL`, `category: SYSTEM_ERROR` saat source offline

```go
// worker/health-check-worker.go
func (w *HealthCheckWorker) RunLoop() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        for _, source := range w.sources {
            go w.ping(source)
        }
    }
}
```

### Incident Management

- [ ] **Incident Auto-Create** — saat source offline, sistem buat incident record otomatis
- [ ] **Incident Timeline** — kapan down, kapan recover, berapa lama downtime
- [x] **Email: "Server Down"** — notifikasi instan saat down dengan estimasi dampak
- [ ] **Email: "Server Recovered"** — notifikasi saat sistem kembali online

### Public / Internal Status Page

- [x] **Public Status Page** — `GET /api/status` endpoint tersedia tanpa autentikasi
- [x] **Status Page UI** — halaman `/status` di dashboard: source cards dengan badge ONLINE/DEGRADED/OFFLINE/MAINTENANCE, auto-refresh 60 detik
- [ ] **URL**: `https://status.ulam.your-domain.com` atau `https://ulam.your-domain.com/status`
- [ ] **Tampilan publik**: uptime % 30 hari terakhir, incident history
- [ ] **Embed Widget** — badge status yang bisa di-embed di README atau halaman lain

---

## ✅ Fase 5 — Error Grouping & Smart Analysis

**Status**: Implemented

Agar tidak pusing membaca ribuan baris log yang sama — kelompokkan dan analisis secara otomatis.

### Error Grouping

- [x] **Auto-Grouping** — log dengan message dan stack trace yang mirip dikelompokkan menjadi satu "Issue" (Via Fingerprinting)
- [x] **Fingerprinting Algorithm** — hash dari `(source_id + normalized_message + stack_trace[:100])` sebagai group key
- [x] **Issue Tracker** — halaman "Issues" di dashboard: daftar group error dengan jumlah occurrence & first/last seen
- [x] **Issue Detail** — lihat semua individual log dalam satu group via modal
- [x] **Issue Status** — tandai issue sebagai `OPEN`, `RESOLVED`, `IGNORED`
- [ ] **Regression Detection** — alert jika issue yang sudah `RESOLVED` muncul lagi

### Error Analytics

- [x] **"Source mana yang paling sering error?"** — CSS progress bar ranking occurrence count per source (Issues › Analytics tab)
- [ ] **"Jam berapa error paling sering terjadi?"** — heatmap per jam dalam seminggu <!-- backlog -->
- [x] **"Error apa yang paling sering muncul?"** — top 10 error messages (Issues › Analytics tab)
- [x] **Level breakdown** — badge count CRITICAL/ERROR/WARN untuk open issues (Issues › Analytics tab)
- [ ] **Error Rate Trend** — persentase log yang error vs total per hari <!-- backlog -->

### AI Integration

- [x] **AI Copilot (Chatbot)** — Floating chat widget (purple, `POST /api/chat`) dengan Groq llama-3.3-70b; system prompt project-aware dengan live stats (total logs, ERROR/CRITICAL count, open issues)
- [ ] **AI Daily Digest** — LLM merangkum error/anomali hari ini dalam bahasa natural (dikirim via email) <!-- backlog -->
- [ ] **AI Error Deduplication** — gunakan embedding untuk mengelompokkan error yang semantically similar (bukan hanya exact match) <!-- backlog -->
- [ ] **Prompt Context** — sertakan framework, bahasa, dan error history agar AI suggestion lebih relevan <!-- backlog -->

```text
Contoh AI Suggestion output:
---
Error: "connection refused at postgres:5432"
Stack: main.connectDB() at db.go:45

💡 Analisis: Database connection pool habis atau PostgreSQL tidak berjalan.
   Penyebab umum:
   1. Max connections di postgresql.conf terlalu rendah
   2. Goroutine leak menyebabkan koneksi tidak dikembalikan ke pool
   Saran: Periksa `db.SetMaxOpenConns()` dan pastikan setiap transaksi ditutup dengan `defer tx.Rollback()`
---
```

---

## ✅ Fase 6 — Centralized Configuration Management

**Status**: Implemented

Kelola semua konfigurasi aplikasi-aplikasi Anda dari satu tempat — tanpa harus SSH ke setiap server.

### Config Storage

- [x] **Config Table** — tabel `source_configs` di PostgreSQL: `source_id`, `key`, `value`, `is_secret`, `updated_at`
- [x] **Secret Management** — nilai sensitif dienkripsi dengan AES-256 (NaCl secretbox) sebelum disimpan
- [x] **Versioning** — setiap perubahan config tersimpan di `source_config_histories`
- [ ] **Environment Namespacing** — config bisa punya namespace: `production`, `staging`, `development`

### Config API

- [x] `GET /api/config/:source_slug` — aplikasi client pull config terbaru saat startup
- [x] `PUT /api/config/:source_slug/:key` — update config value via dashboard
- [x] `GET /api/config/:source_slug/history` — riwayat perubahan config

### Hot Reload

- [ ] **Hot Reload Config** — aplikasi client bisa polling `/api/config/:slug` setiap N detik untuk reload config tanpa restart
- [ ] **Change Notification** — saat config diubah, ULAM kirim notifikasi (webhook atau SSE) ke aplikasi client
- [ ] **Go SDK Config Helper:**

```go
// pkg/ulam/config.go
func WatchConfig(sourceSlug string, interval time.Duration, onUpdate func(map[string]string)) {
    ticker := time.NewTicker(interval)
    for range ticker.C {
        cfg, _ := fetchConfig(sourceSlug)
        onUpdate(cfg)
    }
}
```

### Dashboard Config Management

- [x] **Config Editor UI** — halaman `/config` di dashboard, edit config per source dengan slide-over panel
- [x] **Secret Toggle** — nilai secret ditampilkan sebagai `••••••••`, bisa "reveal" dengan klik Eye icon
- [x] **Rollback** — rollback config ke versi sebelumnya dari tab History dengan satu klik
- [ ] **Config Audit Trail** — siapa yang mengubah config apa, kapan (terintegrasi dengan Fase 2) <!-- backlog -->

---

## 🔄 Fase 7 — Export & Extended Integrations

**Status**: Partially Implemented (CSV export + webhook done; Slack/S3/SDK pending)

Export data dan notifikasi ke platform lain.

- [x] **CSV Export** — `GET /api/logs/export` — export hasil filter log ke CSV
- [ ] **Excel Export** — export ke .xlsx dengan formatting
- [ ] **PDF Audit Report** — export audit trail ke PDF untuk compliance
- [x] **Email Notification** — alert ke email admin (with 5-min throttle)
- [x] **Telegram Bot** — notifikasi ke Telegram group/channel
- [x] **Webhook Support** — generic outgoing webhook ke URL apapun (via `WEBHOOK_URL` env)
- [x] **Rate Limiting** — Token bucket algorithm untuk semua endpoints
- [ ] **Log Archiving** — compress dan archive log lama ke S3/object storage (interface ready)
- [x] **Official Go SDK** — Go SDK untuk integrasi lebih mudah

### Third-Party Audit Integrations (SaaS Connectors)

- [ ] **Google Workspace Tracker** — Background cron-job untuk menarik data dari *Google Admin SDK Reports API* ke ULAM
- [ ] **AWS CloudTrail Sync** — Menarik log IAM dan aktivitas *bucket* S3 Amazon ke dalam Dashboard One-Log
- [ ] **GitHub Org Audit** — Menarik log pendaftaran, penghapusan *repository*, dan perubahan kunci akses dari akun GitHub Enterprise perusahaan
---

## 📊 Feature Priority Matrix

| Feature                              | Value      | Effort     | Priority      |
| ------------------------------------ | ---------- | ---------- | ------------- |
| Ingestion API                        | High       | Low        | 🔴 Must       |
| Email notification                   | High       | Low        | 🔴 Must       |
| Log dashboard + filter               | High       | Medium     | 🔴 Must       |
| **Auth Event Tracking**              | **High**   | **Medium** | **🔴 Must**   |
| Source Management                    | High       | Low        | 🔴 Must       |
| Activity Monitor                     | High       | Medium     | 🟠 High       |
| Audit Trail (immutable)              | High       | Medium     | 🟠 High       |
| Failed login + brute force detection | High       | Low        | 🟠 High       |
| **APM — Response Time Tracking**     | **High**   | **Medium** | **🟠 High**   |
| **APM — Slow Query Detector**        | **High**   | **Low**    | **🟠 High**   |
| **Status Page & Uptime Monitor**     | **High**   | **Medium** | **🟠 High**   |
| Log retention + auto-delete          | Medium     | Low        | 🟡 Should     |
| **Error Grouping (Issues)**          | **High**   | **High**   | **🟡 Should** |
| **Error Analytics**                  | **Medium** | **Medium** | **🟡 Should** |
| CSV/Excel export                     | Medium     | Medium     | 🟡 Should     |
| Slack/Telegram notification          | Medium     | Low        | 🟡 Should     |
| Webhook support                      | Medium     | Low        | 🟡 Should     |
| **AI Stack Trace Analysis**          | **High**   | **High**   | **🟢 Could**  |
| **AI Daily Digest**                  | **Medium** | **High**   | **🟢 Could**  |
| **Centralized Config Management**    | **High**   | **High**   | **🟢 Could**  |
| Geo-IP Login Map                     | Low        | Medium     | 🟢 Could      |
| Multi-admin + RBAC                   | Low        | High       | 🟢 Could      |

---

## 🔧 Technical Debt

| Item                     | Description                            | Target Phase |
| ------------------------ | -------------------------------------- | ------------ |
| In-memory email throttle | Hilang saat restart                    | Fase 3       |
| Hardcoded single admin   | Tidak scalable                         | Fase 7       |
| No log compression       | Disk tumbuh terus                      | Fase 7       |
| No request ID tracing    | Susah debug lintas service             | Fase 3       |
| Auth_method parsing      | Tidak ada validasi standar saat ingest | Fase 2       |
| No health check          | Downtime tidak terdeteksi otomatis     | Fase 4       |
| Error tidak digroup      | Ribuan log duplikat tidak ada konteks  | Fase 5       |
