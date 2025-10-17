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
  // I need to load only when the style is loaded
  // const [ready, setReady] = useState(false);
  // useEffect(() => {
  //   const check = () => {
  //     const styleLoaded = Array.from(document.styleSheets).some(sheet =>
  //       sheet.href?.includes("trieoh") ||
  //       sheet.ownerNode?.textContent?.includes(".trieoh-")
  //     );
  //     if (styleLoaded) setReady(true);
  //   };

  //   check();
  //   const timeout = setTimeout(check, 100);
  //   return () => clearTimeout(timeout);
  // }, []);
  const apiInstance = useMemo(() => new Api(baseURL), [baseURL]);
  const auth = useMemo(() => createAuthService(apiInstance), [apiInstance]);
  // if(!ready) return;
  return (
    <AuthContext.Provider value={{ auth }}>{children}</AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>");
  return ctx;
}
