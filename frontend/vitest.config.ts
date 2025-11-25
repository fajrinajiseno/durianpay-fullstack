// /// <reference types="vitest" />
// import { defineConfig } from 'vitest/config'
// import vue from '@vitejs/plugin-vue'
// import { fileURLToPath, URL } from 'node:url'

// export default defineConfig({
//   plugins: [vue()],

//   test: {
//     environment: 'jsdom',
//     globals: true,
//     setupFiles: ['./vitest.setup.ts'],
//     include: ['tests/**/*.spec.ts', 'tests/**/*.test.ts'],
//     deps: {
//       inline: ['@nuxt/ui', '@vueuse/core'],
//       registerNodeLoader: true
//     },
//     coverage: {
//       reporter: ['text', 'html'],
//       exclude: ['node_modules/']
//     }
//   },

//   resolve: {
//     alias: {
//       // Nuxt aliases
//       '~': fileURLToPath(new URL('./', import.meta.url)),
//       '@': fileURLToPath(new URL('./', import.meta.url)),
//       // if you want alias for composables, stores, components
//       '#imports': fileURLToPath(new URL('./.nuxt/imports.mjs', import.meta.url))
//     }
//   }
// })
import { defineConfig } from 'vitest/config'
import { defineVitestProject } from '@nuxt/test-utils/config'

export default defineConfig({
  test: {
    projects: [
      await defineVitestProject({
        test: {
          name: 'nuxt',
          include: ['tests/**/*.spec.ts', 'tests/**/*.test.ts'],
          environment: 'nuxt',
          globals: true
        }
      })
    ]
  }
})
