import { defineConfig, type Options } from "tsup";

export default defineConfig((options) => {
  const common: Options = {
    entry: { index: "src/index.ts" },
    dts: false,
    minify: true,
    splitting: false,
    sourcemap: true,
    bundle: true,
    clean: true,
    external: ["solid-js", "@solidjs/web"],
    outExtension: () => ({ js: ".js" }),
    esbuildOptions(buildOptions) {
      buildOptions.jsx = "automatic";
      buildOptions.jsxImportSource = "@solidjs/web";
    },
    ...options,
  };

  return [
    { ...common, format: ["esm"], "dts": false, outDir: "dist/esm" },
    { ...common, format: ["cjs"], outDir: "dist/cjs", clean: false },
  ];
});
