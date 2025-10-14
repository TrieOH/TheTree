export function resolveEnv() {
  let viteEnv: ImportMetaEnv = {};
  try {
    viteEnv = (typeof import.meta !== "undefined" && import.meta.env) || {};
  } catch {
    viteEnv = { VITE_TRIEOH_AUTH_API_KEY: undefined };
  }

  return {
    API_KEY:
      // Vite (import.meta.env.VITE_*)
      viteEnv.VITE_TRIEOH_AUTH_API_KEY ??
      // Next (process.env.NEXT_PUBLIC_*)
      (typeof process !== "undefined"
        ? process.env.NEXT_PUBLIC_TRIEOH_AUTH_API_KEY
        : undefined) ??
      // Node (process.env.PUBLIC_TRIEOH_*)
      (typeof process !== "undefined"
        ? process.env.PUBLIC_TRIEOH_AUTH_API_KEY
        : undefined) ??
      "",
    BASE_URL: "https://api.default.com",
  };
}

export const env = resolveEnv();
