# 🗄️ Data Schema

## Unified Log & Activity Monitor (ULAM)

---

## Schema Overview

ULAM menggunakan **2 tabel utama** di PostgreSQL:

```text
┌─────────────────────┐         ┌──────────────────────────┐
│       sources        │         │       log_entries         │
├─────────────────────┤         ├──────────────────────────┤
│ id (PK)             │◄────────│ source_id (FK → slug)   │
│ name                │         │ id (PK)                  │
│ slug (UNIQUE)       │         │ category                 │
│ api_key (UNIQUE)    │         │ level                    │
│ is_active           │         │ message                  │
│ created_at          │         │ stack_trace              │
└─────────────────────┘         │ context (JSONB)          │
                                │ ip_address               │
                                │ created_at               │
                                └──────────────────────────┘
```

---

## Table: `sources`

```sql
CREATE TABLE sources (
    id         SERIAL          PRIMARY KEY,
    name       VARCHAR(100)    NOT NULL,
    slug       VARCHAR(50)     UNIQUE NOT NULL,
    api_key    VARCHAR(64)     UNIQUE NOT NULL,
    is_active  BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

### GORM Model

```go
type Project struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"            json:"id"`
    Name      string    `gorm:"type:varchar(100);not null"          json:"name"`
    Slug      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
    APIKey    string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"-"`
    IsActive  bool      `gorm:"default:true"                        json:"is_active"`
    CreatedAt time.Time `gorm:"autoCreateTime"                      json:"created_at"`
}
```

### Field Descriptions

| Field        | Type         | Description                                         |
| ------------ | ------------ | --------------------------------------------------- |
| `id`         | INT          | Auto-increment primary key                          |
| `name`       | VARCHAR(100) | Nama project human-readable, e.g., "Sistem Absensi" |
| `slug`       | VARCHAR(50)  | Identifier unik, e.g., "absensi-prod"               |
| `api_key`    | VARCHAR(64)  | Hashed API key (bcrypt)                             |
| `is_active`  | BOOLEAN      | Soft-disable tanpa delete data                      |
| `created_at` | TIMESTAMP    | Waktu project didaftarkan (UTC)                     |

---

## Table: `log_entries`

```sql
CREATE TABLE log_entries (
    id          BIGSERIAL    PRIMARY KEY,
    source_id  VARCHAR(50)  NOT NULL,
    category    VARCHAR(30)  NOT NULL,
    level       VARCHAR(20)  NOT NULL,
    message     TEXT         NOT NULL,
    stack_trace TEXT,
    context     JSONB,
    ip_address  VARCHAR(45),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_log_source_id        ON log_entries (source_id);
CREATE INDEX idx_log_level             ON log_entries (level);
CREATE INDEX idx_log_created_at        ON log_entries (created_at DESC);
CREATE INDEX idx_log_source_level_date ON log_entries (source_id, level, created_at DESC);
CREATE INDEX idx_log_context_gin       ON log_entries USING GIN (context);
```

### GORM Model — LogEntry

```go
type LogEntry struct {
    ID         uint           `gorm:"primaryKey;autoIncrement"         json:"id"`
    SourceID  string         `gorm:"index;type:varchar(50);not null"  json:"source_id"`
    Category   string         `gorm:"index;type:varchar(30);not null"  json:"category"`
    Level      string         `gorm:"type:varchar(20);not null"        json:"level"`
    Message    string         `gorm:"type:text;not null"               json:"message"`
    StackTrace string         `gorm:"type:text"                        json:"stack_trace,omitempty"`
    Context    datatypes.JSON `gorm:"type:jsonb"                       json:"context,omitempty"`
    IPAddress  string         `gorm:"type:varchar(45)"                 json:"ip_address,omitempty"`
    CreatedAt  time.Time      `gorm:"autoCreateTime"                   json:"created_at"`
}
```

---

## Enum Values

### Category

| Value           | Description                          |
| --------------- | ------------------------------------ |
| `SYSTEM_ERROR`  | Error teknis (DB fail, 500 errors)   |
| `USER_ACTIVITY` | Aktivitas pengguna (page view, CRUD) |
| `AUTH_EVENT`    | Login, logout, OAuth callback        |

### Level

| Value      | Severity                        | Triggers Email? |
| ---------- | ------------------------------- | --------------- |
| `CRITICAL` | 🔴 Sistem tidak bisa operasi    | ✅ Yes          |
| `ERROR`    | 🟠 Fungsi gagal, user terdampak | ✅ Yes          |
| `WARN`     | 🟡 Potensi masalah              | ❌ No           |
| `INFO`     | 🟢 Operasi normal               | ❌ No           |

