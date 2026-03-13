# 🏗️ Backend Architectural Details

## Pattern: Clean Architecture / Hexagonal
Backend ULAM menggunakan pendekatan **Layered Architecture** untuk memisahkan logika bisnis dari detail infrastruktur.

### 1. Delivery Layer (Internal/Handler)
- Bertanggung jawab untuk HTTP Routing menggunakan **Gin**.
- Melakukan validasi request menggunakan **Gin Validator**.
- Mengubah HTTP request menjadi parameter untuk Service layer.

### 2. Business Logic Layer (Internal/Service)
- Inti dari aplikasi.
- Tidak tahu tentang database atau HTTP.
- Menangani alur kerja seperti: "Jika log level CRITICAL, kirim email dan minta insight AI".

### 3. Data Access Layer (Internal/Repository)
- Berinteraksi langsung dengan **GORM**.
- Menjalankan query database.

### 4. Domain Layer (Internal/Domain)
- Berisi entitas database dan struktur request/response global.

## AI Integration (Groq)
Service AI diletakkan di `pkg/ai` sebagai wrapper untuk Groq API. Alurnya dikontrol oleh Service layer saat log masuk.
