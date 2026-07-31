import { defineConfig } from "tsup";

export default defineConfig({
  entry: ["src/index.ts"],
  format: ["esm"],
  dts: false,
  minify: true,
  bundle: true,
  clean: true,
  sourcemap: true,
  splitting: false,
  treeshake: true,
  shims: true,
  outDir: "dist",
});
