import { SessionI } from "../types/sessions-types";
import { clearAuthTokens, fetchAndSaveClaims, getUserInfo } from "../utils/token-utils";
import type { Api } from "./api";
import { env } from "./env";

export const createAuthService = (apiInstance: Api) => ({
  login: async (email: string, password: string) => {
    const url = env.API_KEY.length > 0 
      ? `/projects/${env.API_KEY}/login` : "/auth/login";

    const res = await apiInstance.post<string>(
      url,
      { email, password }, 
    );
    if(res.code === 200) await fetchAndSaveClaims(apiInstance);
    return res;
  },

  register: (email: string, password: string, flow_id?: string, custom: string[] = [""]) => {
    let url = env.API_KEY.length > 0 
      ? `/projects/${env.API_KEY}/register` : "/auth/register";

    if (env.API_KEY.length > 0) {
      if (!flow_id) {
        return Promise.reject({
          code: 400,
          message: "flow_id is required when a project_id is provided.",
          module: "auth",
          timestamp: new Date().toISOString(),
        });
      }
      const params = new URLSearchParams();
      params.append("flow_id", flow_id);
      params.append("schema_type", "context");
      params.append("version", "0");
      url += `?${params.toString()}`;
    }

    return apiInstance.post<string>(url, { email, password, custom_fields: custom });
  },  
  
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

  refreshProfileInfo: async () => fetchAndSaveClaims(apiInstance),

  profile: () => getUserInfo(),
});
