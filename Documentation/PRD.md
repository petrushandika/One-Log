# 📋 Product Requirements Document (PRD)

## Unified Log & Activity Monitor (ULAM)

| Field       | Details                               |
| ----------- | ------------------------------------- |
| **Project** | Unified Log & Activity Monitor (ULAM) |
| **Status**  | Discovery / Planning                  |
| **Author**  | Petrus Handika                        |
| **Version** | v1.0.0                                |
| **Date**    | March 2026                            |
| **Stack**   | React, Golang, PostgreSQL, GORM, SMTP |

---

## 1. Executive Summary

**ULAM** adalah platform terpusat untuk menangkap, menyimpan, dan menganalisis log teknis serta aktivitas pengguna dari **seluruh aplikasi yang Anda bangun** — baik yang sudah ada maupun yang akan datang.

Sistem dirancang agar **scalable secara horizontal**: setiap aplikasi baru cukup mendaftar, mendapatkan API key, dan mulai mengirim log tanpa perubahan apapun di sisi ULAM.

**Masalah yang diselesaikan:**

- Developer harus SSH ke banyak server untuk membaca log
- Error baru diketahui ketika user sudah melaporkan (reaktif, bukan proaktif)
- Tidak ada gambaran holistik tentang kesehatan semua sistem secara bersamaan
- Stack trace tersebar di berbagai tempat, membuat debugging lambat

---

## 2. Goals & Objectives

| #   | Goal                     | Success Metric                                            |
| --- | ------------------------ | --------------------------------------------------------- |
| 1   | **Centralization**       | 0 server perlu di-SSH untuk membaca log                   |
| 2   | **Real-time Alerting**   | Notifikasi email terkirim dalam < 30 detik setelah event  |
| 3   | **Observability**        | Dapat melihat aktivitas login & audit trail per aplikasi  |
| 4   | **Developer Efficiency** | MTTD (Mean Time to Detect) error kritis < 5 menit         |
| 5   | **Open Integration**     | Aplikasi baru bisa mulai kirim log dalam < 10 menit setup |

---

## 3. Target Audience

### 👨‍💻 Internal Developer (Primary User)

- Memantau kesehatan semua aplikasi secara teknis
- Debugging dengan stack trace lengkap dan konteks tambahan
- Menganalisis tren error dan pola masalah
- **Use case**: "Ada laporan error dari user, saya langsung buka ULAM, filter project + ERROR level, dan temukan stack trace dalam 2 menit"

### 🔐 System Administrator (Secondary User)

- Memantau pola login yang mencurigakan dari berbagai aplikasi
- Audit trail: siapa login, kapan, dari mana
- **Use case**: "Saya lihat ada 50 failed login attempts dari 1 IP dalam 5 menit"

---

## 4. Problem Statement

### Situasi Saat Ini

Setiap aplikasi (e-learning, absensi, CMS, dan aplikasi-aplikasi berikutnya) menyimpan log secara independen. Ini menciptakan masalah yang semakin besar seiring bertambahnya jumlah aplikasi:

```text
Aplikasi A → log di /var/log/app-a/
Aplikasi B → log di server lain
Aplikasi C → log di cloud provider X
Aplikasi D → log di hanya stdout container
...
Aplikasi N → semakin banyak tempat yang harus dicek
```

### Dampak Nyata

| Problem            | Dampak                                            |
| ------------------ | ------------------------------------------------- |
| Fragmentasi log    | Debugging butuh akses ke N server berbeda         |
| Tidak ada alerting | Error diketahui dari laporan user, bukan sistem   |
| Kurang visibilitas | Tidak ada gambaran kesehatan keseluruhan sistem   |
| Audit trail lemah  | Sulit tracking aktivitas pengguna lintas aplikasi |

---

## 5. Proposed Solution

ULAM memperkenalkan model **Push-based Centralized Logging**:

```text
Semua Aplikasi → ULAM API (single endpoint) → PostgreSQL → Dashboard
```

Setiap aplikasi cukup:

