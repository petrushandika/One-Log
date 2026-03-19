# 🔌 API Reference

## Unified Log & Activity Monitor (ULAM)

**Base URL**: `https://api.ulam.your-domain.com`  
**Content-Type**: `application/json`

---

## URL Conventions

```text
/api/ingest                     → Log ingestion (for client apps)
/api/auth/...                   → Authentication
/api/logs/...                   → Log management (dashboard)
/api/logs/export                → CSV export
/api/sources/...                → Source management
/api/stats/...                  → Aggregated statistics
/api/activity/...               → Activity & audit trail (Phase 2)
/api/apm/...                    → APM — endpoint latency (Phase 3)
/api/apm/timeline               → Response time timeline (Phase 3)
/api/issues/...                 → Error grouping & issues (Phase 5)
/api/issues/analytics/trend     → Error rate trend (Phase 5)
/api/issues/analytics/heatmap   → Error heatmap (Phase 5)
/api/incidents/...              → Incident management (Phase 4)
/api/incidents/timeline         → Incident timeline (Phase 4)
/api/config/...                 → Centralized config management (Phase 6)
/api/status                     → Public status page (no auth, Phase 4)
/api/chat                       → AI Copilot chatbot (Phase 5, JWT required)
```

> Note: `/sources` digunakan untuk merepresentasikan "source terdaftar yang mengirim log ke ULAM".
> Tidak menggunakan `/sources` agar lebih semantik dan tidak ambigu.

---

## Authentication

| Context                   | Method          | Mekanisme                                 |
| ------------------------- | --------------- | ----------------------------------------- |
| Client source (ingestion) | API Key         | Header: `X-API-Key: ulam_xxxxxxxx`        |
| Admin dashboard           | httpOnly Cookie | `ulam_access` (24h) + `ulam_refresh` (7d) |

> 🔒 Token **tidak diekspos ke JavaScript** — disimpan sebagai `httpOnly; Secure; SameSite=Strict` cookie untuk mencegah XSS.

---

## 1. Authentication

### `POST /api/auth/login`

Login admin. Response **tidak mengembalikan token di body** — token langsung di-set sebagai httpOnly cookie.

**Auth**: Public

**Request:**

```json
{
  "username": "admin",
  "password": "your_secure_password"
}
```

**Response `200 OK`:**

```json
{
  "message": "Login successful",
  "admin": {
    "id": 1,
    "username": "admin"
  }
}
```

**Set-Cookie headers yang dikirim server:**

```http
Set-Cookie: ulam_access=eyJhbGci...; HttpOnly; Secure; SameSite=Strict; Max-Age=86400; Path=/
Set-Cookie: ulam_refresh=eyJhbGci...; HttpOnly; Secure; SameSite=Strict; Max-Age=604800; Path=/api/auth/refresh
```

| Cookie         | Max-Age           | Scope                                        |
| -------------- | ----------------- | -------------------------------------------- |
| `ulam_access`  | `86400` (24 jam)  | `/` — semua endpoint                         |
| `ulam_refresh` | `604800` (7 hari) | `/api/auth/refresh` — hanya endpoint refresh |

**Response `401 Unauthorized`:**

```json
{
  "status": "error",
  "code": 401,
  "message": "Invalid credentials",
  "errors": null
}
```

---

### `POST /api/auth/refresh`

Generate access token baru menggunakan refresh token. Dipanggil otomatis oleh frontend ketika access token expired (401).

**Auth**: Cookie `ulam_refresh`

**Request:** _(no body — refresh token dibaca otomatis dari cookie)_

**Response `200 OK`:**

```http
Set-Cookie: ulam_access=eyJhbGci...; HttpOnly; Secure; SameSite=Strict; Max-Age=86400; Path=/
```

```json
{ "message": "Token refreshed" }
```

**Response `401 Unauthorized`** _(refresh token expired atau invalid)_:

```json
{
  "error": "Refresh token expired or invalid. Please login again.",
  "code": "REFRESH_TOKEN_EXPIRED"
}
```

---

### `POST /api/auth/logout`

Hapus kedua cookies. Admin harus login ulang.

