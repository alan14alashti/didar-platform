import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  // Allows GitHub Actions to build for a project page (e.g. /<repo-name>/)
  // while local dev and Docker/Nginx deployments keep serving from '/'.
  base: process.env.VITE_BASE_PATH || '/',
  plugins: [vue()],
})