1. Daftarkan diri → dapat API Key
2. Pasang HTTP client sederhana (tersedia helper library)
3. Panggil `LogToULAM()` di titik-titik penting dalam kode

ULAM menangani sisanya: penyimpanan, searching, alerting, dan visualisasi.

---

## 6. Functional Requirements

### 6.1 Application Registration & Management

- Admin bisa mendaftarkan aplikasi baru via dashboard
- Sistem generate API key unik per aplikasi
- Admin bisa menonaktifkan aplikasi (soft-disable) tanpa kehilangan data
- Admin bisa regenerate API key jika terjadi kebocoran

### 6.2 Log Ingestion API

| Spec           | Detail                                                     |
| -------------- | ---------------------------------------------------------- |
| **Endpoint**   | `POST /api/ingest`                                         |
| **Auth**       | `X-API-Key` header                                         |
| **Strategy**   | Async: response < 100ms, proses DB di background goroutine |
| **Format**     | JSON payload                                               |
| **Validation** | Required fields, enum validation untuk category & level    |
| **Rate Limit** | 100 req/menit per API key                                  |

**Payload yang diterima:**

```json
{
  "category": "SYSTEM_ERROR",
  "level": "CRITICAL",
  "message": "Database connection pool exhausted",
  "stack_trace": "goroutine 1 [running]:\nmain.connectDB()...",
  "context": {
    "user_id": "usr_123",
    "endpoint": "/api/attendance",
    "duration_ms": 8000,
    "browser": "Chrome/121"
  }
}
```

> **Catatan**: `source_id` tidak perlu dikirim oleh client karena sudah diketahui dari API key.

### 6.3 Data Privacy & Security (PII Masking)

Untuk menjaga keamanan data, ULAM menerapkan **PII (Personally Identifiable Information) Masking**:

- **Automatic Masking**: Sistem akan mendeteksi dan mensensor (mengganti menjadi `***`) kata kunci sensitif dalam field `message` dan `context` seperti: `password`, `secret`, `api_key`, `access_token`, `auth_token`.
- **Level**: PII Masking dilakukan di tingkat Backend sebelum data ditulis ke database (At-Rest Security).
- **Custom Patterns**: Admin dapat menentukan regex tambahan untuk masking (misal: nomor kartu kredit).

### 6.4 Error Grouping (Log Fingerprinting)

Untuk menghindari "dashboard noise", ULAM mengimplementasikan **Error Grouping**:

- **Mechanism**: Jika fitur ini aktif, sistem akan menghitung hash dari `source_id`, `category`, dan subset dari `message` (50 karakter pertama).
- **Dashboard View**: Log yang identik akan dikelompokkan. User melihat jumlah kejadian (*occurrence count*) dan waktu terakhir muncul (*last seen*) di satu baris yang sama.
- **Alert Suppression**: Email alert hanya dikirim untuk kejadian pertama dalam siklus cooldown (5 menit).

### 6.5 Log Category & Level System

**Categories:**

| Category        | Deskripsi                                                |
| --------------- | -------------------------------------------------------- |
| `SYSTEM_ERROR`  | Error teknis (DB fail, crash, 5xx)                       |
| `USER_ACTIVITY` | Aktivitas pengguna (CRUD, download, navigasi)            |
| `AUTH_EVENT`    | Login, logout, OAuth, token refresh, failed login        |
| `PERFORMANCE`   | Slow query, timeout, resource exhaustion                 |
| `SECURITY`      | Suspicious activity, rate limit hit, unauthorized access |

**Levels:**

| Level      | Severity                                             | Alert?   |
| ---------- | ---------------------------------------------------- | -------- |
| `CRITICAL` | Sistem tidak bisa beroperasi                         | ✅ Email |
| `ERROR`    | Fungsi gagal, user terdampak                         | ✅ Email |
| `WARN`     | Potensi masalah, sistem masih jalan                  | ❌       |
| `INFO`     | Informasi operasional normal                         | ❌       |
| `DEBUG`    | Detail untuk debugging (opsional, di-filter di prod) | ❌       |

