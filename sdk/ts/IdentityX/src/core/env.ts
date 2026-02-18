export interface TrieOHEnv {
  PROJECT_KEY: string;
  API_KEY: string;
  BASE_URL: string;
}

export function resolveEnv(): TrieOHEnv {
  let viteEnv: Partial<ImportMetaEnv> = {};
  
  try {
    const meta = import.meta;
    if (meta && meta.env) viteEnv = meta.env;
  } catch {  }

  const isServer = typeof window === "undefined";
  const safeProcessEnv = typeof process !== "undefined" ? process.env : {};

  return {
    PROJECT_KEY:
      viteEnv.VITE_TRIEOH_AUTH_PROJECT_KEY ??
      safeProcessEnv.NEXT_PUBLIC_TRIEOH_AUTH_PROJECT_KEY ??
      safeProcessEnv.PUBLIC_TRIEOH_AUTH_PROJECT_KEY ??
      "",

    API_KEY: isServer
      ? safeProcessEnv.TRIEOH_AUTH_API_KEY ?? ""
      : "",

    BASE_URL: "https://api.default.com",
  };
}

export const env = resolveEnv();
