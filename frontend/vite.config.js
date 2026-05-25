import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8082', // Your Go backend port
        changeOrigin: true,
        ws: true, // CRITICAL: This allows WebSockets to pass through
      }
    }
  }
})