### 6.6 Metadata Capture

Setiap log entry menangkap:

- **IP Address** dari request header (`X-Forwarded-For` atau `RemoteAddr`)
- **Timestamp** dalam UTC
- **Application ID** dari API key lookup
- **User-defined context** via JSONB field (bebas sesuai kebutuhan)

### 6.7 Notification Engine

| Spec           | Detail                                                            |
| -------------- | ----------------------------------------------------------------- |
| **Trigger**    | `level == ERROR` atau `level == CRITICAL`                         |
| **Delivery**   | SMTP Email ke admin yang terdaftar                                |
| **Throttling** | Max 1 notifikasi per `{source_id}:{message_hash}` per 5 menit     |
| **Template**   | HTML email: Nama App, Level, Message, Stack Trace, Link Dashboard |
| **Timeout**    | SMTP call timeout 10 detik                                        |

### 6.8 Admin Dashboard

| Fitur                 | Deskripsi                                                                                |
| --------------------- | ---------------------------------------------------------------------------------------- |
| **Overview**          | Statistik agregat: total log per level, per aplikasi, dalam 24h/7d/30d                   |
| **Log Table**         | Tabel dengan pagination, filter multi-kriteria, dan search                               |
| **Log Detail**        | Tampilan lengkap: message, stack trace, formatted JSON context                           |
| **AI Insight**        | Tombol "Analyze" untuk mendapatkan ringkasan error & solusi via **Groq API**             |
| **Activity Monitor**  | Halaman khusus untuk melihat login events, auth method breakdown, dan user session trail |
| **Source Management** | Daftar source terdaftar, status aktif/nonaktif, API key management                       |
| **Search**            | Full-text search di field `message`                                                      |

---

### 6.9 AI Insight Engine (via Groq API)

ULAM mengintegrasikan AI untuk membantu interpretasi log teknis:

- **Technology**: Menggunakan **Groq API** dengan model **Llama 3 / Mixtral** (High-speed inference).
- **Error Summarization**: Mengubah stack trace yang kompleks menjadi ringkasan satu kalimat yang mudah dipahami manusia.
- **RCA (Root Cause Analysis)**: Memberikan kemungkinan penyebab utama berdasarkan pesan error.
- **Solution Suggestion**: Memberikan 3 langkah perbaikan yang direkomendasikan.
- **Trigger**:
    - **Automatic**: Untuk log level `CRITICAL`, AI di-trigger otomatis di background.
    - **Manual**: User bisa menekan tombol "Analyze" pada log level apapun di dashboard.

---

### 6.10 User Activity Tracking

Ini adalah fitur **pertama selain error monitoring** — mencatat siapa yang melakukan apa, dari mana, dan dengan cara apa. Tidak hanya error, tapi seluruh jejak aktivitas pengguna lintas semua source yang terdaftar.

#### 6.7.1 Auth Event Tracking

Setiap event autentikasi dari aplikasi manapun harus dicatat ke ULAM dengan `category: AUTH_EVENT`.

**Auth methods yang didukung (via payload):**

| `auth_method` Value | Deskripsi                                    |
| ------------------- | -------------------------------------------- |
| `google_oauth`      | Login via Google OAuth 2.0                   |
| `github_oauth`      | Login via GitHub                             |
| `facebook_oauth`    | Login via Facebook                           |
| `twitter_oauth`     | Login via Twitter/X                          |
| `discord_oauth`     | Login via Discord                            |
| `system_password`   | Login manual dengan username + password      |
| `magic_link`        | Login via magic link (email-based)           |
| `sso`               | Login via SSO / SAML enterprise              |
| `api_token`         | Autentikasi via API token (bukan user login) |

**Auth event types yang dicatat:**

