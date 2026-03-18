Dokumen ini berisi aturan baku untuk menjaga konsistensi kode, struktur, dan desain di seluruh proyek **ULAM**. Semua kontributor **WAJIB** mengikuti panduan ini dengan memprioritaskan prinsip **Clean Code**.

---

## 🏗️ 1. Prinsip Utama: Clean Code & SOLID

Kami menganut prinsip Clean Code untuk memastikan kode mudah dibaca, diuji, dan dipelihara:

- **S - Single Responsibility**: Setiap fungsi atau class hanya boleh mengerjakan SATU hal dengan baik.
- **O - Open/Closed**: Kode harus terbuka untuk ekstensi tapi tertutup untuk modifikasi.
- **D - DRY (Don't Repeat Yourself)**: Hindari duplikasi logika. Gunakan utilitas atau service bersama.
- **K - KISS (Keep It Simple, Stupid)**: Jangan membuat solusi yang terlalu kompleks jika ada cara yang lebih sederhana.
- **Meaningful Names**: Nama variabel harus mendeskripsikan tujuan, bukan tipe data (contoh: `logEntry` bukan `le`).

---

## 📂 2. Aturan Penamaan (Naming Conventions)

### 🖥️ Frontend (React/TypeScript)
- **Folder & Files**: Gunakan `kebab-case` (contoh: `user-profile/`, `log-list.tsx`).
- **Components**: Gunakan `PascalCase` (contoh: `Sidebar.tsx`, `LogTable.tsx`).
- **Styles/CSS**: Gunakan `kebab-case` untuk class names.
- **Variables/Functions**: Gunakan `camelCase`.

### ⚙️ Backend (Golang)
- **Folder & Files**: Gunakan `snake_case` (contoh: `auth_service/`, `db_connection.go`).
- **Packages**: Gunakan huruf kecil semua, satu kata (contoh: `handler`, `repository`).
- **Exported (Public)**: Gunakan `PascalCase` (contoh: `func GetUser()`, `type User struct`).
- **Internal (Private)**: Gunakan `camelCase` (contoh: `func fetchConfig()`).

---

## 🏗️ 2. Arsitektur Folder

### Backend (Clean Architecture)
Struktur harus memisahkan detail teknis dari logika bisnis:
- **cmd/**: Titik masuk aplikasi. Tidak boleh ada logika bisnis di sini.
- **internal/domain/**: "Single Source of Truth" untuk model data dan interface.
- **internal/service/**: Tempat logika bisnis berada.
- **internal/handler/**: Tempat validasi request dan pengiriman response HTTP.

### Frontend (Feature-Based)
Struktur diatur berdasarkan fitur, bukan tipe file:
- **src/features/[feature-name]/**: Berisi semua yang diperlukan fitur tersebut (components, hooks, services).
- **src/shared/**: Asset atau komponen yang digunakan di lebih dari 2 fitur.

---

## 🛠️ 3. Quality Control
- **No Manual Format**: Gunakan `make format` (atau Prettier/GoFmt) sebelum commit.
- **Git Flow**: Push ke `development` untuk testing, `main` hanya melalui Pull Request (PR).
- **Commit Message**: Gunakan standar [Conventional Commits](https://www.conventionalcommits.org/) (contoh: `feat: add ai analysis engine`, `fix: resolve cors issue`).

---

## 🎨 4. Design System (Frontend)
ULAM menggunakan palet warna **"Midnight Obsidian"** yang premium dan modern.

| Token | Warna (Hex) | Kegunaan |
| :--- | :--- | :--- |
| Primary | `#6366F1` | Buttons, Links, Brand Identity (Indigo) |
| Success | `#10B981` | Valid logs, Positive actions |
| Warning | `#F59E0B` | Warning logs, Medium alerts |
| Danger | `#EF4444` | Error logs, Critical alerts |
| Background | `#0F172A` | Main background (Dark Blue/Grey) |
| Card/Surface | `#1E293B` | Components, Card backgrounds |

*Selengkapnya dapat dilihat di [Frontend/docs/DESIGN_SYSTEM.md](../Frontend/docs/DESIGN_SYSTEM.md).*

---

## 🛡️ 5. Backend Communication (API)

Semua API wajib mengembalikan format JSON yang terstruktur. Kita mengikuti standar **JSend-like** agar Frontend mudah mengolah response.

### Success Response (200, 201, 202)
```json
{
  "status": "success",
  "code": 200,
  "message": "Human readable message",
  "data": {
    "items": [],
    "total": 0
  }
}
```

### Error Response (400, 401, 403, 404, 500)
```json
{
  "status": "error",
  "code": 422,
  "message": "Validation failed",
  "errors": [
    {
      "field": "source_id",
      "message": "source_id is a required field"
    }
  ]
}
```

- **Rules**:
  - `status`: Hanya boleh `success` atau `error`.
  - `code`: Wajib menyertakan HTTP Status Code di dalam body untuk kemudahan debugging.
  - `message`: Kalimat singkat yang menjelaskan hasil operasi.
  - `data`: Objek data utama (Hanya ada jika status `success`).
  - `errors`: Array of objects yang mendetailkan kesalahan (Hanya ada jika status `error`).
