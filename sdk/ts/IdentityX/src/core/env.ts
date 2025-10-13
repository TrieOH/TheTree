export function resolveEnv() {
  const viteEnv = (typeof import.meta !== "undefined" && import.meta.env) || {};
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
    BASE_URL: "https://api.default.com", // i need to change later
  };
}

export const env = resolveEnv();