| `event_type`       | Deskripsi                                         |
| ------------------ | ------------------------------------------------- |
| `login_success`    | Login berhasil                                    |
| `login_failed`     | Login gagal (password salah, token invalid, dll.) |
| `logout`           | User logout secara eksplisit                      |
| `token_refresh`    | JWT/session token diperbarui                      |
| `oauth_callback`   | OAuth callback diterima dari provider             |
| `account_linked`   | User menghubungkan akun sosmed ke akun existing   |
| `password_reset`   | User melakukan reset password                     |
| `mfa_challenge`    | Multi-factor authentication diminta               |
| `mfa_success`      | MFA berhasil diverifikasi                         |
| `session_expired`  | Session berakhir karena timeout                   |
| `suspicious_login` | Login dari IP/device baru yang tidak dikenal      |

**Standardized payload untuk AUTH_EVENT:**

```json
{
  "category": "AUTH_EVENT",
  "level": "INFO",
  "message": "User login via Google OAuth",
  "context": {
    "event_type": "login_success",
    "auth_method": "google_oauth",
    "user_id": "usr_abc123",
    "email": "user@example.com",
    "name": "John Doe",
    "avatar_url": "https://googleusercontent.com/...",
    "provider_id": "google_uid_xyz",
    "ip_address": "103.120.45.1",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
    "browser": "Chrome",
    "browser_version": "121.0.0.0",
    "os": "Windows 10",
    "device_type": "desktop",
    "location": {
      "country": "ID",
      "city": "Jakarta"
    },
    "session_id": "sess_xyz789",
    "is_new_user": false
  }
}
```

**Payload untuk login gagal:**

```json
{
  "category": "AUTH_EVENT",
  "level": "WARN",
  "message": "Login failed: invalid password",
  "context": {
    "event_type": "login_failed",
    "auth_method": "system_password",
    "attempted_email": "user@example.com",
    "failure_reason": "invalid_password",
    "attempt_count": 3,
    "ip_address": "103.120.45.1",
    "user_agent": "Mozilla/5.0...",
    "browser": "Chrome",
    "os": "Windows 10"
  }
}
```

> **Level recommendation**: `login_success` → INFO, `login_failed` (1x) → WARN, `login_failed` (5x+ dari IP sama) → ERROR (triggers email alert)

#### 6.7.2 User Activity Trail

Selain auth event, aktivitas umum pengguna di dalam aplikasi juga bisa dikirim ke ULAM dengan `category: USER_ACTIVITY`. Ini membentuk log audit trail yang lengkap.

**Contoh aktivitas yang dicatat:**

| `action`    | Deskripsi                             |
| ----------- | ------------------------------------- |
| `page_view` | User mengunjungi halaman tertentu     |
| `create`    | User membuat data baru                |
| `update`    | User mengubah data                    |
| `delete`    | User menghapus data                   |
| `export`    | User mengekspor data (CSV, PDF, dll.) |
| `download`  | User mengunduh file                   |
| `search`    | User melakukan pencarian              |
| `share`     | User berbagi konten                   |

**Standardized payload untuk USER_ACTIVITY:**

```json
{
  "category": "USER_ACTIVITY",
  "level": "INFO",
  "message": "User exported attendance report",
  "context": {
    "action": "export",
    "user_id": "usr_abc123",
    "resource_type": "attendance_report",
    "resource_id": "rpt_march_2026",
    "metadata": {
      "format": "xlsx",
      "row_count": 245,
      "filters": { "month": "2026-03", "department": "Engineering" }
    },
    "session_id": "sess_xyz789",
    "ip_address": "192.168.1.10"
  }
}
```

#### 6.7.3 Activity Dashboard Page

Halaman khusus di dashboard ULAM untuk memantau aktivitas:

| Sub-Fitur                 | Deskripsi                                                            |
| ------------------------- | -------------------------------------------------------------------- |
| **Auth Method Breakdown** | Pie/bar chart: berapa persen login via Google vs Manual vs lainnya   |
| **Login Timeline**        | Timeline login events per aplikasi \| per user, per hari             |
| **Failed Login Heatmap**  | Visualisasi jam-jam dengan banyak failed login (deteksi brute force) |
| **Recent Sessions**       | Tabel: user, auth_method, IP, browser, timestamp, app                |
| **Top Active Users**      | Siapa user yang paling aktif di tiap aplikasi                        |
| **Geo Map** _(future)_    | Visualisasi peta asal IP login                                       |

