import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/

export default defineConfig({
  plugins: [react()],
  server: {
    host: true,
    allowedHosts: ['pop-os', '.ts.net'],
    proxy: {
      '/api': 'http://localhost:8080',
      '/stream': 'http://localhost:8080',
    },
  },
})
