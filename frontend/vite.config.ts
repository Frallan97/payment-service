import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  base: '/', // Use absolute paths for assets (fixes SPA routing CSS/JS loading)
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
})
