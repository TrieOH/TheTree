import type { CurrentSessionI, SessionI } from "../types/sessions-types";
import {
  clearAuthTokens,
  getUserInfo,
  saveAuthSession,
  type AuthTokens
} from "../utils/token-utils";
import { validateProjectKey } from "../utils/env-validator";
import type { Api } from "./api";
import { env } from "./env";

export const createAuthService = (apiInstance: Api) => ({
  login: async (email: string, password: string) => {
    if (env.PROJECT_ID) validateProjectKey();
    const url = `/auth/login${env.PROJECT_ID ? `?project_id=${env.PROJECT_ID}` : ""}`;
    const res = await apiInstance.post<AuthTokens>(
      url,
      { email, password },
      { requiresAuth: false }
    );

    if (res.success) saveAuthSession(res.data.access_token, res.data.refresh_token);

    return res;
  },

  register: (email: string, password: string) => {
    const options = { requiresAuth: false };
    const url = `/auth/register${env.PROJECT_ID ? `?project_id=${env.PROJECT_ID}` : ""}`;
    if (env.PROJECT_ID) validateProjectKey();

    return apiInstance.post<void>(url, { email, password }, options);
  },

  logout: async (options?: { forceLogout?: boolean }) => {
    const url = `/auth/logout${env.PROJECT_ID ? `?project_id=${env.PROJECT_ID}` : ""}`;
    const res = await apiInstance.post<void>(url);
    if (res.success || options?.forceLogout) clearAuthTokens();
    return res;
  },

  refresh: async () => {
    const res = await apiInstance.post<AuthTokens>(
      "/auth/refresh",
      undefined,
      { skipRefresh: true }
    );

    if (res.success) saveAuthSession(res.data.access_token, res.data.refresh_token);

    return res;
  },

  sessions: async () => {
    return apiInstance.get<SessionI[]>("/sessions");
  },

  currentSession: async () => {
    return apiInstance.get<CurrentSessionI>("/sessions/me");
  },

  revokeASession: async (id: string) => {
    return apiInstance.delete<void>(`/sessions/${id}`);
  },

  revokeSessions: async (revokeAll: boolean = false) => {
    const path = revokeAll ? "/sessions" : "/sessions/others"
    return apiInstance.delete<void>(path);
  },

  profile: () => getUserInfo(),

  sendForgotPassword: async (email: string) => {
    const options = { requiresAuth: false };
    if (env.PROJECT_ID) {
      validateProjectKey();
      return apiInstance.post<void>(
        "/account/forgot-password",
        { email, project_id: env.PROJECT_ID },
        options
      );
    }
    return apiInstance.post<void>("/account/forgot-password", { email }, options);
  },

  resetPassword: async (token: string, new_password: string) => {
    return apiInstance.post<void>(
      `/account/reset-password?token=${token}`,
      { new_password },
      { requiresAuth: false }
    );
  },

  verifyEmail: async () => {
    return apiInstance.get<void>("/account/verify", { requiresAuth: false });
  },

  resendVerifyEmail: async () => {
    return apiInstance.post<void>("/account/verify/resend");
  },

  health: async () => {
    return apiInstance.get<{ service: string; status: string }>("/health", {
      requiresAuth: false,
    });
  },

  authHealth: async () => {
    return apiInstance.get<{ service: string; status: string; user_id: string }>("/protected/health");
  }

});
