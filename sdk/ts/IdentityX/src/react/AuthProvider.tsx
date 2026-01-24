import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { Api } from "../core/api";
import { createAuthService } from "../core/services";
import { getTokenClaims } from "../utils/token-utils";

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
  }, [auth]);

  if (!ready) return null;
  return (
    <AuthContext.Provider value={{ auth, isAuthenticated }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>");
  return ctx;
}