**Auth**: Cookie `ulam_access`

**Response `200 OK`:**

```http
Set-Cookie: ulam_access=; HttpOnly; Secure; Max-Age=0; Path=/
Set-Cookie: ulam_refresh=; HttpOnly; Secure; Max-Age=0; Path=/api/auth/refresh
```

```json
{ "message": "Logged out successfully" }
```

---

## 2. Log Ingestion

### `POST /api/ingest`

Endpoint utama untuk menerima log dari aplikasi client.

**Auth**: `X-API-Key`  
**Rate Limit**: 100 req/menit per API key

> ✅ **Desain**: `source_id` **tidak perlu dikirim** oleh client. Sistem menentukan aplikasi pengirim secara otomatis dari API key.

**Request Body:**

```json
{
  "category": "SYSTEM_ERROR",
  "level": "ERROR",
  "message": "Failed to connect to PostgreSQL",
  "stack_trace": "goroutine 1 [running]:\nmain.connectDB()\n\t/app/db.go:45 +0x1a2",
  "context": {
    "endpoint": "/api/attendance/check-in",
    "method": "POST",
    "user_id": "usr_abc123",
    "duration_ms": 8500,
    "db_host": "localhost:5432"
  }
}
```

**Field Validation:**

| Field         | Type   | Required | Enum / Notes                                                             |
| ------------- | ------ | -------- | ------------------------------------------------------------------------ |
| `category`    | string | ✅       | `SYSTEM_ERROR`, `USER_ACTIVITY`, `AUTH_EVENT`, `PERFORMANCE`, `SECURITY` |
| `level`       | string | ✅       | `CRITICAL`, `ERROR`, `WARN`, `INFO`, `DEBUG`                             |
| `message`     | string | ✅       | Max 5000 chars                                                           |
| `stack_trace` | string | ❌       | Max 50000 chars                                                          |
| `context`     | object | ❌       | Free-form JSON, max 10 fields (MVP)                                      |

**Response `202 Accepted`:**

```json
{
  "status": "accepted",
  "message": "Log received and queued for processing",
  "request_id": "req_7f3a1b2c"
}
```

**Response `400 Bad Request`:**

```json
{
  "error": "Validation failed",
  "code": "INVALID_PAYLOAD",
  "details": [
    "level must be one of: CRITICAL, ERROR, WARN, INFO, DEBUG",
    "message is required"
  ]
}
```

**Response `401 Unauthorized`:**

```json
{
  "error": "Invalid or missing API key",
  "code": "INVALID_API_KEY"
}
```

**Response `429 Too Many Requests`:**

```json
{
  "error": "Rate limit exceeded",
  "code": "RATE_LIMIT_EXCEEDED",
  "retry_after_seconds": 30
}
```

---

## 3. Logs

### `GET /api/logs`

Mengambil daftar log dengan pagination dan filter multi-kriteria.

**Auth**: JWT

**Query Parameters:**

| Param       | Type    | Default | Description                                      |
| ----------- | ------- | ------- | ------------------------------------------------ |
| `source_id` | string  | -       | Filter by application slug, e.g. `absensi-prod`  |
| `level`     | string  | -       | Satu atau lebih, pisahkan koma: `ERROR,CRITICAL` |
| `category`  | string  | -       | Filter by category                               |
| `search`    | string  | -       | Full-text search dalam field `message`           |
| `from`      | ISO8601 | -       | Filter log dari tanggal ini (UTC)                |
| `to`        | ISO8601 | -       | Filter log sampai tanggal ini (UTC)              |
| `page`      | int     | `1`     | Nomor halaman                                    |
| `limit`     | int     | `20`    | Item per halaman (max: `100`)                    |
| `sort`      | string  | `desc`  | `asc` atau `desc` berdasarkan `created_at`       |

**Example Request:**

```http
GET /api/logs?source_id=absensi-prod&level=ERROR,CRITICAL&from=2026-03-01&page=1
Cookie: ulam_access=eyJhbG...
```

**Response `200 OK`:**

