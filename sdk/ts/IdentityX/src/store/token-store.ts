import { sessionBrowserStorage } from "../utils/storage-adapter";

const ACCESS_TOKEN_KEY = "trieoh_access_token";

export const tokenStore = {
  getAccessToken: () => sessionBrowserStorage.getItem(ACCESS_TOKEN_KEY),
  setAccessToken: (token: string | null) => {
    if (token) sessionBrowserStorage.setItem(ACCESS_TOKEN_KEY, token);
    else sessionBrowserStorage.removeItem(ACCESS_TOKEN_KEY);
  },
  clear: () => {
    sessionBrowserStorage.removeItem(ACCESS_TOKEN_KEY);
  }
};
