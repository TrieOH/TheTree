import { defineConfig } from 'vite'
import { devtools } from '@tanstack/devtools-vite'
import tsconfigPaths from 'vite-tsconfig-paths'

import { tanstackStart } from '@tanstack/react-start/plugin/vite'

import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { cloudflare } from "@cloudflare/vite-plugin";

const config = defineConfig({
  plugins: [
    devtools(),
    cloudflare({ viteEnvironment: { name: 'ssr' } }),
    tsconfigPaths({ projects: ['./tsconfig.json'] }),
    tailwindcss(),
    tanstackStart(),
    viteReact({
      babel: {
        plugins: [
          ["babel-plugin-react-compiler", {}],
        ],
      },
    }),
  ],
  build: {
    rollupOptions: {
      output: {
        manualChunks: (id) => {
          if (!id.includes('node_modules')) return

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

          // Internal SDK
          if (
            id.includes('@trieoh/identityx-sdk-ts') ||
            id.includes('@trieoh/envoy-fetch-ts')
          ) {
            return 'vendor-trieoh'
          }
        }
      }
    }
  }
})

export default config