```json
{
  "data": [
    {
      "id": 1042,
      "source_id": "absensi-prod",
      "app_name": "Sistem Absensi Production",
      "category": "SYSTEM_ERROR",
      "level": "ERROR",
      "message": "Failed to connect to PostgreSQL",
      "ip_address": "103.120.45.1",
      "created_at": "2026-03-05T09:30:00Z"
    }
  ],
  "pagination": {
    "total": 47,
    "page": 1,
    "limit": 20,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  }
}
```

---

### `GET /api/logs/:id`

Mengambil detail lengkap satu log entry.

**Auth**: JWT

**Response `200 OK`:**

```json
{
  "id": 1042,
  "source_id": "absensi-prod",
  "app_name": "Sistem Absensi Production",
  "category": "SYSTEM_ERROR",
  "level": "ERROR",
  "message": "Failed to connect to PostgreSQL",
  "stack_trace": "goroutine 1 [running]:\nmain.connectDB()\n\t/app/db.go:45 +0x1a2\n...",
  "context": {
    "endpoint": "/api/attendance/check-in",
    "method": "POST",
    "user_id": "usr_abc123",
    "duration_ms": 8500
  },
  "ip_address": "103.120.45.1",
  "created_at": "2026-03-05T09:30:00Z"
}
```

**Response `404 Not Found`:**

```json
{
  "error": "Log entry not found",
  "code": "NOT_FOUND"
}
```

---

## 4. Sources (Registered Sources)

Menggantikan `/sources`. Merepresentasikan aplikasi apapun yang terdaftar untuk mengirim log ke ULAM.

### `GET /api/sources`

Mengambil semua source yang terdaftar.

**Auth**: JWT

**Query Parameters:**

| Param       | Type   | Default | Description              |
| ----------- | ------ | ------- | ------------------------ |
| `is_active` | bool   | -       | Filter by status aktif   |
| `search`    | string | -       | Search by nama atau slug |

**Response `200 OK`:**

```json
{
  "data": [
    {
      "id": 1,
      "name": "Sistem Absensi Production",
      "slug": "absensi-prod",
      "is_active": true,
      "log_count_24h": 342,
      "error_count_24h": 5,
      "created_at": "2026-03-01T08:00:00Z"
    },
    {
      "id": 2,
      "name": "E-Learning Platform",
      "slug": "elearning-prod",
      "is_active": true,
      "log_count_24h": 1201,
      "error_count_24h": 2,
      "created_at": "2026-03-01T08:05:00Z"
    }
  ],
  "total": 2
}
```

---

### `POST /api/sources`

Mendaftarkan aplikasi baru dan generate API key.

**Auth**: JWT

**Request:**

```json
{
  "name": "Portal CMS Production",
  "slug": "cms-prod"
}
```

| Field  | Type   | Required | Notes                                                 |
| ------ | ------ | -------- | ----------------------------------------------------- |
| `name` | string | ✅       | Nama deskriptif, max 100 chars                        |
| `slug` | string | ✅       | Lowercase, huruf + angka + dash, unique, max 50 chars |

**Response `201 Created`:**

```json
{
  "id": 5,
  "name": "Portal CMS Production",
  "slug": "cms-prod",
  "api_key": "ulam_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "is_active": true,
  "created_at": "2026-03-05T16:31:00Z",
  "_notice": "Store this API key securely. It will NOT be shown again."
}
```

> ⚠️ **Penting**: `api_key` hanya ditampilkan sekali saat pembuatan. Simpan segera ke `.env` aplikasi Anda.

**Response `409 Conflict`:**

```json
{
  "error": "An app with this slug already exists",
  "code": "SLUG_CONFLICT"
}
```

---

### `GET /api/sources/:slug`

Mengambil detail satu aplikasi.

**Auth**: JWT

**Response `200 OK`:**

```json
{
  "id": 1,
  "name": "Sistem Absensi Production",
  "slug": "absensi-prod",
  "is_active": true,
  "stats": {
    "total_logs": 15230,
    "logs_24h": 342,
    "errors_24h": 5,
    "last_log_at": "2026-03-05T16:28:00Z"
  },
  "created_at": "2026-03-01T08:00:00Z"
}
```

