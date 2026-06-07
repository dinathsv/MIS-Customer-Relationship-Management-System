import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  base: '/usermgmt/',
  plugins: [vue()],
  resolve: {
    alias: {
      '@': '/src'
    }
  }
})
