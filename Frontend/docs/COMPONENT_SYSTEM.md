# 🎨 Frontend Component System

## Framework & Style
- **React 19**: Menggunakan fitur terbaru seperti Actions dan Improved Hooks.
- **Tailwind CSS v4**: Menggunakan CSS variable-based configuration.

## Folder Structure (Feature-Based)
Setiap fitur besar (seperti `logs`) memiliki sub-folder:
- `components/`: Komponen UI spesifik fitur (misal: `LogDetailCard.tsx`).
- `hooks/`: Custom hooks untuk fitur tersebut (misal: `useLogAI.ts`).
- `services/`: API calls spesifik fitur menggunakan Axios.

## Shared Components
Komponen atom seperti Button, Input, dan Modal diletakkan di `src/shared/components/ui`.

## State Management
- **Server State**: Menggunakan **TanStack Query (React Query)** untuk caching data dari API.
- **Global UI State**: Menggunakan **Zustand** untuk hal-hal simpel seperti Sidebar open/close atau User Session.

---

## 📡 API Interaction Standards

Semua interaksi dengan Backend wajib mengikuti struktur **JSend-like**:
- **Response Format**: Selalu periksa properti \`status\` dan \`code\` di level root.
- **Success Mapping**: Data utama selalu berada di dalam properti \`data\`.
- **Error Handling**: 
  - Jika \`status === "error"\`, gunakan properti \`message\` untuk user-facing alert (seperti Toast).
  - Gunakan properti \`errors[]\` untuk memetakan error spesifik ke field form.
- **Axios Instance**: Gunakan interceptor untuk standarisasi penanganan error 401 (Unauthorized) dan 422 (Validation Error).