---

### `PATCH /api/sources/:slug`

Update informasi atau status aplikasi.

**Auth**: JWT

**Request (update nama atau status):**

```json
{
  "name": "Sistem Absensi V2",
  "is_active": false
}
```

**Response `200 OK`:**

```json
{
  "id": 1,
  "name": "Sistem Absensi V2",
  "slug": "absensi-prod",
  "is_active": false,
  "updated_at": "2026-03-05T16:31:00Z"
}
```

---

### `POST /api/sources/:slug/rotate-key`

Regenerate API key untuk aplikasi (misalnya jika key bocor).

**Auth**: JWT

**Request:** _(no body)_

**Response `200 OK`:**

```json
{
  "slug": "absensi-prod",
  "api_key": "ulam_newkey_z9y8x7w6v5u4t3s2r1q0...",
  "_notice": "Old API key is now invalidated. Update your app's .env immediately."
}
```

> ⚠️ **Penting**: Key lama **langsung invalid** setelah rotate. Segera update di semua instance aplikasi yang berjalan.

---

## 6. Activity Tracking _(Fase 2)_

Endpoint khusus untuk mengambil dan menganalisis **user activity logs** — login events, auth method breakdown, session trail, dan user audit trail. Semua data ini bersumber dari `log_entries` yang dikirim dengan `category: AUTH_EVENT` atau `category: USER_ACTIVITY`.

> 📌 Endpoint `/api/logs` tetap bisa dipakai untuk query activity dengan filter `category=AUTH_EVENT`. Endpoint di bawah ini menyediakan agregasi yang lebih spesifik.

---

### `GET /api/activity`

Mengambil activity logs (AUTH_EVENT dan USER_ACTIVITY) dengan filter khusus.

**Auth**: JWT

**Query Parameters:**

| Param         | Type    | Default                    | Description                                       |
| ------------- | ------- | -------------------------- | ------------------------------------------------- |
| `source_id`   | string  | -                          | Filter by application slug                        |
| `category`    | string  | `AUTH_EVENT,USER_ACTIVITY` | Filter category                                   |
| `event_type`  | string  | -                          | Filter `context.event_type`, e.g. `login_success` |
| `auth_method` | string  | -                          | Filter `context.auth_method`, e.g. `google_oauth` |
| `user_id`     | string  | -                          | Filter by `context.user_id`                       |
| `from`        | ISO8601 | -                          | Start date (UTC)                                  |
| `to`          | ISO8601 | -                          | End date (UTC)                                    |
| `page`        | int     | `1`                        | Nomor halaman                                     |
| `limit`       | int     | `20`                       | Item per halaman (max: 100)                       |

**Example Request:**

```http
GET /api/activity?source_id=absensi-prod&event_type=login_failed&from=2026-03-01
Authorization: Bearer eyJhbG...
```

**Response `200 OK`:**

```json
{
  "data": [
    {
      "id": 2031,
      "source_id": "absensi-prod",
      "app_name": "Sistem Absensi Production",
      "category": "AUTH_EVENT",
      "level": "WARN",
      "message": "Login failed: invalid password",
      "context": {
        "event_type": "login_failed",
        "auth_method": "system_password",
        "attempted_email": "user@example.com",
        "failure_reason": "invalid_password",
        "attempt_count": 3,
        "ip_address": "45.67.89.10",
        "browser": "Chrome",
        "os": "Android 14"
      },
      "ip_address": "45.67.89.10",
      "created_at": "2026-03-05T10:12:00Z"
    }
  ],
  "pagination": {
    "total": 23,
    "page": 1,
    "limit": 20,
    "total_pages": 2,
    "has_next": true,
    "has_prev": false
  }
}
```

---

### `GET /api/activity/summary`

Ringkasan statistik autentikasi: breakdown per `auth_method` dan `event_type`.

**Auth**: JWT

**Query Parameters:**

| Param       | Default | Description                             |
| ----------- | ------- | --------------------------------------- |
| `source_id` | -       | Filter by app (opsional, default semua) |
| `period`    | `7d`    | `24h`, `7d`, `30d`                      |

