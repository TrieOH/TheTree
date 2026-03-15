import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  useSyncExternalStore
} from "react";
import { Api } from "../core/api";
import { createAuthService } from "../core/services";
import { getTokenClaims, isUpToDate } from "../utils/token-utils";
import { validateProjectKey } from "../utils/env-validator";
import { configure } from "../core/env";
import { authStore } from "../store/auth-store";
import { logger } from "../utils/logger";

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
  exchangeURL,
  isClient = true,
}: {
  children: React.ReactNode;
  baseURL?: string;
  projectId?: string;
  exchangeURL?: string;
  isClient?: boolean;
}) {
  const [ready, setReady] = useState(false);

  const { isAuthenticated, isUpToDate: upToDate } = useSyncExternalStore(
    authStore.subscribe,
    authStore.getSnapshot,
    authStore.getServerSnapshot,
  );

  const prevConfig = useRef({ projectId, baseURL });
  useEffect(() => {
    const prev = prevConfig.current;
    if (prev.projectId === projectId && prev.baseURL === baseURL) return;

    configure({
      ...(projectId ? { PROJECT_ID: projectId } : {}),
      ...(baseURL ? { BASE_URL: baseURL } : {}),
    });

    prevConfig.current = { projectId, baseURL };
  }, [projectId, baseURL]);

  const apiInstance = useMemo(() => new Api(
    baseURL,
    undefined,
    (claims) => authStore.set({
      isAuthenticated: !!claims.access_data,
      isUpToDate: claims.is_up_to_date ?? false,
    }),
    exchangeURL,
  ), [baseURL, exchangeURL]);

  const auth = useMemo(
    () => createAuthService(apiInstance, exchangeURL),
    [apiInstance, exchangeURL],
  );

  useEffect(() => {
    if (isClient) validateProjectKey();

    const restoreSession = async () => {
      if (getTokenClaims()) {
        authStore.set({ isAuthenticated: true, isUpToDate: isUpToDate() });
        setReady(true);
        return;
      }

      logger.log("No cached claims, attempting silent refresh...");
      try {
        const res = await (exchangeURL ? auth.refresh() : auth.refreshProfileInfo());
        if (res.success) {
          authStore.set({ isAuthenticated: true, isUpToDate: isUpToDate() });
          logger.log("Session restored.");
        } else {
          authStore.reset();
          logger.warn("No active session.");
        }
      } catch {
        authStore.reset();
        logger.warn("Could not restore session (offline?).");
      } finally {
        setReady(true);
      }
    };

    restoreSession();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

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