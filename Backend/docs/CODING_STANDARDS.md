# ⚙️ Backend Coding Standards

## 1. Go Project Standard
- Ikuti standar [Effective Go](https://go.dev/doc/effective_go).
- Gunakan `go fmt` dan `go vet` sebelum commit untuk memastikan kebersihan kode.

## 2. API Design & Communication
Aplikasi ini wajib menggunakan standar **JSend-like** untuk seluruh response.

### Struktur Response Utama
Selalu gunakan helper `pkg/utils/response.go` untuk konsistensi:

- **Success**: `{ "status": "success", "code": 200, "message": "...", "data": { ... } }`
- **Error**: `{ "status": "error", "code": 400, "message": "...", "errors": [ { "field": "...", "message": "..." } ] }`

### Aturan Tambahan
- **Status Codes**: 
  - `200 OK`: Request berhasil (GET, PUT).
  - `201 Created`: Resource baru berhasil dibuat.
  - `202 Accepted`: Batch process atau Log Ingestion diterima.
  - `400 Bad Request`: Error logika bisnis umum.
  - `422 Unprocessable Entity`: Error validasi input (body JSON).
- **Versioning**: Gunakan prefix `/api/v1/` untuk menjaga backward compatibility.
- **Validation**: Selalu validasi input menggunakan tag `binding` di level struct (Gin Validator).

## 3. Database Interactivity (GORM)
- **AutoMigrate**: Hanya jalankan di lingkungan development. Produksi harus menggunakan manual migration atau tool terkontrol.
- **Efficiency**: Selalu gunakan **Preloads** secara efisien untuk menghindari masalah query N+1.
- **Transaction**: Gunakan **Transaction** (`db.Transaction`) untuk operasi yang melibatkan lebih dari satu tabel untuk menjaga integritas data.

## 4. Error Handling
- Berikan error yang deskriptif namun tidak membocorkan detail teknis sistem (seperti stack trace internal DB) ke user.
- Hindari `panic()`, selalu tangani error dengan mengembalikan nilai `error`.

## 5. Logging & Observability
- Backend harus mampu me-log aktivitasnya sendiri secara internal menggunakan library logging standar atau terstruktur (seperti Zerolog atau Zap).
- Gunakan level log yang tepat (DEBUG, INFO, WARN, ERROR) agar mudah difilter.

## 6. Naming Consistency
- Database Columns: `snake_case` (contoh: `created_at`).
- JSON Keys: `snake_case` (contoh: `source_id`).
