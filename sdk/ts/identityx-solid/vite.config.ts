import { defineConfig } from "vite";
import solid from "@solidjs/vite-plugin";

export default defineConfig({
  plugins: [solid()],
  build: {
    lib: {
      entry: "src/index.ts",
      formats: ["es", "cjs"],
      fileName: (format) => format === "es" ? "esm/index.js" : "cjs/index.cjs",
    },
    rollupOptions: {
      external: ["solid-js", "@solidjs/web", "@trieoh/envoy-fetch-ts", "@trieoh/identityx-sdk-ts"],
    },
  },
});
