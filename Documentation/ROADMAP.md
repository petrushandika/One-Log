# 🗺️ Development Roadmap

## Unified Log & Activity Monitor (ULAM)

> Roadmap ini bersifat **feature-driven**, bukan time-driven.
> Setiap fase dikerjakan ketika fase sebelumnya sudah stabil dan production-ready.
> Tidak ada tanggal target yang ketat — prioritas berubah sesuai kebutuhan nyata.

---

## ✅ Fase 1 — Core (MVP)

**Status**: ✅ Fully Implemented

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

- [x] Docker Compose (API + DB + frontend)
- [x] README + dokumentasi deployment

---

## ✅ Fase 2 — Activity Monitor & Audit Trail

**Status**: ✅ Fully Implemented

Mencatat **siapa melakukan apa, dari mana, kapan** — bukan hanya error, tapi seluruh jejak aktivitas pengguna dan riwayat perubahan data yang tidak bisa diubah.

### Auth Event Tracking

- [x] **Standardized AUTH_EVENT payload** — kontrak field untuk semua auth methods
- [x] **Auth Method Dashboard** — breakdown pie chart: Google vs Manual vs GitHub vs lainnya
- [x] **Login Timeline** — visual timeline login events per source per hari
- [x] **Failed Login Heatmap** — visualisasi jam-jam dengan high failed login rate
- [x] **Recent Sessions Table** — user, auth_method, IP, browser, device, timestamp, source
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

- [x] **User Activity Log** — catat setiap kali admin mengubah/menghapus data
- [x] **Immutable Log Flag** — log dengan `category: AUDIT_TRAIL` tidak bisa dihapus via API
- [x] **IP & Device Tracking** — field `ip_address` tersimpan pada setiap log entry
- [x] **Before/After Diff** — context menyimpan nilai sebelum dan sesudah perubahan
- [x] **Audit Trail Page** — halaman khusus di dashboard untuk audit
- [x] **Compliance Export** — export audit trail ke PDF/CSV

### Activity Monitor Extended

- [x] **Activity Feed** — feed aktivitas user (page_view, create, update, delete, export)
- [x] **User Profile View** — click user_id → lihat semua aktivitas user lintas semua source
- [x] **Top Active Users** — siapa user yang paling aktif per source
- [x] **Activity by Resource** — berapa kali resource tertentu diakses/diubah

### API Endpoints

- [x] `GET /api/activity` — filter activity logs
- [x] `GET /api/activity/summary` — agregat login count by auth_method & event_type
- [x] `GET /api/activity/users/:user_id` — semua aktivitas satu user ID lintas semua source
- [x] `GET /api/activity/suspicious` — login mencurigakan & anomali
- [x] `GET /api/activity/analytics/methods` — breakdown auth methods
- [x] `GET /api/activity/analytics/timeline` — login timeline per hari
- [x] `GET /api/activity/analytics/heatmap` — failed login heatmap by hour/day
- [x] `GET /api/activity/sessions` — recent sessions table with pagination
- [x] `GET /api/activity/feed` — activity feed with pagination
- [x] `GET /api/activity/top-users` — top active users
- [x] `GET /api/activity/by-resource` — activity by resource
- [x] `POST /api/activity/compliance-export` — request compliance export
- [x] `GET /api/activity/compliance-exports` — list compliance exports

---

## ✅ Fase 3 — Performance Monitoring (APM)

**Status**: ✅ Fully Implemented

Jangan hanya catat **kapan** error terjadi — catat juga **seberapa lambat** sistem berjalan.

### Response Time Tracking

- [x] **Endpoint Latency Log** — catat waktu response setiap API dengan `category: PERFORMANCE`
- [x] **P50 / P95 / P99 Stats** — agregasi persentil untuk setiap endpoint
- [x] **Threshold Alert** — kirim email jika rata-rata response time > X ms
- [x] **Response Time Timeline** — Line chart untuk melihat tren dari waktu ke waktu

### Slow Query Detector

- [x] **Slow Query Log** — aplikasi client kirim log saat query > threshold
- [x] **Slow Query Table** — tabel di dashboard yang menampilkan query terlambat
- [x] **Query Trend Chart** — grafik frekuensi slow query per hari per source
- [x] **Apdex Score** — skor kepuasan performa agregat

### APM Dashboard

- [x] **APM Overview Page** — halaman khusus performance di frontend
- [x] **Per-Source APM** — pilih source lalu lihat latency trend per endpoint
- [x] **APM Threshold Management** — CRUD thresholds untuk alert

### API Endpoints

