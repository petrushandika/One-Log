# 🛠️ Frontend Development Workflow

## Standar Kode
Kami menggunakan **Prettier** dan **ESLint** untuk menjaga kerapihan. Jangan lewatkan proses Lint agar build tidak gagal di CI/CD.

## Menjalankan Proyek
1. **Instalasi**:
   ```bash
   npm install --legacy-peer-deps
   ```
2. **Dev Server**:
   ```bash
   npm run dev
   ```
3. **Lint & Format**:
   ```bash
   npm run lint
   npm run format
   ```

## Branching & Commit
Gunakan **Husky** yang sudah terinstal. Jika ada error lint saat commit, perbaiki terlebih dahulu sebelum mencoba commit ulang.
