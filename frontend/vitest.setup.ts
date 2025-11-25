// tests/setup.ts
import '@testing-library/jest-dom/vitest'

// // Optional: stub Nuxt runtime config
// import { vi } from 'vitest'

// vi.mock('#app', async () => {
//   const actual = await vi.importActual<any>('#app')
//   return {
//     ...actual,
//     // mock navigateTo for tests
//     navigateTo: vi.fn()
//   }
// })

// // Optional: mock useRouter, useRoute if needed
// vi.mock('vue-router', () => ({
//   useRouter: () => ({ push: vi.fn() }),
//   useRoute: () => ({ query: {} })
// }))
