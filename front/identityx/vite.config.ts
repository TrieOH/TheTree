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
    tanstackStart({
      serverFns: {
        // Readable /_serverFn/<file>.<fn> URLs in devtools instead of hashes;
        // the code is source-available, so the hashes only hide it from us.
        generateFunctionId: ({ filename, functionName }) => {
          const file = filename.split('/').pop()?.replace(/\.[jt]sx?$/, '')
          return `${file}.${functionName.replace(/(ServerFn)?_createServerFn_handler$/, '')}`
        },
      },
    }),
    viteReact(),
  ],
  build: {
    chunkSizeWarningLimit: 1000,
    // real file names and lines in the Sources tab and stack traces
    sourcemap: true,
  },
})

export default config
