# 🚀 Backend Detailed Setup Guide

## Lingkungan Pengembangan
1. **Instalasi Go**: Pastikan menggunakan Go 1.26+.
2. **PostgreSQL**: Siapkan database kosong bernama `ulam_db`.

## Langkah-langkah
1. **Clone & Install**:
   ```bash
   go mod tidy
   ```
2. **Konfigurasi Environment**:
   Salin `.env.example` ke `.env` dan sesuaikan kredensial database & API Key.
3. **Migrasi Database**:
   ```bash
   go run cmd/api/main.go --migrate
   ```
4. **Menjalankan Server**:
   ```bash
   go run cmd/api/main.go
   ```

## Tips Debugging
- Gunakan `GIN_MODE=debug` untuk melihat log HTTP secara detail.
- Cek tabel `log_entries` menggunakan tool seperti TablePlus/DBeaver untuk memastikan data masuk.
