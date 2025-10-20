import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { Api } from "../core/api";
import { createAuthService } from "../core/services";

type AuthContextType = {
  auth: ReturnType<typeof createAuthService>;
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

  useEffect(() => {
    requestAnimationFrame(() => requestAnimationFrame(() => setReady(true)));
  }, []);
  const apiInstance = useMemo(() => new Api(baseURL), [baseURL]);
  const auth = useMemo(() => createAuthService(apiInstance), [apiInstance]);
  if (!ready) return <div aria-hidden style={{ minHeight: 40, minWidth: 120 }} />;
  return (
    <AuthContext.Provider value={{ auth }}>{children}</AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>");
  return ctx;
}
