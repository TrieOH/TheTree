import type { ProjectFieldDefinitionResultI, FieldValue } from "../types/fields-types";
import type { SessionI } from "../types/sessions-types";
import {
  clearAuthTokens,
  exchangeAndSaveClaims,
  fetchAndSaveClaims,
  getUserInfo,
  isUpToDate,
  TokenClaims
} from "../utils/token-utils";
import { validateApiKey, validateProjectKey } from "../utils/env-validator";
import type { Api, ApiResponse } from "./api";
import { env } from "./env";

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  is_up_to_date: boolean;
}

export const createAuthService = (apiInstance: Api, exchangeURL?: string) => ({
  login: async (email: string, password: string) => {
    const options = { requiresAuth: false };
    let res: ApiResponse<AuthTokens>;

    if (env.PROJECT_ID) {
      validateProjectKey();
      const url = `/projects/${env.PROJECT_ID}/login`;
      res = await apiInstance.post<AuthTokens>(
        url,
        { email, password },
        options
      );
    } else {
      res = await apiInstance.post<AuthTokens>(
        "/auth/login",
        { email, password },
        options
      );
    }

    if (res.success) {
      try {
        await exchangeAndSaveClaims(
          apiInstance,
          res.data.access_token,
          res.data.refresh_token,
          res.data.is_up_to_date,
          exchangeURL
        );
        return res;
      } catch (error) {
        console.error("[TRIEOH SDK] Exchange failed during login:", error);
        clearAuthTokens();
        return {
          success: false,
          code: 500,
          message: error instanceof Error ? error.message : "Authentication failed during exchange"
        } as ApiResponse<AuthTokens>;
      }
    }

    return res;
  },

  register: (email: string, password: string, flow_id?: string, custom: Record<string, FieldValue> = {}) => {
    const options = { requiresAuth: false };
    if (env.PROJECT_ID) {
      validateProjectKey();

      const params = new URLSearchParams();
      if (flow_id) {
        params.append("flow_id", flow_id);
        params.append("schema_type", "context");
        params.append("version", "1");
      }
      const bodyData = { ...{ email, password }, ...flow_id && { custom_fields: custom } };
      const paramsUrl = params.toString() ? `?${params.toString()}` : "";
      const url = `/projects/${env.PROJECT_ID}/register${paramsUrl}`;
      return apiInstance.post<string>(url, bodyData, options);
    }

    return apiInstance.post<string>("/auth/register", { email, password }, options);
  },

  logout: async () => {
    const res = await apiInstance.post<string>("/auth/logout");
    if (res.success) clearAuthTokens();
    return res;
  },

  refresh: async () => {
    const res = await apiInstance.post<AuthTokens>(
      "/auth/refresh",
      undefined,
      { skipRefresh: true }
    );

    if (res.success) {
      try {
        await exchangeAndSaveClaims(
          apiInstance,
          res.data.access_token,
          res.data.refresh_token,
          res.data.is_up_to_date,
          exchangeURL
        );
        return res;
      } catch (error) {
        console.error("[TRIEOH SDK] Exchange failed during refresh:", error);
        clearAuthTokens();
        return {
          success: false,
          code: 500,
          message: error instanceof Error ? error.message : "Refresh failed during exchange"
        } as ApiResponse<AuthTokens>;
      }
    }

    return res;
  },

  sessions: async () => {
    return apiInstance.get<SessionI[]>("/sessions");
  },

  revokeASession: async (id: string) => {
    return apiInstance.delete<string>(`/sessions/${id}`);
  },

  revokeSessions: async (revokeAll: boolean = false) => {
    const path = revokeAll ? "/sessions" : "/sessions/others"
    return apiInstance.delete<string>(path);
  },

  /**
   * Only for non Project
   * @returns APIResponse Claim
   */
  refreshProfileInfo: async () => fetchAndSaveClaims(apiInstance),

  profile: () => getUserInfo(),

  isUpToDate: () => isUpToDate(),

  getProfileUpgradeForms: async () => {
    validateProjectKey();
    const url = `/projects/${env.PROJECT_ID}/upgrade-form`;
    return apiInstance.get<ProjectFieldDefinitionResultI>(url);
  },

  updateProfile: async (custom: Record<string, FieldValue>) => {
    validateProjectKey();
    const url = `/projects/${env.PROJECT_ID}/metadata`;
    return apiInstance.post<string>(url, { custom_fields: custom });
  },

  sendForgotPassword: async (email: string) => {
    const options = { requiresAuth: false };
    if (env.PROJECT_ID) {
      validateProjectKey();
      return apiInstance.post<string>(
        "/auth/forgot-password",
        { email, project_id: env.PROJECT_ID },
        options
      );
    }
    return apiInstance.post<string>("/auth/forgot-password", { email }, options);
  },

  resetPassword: async (token: string, new_password: string) => {
    return apiInstance.post<string>(
      `/auth/reset-password?token=${token}`,
      { new_password },
      { requiresAuth: false }
    );
  },

  resendVerifyEmail: async () => {
    return apiInstance.post<string>("/auth/verify/resend");
  },

  verifyEmail: async () => {
    return apiInstance.get<string>(
      "/auth/verify",
      { requiresAuth: false }
    );
  },

  addSubContext: async (user_id: string, data: Record<string, unknown>) => {
    validateProjectKey();
    return apiInstance.post<void>(
      `/projects/${env.PROJECT_ID}/sub-context`,
      { data, user_id }
    );
  },

  removeSubContext: async (user_id: string, keys: string[]) => {
    validateProjectKey();
    return apiInstance.delete<void>(
      `/projects/${env.PROJECT_ID}/sub-context`,
      { keys, user_id }
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
