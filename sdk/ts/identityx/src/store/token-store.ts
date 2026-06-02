import { sessionBrowserStorage } from "../utils/storage-adapter";
import { obfuscate, deobfuscate } from "../utils/obfuscation-utils";

const ACCESS_TOKEN_KEY = "trieoh_access_token";

export const tokenStore = {
  getAccessToken: () => deobfuscate<string>(sessionBrowserStorage.getItem(ACCESS_TOKEN_KEY)),
  setAccessToken: (token: string | null) => {
    if (token) sessionBrowserStorage.setItem(ACCESS_TOKEN_KEY, obfuscate(token));
    else sessionBrowserStorage.removeItem(ACCESS_TOKEN_KEY);
  },
  clear: () => {
    sessionBrowserStorage.removeItem(ACCESS_TOKEN_KEY);
  }
};
