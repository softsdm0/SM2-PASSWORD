import { fileURLToPath, URL } from 'url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { nodePolyfills } from 'vite-plugin-node-polyfills'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  build: {
    outDir: 'dist',
  },
  plugins: [
    vue(),
    vueDevTools(),
    nodePolyfills({
      globals: {
        Buffer: true,
      },
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      // '@': path.resolve(__dirname, 'src'),
    },
  },
  server: {
    host: '0.0.0.0', // 添加bind地址
    port: 5173,
    proxy: {
      '/public/api': {
        target: 'http://10.0.0.100:30005',
        changeOrigin: true,
      },
      // casdoor
      '/api': {
        target: 'http://10.0.0.200:30003',
        changeOrigin: true,
      },
      '/static': {
        target: 'http://10.0.0.200:30003',
        changeOrigin: true,
      },
      '/login/oauth': {
        target: 'http://10.0.0.200:30003',
        changeOrigin: true,
      },
      '/files/resource': {
        target: 'http://10.0.0.200:30003',
        changeOrigin: true,
      },
    },
  },
})