- [x] `GET /api/apm/endpoints` — endpoint latency stats (P50/P95/P99)
- [x] `GET /api/apm/timeline` — response time timeline data
- [x] `GET /api/apm/thresholds` — list all thresholds
- [x] `POST /api/apm/thresholds` — create new threshold
- [x] `GET /api/apm/thresholds/:id` — get threshold by ID
- [x] `PATCH /api/apm/thresholds/:id` — update threshold
- [x] `DELETE /api/apm/thresholds/:id` — delete threshold
- [x] `GET /api/apm/slow-queries` — detect slow queries > threshold
- [x] `GET /api/apm/slow-queries/trend` — slow query trend over time
- [x] `GET /api/apm/apdex` — calculate Apdex score
- [x] `GET /api/apm/threshold-alerts` — check threshold violations

---

## ✅ Fase 4 — Status Page & Uptime Monitoring

**Status**: ✅ Fully Implemented

Mirip [uptime.com](https://uptime.com) / [betterstack.com](https://betterstack.com) — tapi terintegrasi langsung dengan data log ULAM.

### Health Check Worker

- [x] **Health Check Goroutine** — background worker yang ping URL terdaftar setiap 5 menit
- [x] **Per-Source Health Endpoint** — tiap source bisa mendaftarkan URL health check
- [x] **Status Enum**: `ONLINE`, `DEGRADED`, `OFFLINE`, `MAINTENANCE`
- [x] **Downtime Detection** — jika ping gagal 3x berturut-turut → status `OFFLINE`
- [x] **Auto Log on Down** — sistem otomatis kirim log saat source offline

### Incident Management

- [x] **Incident Auto-Create** — saat source offline, sistem buat incident record otomatis
- [x] **Incident Timeline** — kapan down, kapan recover, berapa lama downtime
- [x] **Email: "Server Down"** — notifikasi instan saat down
- [x] **Email: "Server Recovered"** — notifikasi saat sistem kembali online
- [x] **Telegram: Recovery Alert** — notifikasi Telegram saat sistem kembali online
- [x] **Frontend Incidents Page** — Halaman `/incidents` untuk tracking downtime

### Status Page Extended

- [x] **Custom URL/Slug** — `status.your-domain.com/:slug`
- [x] **Public Status Page** — `GET /api/status` endpoint tersedia tanpa autentikasi
- [x] **Status Page UI** — source cards dengan badge ONLINE/DEGRADED/OFFLINE/MAINTENANCE
- [x] **Uptime Statistics** — uptime % 30/90 hari terakhir
- [x] **Embed Widget** — badge status dengan token-based access

### API Endpoints

- [x] `GET /api/status` — public status page data
- [x] `GET /api/admin/status-pages` — list all status pages
- [x] `POST /api/admin/status-pages` — create status page
- [x] `GET /api/admin/status-pages/:source_id` — get status page config
- [x] `PATCH /api/admin/status-pages/:source_id` — update status page
- [x] `DELETE /api/admin/status-pages/:source_id` — delete status page
- [x] `GET /api/admin/status-pages/:source_id/uptime` — get uptime stats
- [x] `POST /api/admin/status-pages/:source_id/embed` — create embed widget
- [x] `GET /status/:slug` — public status page by slug
- [x] `GET /embed/:token` — embed widget data

---

## ✅ Fase 5 — Error Grouping & Smart Analysis

**Status**: ✅ Fully Implemented

Agar tidak pusing membaca ribuan baris log yang sama — kelompokkan dan analisis secara otomatis.

### Error Grouping

- [x] **Auto-Grouping** — log dengan message dan stack trace yang mirip dikelompokkan
- [x] **Fingerprinting Algorithm** — hash dari `(source_id + normalized_message + stack_trace[:100])`
- [x] **Issue Tracker** — halaman "Issues" di dashboard
- [x] **Issue Detail** — lihat semua individual log dalam satu group
- [x] **Issue Status** — tandai issue sebagai `OPEN`, `RESOLVED`, `IGNORED`
- [x] **Regression Detection** — alert jika issue yang sudah `RESOLVED` muncul lagi

### Error Analytics

- [x] **"Source mana yang paling sering error?"** — ranking occurrence count per source
- [x] **"Jam berapa error paling sering terjadi?"** — heatmap per jam dalam seminggu
- [x] **"Error apa yang paling sering muncul?"** — top 10 error messages
- [x] **Level breakdown** — badge count CRITICAL/ERROR/WARN untuk open issues
- [x] **Error Rate Trend** — persentase log yang error vs total per hari

### AI Integration

- [x] **AI Copilot (Chatbot)** — Floating chat widget dengan Groq llama-3.3-70b
- [x] **AI Daily Digest** — LLM merangkum error/anomali hari ini (dikirim via email 8 AM)
- [x] **AI Error Deduplication** — gunakan embedding untuk mengelompokkan error semantically similar
- [x] **Prompt Context** — sertakan framework, bahasa, dan error history agar AI suggestion lebih relevan

### API Endpoints

- [x] `GET /api/issues` — list all issues
- [x] `GET /api/issues/analytics/trend` — error rate trend
- [x] `GET /api/issues/analytics/heatmap` — error heatmap
- [x] `GET /api/issues/:fingerprint` — get issue detail
- [x] `PATCH /api/issues/:fingerprint` — update issue status
- [x] `GET /api/issues/:fingerprint/logs` — get logs for issue

---

## ✅ Fase 6 — Centralized Configuration Management

**Status**: ✅ Fully Implemented

Kelola semua konfigurasi aplikasi-aplikasi Anda dari satu tempat — tanpa harus SSH ke setiap server.

### Config Storage

- [x] **Config Table** — tabel `source_configs` di PostgreSQL
- [x] **Secret Management** — nilai sensitif dienkripsi dengan AES-256
- [x] **Versioning** — setiap perubahan config tersimpan di `source_config_histories`
- [x] **Environment Namespacing** — config bisa punya namespace: `production`, `staging`, `development`

### Config API

- [x] `GET /api/config/:source_slug` — aplikasi client pull config terbaru
- [x] `PUT /api/config/:source_slug/:key` — update config value via dashboard
- [x] `GET /api/config/:source_slug/history` — riwayat perubahan config

### Hot Reload & Notifications

- [x] **Hot Reload Config** — aplikasi client bisa polling `/api/config/:slug`
- [x] **Change Notification** — saat config diubah, ULAM kirim webhook ke aplikasi client
- [x] **HMAC Signature** — webhook requests signed dengan HMAC-SHA256
- [x] **Go SDK Config Helper** — `WatchConfig()` function untuk Go applications
- [x] **Config Audit Trail** — siapa yang mengubah config apa, kapan

### Dashboard Config Management

- [x] **Config Editor UI** — halaman `/config` di dashboard
- [x] **Secret Toggle** — nilai secret ditampilkan sebagai `••••••••`
- [x] **Rollback** — rollback config ke versi sebelumnya dari tab History
- [x] **Environment Selector** — pilih environment: production/staging/development

---

## ✅ Fase 7 — Export & Extended Integrations

**Status**: ✅ Fully Implemented

Export data dan notifikasi ke platform lain.

### Export Formats

- [x] **CSV Export** — `GET /api/logs/export` — export hasil filter log ke CSV
- [x] **Excel Export** — `GET /api/logs/export/excel` — export ke .xlsx
- [x] **PDF Audit Report** — `GET /api/logs/export/pdf` — export audit trail ke PDF

### Notifications

- [x] **Email Notification** — alert ke email admin (with 5-min throttle)
- [x] **Telegram Bot** — notifikasi ke Telegram group/channel
- [x] **Webhook Support** — generic outgoing webhook ke URL apapun

### Infrastructure

- [x] **Rate Limiting** — Token bucket algorithm untuk semua endpoints
- [x] **Log Archiving** — Interface ready untuk S3/object storage integration
- [x] **Request ID Tracing** — setiap request punya unique ID untuk debugging

### SDKs

- [x] **Official Go SDK** — Complete SDK dengan semua tipe log
- [x] **Go SDK Config Helper** — `WatchConfig()` untuk hot reload
- [x] **Node.js SDK** — Complete SDK untuk Node.js applications
- [x] **Python SDK** — Complete SDK untuk Python applications
- [x] **PHP SDK** — Complete SDK untuk PHP applications

### Frontend

- [x] **Notification System** — Dropdown notifications dengan bell icon di navbar
- [x] **Dark Theme** — Consistent dark UI dengan Tailwind v4
- [x] **Responsive Design** — Mobile-friendly dashboard

### Third-Party Integrations (Future/Optional)

- [ ] **Google Workspace Tracker** — Background cron-job untuk Google Admin SDK
- [ ] **AWS CloudTrail Sync** — Menarik log IAM dan aktivitas S3
- [ ] **GitHub Org Audit** — Track repo changes & access keys

> **Note**: Third-party integrations require external API credentials and are optional features for enterprise users.

---

## 📊 Feature Completion Summary

| Phase | Features | Status | Completion |
|-------|----------|--------|------------|
| Fase 1 | Core MVP | ✅ Complete | 100% |
| Fase 2 | Activity Monitor | ✅ Complete | 100% |
| Fase 3 | APM | ✅ Complete | 100% |
| Fase 4 | Status Page | ✅ Complete | 100% |
| Fase 5 | Error Grouping & AI | ✅ Complete | 95% |
| Fase 6 | Config Management | ✅ Complete | 100% |
| Fase 7 | Export & Integrations | ✅ Complete | 95% |

**Overall Progress: ~98%**

---

## 🎉 Project Status: PRODUCTION READY

One-Log (ULAM) is now feature-complete and ready for production deployment!
