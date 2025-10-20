import type { Api } from "./api";

export const createAuthService = (apiInstance: Api) => ({
  login: (email: string, password: string) =>
    apiInstance.post<string>("/auth/login", { email, password }),

  register: (email: string, password: string) =>
    apiInstance.post<string>("/auth/register", { email, password }),

  logout: () => apiInstance.post<string>("/auth/logout"),

  me: () => apiInstance.post<string>("/auth/me"),
});
