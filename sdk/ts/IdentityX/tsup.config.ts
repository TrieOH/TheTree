import { defineConfig, Options } from "tsup";

export default defineConfig((options) => {
  const commonOptions: Options = {
    entry: ["src/next.ts"],
    dts: false,
    minify: true,
    splitting: false,
    sourcemap: true,
    bundle: true,
    clean: true,
    external: ["react", "react-dom"],
    outExtension({ format }) {
      return format === "esm" ? { js: ".js" } : { js: ".js" };
    },
    ...options,
  };

  return [
    // ESM build + declarations
    {
      ...commonOptions,
      format: ["esm"],
      dts: true,
      outDir: "dist/esm",
    },
    // CJS build
    {
      ...commonOptions,
      format: ["cjs"],
      outDir: "dist/cjs",
    },
  ];
});