**Response `200 OK`:**

```json
{
  "period": "7d",
  "total_auth_events": 3421,
  "by_auth_method": [
    { "method": "google_oauth", "count": 2100, "percentage": 61.4 },
    { "method": "system_password", "count": 980, "percentage": 28.6 },
    { "method": "github_oauth", "count": 200, "percentage": 5.8 },
    { "method": "facebook_oauth", "count": 141, "percentage": 4.2 }
  ],
  "by_event_type": [
    { "event_type": "login_success", "count": 3100 },
    { "event_type": "login_failed", "count": 280 },
    { "event_type": "logout", "count": 41 }
  ],
  "failed_logins_trend": [
    { "date": "2026-02-28", "failed": 35 },
    { "date": "2026-03-01", "failed": 42 },
    { "date": "2026-03-05", "failed": 61 }
  ],
  "suspicious_logins": 3
}
```

---

### `GET /api/activity/users/:user_id`

Semua aktivitas satu user ID, **lintas semua aplikasi** yang terdaftar di ULAM.

**Auth**: JWT

**Path Parameter:**

| Param     | Description                                        |
| --------- | -------------------------------------------------- |
| `user_id` | User ID yang dicari (dari field `context.user_id`) |

**Query Parameters:**

| Param      | Default | Description                   |
| ---------- | ------- | ----------------------------- |
| `period`   | `30d`   | `7d`, `30d`, `90d`            |
| `category` | all     | `AUTH_EVENT`, `USER_ACTIVITY` |

**Response `200 OK`:**

```json
{
  "user_id": "usr_abc123",
  "period": "30d",
  "summary": {
    "total_events": 142,
    "auth_events": 31,
    "activity_events": 111,
    "apps_accessed": ["absensi-prod", "elearning-prod"],
    "last_seen": "2026-03-05T15:50:00Z",
    "auth_methods_used": ["google_oauth", "system_password"]
  },
  "recent_logins": [
    {
      "source_id": "absensi-prod",
      "auth_method": "google_oauth",
      "event_type": "login_success",
      "ip_address": "103.120.45.1",
      "browser": "Chrome",
      "os": "Windows 10",
      "created_at": "2026-03-05T08:00:00Z"
    }
  ],
  "recent_activity": [
    {
      "source_id": "absensi-prod",
      "action": "export",
      "resource_type": "attendance_report",
      "created_at": "2026-03-05T09:30:00Z"
    }
  ]
}
```

---

### `GET /api/activity/suspicious`

Daftar aktivitas yang terdeteksi mencurigakan (level ERROR pada AUTH_EVENT, atau `event_type: suspicious_login`).

**Auth**: JWT

**Query Parameters:**

| Param       | Default | Description      |
| ----------- | ------- | ---------------- |
| `source_id` | -       | Filter by app    |
| `period`    | `24h`   | `24h`, `7d`      |
| `page`      | `1`     | Nomor halaman    |
| `limit`     | `20`    | Item per halaman |

**Response `200 OK`:**

```json
{
  "data": [
    {
      "id": 5012,
      "source_id": "absensi-prod",
      "level": "ERROR",
      "message": "Suspicious login detected",
      "context": {
        "event_type": "suspicious_login",
        "auth_method": "google_oauth",
        "user_id": "usr_abc123",
        "email": "user@example.com",
        "ip_address": "185.220.101.1",
        "reason": "New IP address not seen before",
        "previous_ip": "103.120.45.1"
      },
      "created_at": "2026-03-05T03:12:00Z"
    }
  ],
  "pagination": { "total": 3, "page": 1, "limit": 20, "total_pages": 1 }
}
```

---

## 7. Statistics

### `GET /api/stats/overview`

Statistik agregat untuk dashboard overview.

**Auth**: JWT

**Query Parameters:**

| Param    | Default | Options            |
| -------- | ------- | ------------------ |
| `period` | `24h`   | `24h`, `7d`, `30d` |

**Response `200 OK`:**

