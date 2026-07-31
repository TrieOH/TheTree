import { defineConfig } from 'vite'
import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import viteReact, { reactCompilerPreset } from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { cloudflare } from '@cloudflare/vite-plugin'

import babel from '@rolldown/plugin-babel'

const config = defineConfig({
  resolve: { tsconfigPaths: true },
  plugins: [
    cloudflare({ viteEnvironment: { name: 'ssr' } }),
    tailwindcss(),
    {
      ...babel({
        presets: [reactCompilerPreset()],
        include: /\.[jt]sx?$/,
        exclude: /node_modules/,
      }),
      enforce: 'post'
    },
    tanstackStart(),
    viteReact(),
  ],
  build: {
    chunkSizeWarningLimit: 1000,
  },
})

export default config
