import type { ProjectFieldDefinitionResultI, FieldValue } from "../types/fields-types";
import type { SessionI } from "../types/sessions-types";
import type {
  CheckPermissionResponse,
  CheckPermissionRequest,
  CompletePermissionBuilder
} from "../types/permission-types";
import {
  clearAuthTokens,
  getUserInfo,
  saveAuthSession,
  type AuthTokens
} from "../utils/token-utils";
import { validateApiKey, validateProjectKey } from "../utils/env-validator";
import type { Api } from "./api";
import { env } from "./env";

export const createAuthService = (apiInstance: Api) => ({
  login: async (email: string, password: string) => {
    if (env.PROJECT_ID) validateProjectKey();
    const url = env.PROJECT_ID ? `/projects/${env.PROJECT_ID}/login` : "/auth/login";
    const res = await apiInstance.post<AuthTokens>(
      url,
      { email, password },
      { requiresAuth: false }
    );

    if (res.success) {
      saveAuthSession(res.data.access_token, res.data.refresh_token);
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

  logout: async (options?: { forceLogout?: boolean }) => {
    const url = env.PROJECT_ID ? `/projects/${env.PROJECT_ID}/logout` : "/auth/logout";
    const res = await apiInstance.post<null>(url);
    if (res.success || options?.forceLogout) clearAuthTokens();
    return res;
  },

  refresh: async () => {
    const res = await apiInstance.post<AuthTokens>(
      "/auth/refresh",
      undefined,
      { skipRefresh: true }
    );

    if (res.success) {
      saveAuthSession(res.data.access_token, res.data.refresh_token);
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

  profile: () => getUserInfo(),

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
          "Authorization": env.API_KEY,
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
          "Authorization": env.API_KEY,
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
          "Authorization": env.API_KEY,
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
          "Authorization": env.API_KEY,
          "Content-Type": "application/json"
        }
      }
    );
  },

  /**
   * Checks if a user has a specific permission.
   * @param permission PermissionBuilder instance or raw CheckPermissionRequest object.
   */
  checkPermission: async (permission: CompletePermissionBuilder | CheckPermissionRequest) => {
    validateApiKey();
    const req = ("toJSON" in permission) ? permission.toJSON() : permission;
    const payload = {
      project_id: env.PROJECT_ID,
      ...req
    };
    return apiInstance.post<CheckPermissionResponse>(
      "/authz/check",
      payload,
      {
        headers: {
          "X-API-Key": env.API_KEY,
          "Content-Type": "application/json"
        }
      }
    );
  },
});