```json
{
  "period": "24h",
  "summary": {
    "total_logs": 3521,
    "critical": 3,
    "errors": 47,
    "warnings": 412,
    "infos": 3059
  },
  "by_app": [
    {
      "source_id": "absensi-prod",
      "app_name": "Sistem Absensi Production",
      "critical": 1,
      "errors": 12,
      "warnings": 89,
      "infos": 240
    },
    {
      "source_id": "elearning-prod",
      "app_name": "E-Learning Platform",
      "critical": 2,
      "errors": 35,
      "warnings": 323,
      "infos": 2819
    }
  ],
  "trend": [
    { "timestamp": "2026-03-05T08:00:00Z", "count": 120 },
    { "timestamp": "2026-03-05T09:00:00Z", "count": 145 },
    { "timestamp": "2026-03-05T10:00:00Z", "count": 98 }
  ]
}
```

---

### `GET /api/stats/sources/:slug`

Statistik spesifik untuk satu aplikasi.

**Auth**: JWT

**Query Parameters:**

| Param    | Default | Options            |
| -------- | ------- | ------------------ |
| `period` | `24h`   | `24h`, `7d`, `30d` |

**Response `200 OK`:**

```json
{
  "source_id": "absensi-prod",
  "app_name": "Sistem Absensi Production",
  "period": "7d",
  "total_logs": 4230,
  "by_level": {
    "CRITICAL": 2,
    "ERROR": 45,
    "WARN": 318,
    "INFO": 3865
  },
  "by_category": {
    "SYSTEM_ERROR": 47,
    "AUTH_EVENT": 892,
    "USER_ACTIVITY": 3291
  },
  "trend_daily": [
    { "date": "2026-02-27", "count": 601 },
    { "date": "2026-02-28", "count": 589 }
  ],
  "top_errors": [
    {
      "message": "DB connection pool exhausted",
      "count": 23,
      "last_seen": "2026-03-05T14:22:00Z"
    }
  ]
}
```

---

---

## 6. Activity & Audit Trail

### `GET /api/activity`

**Auth**: JWT

**Query Parameters:**

| Param       | Default | Description                                     |
| ----------- | ------- | ----------------------------------------------- |
| `category`  | —       | `AUTH_EVENT`, `USER_ACTIVITY`, `AUDIT_TRAIL`    |
| `source_id` | —       | Filter by source UUID                           |
| `page`      | `1`     | Page number                                     |
| `limit`     | `20`    | Items per page                                  |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "items": [
      {
        "id": 101,
        "source_id": "uuid",
        "category": "AUTH_EVENT",
        "level": "INFO",
        "message": "Admin logged in successfully",
        "ip_address": "192.168.1.5",
        "context": { "email": "admin@onelog.com" },
        "created_at": "2026-03-18T09:00:00Z"
      }
    ],
    "meta": { "total": 500, "page": 1, "limit": 20 }
  }
}
```

### `GET /api/activity/summary`

**Auth**: JWT — Returns count breakdown by category/level.

### `GET /api/activity/users/:user_id`

**Auth**: JWT — All activity events for a specific user ID.

### `GET /api/activity/suspicious`

**Auth**: JWT — Events flagged as suspicious (brute force, anomaly).

---

## 7. APM — Performance Monitoring

### `GET /api/apm/endpoints`

**Auth**: JWT

**Query Parameters:**

| Param       | Default | Description            |
| ----------- | ------- | ---------------------- |
| `source_id` | —       | Filter by source UUID  |
| `period`    | `24h`   | `24h`, `7d`, `30d`     |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "endpoint": "/api/users",
      "method": "GET",
      "p50_ms": 45,
      "p95_ms": 120,
      "p99_ms": 340,
      "count": 4200
    }
  ]
}
```

### `GET /api/apm/timeline`

Time-series response time data for APM charts.

**Auth**: JWT

**Query Parameters:**

