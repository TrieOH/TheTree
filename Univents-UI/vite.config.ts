import { defineConfig } from 'vite'
import { devtools } from '@tanstack/devtools-vite'
import tsconfigPaths from 'vite-tsconfig-paths'

import { tanstackStart } from '@tanstack/react-start/plugin/vite'

import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { cloudflare } from '@cloudflare/vite-plugin'

const config = defineConfig({
  plugins: [
    devtools(),
    cloudflare({ viteEnvironment: { name: 'ssr' } }),
    tsconfigPaths({ projects: ['./tsconfig.json'] }),
    tailwindcss(),
    tanstackStart(),
    viteReact(),
  ],
  build: {
    rollupOptions: {
      output: {
        manualChunks: (id) => {
          if (!id.includes('node_modules')) return

          // React core + TanStack Store
          if (
            id.includes('/react/') ||
            id.includes('/react-dom/') ||
            id.includes('/scheduler/') ||
            id.includes('@tanstack/react-store') ||
            id.includes('@tanstack/store')
          ) {
            return 'vendor-react'
          }

          // TanStack Router + SSR
          if (
            id.includes('@tanstack/react-router') ||
            id.includes('@tanstack/react-start') ||
            id.includes('@tanstack/router-plugin') ||
            id.includes('@tanstack/react-router-ssr-query')
          ) {
            return 'vendor-tanstack-router'
          }

          // TanStack Query
          if (id.includes('@tanstack/react-query')) return 'vendor-tanstack-query'

          // Validations - Zod
          if (id.includes('/zod/')) return 'vendor-zod'

          // Forms
          if (id.includes('react-hook-form') || id.includes('@hookform/resolvers'))
            return 'vendor-forms'

          // Analytics / tracking
          if (id.includes('posthog-js') || id.includes('@posthog/react'))
            return 'vendor-analytics'

          // Animations
          if (id.includes('/motion/')) return 'vendor-animations'

          // Icons
          if (id.includes('lucide-react')) return 'vendor-icons'

          // UI: components and primitives
          if (
            id.includes('@base-ui/react') ||
            id.includes('vaul') ||
            id.includes('sonner') ||
            id.includes('next-themes')
          ) {
            return 'vendor-ui-primitives'
          }

          // Utils CSS
          if (
            id.includes('class-variance-authority') ||
            id.includes('tailwind-merge') ||
            id.includes('clsx')
          ) {
            return 'vendor-css-utils'
          }

          // Payments
          if (id.includes('@mercadopago/sdk-js')) return 'vendor-payments'

          // Internal SDK
          if (
            id.includes('@soramux/node-auth-sdk') ||
            id.includes('@soramux/node-fetch-sdk') ||
            id.includes('@soramux/node-payments-sdk')
          ) {
            return 'vendor-soramux'
          }

          // Utils Network and Infra
          if (
            id.includes('aws4fetch') ||
            id.includes('@microsoft/fetch-event-source') ||
            id.includes('@t3-oss/env-core')
          ) {
            return 'vendor-infra'
          }
        },
      },
    },
  },
})

export default config
