import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    // Raise the chunk size warning threshold (react-markdown + recharts make the
    // bundle legitimately large; consider lazy-loading pages in the future).
    chunkSizeWarningLimit: 1500,
  },
})
