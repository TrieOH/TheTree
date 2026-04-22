import { defineConfig } from 'vite'
import viteTsConfigPaths from 'vite-tsconfig-paths'
import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { cloudflare } from '@cloudflare/vite-plugin'

const config = defineConfig({
  plugins: [
    cloudflare({ viteEnvironment: { name: 'ssr' } }),
    viteTsConfigPaths({ projects: ['./tsconfig.json'] }),
    tailwindcss(),
    tanstackStart(),
    viteReact(),
  ],
  build: {
    rollupOptions: {
      output: {
        manualChunks: (id) => {
          if (!id.includes('node_modules')) return

          // React Core e TanStack Store
          if (
            id.includes('/react/') ||
            id.includes('/react-dom/') ||
            id.includes('/scheduler/') ||
            id.includes('@tanstack/react-store') ||
            id.includes('@tanstack/store')
          ) {
            return 'vendor-react'
          }

          // TanStack Router, Start e SSR
          if (
            id.includes('@tanstack/react-router') ||
            id.includes('@tanstack/router-plugin') ||
            id.includes('@tanstack/react-router-ssr-query')
          ) {
            return 'vendor-tanstack-router'
          }

          // TanStack Query
          if (id.includes('@tanstack/react-query')) return 'vendor-tanstack-query'

          // Validação (Zod)
          if (id.includes('/zod/')) return 'vendor-zod'

          // TanStack Form
          if (
            id.includes('@tanstack/react-form') || 
            id.includes('@tanstack/react-form-start')
          ) {
            return 'vendor-forms'
          }
          
          // Analytics / tracking
          if (id.includes('posthog-js')) return 'vendor-analytics'

          // Animations
          if (id.includes('/motion/')) return 'vendor-animations'
          
          // Icons
          if (id.includes('lucide-react')) return 'vendor-icons'

          // UI: components and primitives
          if (
            id.includes('radix-ui') ||
            id.includes('sonner') ||
            id.includes('next-themes')
          ) {
            return 'vendor-ui-primitives'
          }

          // Utils CSS
          if (
            id.includes('class-variance-authority') ||
            id.includes('tailwind-merge') ||
            id.includes('clsx') ||
            id.includes('tw-animate-css')
          ) {
            return 'vendor-css-utils'
          }

          // Internal SDK
          if (id.includes('@soramux/identityx-sdk-ts')) return 'vendor-identityx'

          // Utils Network and Infra
          if (id.includes('@t3-oss/env-core'))  return 'vendor-infra'
        },
      },
    },
  },
})

export default config
