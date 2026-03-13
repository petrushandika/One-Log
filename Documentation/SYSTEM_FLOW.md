# ⚡ System Flow

## Unified Log & Activity Monitor (ULAM)

---

## Overview

Dokumen ini menjelaskan bagaimana data mengalir dalam sistem ULAM dari berbagai sudut pandang.

---

## A. Alur Ingesti Log (Log Ingestion Flow)

Terjadi ketika project eksternal mendeteksi event dan mengirimkan log ke ULAM.

```
1. Project Client (e.g., Absensi) mendeteksi event
   (e.g., Google OAuth login, DB error, user action)
         │
         ▼
2. Client membuat payload JSON:
   {
     "source_id": "absensi-prod",
     "category": "AUTH_EVENT",
     "level": "INFO",
     "message": "User login via Google",
     "context": { "user_id": "usr_123", "method": "google_oauth" }
   }
         │
         ▼
3. Client kirim POST /v1/logs
   Header: X-API-Key: ulam_abcdef...
         │
         ▼
4. ULAM API — Middleware: Rate Limit Check
   └─ Jika melebihi limit → 429 Too Many Requests
         │
         ▼
5. ULAM API — Middleware: API Key Validation
   └─ Hash key → Query tabel sources
   └─ Jika tidak valid → 401 Unauthorized
         │
         ▼
6. ULAM API — Handler: Payload Validation
   └─ Check required fields (source_id, category, level, message)
   └─ Jika invalid → 400 Bad Request
          │
          ▼
7. ULAM API — Middleware: Data Masking (PII)
   └─ Scan field `message` & `context`
   └─ Replace sensitive keys with `***`
          │
          ├─────────────────────────────────────────────► 7a. Return 202 Accepted to Client
         │                                                   (< 100ms)
         ▼ (goroutine — async)
8. Background Worker: Save to Database
   └─ db.Create(&logEntry) via GORM → PostgreSQL
         │
         ▼
9. Check: level == ERROR || level == CRITICAL?
   ├─ NO  → Done
   └─ YES ▼
         │
10. Email Worker: Throttle Check
    └─ Key: "source_id:message_hash"
    └─ Cek apakah sudah ada email dalam 5 menit terakhir
    ├─ THROTTLED → Skip (Done)
    └─ OK ▼
         │
11. Email Worker: Send Notification
    └─ Render HTML template
    └─ smtp.SendMail() ke ALERT_EMAIL
    └─ Update throttle map timestamp

---

## F. Log Retention Flow (Cleanup)

Proyek dijalankan secara berkala untuk menjaga kebersihan data.

```
1. Cron Job: Triggered at 02:00 UTC Daily
          │
          ▼
2. Retention Worker: Build Delete Queries
   ├─ DELETE FROM log_entries WHERE level != 'CRITICAL' AND created_at < NOW() - 30 days
   └─ DELETE FROM log_entries WHERE level = 'CRITICAL'  AND created_at < NOW() - 90 days
          │
          ▼
3. PostgreSQL: Execute Batch Deletion
          │
          ▼
4. Worker: VACUUM (Optional) & Log Summary to stdout

---

## G. AI Insight Flow (via Groq API)

Bagaimana sistem memproses analisis log teknis menggunakan AI.

```
1. Trigger: Level == CRITICAL (Auto) atau Klik Tombol "Analyze" (Manual)
          │
          ▼
2. API Service: Collect Log Context
   └─ Ambil Message + Stack Trace + Context (User ID, Endpoint, dll.)
   └─ Gabungkan ke dalam prompt template
          │
          ▼
3. Groq Worker: Call Groq API
   └─ Endpoint: https://api.groq.com/openai/v1/chat/completions
   └─ Model: llama-3.3-70b-versatile
   └─ Header: Authorization: Bearer <GROQ_API_KEY>
          │
          ▼
4. Groq Response Processing
   └─ Parse JSON response (Summary, RCA, Solution)
   └─ Jika gagal → Return error / Silently skip
          │
          ▼
5. Persistence & Delivery
   └─ Update record log_entries (kolom: ai_insight)
   └─ (Jika Auto) Sertakan insight dalam Email Notification
          │
          ▼
6. Dashboard: Display Result
   └─ Tampilkan markdown formatted insight pada detail log
```
```
```

**Karakteristik penting:**

- Step 7a (Response ke client) terjadi **sebelum** DB write selesai → Non-blocking
- Email dikirim dalam **goroutine terpisah** → Tidak memblock proses lain
- Jika DB write gagal → Log error ke stdout (tidak retry di MVP)

---

## B. Alur Monitoring Dashboard

Terjadi ketika admin membuka dashboard React untuk melihat log.

```
1. Admin buka browser → https://dashboard.ulam.dev
         │
         ▼
2. React App load → Check `ulam_access` cookie (httpOnly, dikirim otomatis oleh browser)
   ├─ Token tidak ada → Redirect ke /login
   └─ Token ada → Continue
         │
         ▼
3. Login Flow (jika belum login):
   ├─ Admin input username + password
   ├─ POST /v1/auth/login
   ├─ API verify credentials → Return JWT
   └─ Server set httpOnly cookie `ulam_access` (24h) + `ulam_refresh` (7d) → Redirect ke dashboard
         │
         ▼
4. Dashboard Load:
   React fetch parallel:
   ├─ GET /v1/stats/overview?period=24h → Overview cards
   └─ GET /v1/logs?page=1&limit=20     → Log table
         │
         ▼
5. User Apply Filter:
   User pilih: project=absensi-prod, level=ERROR, date=today
         │
         ▼
6. React kirim request:
   GET /v1/logs?source_id=absensi-prod&level=ERROR&from=2026-03-05
   Header: Authorization: Bearer <jwt>
         │
         ▼
7. ULAM API — Validate JWT
   └─ Jika expired/invalid → 401 → React redirect ke login
         │
         ▼
8. ULAM API — Build Query:
   SELECT * FROM log_entries
   WHERE source_id = 'absensi-prod'
     AND level = 'ERROR'
     AND created_at >= '2026-03-05'
   ORDER BY created_at DESC
   LIMIT 20 OFFSET 0
         │
         ▼
9. PostgreSQL execute query → Return rows
         │
         ▼
10. API serialize ke JSON → Return response
    {
      "data": [...],
      "pagination": { "total": 47, "page": 1, ... }
    }
         │
         ▼
11. React renders LogTable dengan data
    User bisa klik row → Fetch detail LOG → Show modal
```

