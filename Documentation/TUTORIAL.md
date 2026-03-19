# 📖 One Log (ULAM) — Tutorial Lengkap

> **ULAM** — Unified Log & Activity Monitor  
> Platform observabilitas terpusat untuk memantau log, error, performa, dan aktivitas sistem dari berbagai aplikasi dalam satu dashboard.

---

## Daftar Isi

1. [Prerequisites](#1-prerequisites)
2. [Clone & Struktur Project](#2-clone--struktur-project)
3. [Setup Backend](#3-setup-backend)
4. [Setup Frontend](#4-setup-frontend)
5. [Menjalankan Aplikasi](#5-menjalankan-aplikasi)
6. [Login ke Dashboard](#6-login-ke-dashboard)
7. [Mendaftarkan Source (Aplikasi)](#7-mendaftarkan-source-aplikasi)
8. [Mengirim Log dari Aplikasi](#8-mengirim-log-dari-aplikasi)
9. [Fitur Dashboard](#9-fitur-dashboard)
   - [Overview](#91-overview)
   - [Log Explorer](#92-log-explorer)
   - [Issues](#93-issues)
   - [APM — Performa Endpoint](#94-apm--performa-endpoint)
   - [Incidents](#95-incidents)
   - [Audit Trail](#96-audit-trail)
   - [Status Page](#97-status-page)
   - [Config Manager](#98-config-manager)
   - [One Log AI (Chatbot)](#99-one-log-ai-chatbot)
10. [Contoh Integrasi Lengkap](#10-contoh-integrasi-lengkap)
11. [Troubleshooting](#11-troubleshooting)

---

## 1. Prerequisites

Pastikan semua tools berikut sudah terinstall di sistem kamu sebelum memulai.

| Tool | Versi Minimum | Cek |
|---|---|---|
| Go | 1.22+ | `go version` |
| Node.js | 20+ | `node --version` |
| npm | 9+ | `npm --version` |
| PostgreSQL | 15+ | `psql --version` |
| Git | - | `git --version` |

> **Opsional:** Docker & Docker Compose (untuk menjalankan database secara praktis)

### Install PostgreSQL via Docker (direkomendasikan)

Cara termudah menjalankan PostgreSQL tanpa install lokal:

```bash
docker run --name ulam-db \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=one-log \
  -p 5432:5432 \
  -d postgres:17-alpine
```

Atau gunakan `compose.yml` yang sudah disediakan:

```bash
docker compose up -d database
```

---

## 2. Clone & Struktur Project

```bash
git clone https://github.com/petrushandika/One-Log.git
cd One-Log
```

Struktur folder utama:

```
One-Log/
├── Backend/            ← Go API server
│   ├── cmd/api/        ← Entry point (main.go)
│   ├── internal/       ← Handler, Service, Repository, Domain
│   ├── migrations/     ← SQL migration files
│   ├── pkg/            ← Database, AI, Utils
│   └── .env.example    ← Template environment variables
├── Frontend/           ← React 19 dashboard
│   ├── src/pages/      ← Halaman aplikasi
│   └── src/shared/     ← Components, Layout, API client
├── Documentation/      ← Docs, Roadmap, API Reference
├── Makefile            ← Shortcut commands
└── compose.yml         ← Docker Compose
```

---

## 3. Setup Backend

### 3.1 Install Dependencies Go

```bash
cd Backend
go mod download
```

### 3.2 Buat File `.env`

Copy dari template yang sudah disediakan:

```bash
cp .env.example .env
```

Buka `.env` dan isi setiap variabel:

```env
# ── Server ──────────────────────────────────────
SERVER_PORT=8080
GIN_MODE=debug          # Ganti ke "release" di production

# ── Database ─────────────────────────────────────
DATABASE_URL=postgres://user:password@localhost:5432/one-log?sslmode=disable

# ── Security ─────────────────────────────────────
# Generate JWT_SECRET yang kuat:
#   openssl rand -hex 64
JWT_SECRET=isi_dengan_string_random_panjang_min_32_karakter
JWT_EXPIRY_HOURS=24
REFRESH_TOKEN_EXPIRY_DAYS=7

# ── Admin Credentials ────────────────────────────
ADMIN_EMAIL=admin@onelog.com
ADMIN_PASSWORD=ganti_password_kuat_disini
ADMIN_NAME=System Admin

# ── CORS (wajib sesuai URL frontend) ─────────────
CORS_ALLOWED_ORIGIN=http://localhost:5173

# ── Groq AI (untuk AI analysis & chatbot) ────────
GROQ_API_KEY=gsk_xxxxxxxxxxxx   # Dapatkan di console.groq.com (gratis)
AI_MODEL=llama-3.3-70b-versatile

# ── Log Retention ────────────────────────────────
LOG_RETENTION_DAYS=30

# ── Email / SMTP (opsional — untuk notifikasi) ───
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USER=email-kamu@gmail.com
MAIL_PASSWORD=app-password-gmail
ALERT_EMAIL=admin@yourdomain.com
```

> **Cara mendapatkan Groq API Key (gratis):**
> 1. Buka [console.groq.com](https://console.groq.com)
> 2. Daftar akun gratis
> 3. Buka **API Keys** → **Create API Key**
> 4. Copy key yang dihasilkan ke `GROQ_API_KEY`

### 3.3 Jalankan Database Migration

Perintah ini membuat semua tabel yang dibutuhkan di PostgreSQL:

```bash
# Dari folder Backend/
go run cmd/api/main.go --migrate
```

Output yang diharapkan:
```
Database connection established
Running AutoMigrate...
Migration completed successfully
```

### 3.4 Seed Admin User

Perintah ini membuat akun admin pertama berdasarkan `ADMIN_EMAIL` dan `ADMIN_PASSWORD` di `.env`:

```bash
go run cmd/api/main.go --seed
```

Output yang diharapkan:
```
Database seeding started...
Admin user created: admin@onelog.com
Seeding completed
```

> ⚠️ **Penting:** Jika kamu mengubah `ADMIN_EMAIL`/`ADMIN_PASSWORD` di `.env` setelah seed, jalankan `--seed` lagi. Admin lama tidak akan dihapus secara otomatis.

---

## 4. Setup Frontend

### 4.1 Install Dependencies Node.js

```bash
cd ../Frontend     # atau dari root: cd Frontend
npm install --legacy-peer-deps
```

### 4.2 Konfigurasi API URL (opsional)

Secara default frontend mengarah ke `http://localhost:8080`. Jika backend kamu berjalan di port/host berbeda, buat file `.env.local` di dalam folder `Frontend/`:

```env
VITE_API_BASE_URL=http://localhost:8080
```

---

## 5. Menjalankan Aplikasi

Buka **dua terminal** secara bersamaan.

### Terminal 1 — Backend

```bash
cd Backend
go run cmd/api/main.go
```

Output sukses:
```
Server starting on port 8080...
[GIN-debug] Listening and serving HTTP on :8080
```

### Terminal 2 — Frontend

```bash
cd Frontend
npm run dev
```

Output sukses:
```
  VITE v8.x.x  ready in 300 ms
  ➜  Local:   http://localhost:5173/
```

### Shortcut via Makefile (opsional)

```bash
# Dari root folder
make dev       # Jalankan backend + frontend sekaligus
make migrate   # Jalankan migration
make build     # Build production binary + frontend
```

### Cek Health Backend

Pastikan backend berjalan dengan membuka:

```
http://localhost:8080/health
```

Response:
```json
{ "status": "success", "data": { "app": "ULAM API" } }
```

---

## 6. Login ke Dashboard

1. Buka browser → [http://localhost:5173](http://localhost:5173)
2. Masukkan kredensial dari `.env`:
   - **Email**: nilai `ADMIN_EMAIL` (default: `admin@onelog.com`)
   - **Password**: nilai `ADMIN_PASSWORD`
3. Klik **Sign In**

Setelah login berhasil, kamu akan diarahkan ke halaman **Overview**.

---

## 7. Mendaftarkan Source (Aplikasi)

**Source** adalah aplikasi yang akan mengirim log ke One Log. Setiap source mendapat API key unik.

### Langkah-langkah:

1. Klik **Sources** di sidebar kiri
2. Klik tombol **Register Source** (pojok kanan atas)
3. Isi form:
   - **Source Name**: Nama aplikasimu (misal: `Payment Service`, `Auth API`)
   - **Health Check URL** (opsional): URL endpoint health check aplikasimu (misal: `https://api.example.com/health`)
4. Klik **Register**

### Menyimpan API Key

Setelah source terdaftar, **API key akan ditampilkan sekali saja** dengan format `ulam_live_xxxxxxxxxx`.

> ⚠️ **Sangat penting:** Copy API key ini sekarang menggunakan tombol **Copy** karena setelah kamu menyembunyikannya, key tidak bisa dilihat lagi secara penuh. Jika lupa, kamu harus **Rotate & Reveal** untuk mendapatkan key baru.

Tampilan key:
- **Visible**: `ulam_live_a1b2c3d4e5f6` (gunakan tombol mata 👁 untuk hide)
- **Hidden**: `ulam_live_••••••••••••`

### Merotasi API Key

Jika API key bocor atau kamu ingin menggantinya:

1. Pergi ke **Sources**
2. Di kartu source yang bersangkutan, klik **Rotate & Reveal**
3. Key baru akan ditampilkan — copy segera

---

## 8. Mengirim Log dari Aplikasi

Setelah punya API key, kamu bisa mulai mengirim log dari aplikasimu.

### Endpoint Ingestion

```
POST http://localhost:8080/api/ingest
Header: X-API-Key: <api_key_source_kamu>
Content-Type: application/json
```

### Struktur Payload

```json
{
  "category": "SYSTEM_ERROR",
  "level": "ERROR",
  "message": "Database connection timeout",
  "stack_trace": "Error: connect ETIMEDOUT\n  at TCPConnectWrap.afterConnect",
  "context": {
    "endpoint": "/api/users",
    "method": "GET",
    "user_id": "usr_123",
    "duration_ms": 5000
  }
}
```

### Field yang Tersedia

| Field | Wajib | Nilai yang Valid |
|---|---|---|
| `category` | ✅ | `SYSTEM_ERROR`, `AUTH_EVENT`, `USER_ACTIVITY`, `SECURITY`, `PERFORMANCE`, `AUDIT_TRAIL` |
| `level` | ✅ | `CRITICAL`, `ERROR`, `WARN`, `INFO`, `DEBUG` |
| `message` | ✅ | Teks bebas, maks 5000 karakter |
| `stack_trace` | ❌ | Stack trace error, maks 50000 karakter |
| `context` | ❌ | Object JSON bebas (maks 10 field) |
| `ip_address` | ❌ | IP address pengirim |

### Contoh Pengiriman

**Menggunakan `curl`:**

```bash
curl -X POST http://localhost:8080/api/ingest \
  -H "X-API-Key: ulam_live_a1b2c3d4e5f6" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "SYSTEM_ERROR",
    "level": "ERROR",
    "message": "Failed to connect to Redis",
    "stack_trace": "Error: connect ECONNREFUSED 127.0.0.1:6379",
    "context": { "service": "cache", "retry_count": 3 }
  }'
```

**Menggunakan JavaScript/TypeScript:**

```typescript
async function sendLog(apiKey: string, payload: object) {
  const response = await fetch('http://localhost:8080/api/ingest', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': apiKey,
    },
    body: JSON.stringify(payload),
  });
  return response.json();
}

// Contoh penggunaan
await sendLog('ulam_live_xxxxxxxxxx', {
  category: 'AUTH_EVENT',
  level: 'WARN',
  message: 'Failed login attempt',
  context: {
    email: 'user@example.com',
    ip_address: '203.0.113.42',
    event_type: 'login_failed',
    auth_method: 'system_password',
  },
});
```

**Menggunakan Go:**

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func sendLog(apiKey string, payload map[string]interface{}) error {
    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "http://localhost:8080/api/ingest", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", apiKey)
    
    client := &http.Client{}
    _, err := client.Do(req)
    return err
}

// Contoh
sendLog("ulam_live_xxxxxxxxxx", map[string]interface{}{
    "category": "PERFORMANCE",
    "level":    "WARN",
    "message":  "Slow database query detected",
    "context": map[string]interface{}{
        "endpoint":    "/api/orders",
        "duration_ms": 2500,
        "query":       "SELECT * FROM orders WHERE ...",
    },
})
```

**Menggunakan Python:**

```python
import requests

def send_log(api_key: str, payload: dict):
    response = requests.post(
        "http://localhost:8080/api/ingest",
        json=payload,
        headers={
            "X-API-Key": api_key,
            "Content-Type": "application/json",
        }
    )
    return response.json()

# Contoh
send_log("ulam_live_xxxxxxxxxx", {
    "category": "SECURITY",
    "level": "CRITICAL",
    "message": "Brute force attack detected",
    "context": {
        "ip_address": "198.51.100.23",
        "attempts": 50,
        "blocked": True,
    }
})
```

### Response Sukses

```json
{
  "status": "accepted",
  "message": "Log received and queued for processing"
}
```

---

## 9. Fitur Dashboard

### 9.1 Overview

**Navigasi:** Sidebar → **Overview**

Halaman utama yang menampilkan ringkasan kondisi sistem secara real-time:

- **Total Logs** — jumlah seluruh log yang tersimpan
- **Errors & Critical** — jumlah log dengan level ERROR atau CRITICAL
- **Online Sources** — jumlah source dengan status ONLINE
- **Security Alerts** — jumlah log kategori SECURITY

Tersedia juga:
- **Bar chart** distribusi log berdasarkan level (CRITICAL, ERROR, WARN, INFO)
- **Area chart** tren ingestion log sepanjang hari

Klik tombol **Refresh** (pojok kanan atas) untuk memperbarui data.

---

### 9.2 Log Explorer

**Navigasi:** Sidebar → **Log Explorer**

Temukan dan analisis log dengan filter yang lengkap.

#### Filter yang Tersedia

Klik tombol **Filters** untuk membuka panel filter:

| Filter | Keterangan |
|---|---|
| **Level** | Pilih CRITICAL, ERROR, WARN, INFO, atau DEBUG |
| **Source** | Filter berdasarkan aplikasi pengirim |
| **Category** | Filter berdasarkan kategori (System Error, Auth Event, dll.) |
| **From** | Waktu mulai (datetime picker) |
| **To** | Waktu akhir (datetime picker) |

#### Melihat Detail Log

Klik baris mana saja di tabel untuk membuka **Log Detail Modal** yang menampilkan:
- Message lengkap
- Timestamp, IP Address, Log ID, Fingerprint
- Stack trace (jika ada)
- Hasil analisis AI (jika sudah dianalisis)

#### Analisis AI pada Log

Dari dalam Log Detail Modal:
1. Klik tombol **Run AI Analysis** di bagian bawah
2. Tunggu beberapa detik
3. AI akan menghasilkan:
   - **Root Cause Analysis** — penyebab utama error
   - **Impact** — dampak yang mungkin terjadi
   - **Suggested Fix** — cara memperbaikinya

#### Export CSV

Klik tombol **Export CSV** (pojok kanan atas) untuk mengunduh log yang sedang difilter dalam format spreadsheet.

---

### 9.3 Issues

**Navigasi:** Sidebar → **Issues**

Log dengan error/fingerprint yang sama otomatis dikelompokkan menjadi satu **Issue**.

> **Fingerprint** dihitung dari: `SHA-256(source_id + message + stack_trace[:100])`  
> Artinya, setiap kali error yang sama terjadi, hitungan occurrences-nya bertambah — bukan membuat log baru yang terpisah.

#### Tab Issues

Menampilkan tabel semua issue dengan kolom:
- **Issue** — pesan error + fingerprint singkat + kategori
- **Level** — CRITICAL / ERROR / WARN
- **Status** — OPEN / RESOLVED / IGNORED
- **Occurrences** — berapa kali error ini terjadi
- **Last Seen** — kapan terakhir kali muncul

#### Aksi pada Issue

Klik baris issue untuk membuka detail, kemudian:

| Tombol | Fungsi |
|---|---|
| **Mark Resolved** | Tandai issue sudah diperbaiki |
| **Ignore** | Sembunyikan dari tampilan utama (false positive) |
| **Reopen** | Buka kembali issue yang sudah resolved/ignored |

#### Tab Analytics

Visualisasi agregat semua issues:
- **Level breakdown** — jumlah open issues per level (CRITICAL / ERROR / WARN)
- **Top 10 Most Frequent Errors** — error yang paling sering terjadi
- **Total Occurrences by Source** — source mana yang paling banyak error

---

### 9.4 APM — Performa Endpoint

**Navigasi:** Sidebar → **APM**

Pantau latensi endpoint API berdasarkan log dengan kategori `PERFORMANCE`.

> **Cara mengaktifkan APM:** Kirim log dengan:
> ```json
> {
>   "category": "PERFORMANCE",
>   "level": "INFO",
>   "message": "API request completed",
>   "context": {
>     "endpoint": "/api/orders",
>     "duration_ms": 145
>   }
> }
> ```
> Field `endpoint` dan `duration_ms` di dalam `context` **wajib ada** untuk APM.

#### Filter yang Tersedia

- **Period**: 1 Hour / 24 Hours / 7 Days
- **Source**: Filter per aplikasi

#### Metrics yang Ditampilkan

| Kolom | Keterangan |
|---|---|
| **Endpoint** | Nama path API |
| **Requests** | Jumlah request dalam periode |
| **P50** | Latensi median (50% request lebih cepat dari ini) |
| **P95** | 95% request selesai sebelum nilai ini |
| **P99** | 99% request selesai sebelum nilai ini |

**Kode warna latensi:**
- 🟢 `< 500ms` — Fast
- 🟡 `500–999ms` — Acceptable
- 🟠 `1000–1999ms` — Slow
- 🔴 `≥ 2000ms` — Critical

---

### 9.5 Incidents

**Navigasi:** Sidebar → **Reliability** → **Incidents**

Pantau dan tracking downtime sistem secara otomatis. Setiap kali source berstatus OFFLINE, sistem akan:

1. Membuat incident record otomatis
2. Mengirim notifikasi Email + Telegram
3. Menghitung downtime duration
4. Update status saat source kembali ONLINE

#### Stats Overview

- **Open Incidents** — Jumlah incident yang sedang berlangsung
- **Resolved (30d)** — Total incident yang sudah resolved dalam 30 hari
- **Total Downtime** — Total waktu downtime dalam 30 hari
- **Uptime %** — Persentase uptime dalam 30 hari terakhir

#### Timeline Chart

Visualisasi incident per hari (bar chart):
- **Blue bars** — Incident yang dibuka
- **Green bars** — Incident yang resolved

#### Incident Table

| Kolom | Keterangan |
|---|---|
| **Status** | OPEN atau RESOLVED dengan badge warna |
| **Message** | Deskripsi incident (misal: "Source X is DOWN") |
| **Started** — Waktu incident dimulai (relative time) |
| **Resolved** — Waktu incident selesai (jika resolved) |
| **Duration** — Durasi downtime dalam format human-readable |

> **Notifikasi:**
> - Email dan Telegram terkirim saat source DOWN
> - Email dan Telegram terkirim saat source kembali UP dengan informasi downtime duration

---

### 9.6 Audit Trail

**Navigasi:** Sidebar → **Compliance** → **Audit Trail**

Catatan aktivitas yang **tidak bisa dihapus** — cocok untuk keperluan compliance dan keamanan.

Log dengan kategori `AUDIT_TRAIL` akan muncul di halaman ini. Gunakan kolom `Search` untuk mencari berdasarkan message, kategori, atau IP address.

> **Cara mengirim log audit:**
> ```json
> {
>   "category": "AUDIT_TRAIL",
>   "level": "INFO",
>   "message": "User updated billing information",
>   "context": {
>     "actor_id": "usr_456",
>     "resource_type": "billing",
>     "action": "update",
>     "before": { "plan": "free" },
>     "after": { "plan": "pro" }
>   }
> }
> ```

---

### 9.7 Status Page

**Navigasi:** Sidebar → **Observe** → **Status**

Halaman monitoring uptime semua source yang terdaftar. Data **diperbarui otomatis setiap 60 detik**.

Status yang mungkin muncul:

| Status | Warna | Keterangan |
|---|---|---|
| **Online** | 🟢 Hijau | Source berjalan normal |
| **Degraded** | 🟡 Kuning | Response lambat atau sebagian tidak berfungsi |
| **Offline** | 🔴 Merah | Source tidak bisa diakses |
| **Maintenance** | ⚫ Abu | Source sedang dalam pemeliharaan |

> **Cara mengaktifkan health monitoring:** Saat mendaftarkan source, isi field **Health Check URL** dengan URL endpoint yang mengembalikan HTTP 200 (misal: `https://api.example.com/health`). Backend akan melakukan ping secara berkala ke URL tersebut.

---

### 9.8 Config Manager

**Navigasi:** Sidebar → **Manage** → **Config**

Simpan dan kelola konfigurasi per source secara terpusat — tanpa perlu deploy ulang.

#### Cara Menggunakan

1. Pilih **Source** dari dropdown di bagian atas
2. Tab **Config** menampilkan semua key/value config yang tersimpan
3. Nilai yang bertipe secret akan ditampilkan sebagai `••••••••` — klik **Reveal** untuk melihatnya

#### Menambah/Update Config

1. Klik **Edit** pada baris config yang ingin diubah
2. Ubah value-nya
3. Centang **Secret** jika nilainya sensitif (akan dienkripsi di database)
4. Klik **Save**

#### History & Rollback

Buka tab **History** untuk melihat semua perubahan config sebelumnya. Klik **Rollback** untuk mengembalikan ke versi tertentu.

---

### 9.8 One Log AI (Chatbot)

**Lokasi:** Ikon chat di pojok kanan bawah layar (selalu terlihat di semua halaman)

Chatbot AI yang mengetahui kondisi sistem kamu secara real-time dan bisa menjawab pertanyaan teknis.

#### Cara Memulai

1. Klik ikon chat **🟣** di pojok kanan bawah
2. Panel chat akan muncul di atas tombol tersebut
3. Ketik pertanyaanmu dan tekan **Enter**

#### Yang Bisa Ditanyakan

**Tentang sistem kamu:**
- _"How many critical errors do I have?"_
- _"What are the most frequent errors?"_
- _"Which endpoint has the highest P99 latency?"_
- _"Explain this error: connection refused ECONNREFUSED"_

**Tentang platform One Log:**
- _"What is the AUDIT_TRAIL category used for?"_
- _"How do I send a PERFORMANCE log?"_
- _"How does error fingerprinting work?"_

**Pertanyaan teknis umum (tetap relevan ke project):**
- _"How do I optimize slow PostgreSQL queries?"_
- _"Best practices for JWT authentication?"_
- _"How to set up Docker for a Go application?"_

> **Bilingual:** Chatbot akan menjawab dalam bahasa yang sama dengan pertanyaanmu. Tanya dalam Bahasa Indonesia → dijawab Bahasa Indonesia. Tanya dalam English → dijawab English.

#### Quick Questions

Saat pertama kali membuka chatbot, tersedia **4 pertanyaan cepat** yang bisa diklik langsung tanpa mengetik:
- How many ERROR logs do I have today?
- What is the AUDIT_TRAIL category?
- How do I send a PERFORMANCE log?
- How do I debug a SYSTEM_ERROR?

---

## 10. Contoh Integrasi Lengkap

Berikut contoh lengkap mengintegrasikan One Log ke aplikasi Node.js/Express:

```typescript
// logger.ts — Helper untuk mengirim log ke One Log
const ULAM_URL = process.env.ULAM_URL ?? 'http://localhost:8080/api/ingest';
const ULAM_API_KEY = process.env.ULAM_API_KEY ?? '';

type LogLevel = 'CRITICAL' | 'ERROR' | 'WARN' | 'INFO' | 'DEBUG';
type LogCategory = 'SYSTEM_ERROR' | 'AUTH_EVENT' | 'USER_ACTIVITY' | 'SECURITY' | 'PERFORMANCE' | 'AUDIT_TRAIL';

interface LogPayload {
  category: LogCategory;
  level: LogLevel;
  message: string;
  stack_trace?: string;
  context?: Record<string, unknown>;
  ip_address?: string;
}

export async function log(payload: LogPayload): Promise<void> {
  try {
    await fetch(ULAM_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': ULAM_API_KEY,
      },
      body: JSON.stringify(payload),
    });
  } catch (err) {
    // Jangan biarkan kegagalan logging mengganggu aplikasi utama
    console.error('[ULAM] Failed to send log:', err);
  }
}
```

```typescript
// app.ts — Contoh pemakaian di Express
import express from 'express';
import { log } from './logger';

const app = express();

// Middleware: log setiap request + duration (untuk APM)
app.use(async (req, res, next) => {
  const start = Date.now();
  res.on('finish', () => {
    const duration = Date.now() - start;
    log({
      category: 'PERFORMANCE',
      level: duration > 2000 ? 'WARN' : 'INFO',
      message: `${req.method} ${req.path} — ${res.statusCode}`,
      context: {
        endpoint: req.path,
        method: req.method,
        status_code: res.statusCode,
        duration_ms: duration,
      },
    });
  });
  next();
});

// Route dengan error handling
app.get('/api/users/:id', async (req, res) => {
  try {
    const user = await getUserFromDB(req.params.id);
    
    // Log aktivitas user
    await log({
      category: 'USER_ACTIVITY',
      level: 'INFO',
      message: `User profile viewed`,
      context: {
        action: 'view',
        resource_type: 'user',
        resource_id: req.params.id,
        actor_id: req.user?.id,
      },
    });
    
    res.json(user);
  } catch (err: any) {
    // Log error ke One Log
    await log({
      category: 'SYSTEM_ERROR',
      level: 'ERROR',
      message: err.message,
      stack_trace: err.stack,
      context: {
        endpoint: req.path,
        user_id: req.params.id,
      },
    });
    res.status(500).json({ error: 'Internal server error' });
  }
});

// Contoh: Log auth event
app.post('/api/login', async (req, res) => {
  const { email, password } = req.body;
  const ip = req.ip;
  
  try {
    const user = await authenticateUser(email, password);
    await log({
      category: 'AUTH_EVENT',
      level: 'INFO',
      message: 'Successful login',
      ip_address: ip,
      context: {
        event_type: 'login_success',
        auth_method: 'system_password',
        user_id: user.id,
      },
    });
    res.json({ token: generateToken(user) });
  } catch {
    await log({
      category: 'AUTH_EVENT',
      level: 'WARN',
      message: 'Failed login attempt',
      ip_address: ip,
      context: {
        event_type: 'login_failed',
        auth_method: 'system_password',
        email,
      },
    });
    res.status(401).json({ error: 'Invalid credentials' });
  }
});
```

---

## 11. Troubleshooting

### ❌ "CRITICAL SECURITY ERROR: ADMIN_EMAIL... is missing"

**Penyebab:** File `.env` tidak ditemukan atau variabel wajib kosong.

**Solusi:**
```bash
# Pastikan .env ada di folder Backend/
ls Backend/.env

# Atau buat dari template
cp Backend/.env.example Backend/.env
# Lalu isi ADMIN_EMAIL, ADMIN_PASSWORD, dan JWT_SECRET
```

---

### ❌ Login gagal / CORS error

**Penyebab:** `CORS_ALLOWED_ORIGIN` di `.env` tidak sesuai dengan URL frontend.

**Solusi:**
```env
# Di Backend/.env — harus persis sama dengan URL browser kamu
CORS_ALLOWED_ORIGIN=http://localhost:5173
```

Restart backend setelah mengubah `.env`.

---

### ❌ "Failed to update status" saat disable source

**Penyebab:** Method `PATCH` tidak diizinkan di CORS (sudah diperbaiki di versi terbaru).

**Solusi:** Pastikan kamu menggunakan versi terbaru dari repository.

---

### ❌ APM page kosong / "No performance data found"

**Penyebab:** Belum ada log dengan `category: PERFORMANCE` yang memiliki `context.duration_ms` dan `context.endpoint`.

**Solusi:** Kirim log seperti contoh berikut:
```bash
curl -X POST http://localhost:8080/api/ingest \
  -H "X-API-Key: <api_key_kamu>" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "PERFORMANCE",
    "level": "INFO",
    "message": "Test endpoint",
    "context": { "endpoint": "/api/test", "duration_ms": 200 }
  }'
```

---

### ❌ Chatbot error: "Connection failed. Make sure backend is running..."

**Penyebab:** `GROQ_API_KEY` tidak diset atau salah.

**Solusi:**
1. Buka [console.groq.com](https://console.groq.com) → buat API key
2. Tambahkan ke `Backend/.env`:
   ```env
   GROQ_API_KEY=gsk_xxxxxxxxxxxxxxxxxxxx
   ```
3. Restart backend

---

### ❌ "pq: password authentication failed for user"

**Penyebab:** Credential PostgreSQL di `DATABASE_URL` salah.

**Solusi:**
```env
# Sesuaikan dengan konfigurasi PostgreSQL kamu
DATABASE_URL=postgres://USER:PASSWORD@localhost:5432/DBNAME?sslmode=disable
```

---

### ❌ "port 8080 already in use"

**Solusi:**
```bash
# Cari dan kill proses yang menggunakan port 8080
lsof -ti:8080 | xargs kill -9

# Atau ganti port di .env
SERVER_PORT=8081
```

---

## Selamat! 🎉

Kamu sudah berhasil menjalankan **One Log (ULAM)** dan siap memantau sistem kamu secara terpusat. Jika ada pertanyaan, gunakan fitur **One Log AI** di dalam dashboard untuk mendapatkan bantuan langsung!
