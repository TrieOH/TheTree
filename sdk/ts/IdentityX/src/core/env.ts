export interface TrieOHEnv {
  PROJECT_ID: string;
  API_KEY: string;
  BASE_URL: string;
}

export function resolveEnv(): TrieOHEnv {
  const isServer = typeof window === "undefined";

  let viteEnv: ImportMetaEnv | Record<string, never> = {};

  if (typeof import.meta !== "undefined") viteEnv = import.meta.env;

  const safeProcessEnv: NodeJS.ProcessEnv = 
    typeof process !== "undefined" ? process.env : {};

  const resolvedProjectId =
    viteEnv.VITE_TRIEOH_AUTH_PROJECT_ID ||
    safeProcessEnv.NEXT_PUBLIC_TRIEOH_AUTH_PROJECT_ID ||
    safeProcessEnv.PUBLIC_TRIEOH_AUTH_PROJECT_ID ||
    "";

  const resolvedApiKey = isServer
    ? (safeProcessEnv.TRIEOH_AUTH_API_KEY || "")
    : "";

  return {
    PROJECT_ID: resolvedProjectId,
    API_KEY: resolvedApiKey,
    BASE_URL: "https://api.default.com",
  };
}

export const env = resolveEnv();
