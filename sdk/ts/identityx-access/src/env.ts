export interface IdentityXAccessEnv {
  BASE_URL: string;
  API_KEY: string;
}

const DEFAULT_BASE_URL = "https://api.identityx.trieoh.com";

export function resolveEnv(): IdentityXAccessEnv {
  const env = process.env;

  return {
    BASE_URL: env.TRIEOH_IDENTITYX_ACCESS_BASE_URL ?? DEFAULT_BASE_URL,
    API_KEY: env.TRIEOH_IDENTITYX_ACCESS_API_KEY ?? "",
  };
}

export const env = resolveEnv();
