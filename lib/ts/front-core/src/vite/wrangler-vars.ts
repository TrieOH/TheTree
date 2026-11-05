import { unstable_readConfig } from "wrangler";

// A front's production config lives in its wrangler.jsonc `vars`, reviewed
// with the code like the backends', instead of in CI variables. Vite inlines
// `import.meta.env.VITE_*` at build time and only reads process.env and .env
// files, so a build copies the VITE_* vars into process.env first; the Worker
// gets the same vars at runtime as bindings. Dev keeps reading .env (only
// `vite build` loads them), and a variable that is already set still wins.
//
//   export default defineConfig(({ command }) => {
//     loadWranglerVars(command);
//     return { ... };
//   });
export function loadWranglerVars(
  command: "build" | "serve",
  configPath = "wrangler.jsonc",
): void {
  if (command !== "build") return;
  const { vars } = unstable_readConfig({ config: configPath });
  for (const [key, value] of Object.entries(vars)) {
    if (key.startsWith("VITE_") && process.env[key] === undefined) {
      process.env[key] =
        typeof value === "string" ? value : JSON.stringify(value);
    }
  }
}
