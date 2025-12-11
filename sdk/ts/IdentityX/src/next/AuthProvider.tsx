import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { Api } from "../core/api";
import { createAuthService } from "../core/services";
import { getTokenClaims, isRefreshSessionExpired } from "../utils/token-utils";

type AuthContextType = {
  auth: ReturnType<typeof createAuthService>;
  isAuthenticated: boolean;
};

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({
  children,
  baseURL,
}: {
  children: React.ReactNode;
  baseURL?: string;
}) {
  const [ready, setReady] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const apiInstance = useMemo(() => new Api(baseURL), [baseURL]);
  const auth = useMemo(() => createAuthService(apiInstance), [apiInstance]);

  useEffect(() => {
    const loadAuthStatus = async () => {
      if(getTokenClaims()) {
        setIsAuthenticated(true);
        setReady(true);
        return;
      }
      if(isRefreshSessionExpired()) {
        console.warn("[AuthProvider] Persistent session is expired. Skipping refresh request.");
        setReady(true);
        return;
      }
      console.log("[TRIEOH SDK] Attempting to refresh session...");
      try {
        const res = await auth.refresh();
        if (res.code === 200) {
          setIsAuthenticated(true);
          console.log("[TRIEOH SDK] Session restored successfully.");
        } else {
          setIsAuthenticated(false);
          console.warn("[TRIEOH SDK] Session restoration failed/no session.");
        }
      } catch(error) {
        setIsAuthenticated(false);
        console.error("[TRIEOH SDK] Session bootstrap error:", error);
      } finally {
        setReady(true);
      }
    }
    loadAuthStatus();
  }, [auth]);

  if (!ready) return <div aria-hidden style={{ minHeight: 40, minWidth: 120 }} />;
  return (
    <AuthContext.Provider value={{ auth, isAuthenticated }}>{children}</AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>");
  return ctx;
}