| Param       | Default | Description                          |
| ----------- | ------- | ------------------------------------ |
| `source_id` | —       | Filter by source UUID                |
| `endpoint`  | —       | Filter by specific endpoint path     |
| `period`    | `24h`   | `24h`, `7d`, `30d`                   |
| `interval`  | `1h`    | Time bucket interval: `1h`, `1d`     |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "timestamp": "2026-03-19T10:00:00Z",
      "request_count": 150,
      "avg_duration": 45.2,
      "p50": 30,
      "p95": 120,
      "p99": 250
    }
  ]
}
```

---

## 8. Issues — Error Grouping

### `GET /api/issues`

**Auth**: JWT

**Query Parameters:**

| Param       | Default | Options                         |
| ----------- | ------- | ------------------------------- |
| `status`    | —       | `OPEN`, `RESOLVED`, `IGNORED`   |
| `source_id` | —       | Filter by source UUID           |
| `page`      | `1`     |                                 |
| `limit`     | `20`    |                                 |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "items": [
      {
        "fingerprint": "sha256hash...",
        "source_id": "uuid",
        "category": "SYSTEM_ERROR",
        "level": "ERROR",
        "status": "OPEN",
        "message_sample": "connection refused at postgres:5432",
        "occurrence_count": 43,
        "first_seen_at": "2026-03-10T12:00:00Z",
        "last_seen_at": "2026-03-18T08:55:00Z"
      }
    ],
    "meta": { "total": 12, "page": 1, "limit": 20 }
  }
}
```

### `GET /api/issues/:fingerprint`

**Auth**: JWT — Full detail for one issue.

### `PATCH /api/issues/:fingerprint`

**Auth**: JWT — Update issue status.

**Request:**

```json
{ "status": "RESOLVED" }
```

### `GET /api/issues/:fingerprint/logs`

**Auth**: JWT — Individual log entries belonging to this issue (paginated).

### `GET /api/issues/analytics/trend`

Daily error rate trend for the last N days.

**Auth**: JWT

**Query Parameters:**

| Param       | Default | Description          |
| ----------- | ------- | -------------------- |
| `source_id` | —       | Filter by source     |
| `days`      | `30`    | Number of days (1-90)|

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "date": "2026-03-18",
      "total_logs": 5000,
      "error_count": 25,
      "error_rate": 0.5
    }
  ]
}
```

### `GET /api/issues/analytics/heatmap`

Error frequency heatmap by hour of day and day of week.

**Auth**: JWT

**Query Parameters:**

| Param       | Default | Description          |
| ----------- | ------- | -------------------- |
| `source_id` | —       | Filter by source     |
| `days`      | `30`    | Number of days (1-90)|

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "day_of_week": 1,
      "hour_of_day": 14,
      "error_count": 45
    }
  ]
}
```

---

## 9. Incidents — Downtime Tracking

Incident management for tracking source downtime and uptime.

### `GET /api/incidents`

List all incidents with pagination.

**Auth**: JWT

**Query Parameters:**

| Param       | Default | Options                |
| ----------- | ------- | ---------------------- |
| `status`    | —       | `OPEN`, `RESOLVED`     |
| `source_id` | —       | Filter by source       |
| `page`      | `1`     | Page number            |
| `limit`     | `20`    | Items per page         |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "source_id": "uuid",
        "status": "RESOLVED",
        "started_at": "2026-03-19T08:00:00Z",
        "resolved_at": "2026-03-19T08:15:00Z",
        "duration_sec": 900,
        "message": "Source 'Payment Gateway' is DOWN. Health check failed."
      }
    ],
    "meta": { "total": 5, "page": 1, "limit": 20 }
  }
}
```

### `GET /api/incidents/timeline`

Daily incident statistics for uptime reporting.

**Auth**: JWT

**Query Parameters:**

| Param       | Default | Description           |
| ----------- | ------- | --------------------- |
| `source_id` | —       | Filter by source      |
| `days`      | `30`    | Number of days (1-90) |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "date": "2026-03-18",
      "open_count": 1,
      "resolved_count": 2,
      "total_downtime_sec": 1800
    }
  ]
}
```

> **Note**: Incidents are automatically created when a source goes OFFLINE and resolved when it comes back ONLINE. Recovery emails are sent automatically.

---

## 10. Config Management

### `GET /api/config/:source_slug`

**Auth**: JWT — Returns all non-secret config values for a source.

### `PUT /api/config/:source_slug/:key`