---

## 7. Non-Functional Requirements

### 7.1 Performance

| Metric                      | Target                   |
| --------------------------- | ------------------------ |
| Ingestion API response time | < 100ms (P95)            |
| Dashboard log list load     | < 500ms (P95)            |
| Dashboard overview stats    | < 1 detik                |
| Email notification delivery | < 30 detik setelah event |

### 7.2 Reliability & Availability

| Metric             | Target                               |
| ------------------ | ------------------------------------ |
| Uptime             | 99.5% per bulan                      |
| DB connection pool | Min 5, Max 20 concurrent connections |
| Graceful shutdown  | Drain goroutines sebelum shutdown    |

### 7.3 Security

| Aspek            | Implementasi                                                 |
| ---------------- | ------------------------------------------------------------ |
| Transport        | HTTPS wajib (TLS 1.2+)                                       |
| Ingestion Auth   | API Key (hashed bcrypt di DB)                                |
| Dashboard Auth   | JWT via **httpOnly Cookie** (`ulam_access` + `ulam_refresh`) |
| Access Token     | HS256, expire **24 jam** — disimpan di httpOnly cookie       |
| Refresh Token    | HS256, expire **7 hari** — scope hanya `/api/auth/refresh`   |
| API Key format   | `ulam_` prefix + 32 random hex chars                         |
| Password storage | bcrypt cost factor 12                                        |
| Rate limiting    | Per API Key dan per IP                                       |
| Cookie flags     | `HttpOnly`, `Secure`, `SameSite=Strict`                      |
| CORS             | Whitelist origin dashboard saja                              |

### 7.5 Log Retention Policy

Untuk menjaga performa database dan efisiensi penyimpanan:

| Retention Rule       | Detail                                          |
| -------------------- | ----------------------------------------------- |
| **Max Retention**    | Log akan dihapus otomatis setelah 30 hari       |
| **Critical Logs**    | Log level `CRITICAL` disimpan selama 90 hari    |
| **Cleanup Schedule** | Background job dijalankan setiap hari pada 02:00 UTC |
| **Auto-Archive**     | (Post-MVP) Opsi backup ke S3 sebelum dihapus    |

### 7.4 Scalability

| Aspek                   | Kapasitas MVP                             |
| ----------------------- | ----------------------------------------- |
| Jumlah source terdaftar | Unlimited (tidak ada batasan di DB)       |
| Log throughput          | 1.000 log/menit (1 instance)              |
| Concurrent connections  | 20 DB connections, 1.000 goroutines       |
| Storage growth          | ~1 KB/log → 1.000 log/hari = ~30 MB/bulan |

---

## 8. Integration Guide (untuk Aplikasi Client)

### Setup (One-time, < 10 menit)

1. Admin buka ULAM Dashboard → **Apps** → **Register New App**
2. Isi nama dan slug aplikasi → Klik **Generate**
3. Salin API Key yang muncul (hanya tampil sekali)
4. Simpan ke `.env` aplikasi:

```env
ULAM_ENDPOINT=https://api.ulam.your-domain.com/api/ingest
ULAM_API_KEY=ulam_a1b2c3d4e5f6g7h8i9j0...
```

### Integrasi di Kode (Golang)

```go
// pkg/ulam/client.go — Helper yang dipasang di tiap project
package ulam

import (
    "bytes"
    "encoding/json"
    "net/http"
    "os"
)

type LogPayload struct {
    Category   string                 `json:"category"`
    Level      string                 `json:"level"`
    Message    string                 `json:"message"`
    StackTrace string                 `json:"stack_trace,omitempty"`
    Context    map[string]interface{} `json:"context,omitempty"`
}

func Send(p LogPayload) {
    go func() {
        body, _ := json.Marshal(p)
        req, _ := http.NewRequest("POST", os.Getenv("ULAM_ENDPOINT"), bytes.NewBuffer(body))
        req.Header.Set("X-API-Key", os.Getenv("ULAM_API_KEY"))
        req.Header.Set("Content-Type", "application/json")
        http.DefaultClient.Do(req) // Fire and forget
    }()
}
```

