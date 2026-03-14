import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { Api } from "../core/api";
import { createAuthService } from "../core/services";
import { getTokenClaims, isUpToDate } from "../utils/token-utils";
import { validateProjectKey } from "../utils/env-validator";
import { configure } from "../core/env";

type AuthContextType = {
  auth: ReturnType<typeof createAuthService>;
  isAuthenticated: boolean;
  isUpToDate: boolean;
  isClient?: boolean;
};

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({
  children,
  baseURL,
  projectId,
  isClient = true,
}: {
  children: React.ReactNode;
  baseURL?: string;
  projectId?: string;
  isClient?: boolean;
}) {
  const [ready, setReady] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [upToDate, setUpToDate] = useState(false);

  // Apply manual configuration if provided
  useMemo(() => {
    if (projectId || baseURL) {
      configure({
        ...(projectId ? { PROJECT_ID: projectId } : {}),
        ...(baseURL ? { BASE_URL: baseURL } : {}),
      });
    }
  }, [projectId, baseURL]);

  const apiInstance = useMemo(() => new Api(baseURL, undefined, (claims) => {
    setUpToDate(claims.is_up_to_date || false);
  }), [baseURL]);

  const rawAuth = useMemo(() => createAuthService(apiInstance), [apiInstance]);

  const auth = useMemo(() => ({
    ...rawAuth,
    login: async (...args: Parameters<typeof rawAuth.login>) => {
      const res = await rawAuth.login(...args);
      if (res.success) {
        setIsAuthenticated(true);
        setUpToDate(isUpToDate());
      }
      return res;
    },
    logout: async (...args: Parameters<typeof rawAuth.logout>) => {
      const res = await rawAuth.logout(...args);
      if (res.success) {
        setIsAuthenticated(false);
        setUpToDate(false);
      }
      return res;
    },
    refresh: async (...args: Parameters<typeof rawAuth.refresh>) => {
      const res = await rawAuth.refresh(...args);
      if (res.success) {
        setIsAuthenticated(true);
        setUpToDate(isUpToDate());
      }
      return res;
    },
  }), [rawAuth]);

  useEffect(() => {
    if (isClient) validateProjectKey();

    const loadAuthStatus = async () => {
      const claims = getTokenClaims();
      if (claims) {
        setIsAuthenticated(true);
        setUpToDate(isUpToDate());
        setReady(true);
        return;
      }
      console.log("[TRIEOH SDK] Attempting to refresh session...");
      try {
        const res = await auth.refreshProfileInfo();
        if (res.success) {
          setIsAuthenticated(true);
          setUpToDate(isUpToDate());
          console.log("[TRIEOH SDK] Session restored successfully.");
        } else {
          setIsAuthenticated(false);
          setUpToDate(false);
          console.warn("[TRIEOH SDK] Session restoration failed/no session.");
        }
      } catch (error) {
        console.warn("[TRIEOH SDK] Unable to verify session (offline?)");
        setIsAuthenticated(false);
        setUpToDate(false);
      } finally { setReady(true); }
    }
    loadAuthStatus();
  }, [auth, isClient]);

  if (!ready) return null;
  return (
    <AuthContext.Provider value={{ auth, isAuthenticated, isUpToDate: upToDate, isClient }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>");
  return ctx;
}