**Auth**: JWT — Create or update a config key.

**Request:**

```json
{ "value": "my-value", "is_secret": false }
```

### `GET /api/config/:source_slug/history`

**Auth**: JWT — Returns version history for a source's config.

---

## 11. Logs — CSV Export

### `GET /api/logs/export`

**Auth**: JWT  
**Response Content-Type**: `text/csv`

**Query Parameters**: same as `GET /api/logs` (source_id, level, category).

---

## 12. Public Status Page

### `GET /api/status`

**Auth**: Public (no authentication required)

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "id": "uuid",
      "name": "Payment Gateway",
      "status": "ONLINE",
      "health_url": "https://api.example.com/health",
      "updated_at": "2026-03-18T09:00:00Z"
    }
  ]
}
```

---

## Error Response Format

Semua error menggunakan format konsisten:

```json
{
  "status": "error",
  "code": 400,
  "message": "Human-readable error message",
  "errors": [
    { "field": "email", "message": "invalid email format" }
  ]
}
```

### Standard Error Codes

| Code                  | HTTP Status | Meaning                              |
| --------------------- | ----------- | ------------------------------------ |
| `INVALID_CREDENTIALS` | 401         | Username/password salah              |
| `INVALID_API_KEY`     | 401         | API key tidak valid atau tidak aktif |
| `TOKEN_EXPIRED`       | 401         | JWT kadaluarsa                       |
| `FORBIDDEN`           | 403         | Tidak punya akses                    |
| `NOT_FOUND`           | 404         | Resource tidak ditemukan             |
| `INVALID_PAYLOAD`     | 400         | Validasi request gagal               |
| `SLUG_CONFLICT`       | 409         | Slug sudah dipakai                   |
| `RATE_LIMIT_EXCEEDED` | 429         | Terlalu banyak request               |
| `INTERNAL_ERROR`      | 500         | Error tidak terduga di server        |

---

---

## 13. AI Copilot — Chat

### POST /api/chat

Send a natural-language question to the AI Copilot. Requires JWT authentication.

**Auth**: `Authorization: Bearer <token>` (JWT)

**Request:**

```json
{
  "message": "How many CRITICAL errors do I have?"
}
```

**Response `200`:**

```json
{
  "status": "success",
  "message": "Reply generated",
  "data": {
    "reply": "## Critical Errors\n\nBased on the current snapshot, you have **3 CRITICAL** log entries..."
  }
}
```

The AI has access to:
- Total log count
- ERROR and CRITICAL log counts
- Number of open Issues
- Full project context (categories, levels, APM structure, etc.)

Responses are formatted in Markdown and rendered in the chat widget.

---

## Rate Limits Summary

| Endpoint                             | Limit   | Window                 | Status         | Algorithm       |
| ------------------------------------ | ------- | ---------------------- | -------------- | --------------- |
| `POST /api/ingest`                   | 100 req | Per menit, per API key | ✅ Implemented | Token Bucket    |
| `GET /api/logs`                      | 60 req  | Per menit, per JWT     | ✅ Implemented | Token Bucket    |
| `GET /api/stats/*`                   | 60 req  | Per menit, per JWT     | ✅ Implemented | Token Bucket    |
| `POST /api/auth/login`               | 10 req  | Per menit, per IP      | ✅ Implemented | Token Bucket    |
| `POST /api/sources/:slug/rotate-key` | 5 req   | Per jam, per JWT       | ✅ Implemented | Token Bucket    |
| Email Notifications                  | 1 email | Per 5 menit, per error | ✅ Implemented | Time-based      |

**Rate Limit Response (429 Too Many Requests):**

```json
{
  "status": "error",
  "code": 429,
  "message": "Rate limit exceeded. Please try again later.",
  "errors": {
    "retry_after_seconds": 30,
    "limit": 100,
    "window": "1m0s"
  }
}
```

**Headers:**
- `Retry-After`: Seconds to wait before retry
- `X-RateLimit-Limit`: Request limit
- `X-RateLimit-Window`: Time window

> **Algorithm**: Token Bucket - Smooth rate limiting with burst capability

