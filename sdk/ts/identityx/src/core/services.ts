import type { CurrentSessionI, SessionI } from "../types/sessions-types";
import { clearAuthTokens, getUserInfo, saveAuthSession } from "../utils/token-utils";
import { validateProjectKey } from "../utils/env-validator";
import type { Api, ApiResponse } from "./api";
import { env } from "./env";
import type { IntrospectResponse } from "../types/instropect-types";
import { AuthTokens } from "../types/token-types";
import type { ProviderI } from "../types/common-types";

export interface AuthCallbacks {
  onLogin?: (res: ApiResponse<AuthTokens>) => void;
  onSetup?: (res: ApiResponse<AuthTokens>) => void;
  onRegister?: (res: ApiResponse<void>) => void;
  onVerify?: (res: ApiResponse<void>) => void;
  onResetPassword?: (res: ApiResponse<void>) => void;
  onRefresh?: (res?: ApiResponse<AuthTokens>) => void;
}

export const createAuthService = (apiInstance: Api, callbacks?: AuthCallbacks) => ({
  isSetupDone: async () => {
    return apiInstance.get<void>("/auth/setup", { requiresAuth: false });
  },

  setup: async (email: string, password: string) => {
    if (env.PROJECT_ID) validateProjectKey();
    const url = `/auth/setup${env.PROJECT_ID ? `?project_id=${env.PROJECT_ID}` : ""}`;
    const res = await apiInstance.post<AuthTokens>(
      url,
      { email, password },
      { requiresAuth: false }
    );

    if (res.success) {
      saveAuthSession(res.data);
      callbacks?.onSetup?.(res);
    }

    return res;
  },

  login: async (email: string, password: string) => {
    if (env.PROJECT_ID) validateProjectKey();
    const url = `/auth/login${env.PROJECT_ID ? `?project_id=${env.PROJECT_ID}` : ""}`;
    const res = await apiInstance.post<AuthTokens>(
      url,
      { email, password },
      { requiresAuth: false }
    );

    if (res.success) {
      saveAuthSession(res.data);
      callbacks?.onLogin?.(res);
    }

    return res;
  },

  loginWithProvider: async (provider: ProviderI) => {
    const url = `/auth/${provider}/connect`;
    const res = await apiInstance.get<{ url: string }>(url, { requiresAuth: false });
    return res;
  },

  completeProviderLogin: async (provider: ProviderI, code: string) => {
    const url = `/auth/${provider}/callback?code=${code}`;
    const res = await apiInstance.get<AuthTokens>(url, { requiresAuth: false });
    if (res.success) {
      saveAuthSession(res.data);
      callbacks?.onLogin?.(res);
    }
    return res;
  },

  register: async (email: string, password: string) => {
    const options = { requiresAuth: false };
    const url = `/auth/register${env.PROJECT_ID ? `?project_id=${env.PROJECT_ID}` : ""}`;
    if (env.PROJECT_ID) validateProjectKey();

    const res = await apiInstance.post<void>(url, { email, password }, options);
    if (res.success) callbacks?.onRegister?.(res);
    return res;
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

    if (res.success) {
      saveAuthSession(res.data);
      callbacks?.onRefresh?.(res);
    }

    return res;
  },

  // FIXME: This is not being used for now
  sessions: async () => apiInstance.get<SessionI[]>("/sessions"),

  currentSession: async () => apiInstance.get<CurrentSessionI>("/sessions/me"),

  revokeASession: async (id: string) => apiInstance.delete<void>(`/sessions/${id}`),

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
    const res = await apiInstance.post<void>(
      `/account/reset-password?token=${token}`,
      { new_password },
      { requiresAuth: false }
    );
    if (res.success) callbacks?.onResetPassword?.(res);
    return res;
  },

  verifyEmail: async (token: string) => {
    const url = `/account/verify?token=${token}`;
    const res = await apiInstance.post<void>(url);
    if (res.success) callbacks?.onVerify?.(res);
    return res;
  },

  resendVerifyEmail: async () => apiInstance.post<void>("/account/verify/resend"),

  introspect: async () => apiInstance.get<IntrospectResponse>("/auth/introspect"),

  health: async () => {
    return apiInstance.get<{ service: string; status: string }>("/health", {
      requiresAuth: false,
    });
  },

  authHealth: async () => {
    return apiInstance.get<{ service: string; status: string; user_id: string }>("/protected/health");
  }
});
