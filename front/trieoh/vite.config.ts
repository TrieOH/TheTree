import { defineConfig } from "vite";
import solid from "@solidjs/vite-plugin";
import tailwindcss from "@tailwindcss/vite";
import Icons from "unplugin-icons/vite";

export default defineConfig({
  plugins: [tailwindcss(), Icons({ compiler: "solid" }), solid()],
  server: { port: 3005 },
  build: { target: "esnext" },
});
