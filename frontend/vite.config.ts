import path from 'path'
import fs from 'fs'

import { defineConfig, Plugin } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

import Components from 'unplugin-vue-components/vite'
import AutoImport from 'unplugin-auto-import/vite'

import tailwindcss from '@tailwindcss/vite'

function fontPreloadInject(): Plugin {
  return {
    name: 'font-preload-inject',
    transformIndexHtml(html) {
      return html
    },
    closeBundle() {
      const distPath = path.resolve(__dirname, 'dist')
      const assetsPath = path.resolve(distPath, 'assets')

      if (!fs.existsSync(assetsPath)) return

      const fontFiles = fs.readdirSync(assetsPath).filter((file) => file.endsWith('.woff2'))

      if (fontFiles.length === 0) return

      const preloadLinks = fontFiles
        .map((fileName) => {
          return `    <link rel="preload" href="/assets/${fileName}" as="font" type="font/woff2" crossorigin>`
        })
        .join('\n')

      const htmlPath = path.resolve(distPath, 'index.html')
      if (!fs.existsSync(htmlPath)) return

      const originalHtml = fs.readFileSync(htmlPath, 'utf-8')
      const insertPoint = originalHtml.indexOf('</head>')

      if (insertPoint !== -1 && !originalHtml.includes('font/woff2')) {
        const modifiedHtml =
          originalHtml.slice(0, insertPoint) +
          `\n${preloadLinks}\n` +
          originalHtml.slice(insertPoint)
        fs.writeFileSync(htmlPath, modifiedHtml)
      }
    },
  }
}

// https://vite.dev/config/
export default defineConfig({
  base: process.env.BASE_URL || '/',
  plugins: [
    vue(),
    vueDevTools(),
    tailwindcss(),
    Components({
      dts: true,
    }),
    AutoImport({
      imports: ['vue', 'vue-router'],
    }),
    fontPreloadInject(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
  },
  define: {
    'process.env.VITE_API_URL': JSON.stringify(process.env.VITE_API_URL || 'http://localhost:8080'),
  },
})
