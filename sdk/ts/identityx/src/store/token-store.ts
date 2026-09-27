import {
  browserStorage,
  sessionBrowserStorage,
  type StorageAdapter,
} from "../utils/storage-adapter";
import type { AuthTokenClaims } from "../types/token-types";

/**
 * Where a session's tokens live. Split in two because the browser keeps the
 * access token per tab (sessionStorage) while the refresh token has to survive
 * reloads (localStorage). On a server both are usually the same backing store.
 */
export interface TokenStorage {
  access: StorageAdapter;
  persistent: StorageAdapter;
}

/** Browser default: access token per tab, refresh token across reloads. */
export const browserTokenStorage: TokenStorage = {
  access: sessionBrowserStorage,
  persistent: browserStorage,
};

/**
 * The token layout plus the decoded-claims memo, bound to one session.
 *
 * Inject one to keep the session somewhere other than the browser — a cookie on
 * a BFF, memory in tests — instead of the module-level browser storages.
 */
export interface TokenStore {
  getAccessToken(): string | null;
  setAccessToken(token: string | null): void;
  getRefreshToken(): string | null;
  setRefreshToken(token: string | null): void;
  getAccessExpiry(): number | null;
  setAccessExpiry(timestamp: number | null): void;
  getRefreshExpiry(): number | null;
  setRefreshExpiry(timestamp: number | null): void;
  getClaims(): AuthTokenClaims | null;
  setClaims(claims: AuthTokenClaims | null): void;
  clear(): void;
}

const ACCESS_TOKEN_KEY = "trieoh_access_token";
const ACCESS_EXPIRY_KEY = "trieoh_access_expiry";
const REFRESH_EXPIRY_KEY = "trieoh_refresh_expiry";
const REFRESH_TOKEN_KEY = "trieoh_refresh_token";

function readNumber(storage: StorageAdapter, key: string): number | null {
  const raw = storage.getItem(key);
  if (raw === null) return null;
  const value = Number(raw);
  return Number.isFinite(value) ? value : null;
}

function writeNumber(
  storage: StorageAdapter,
  key: string,
  value: number | null,
): void {
  if (value === null) storage.removeItem(key);
  else storage.setItem(key, String(value));
}

/** Binds the token layout to a pair of storages, isolated per store. */
export function createTokenStore(
  storage: TokenStorage = browserTokenStorage,
): TokenStore {
  let claims: AuthTokenClaims | null = null;

  return {
    getAccessToken: () => storage.access.getItem(ACCESS_TOKEN_KEY),
    setAccessToken: (token) => {
      if (token) storage.access.setItem(ACCESS_TOKEN_KEY, token);
      else storage.access.removeItem(ACCESS_TOKEN_KEY);
    },

    getRefreshToken: () => storage.persistent.getItem(REFRESH_TOKEN_KEY),
    setRefreshToken: (token) => {
      if (token) storage.persistent.setItem(REFRESH_TOKEN_KEY, token);
      else storage.persistent.removeItem(REFRESH_TOKEN_KEY);
    },

    getAccessExpiry: () => readNumber(storage.persistent, ACCESS_EXPIRY_KEY),
    setAccessExpiry: (timestamp) =>
      writeNumber(storage.persistent, ACCESS_EXPIRY_KEY, timestamp),

    getRefreshExpiry: () => readNumber(storage.persistent, REFRESH_EXPIRY_KEY),
    setRefreshExpiry: (timestamp) =>
      writeNumber(storage.persistent, REFRESH_EXPIRY_KEY, timestamp),

    getClaims: () => claims,
    setClaims: (next) => {
      claims = next;
    },

    clear: () => {
      claims = null;
      storage.access.removeItem(ACCESS_TOKEN_KEY);
      storage.persistent.removeItem(ACCESS_EXPIRY_KEY);
      storage.persistent.removeItem(REFRESH_EXPIRY_KEY);
      storage.persistent.removeItem(REFRESH_TOKEN_KEY);
    },
  };
}

/** Used when a consumer does not inject a store: the browser's own storages. */
export const defaultTokenStore: TokenStore = createTokenStore();
