export interface TrieOHEnv {
  PROJECT_ID: string;
  API_KEY: string;
  BASE_URL: string;
}

export function resolveEnv(): TrieOHEnv {
  const isServer = typeof window === "undefined";

  const viteEnv = (
    typeof import.meta !== "undefined" && import.meta.env
      ? import.meta.env
      : {}
  ) as Partial<ImportMetaEnv>;

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
let memoizedEnv: TrieOHEnv | null = null;
let overrides: Partial<TrieOHEnv> = {};

/**
 * Configure the SDK manually. This will override any environment variables.
 */
export function configure(config: Partial<TrieOHEnv>) {
  overrides = { ...overrides, ...config };
  memoizedEnv = null; // Reset memoization to apply new config
}

function getEnv(): TrieOHEnv {
  if (!memoizedEnv) {
    const resolved = resolveEnv();
    memoizedEnv = {
      ...resolved,
      ...overrides,
    };
  }
  return memoizedEnv;
}
export const env: TrieOHEnv = {
  get PROJECT_ID() {
    return getEnv().PROJECT_ID;
  },
  get API_KEY() {
    return getEnv().API_KEY;
  },
  get BASE_URL() {
    return getEnv().BASE_URL;
  },
};
