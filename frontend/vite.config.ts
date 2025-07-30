import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import Unocss from 'unocss/vite'
import Components from 'unplugin-vue-components/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'
import { fileURLToPath, URL } from 'node:url'
import svgLoader from 'vite-svg-loader'
import AutoImport from 'unplugin-auto-import/vite'
import vuetify from 'vite-plugin-vuetify'

// https://vite.dev/config/
export default defineConfig({
  build: {
    outDir: 'dist',
    rollupOptions: {
      output: {
        advancedChunks: {
          groups: [
            {
              test: /node_modules\/(?:vue|vue-router|@vueuse\/core|pinia)/,
              name: 'common',
              priority: 1,
            },
            {
              test: /node_modules\/(?:naive-ui)/,
              name: 'ui',
              priority: 2,
            },
          ],
        },
      },
    },
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  plugins: [
    vue(),
    svgLoader(),
    Unocss(),
    vuetify({
      autoImport: true,
    }),
    AutoImport({
      imports: ['vue', '@vueuse/core', 'vue-router', 'pinia'],
      dts: 'src/auto-imports.d.ts',
      eslintrc: {
        enabled: true,
      },
      vueTemplate: true,
    }),
    Components({
      resolvers: [NaiveUiResolver()],
      dts: true,
    }),
  ],
  server: {
    host: '0.0.0.0',
    port: 5173,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
      Authorization: 'Basic YWRtaW46Q2lkc2ljLXNpc2phZC1yeXptdTE=',
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8088',
        changeOrigin: true,
      },
    },
  },
})
