import { type AuthTokenClaims, clearAuthTokens, saveTokenClaims } from "../utils/token-utils";
import type { Api } from "./api";

export const createAuthService = (apiInstance: Api) => ({
  login: async (email: string, password: string) => {
    const res = await apiInstance.post<AuthTokenClaims>(
      "/auth/login",
      { email, password }, 
      { requiresAuth: true }
    );
    if(res.code === 200 && res.data) saveTokenClaims(res.data);
    return res;
  },

  register: (email: string, password: string) =>
    apiInstance.post<string>("/auth/register", { email, password }),

  logout: async () => {
    const res = await apiInstance.post<string>(
      "/auth/logout",
      undefined, 
      { requiresAuth: true }
    );
    if(res.code === 200) clearAuthTokens();
    return res;
  },

  refresh: async () => {
    const res = await apiInstance.post<AuthTokenClaims>(
      "/auth/refresh",
      undefined,
      { requiresAuth: true, skipRefresh: true }
    );
    if(res.code === 200 && res.data) saveTokenClaims(res.data);
    return res;
  }
});
