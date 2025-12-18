import { SessionI } from "../types/sessions-types";
import { clearAuthTokens, fetchAndSaveClaims, getUserInfo } from "../utils/token-utils";
import type { Api } from "./api";

export const createAuthService = (apiInstance: Api) => ({
  login: async (email: string, password: string) => {
    const res = await apiInstance.post<string>(
      "/auth/login",
      { email, password }, 
    );
    if(res.code === 200) await fetchAndSaveClaims(apiInstance);
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
    const res = await apiInstance.post<string>(
      "/auth/refresh",
      undefined,
      { requiresAuth: true, skipRefresh: true }
    );
    if(res.code === 200) await fetchAndSaveClaims(apiInstance);
    return res;
  },

  sessions: async () => {
    const res = await apiInstance.get<SessionI[]>("/sessions", { requiresAuth: true });
    return res;
  },

  revokeASession: async (id: string) => {
    const res = await apiInstance.delete<string>(`/sessions/${id}`, { requiresAuth: true });
    return res;
  },

  revokeSessions: async (revokeAll: boolean = false) => {
    const path = revokeAll ? "/sessions" : "/sessions/others"
    const res = await apiInstance.delete<string>(path, { requiresAuth: true });
    return res;
  },

  profile: () => getUserInfo(),
});