### Penggunaan

```go
// Saat error database
ulam.Send(ulam.LogPayload{
    Category: "SYSTEM_ERROR",
    Level:    "ERROR",
    Message:  "DB connection failed: " + err.Error(),
    StackTrace: debug.Stack(),
    Context: map[string]interface{}{
        "endpoint": r.URL.Path,
        "user_id":  currentUser.ID,
    },
})

// Saat user login via Google
ulam.Send(ulam.LogPayload{
    Category: "AUTH_EVENT",
    Level:    "INFO",
    Message:  "User login via Google OAuth",
    Context: map[string]interface{}{
        "user_id": user.ID,
        "email":   user.Email,
        "method":  "google_oauth",
        "browser": r.Header.Get("User-Agent"),
    },
})
```

---

## 9. Data Model Summary

```text
LogEntry {
  id          uint       — Auto-increment PK
  source_id      string     — Slug aplikasi pengirim (from API key)
  category    string     — SYSTEM_ERROR | USER_ACTIVITY | AUTH_EVENT | ...
  level       string     — CRITICAL | ERROR | WARN | INFO | DEBUG
  message     text       — Pesan utama log
  stack_trace text?      — Stack trace (opsional)
  context     jsonb?     — Metadata bebas
  ai_insight  jsonb?     — Analisis AI dari Groq
  ip_address  string?    — IP pengirim
  created_at  timestamp  — Waktu diterima (UTC)
}

Application {
  id         uint
  name       string     — "Sistem Absensi Production"
  slug       string     — "absensi-prod" (unique, dipakai sebagai source_id)
  api_key    string     — Hashed bcrypt
  is_active  bool
  created_at timestamp
}
```

---

## 10. Out of Scope (MVP)

| Feature                    | Alasan Ditunda                  |
| -------------------------- | ------------------------------- |
| CSV/Excel export           | Tidak urgent                    |
| Slack/Telegram notifikasi  | SMTP cukup                      |
| Custom alert rules per app | Hardcoded trigger cukup         |
| Multi-admin dengan RBAC    | Single admin cukup              |
| SDK resmi per bahasa       | Helper function sederhana cukup |
| Log streaming (WebSocket)  | Polling cukup untuk MVP         |

---

## 11. Future Roadmap

| Version | Feature                                           | Priority |
| ------- | ------------------------------------------------- | -------- |
| v1.1    | Log Retention Policy (auto-delete setelah N hari) | High     |
| v1.2    | Export CSV/Excel                                  | Medium   |
| v1.3    | Slack / Telegram Integration                      | Medium   |
| v1.4    | Webhook support (generic outgoing)                | Medium   |
| v2.0    | AI Insight — LLM summarize daily errors           | Low      |
| v2.0    | Multi-admin + RBAC                                | Low      |
| v2.0    | Official SDK (Go, Node.ts, Python)                | Low      |

---

## 12. Risks & Mitigations

| Risk                      | Probability | Impact | Mitigation                                         |
| ------------------------- | ----------- | ------ | -------------------------------------------------- |
| SMTP rate limit (Gmail)   | Medium      | Medium | Gunakan App Password; switch ke Resend jika perlu  |
| API Key bocor dari client | Medium      | High   | Dokumentasikan cara rotate key; throttling per key |
| DB disk penuh             | Low         | High   | Monitor storage; alert manual jika > 80%           |
| Goroutine leak            | Low         | Medium | Timeout context di setiap goroutine                |
| Email spam (no throttle)  | Low         | Medium | Throttle map dengan 5 menit cooldown               |

---

_Dokumen ini adalah dokumen hidup yang diperbarui seiring perkembangan project._