---

## Context JSONB — Standardized Schemas

Field `context` adalah JSONB sehingga bebas, namun berikut adalah **kontrak yang direkomendasikan** agar data konsisten dan bisa divisualisasikan dengan baik di dashboard.

---

### 📋 AUTH_EVENT Context

Digunakan untuk semua event autentikasi. Field `event_type` dan `auth_method` adalah **wajib** untuk kategori ini.

**Login berhasil (Google OAuth):**

```json
{
  "event_type": "login_success",
  "auth_method": "google_oauth",
  "user_id": "usr_abc123",
  "email": "user@example.com",
  "name": "John Doe",
  "provider_id": "google_uid_xyz789",
  "ip_address": "103.120.45.1",
  "browser": "Chrome",
  "browser_version": "121.0.0.0",
  "os": "Windows 10",
  "device_type": "desktop",
  "session_id": "sess_xyz789",
  "is_new_user": false
}
```

**Login berhasil (Manual / System Password):**

```json
{
  "event_type": "login_success",
  "auth_method": "system_password",
  "user_id": "usr_def456",
  "email": "admin@company.com",
  "role": "admin",
  "ip_address": "192.168.1.5",
  "browser": "Firefox",
  "os": "Ubuntu 22.04",
  "device_type": "desktop",
  "session_id": "sess_abc123"
}
```

**Login gagal:**

```json
{
  "event_type": "login_failed",
  "auth_method": "system_password",
  "attempted_email": "user@example.com",
  "failure_reason": "invalid_password",
  "attempt_count": 3,
  "ip_address": "45.67.89.10",
  "browser": "Chrome",
  "os": "Android 14",
  "device_type": "mobile"
}
```

**OAuth social login (GitHub, Facebook, Twitter/X, Discord):**

```json
{
  "event_type": "login_success",
  "auth_method": "github_oauth",
  "user_id": "usr_ghi789",
  "email": "dev@email.com",
  "name": "Jane Dev",
  "provider_id": "github_uid_12345",
  "provider_username": "janedev",
  "provider_avatar": "https://avatars.githubusercontent.com/...",
  "ip_address": "103.120.45.2",
  "browser": "Safari",
  "os": "macOS Sonoma",
  "device_type": "desktop",
  "session_id": "sess_mno456",
  "is_new_user": true
}
```

**Suspicious login (login dari lokasi baru):**

```json
{
  "event_type": "suspicious_login",
  "auth_method": "google_oauth",
  "user_id": "usr_abc123",
  "email": "user@example.com",
  "ip_address": "185.220.101.1",
  "browser": "Chrome",
  "os": "Windows 10",
  "device_type": "desktop",
  "reason": "New IP address not seen before",
  "previous_ip": "103.120.45.1",
  "session_id": "sess_pqr789"
}
```

**Enum values untuk `auth_method`:**

| Value             | Provider              |
| ----------------- | --------------------- |
| `google_oauth`    | Google                |
| `github_oauth`    | GitHub                |
| `facebook_oauth`  | Facebook              |
| `twitter_oauth`   | Twitter/X             |
| `discord_oauth`   | Discord               |
| `system_password` | Manual login          |
| `magic_link`      | Email magic link      |
| `sso`             | Enterprise SSO / SAML |
| `api_token`       | API token (non-user)  |

**Enum values untuk `event_type`:**

| Value              | Level Rekomendasi                   |
| ------------------ | ----------------------------------- |
| `login_success`    | INFO                                |
| `login_failed`     | WARN (ERROR jika > 5x dari IP sama) |
| `logout`           | INFO                                |
| `token_refresh`    | INFO                                |
| `oauth_callback`   | INFO                                |
| `account_linked`   | INFO                                |
| `password_reset`   | WARN                                |
| `mfa_challenge`    | INFO                                |
| `mfa_success`      | INFO                                |
| `session_expired`  | INFO                                |
| `suspicious_login` | ERROR                               |

---

### 🚶 USER_ACTIVITY Context

Digunakan untuk mencatat aktivitas pengguna di dalam aplikasi.

**Export data:**

```json
{
  "action": "export",
  "user_id": "usr_abc123",
  "resource_type": "attendance_report",
  "resource_id": "rpt_march_2026",
  "metadata": {
    "format": "xlsx",
    "row_count": 245,
    "filters": { "month": "2026-03", "department": "Engineering" }
  },
  "session_id": "sess_xyz789"
}
```

