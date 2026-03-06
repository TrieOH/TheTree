import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { Api } from "../core/api";
import { createAuthService } from "../core/services";
import { getTokenClaims } from "../utils/token-utils";
import { validateProjectKey } from "../utils/env-validator";
import { configure } from "../core/env";

type AuthContextType = {
  auth: ReturnType<typeof createAuthService>;
  isAuthenticated: boolean;
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

  // Apply manual configuration if provided
  useMemo(() => {
    if (projectId || baseURL) {
      configure({
        ...(projectId ? { PROJECT_ID: projectId } : {}),
        ...(baseURL ? { BASE_URL: baseURL } : {}),
      });
    }
  }, [projectId, baseURL]);

  const apiInstance = useMemo(() => new Api(baseURL), [baseURL]);
  const auth = useMemo(() => createAuthService(apiInstance), [apiInstance]);

  useEffect(() => {
    if (isClient) validateProjectKey();

    const loadAuthStatus = async () => {
      if(getTokenClaims()) {
        setIsAuthenticated(true);
        setReady(true);
        return;
      }
      console.log("[TRIEOH SDK] Attempting to refresh session...");
      try {
        const res = await auth.refreshProfileInfo();
        if (res.code === 200) {
          setIsAuthenticated(true);
          console.log("[TRIEOH SDK] Session restored successfully.");
        } else {
          setIsAuthenticated(false);
          console.warn("[TRIEOH SDK] Session restoration failed/no session.");
        }
      } catch(error) {
        console.warn("[TRIEOH SDK] Unable to verify session (offline?)");
        setIsAuthenticated(false);
      } finally { setReady(true); }
    }
    loadAuthStatus();
  }, [auth, isClient]);

  if (!ready) return null;
  return (
    <AuthContext.Provider value={{ auth, isAuthenticated, isClient }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>");
  return ctx;
}
