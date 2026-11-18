import { defineConfig } from "vitest/config";
import solid from "@solidjs/vite-plugin";
import Icons from "unplugin-icons/vite";

export default defineConfig({
  plugins: [Icons({ compiler: "solid" }), solid()],
  test: {
    environment: "jsdom",
    include: ["tests/**/*.test.tsx"],
    restoreMocks: true,
  },
});
