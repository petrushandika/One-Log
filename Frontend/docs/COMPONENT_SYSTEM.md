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