**CRUD operation:**

```json
{
  "action": "delete",
  "user_id": "usr_abc123",
  "resource_type": "course",
  "resource_id": "course_456",
  "metadata": {
    "course_title": "Introduction to Go",
    "enrolled_students": 38
  },
  "session_id": "sess_xyz789"
}
```

**Enum values untuk `action`:**

| Value       | Deskripsi           |
| ----------- | ------------------- |
| `page_view` | Mengunjungi halaman |
| `create`    | Membuat data baru   |
| `update`    | Mengubah data       |
| `delete`    | Menghapus data      |
| `export`    | Mengekspor data     |
| `download`  | Mengunduh file      |
| `search`    | Melakukan pencarian |
| `share`     | Berbagi konten      |

---

### ⚙️ SYSTEM_ERROR Context

```json
{
  "endpoint": "/api/attendance/check-in",
  "http_method": "POST",
  "http_status": 500,
  "user_id": "usr_abc123",
  "request_id": "req_def456",
  "db_error": "connection refused",
  "duration_ms": 5023
}
```

### 🐢 PERFORMANCE Context

```json
{
  "endpoint": "/api/reports/generate",
  "duration_ms": 8500,
  "threshold_ms": 3000,
  "query_count": 47,
  "slowest_query_ms": 4200,
  "user_id": "usr_abc123"
}
```

### 🔒 SECURITY Context

```json
{
  "event": "rate_limit_exceeded",
  "ip_address": "185.220.101.1",
  "endpoint": "/api/ingest",
  "request_count": 150,
  "window_seconds": 60,
  "blocked": true
}
```

---

## Sample Queries

### Filter log per app + level

```sql
SELECT * FROM log_entries
WHERE source_id = 'absensi-prod'
  AND level IN ('ERROR', 'CRITICAL')
  AND created_at BETWEEN '2026-03-01' AND '2026-03-31'
ORDER BY created_at DESC
LIMIT 20;
```

### Overview: Error count per app

```sql
SELECT
    source_id,
    COUNT(*) FILTER (WHERE level = 'CRITICAL') AS critical_count,
    COUNT(*) FILTER (WHERE level = 'ERROR')    AS error_count,
    COUNT(*) FILTER (WHERE level = 'WARN')     AS warn_count,
    COUNT(*) FILTER (WHERE level = 'INFO')     AS info_count
FROM log_entries
WHERE created_at >= NOW() - INTERVAL '24 hours'
GROUP BY source_id
ORDER BY critical_count DESC;
```

### Activity: Login events per auth method (last 7 days)

```sql
SELECT
    context->>'auth_method' AS auth_method,
    context->>'event_type'  AS event_type,
    COUNT(*)                AS total
FROM log_entries
WHERE category = 'AUTH_EVENT'
  AND created_at >= NOW() - INTERVAL '7 days'
GROUP BY context->>'auth_method', context->>'event_type'
ORDER BY total DESC;
```

### Activity: Semua login dari satu user (lintas semua app)

```sql
SELECT
    source_id,
    context->>'auth_method' AS method,
    context->>'event_type'  AS event,
    ip_address,
    created_at
FROM log_entries
WHERE category = 'AUTH_EVENT'
  AND context->>'user_id' = 'usr_abc123'
ORDER BY created_at DESC;
```

### Activity: Deteksi brute force (> 5 login failed dari 1 IP dalam 10 menit)

```sql
SELECT
    ip_address,
    COUNT(*) AS failed_attempts,
    MIN(created_at) AS first_attempt,
    MAX(created_at) AS last_attempt
FROM log_entries
WHERE category = 'AUTH_EVENT'
  AND context->>'event_type' = 'login_failed'
  AND created_at >= NOW() - INTERVAL '10 minutes'
GROUP BY ip_address
HAVING COUNT(*) > 5
ORDER BY failed_attempts DESC;
```

### Activity: Top active users per app

```sql
SELECT
    source_id,
    context->>'user_id' AS user_id,
    COUNT(*) AS activity_count
FROM log_entries
WHERE category = 'USER_ACTIVITY'
  AND created_at >= NOW() - INTERVAL '24 hours'
GROUP BY source_id, context->>'user_id'
ORDER BY activity_count DESC
LIMIT 20;
```

---

## Auto-Migration

```go
func ConnectDB(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    err = db.AutoMigrate(&models.Source{}, &models.LogEntry{})
    return db, err
}
```
