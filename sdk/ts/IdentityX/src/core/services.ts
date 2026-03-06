import type { ProjectFieldDefinitionResultI, FieldValue } from "../types/fields-types";
import type { SessionI } from "../types/sessions-types";
import { clearAuthTokens, fetchAndSaveClaims, getUserInfo } from "../utils/token-utils";
import { validateApiKey, validateProjectKey } from "../utils/env-validator";
import type { Api } from "./api";
import { env } from "./env";

export const createAuthService = (apiInstance: Api) => ({
  login: async (email: string, password: string) => {
    if (env.PROJECT_ID) {
      validateProjectKey();
      const url = `/projects/${env.PROJECT_ID}/login`;
      const res = await apiInstance.post<{is_up_to_date: boolean}>(url, { email, password });
      if(res.code === 200) await fetchAndSaveClaims(apiInstance);
      return res;
    }

    const res = await apiInstance.post<{is_up_to_date: boolean}>("/auth/login", { email, password });
    if(res.code === 200) await fetchAndSaveClaims(apiInstance);
    return res;
  },

  register: (email: string, password: string, flow_id?: string, custom: Record<string, FieldValue> = {}) => {
    if (env.PROJECT_ID) {
      validateProjectKey();
      
      const params = new URLSearchParams();
      params.append("flow_id", flow_id || "none");
      params.append("schema_type", "context");
      params.append("version", "1");
      const url = `/projects/${env.PROJECT_ID}/register?${params.toString()}`;
      return apiInstance.post<string>(url, { email, password, custom_fields: custom });
    }

    return apiInstance.post<string>("/auth/register", { email, password });
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
    const res = await apiInstance.post<{is_up_to_date: boolean}>(
      "/auth/refresh",
      undefined,
      { requiresAuth: true, skipRefresh: true }
    );
    if(res.code === 200) await fetchAndSaveClaims(apiInstance);
    return res;
  },

  sessions: async () => {
    return apiInstance.get<SessionI[]>("/sessions", { requiresAuth: true });
  },

  revokeASession: async (id: string) => {
    return apiInstance.delete<string>(`/sessions/${id}`, { requiresAuth: true });
  },

  revokeSessions: async (revokeAll: boolean = false) => {
    const path = revokeAll ? "/sessions" : "/sessions/others"
    return apiInstance.delete<string>(path, { requiresAuth: true });
  },

  refreshProfileInfo: async () => fetchAndSaveClaims(apiInstance),

  profile: () => getUserInfo(),

  getProfileUpgradeForms: async () => {
    validateProjectKey();
    const url = `/projects/${env.PROJECT_ID}/upgrade-form`;
    return apiInstance.get<ProjectFieldDefinitionResultI>(url, { requiresAuth: true });
  },

  updateProfile: async (custom: Record<string, FieldValue>) => {
    validateProjectKey();
    const url = `/projects/${env.PROJECT_ID}/metadata`;
    return apiInstance.post<string>(url, { custom_fields: custom }, { requiresAuth: true });
  },

  sendForgotPassword: async (email: string) => {
    if (env.PROJECT_ID) {
      validateProjectKey();
      return apiInstance.post<string>(
        "/auth/forgot-password",
        {email, project_id: env.PROJECT_ID}, 
      );
    }
    return apiInstance.post<string>("/auth/forgot-password", {email});
  },

  resetPassword: async (token: string, new_password: string) => {
    return apiInstance.post<string>(
      `/auth/reset-password?token=${token}`,
      {new_password}, 
    );
  },

  resendVerifyEmail: async () => {
    return apiInstance.post<string>(
      "/auth/verify/resend",
      undefined,
      { requiresAuth: true } 
    );
  },

  verifyEmail: async () => {
    return apiInstance.get<string>(
      "/auth/verify",
      undefined,
    );
  },
});

export const createServerAuthService = (apiInstance: Api) => ({
  getProjectLatestRegisterFields: async (flow_id: string) => {
    validateProjectKey();
    validateApiKey();

    let url = `/projects/${env.PROJECT_ID}/schemas/lookup/latest`
    const params = new URLSearchParams();
    params.append("flow_id", flow_id);
    params.append("schema_type", "context");
    url += `?${params.toString()}`;

    return apiInstance.get<ProjectFieldDefinitionResultI>(
      url,
      {
        headers: {
          "Authorization": `Bearer ${env.API_KEY}`,
          "Content-Type": "application/json"
        }
      }
    );
  },

  getProjectRegisterFields: async (flow_id: string) => {
    validateProjectKey();
    validateApiKey();

    const version = 1;
    let url = `/projects/${env.PROJECT_ID}/schemas/lookup/v${version}`
    const params = new URLSearchParams();
    params.append("flow_id", flow_id);
    params.append("schema_type", "context");
    url += `?${params.toString()}`;

    return apiInstance.get<ProjectFieldDefinitionResultI>(
      url,
      {
        headers: {
          "Authorization": `Bearer ${env.API_KEY}`,
          "Content-Type": "application/json"
        }
      }
    );
  },

  assignRoleByNameToUser: async (user_id: string, role_name: string, scope_id: string | null) => {
    validateProjectKey();
    validateApiKey();

    const url = `/projects/${env.PROJECT_ID}/identities/${user_id}/roles/by-name`

    return apiInstance.post<void>(
      url,
      { role_name, scope_id },
      {
        headers: {
          "Authorization": `Bearer ${env.API_KEY}`,
          "Content-Type": "application/json"
        }
      }
    );
  },

  removeRoleByNameFromUser: async (user_id: string, role_name: string, scope_id: string | null) => {
    validateProjectKey();
    validateApiKey();

    const url = `/projects/${env.PROJECT_ID}/identities/${user_id}/roles/by-name`

    return apiInstance.delete<void>(
      url,
      { role_name, scope_id },
      {
        headers: {
          "Authorization": `Bearer ${env.API_KEY}`,
          "Content-Type": "application/json"
        }
      }
    );
  },
});
