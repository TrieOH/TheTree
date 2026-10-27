import { defineConfig } from "vite";
import solid from "@solidjs/vite-plugin";
import tailwindcss from "@tailwindcss/vite";
import Icons from "unplugin-icons/vite";
import { apiReference } from "./plugins/api-reference.ts";

export default defineConfig({
  plugins: [
    apiReference(),
    tailwindcss(),
    Icons({ compiler: "solid" }),
    solid(),
  ],
  server: {
    port: 3004,
  },
  build: {
    target: "esnext",
  },
});