---

## C. Alur Integrasi Client Project

Cara project eksternal mengintegrasikan diri dengan ULAM.

```
PROJECT SETUP (One-time):
━━━━━━━━━━━━━━━━━━━━━━━━

Admin → POST /api/sources { name: "Absensi", slug: "absensi-prod" }
     → Dapat API Key: "ulam_abc123..."
     → Simpan ke .env project:
          ULAM_API_KEY=ulam_abc123...
          ULAM_ENDPOINT=https://api.ulam.dev/v1/logs
          ULAM_source_id=absensi-prod

RUNTIME INTEGRATION:
━━━━━━━━━━━━━━━━━━━

// Contoh helper function di project Go
func LogToULAM(level, category, message, stackTrace string, context map[string]interface{}) {
    go func() {  // Fire and forget
        payload := map[string]interface{}{
            "source_id": os.Getenv("ULAM_source_id"),
            "category":   category,
            "level":      level,
            "message":    message,
            "stack_trace": stackTrace,
            "context":    context,
        }
        body, _ := json.Marshal(payload)
        req, _ := http.NewRequest("POST", os.Getenv("ULAM_ENDPOINT"), bytes.NewBuffer(body))
        req.Header.Set("X-API-Key", os.Getenv("ULAM_API_KEY"))
        req.Header.Set("Content-Type", "application/json")
        http.DefaultClient.Do(req)  // Ignore error (best effort)
    }()
}

// Penggunaan:
// Saat login Google berhasil:
LogToULAM("INFO", "AUTH_EVENT", "User login via Google", "", map[string]interface{}{
    "user_id": user.ID,
    "email":   user.Email,
    "method":  "google_oauth",
})

// Saat terjadi database error:
LogToULAM("ERROR", "SYSTEM_ERROR", "DB connection failed", err.Error(), map[string]interface{}{
    "endpoint": r.URL.Path,
})
```

---

## D. Alur Email Notifikasi

Detail bagaimana email dikirim dan throttling bekerja.

```
Log masuk dengan level ERROR/CRITICAL
         │
         ▼
Generate throttle key:
key = SHA256(source_id + ":" + message[:50])
         │
         ▼
Check throttle map:
throttleMap[key] → last_sent_time
         │
         ├─ Tidak ada / > 5 menit → KIRIM EMAIL
         │         │
         │         ▼
         │   Render HTML Template:
         │   - Subject: [ULAM] 🔴 ERROR di absensi-prod
         │   - Body:
         │     Project  : Sistem Absensi Production
         │     Level    : ERROR
         │     Message  : Failed to connect to database
         │     Time     : 2026-03-05 10:30:00 UTC
         │     Stack    : goroutine 1 [running]...
         │     Link     : https://dashboard.ulam.dev/logs/42
         │         │
         │         ▼
         │   smtp.SendMail(
         │       host:SMTP_HOST, port:SMTP_PORT,
         │       from:SMTP_USER, to:ALERT_EMAIL,
         │       body:htmlTemplate
         │   )
         │         │
         │         ▼
         │   Update throttleMap[key] = time.Now()
         │
         └─ ≤ 5 menit → SKIP (tidak kirim, hindari spam)
```

---

## E. Error Handling Flow

Bagaimana sistem menangani berbagai kondisi gagal.

| Kondisi             | Handling                                  | Impact                                |
| ------------------- | ----------------------------------------- | ------------------------------------- |
| Invalid API Key     | Return 401 langsung                       | Client tahu token salah               |
| Payload invalid     | Return 400 + details                      | Client tahu field mana yang salah     |
| DB write gagal      | Log error ke stdout                       | Log hilang, client tidak tahu         |
| Email gagal         | Log error ke stdout, skip throttle update | Email tidak terkirim, retry next time |
| DB down (dashboard) | Return 500                                | Dashboard tidak bisa load data        |
| JWT expired         | Return 401                                | Frontend redirect ke login            |

> **MVP Note**: DB write failure tidak di-retry. Jika reliability dibutuhkan, pertimbangkan Redis queue di Post-MVP.

---

## Sequence Diagram — Happy Path

```
Client          ULAM API        PostgreSQL      EmailSMTP
  │                │                │               │
  │  POST /v1/logs │                │               │
  │───────────────►│                │               │
  │                │ Validate Token │               │
  │                │───────────────►│               │
  │                │◄───────────────│               │
  │  202 Accepted  │                │               │
  │◄───────────────│                │               │
  │                │                │               │
  │            [goroutine]          │               │
  │                │ db.Create()    │               │
  │                │───────────────►│               │
  │                │◄───────────────│               │
  │                │                │               │
  │           [if ERROR]            │               │
  │                │ smtp.Send()    │               │
  │                │──────────────────────────────► │
  │                │◄───────────────────────────────│
```
