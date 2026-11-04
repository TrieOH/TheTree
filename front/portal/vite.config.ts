import { tanstackRouter } from "@tanstack/router-plugin/vite";
import { defineConfig } from "vitest/config";
import mdx from "@mdx-js/rollup";
import solid from "@solidjs/vite-plugin";
import tailwindcss from "@tailwindcss/vite";
import Icons from "unplugin-icons/vite";
import { cloudflare } from "@cloudflare/vite-plugin";
import { fileURLToPath } from "node:url";
import rehypeShiki from "@shikijs/rehype";
import remarkGfm from "remark-gfm";
import { problems } from "./mdx/problems-plugin.ts";
import { rehypeHeadingIds } from "./mdx/rehype-heading-ids.ts";
import { remarkAdr } from "./mdx/remark-adr.ts";
import { remarkMermaid } from "./mdx/remark-mermaid.ts";

export default defineConfig(({ mode }) => ({
  resolve: { tsconfigPaths: true },

  // Same turnkey client mode as front/univents-solid: the plugin generates the
  // entries around src/App.tsx, wrapped in src/Document.tsx, and `vite build`
  // emits a static dist/client that the Worker serves through ASSETS.
  plugins: [
    // The Cloudflare plugin owns worker environments and rejects the
    // `resolve.external` that Vitest injects, so keep it out of test runs.
    ...(mode === "test" ? [] : [cloudflare()]),
    tailwindcss(),
    Icons({ compiler: "solid" }),
    // Docs are Markdown/MDX compiled to plain JS calls into our own
    // jsx-runtime (src/shared/mdx), not to JSX: MDX renders tags as
    // `<_components.h1>` with string defaults, which Solid's compiler would
    // treat as components. The repo's docs/*.md go through the same pipeline;
    // they sit outside this project, where the `@/` alias doesn't apply, hence
    // the absolute path.
    {
      enforce: "pre" as const,
      ...mdx({
        jsxImportSource: fileURLToPath(new URL("./src/shared/mdx", import.meta.url)),
        mdExtensions: [".md"],
        mdxExtensions: [".mdx"],
        remarkPlugins: [remarkAdr, remarkGfm, remarkMermaid],
        // Highlighted at build time, in both themes: App.css picks one.
        rehypePlugins: [
          rehypeHeadingIds,
          [rehypeShiki, { themes: { light: "github-light", dark: "github-dark-dimmed" }, defaultColor: false }],
        ],
      }),
    },
    problems(),
    // Scans src/routes and generates src/routeTree.gen.ts — the typed route
    // tree — on dev and build. Must be registered before solid().
    tanstackRouter({ target: "solid", autoCodeSplitting: true }),
    solid({ start: true, diagnostics: true }),
  ],
  server: {
    port: 3010,
  },
  test: {
    environment: "jsdom",
    environmentOptions: {
      jsdom: {
        url: "http://localhost:3010",
      },
    },
    globals: false,
    setupFiles: ["./tests/ui/setup.ts"],
    isolate: true,
    // UI/component tests live here; Worker and handler tests run inside workerd
    // through vitest.workers.config.ts.
    include: ["tests/ui/**/*.test.{ts,tsx}"],
  },
  build: {
    target: "esnext",
    assetsInlineLimit: 0,
  },
}));
