import { defineConfig } from 'waku/config';
import mdx from 'fumadocs-mdx/vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  vite: {
    resolve: {
      tsconfigPaths: true,
      external: ['@takumi-rs/image-response'],
    },
    optimizeDeps: {
      include: ['ajv/dist/2020.js', 'xml-js', 'xml-js/lib/js2xml.js'],
    },
    ssr: {
      noExternal: ['fumadocs-openapi'],
      external: ['xml-js'],
    },
    plugins: [tailwindcss(), mdx()],
  },
});
