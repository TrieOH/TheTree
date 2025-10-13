import React, { createContext, useContext, useMemo } from "react";
import { Api } from "../core/api";
import { AuthService } from "../core/services";

type AuthContextType = {
  auth: typeof AuthService;
};

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({
  children,
  baseURL,
}: {
  children: React.ReactNode;
  baseURL?: string;
}) {
  const apiInstance = useMemo(() => new Api(baseURL), [baseURL]);
  const auth = useMemo(
    () => ({
      login: (email: string, password: string) =>
        apiInstance.post("/auth/login", { email, password }),
      register: (email: string, password: string) =>
        apiInstance.post("/auth/register", { email, password }),
    }),
    [apiInstance]
  );

  return (
    <AuthContext.Provider value={{ auth }}>{children}</AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>");
  return ctx;
}
