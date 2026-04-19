import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useSyncExternalStore
} from "react";
import { Api } from "../core/api";
import { createAuthService } from "../core/services";
import { getTokenClaims, type AuthTokenClaims } from "../utils/token-utils";
import { validateProjectKey } from "../utils/env-validator";
import { configure } from "../core/env";
import { authStore } from "../store/auth-store";
import { logger, type DefaultFetchClientConfig } from "@soramux/node-fetch-sdk";

type AuthContextType = {
  auth: ReturnType<typeof createAuthService>;
  isAuthenticated: boolean;
  isInitializing: boolean;
  isClient?: boolean;
};

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({
  children,
  baseURL,
  projectId,
  isClient = true,
  fallback,
  waitSession = true,
  clientConfig,
}: {
  children: React.ReactNode;
  baseURL?: string;
  projectId?: string;
  isClient?: boolean;
  /** Component to show while initial auth check is in progress */
  fallback?: React.ReactNode;
  /** Whether to wait for the session restoration before rendering children. Defaults to true. */
  waitSession?: boolean;
  /** Extra config forwarded to the API client (e.g. timeout) */
  clientConfig?: Omit<DefaultFetchClientConfig, "adapter">;
}) {
  const isRestoring = useRef(false);

  const { isAuthenticated, isInitializing } = useSyncExternalStore(
    authStore.subscribe,
    authStore.getSnapshot,
    authStore.getServerSnapshot,
  );

  useEffect(() => {
    configure({
      ...(projectId ? { PROJECT_ID: projectId } : {}),
      ...(baseURL ? { BASE_URL: baseURL } : {}),
    });
  }, [projectId, baseURL]);

  const onTokenRefreshed = useCallback((claims: AuthTokenClaims) => {
    authStore.set({
      isAuthenticated: !!claims.access_data,
      isInitializing: false,
    });
  }, []);

  const apiInstance = useMemo(() => new Api(
    baseURL,
    undefined,
    onTokenRefreshed,
    clientConfig,
  ), [baseURL, onTokenRefreshed, clientConfig]);

  const auth = useMemo(
    () => createAuthService(apiInstance),
    [apiInstance],
  );

  useEffect(() => {
    if (isClient) validateProjectKey();

    const restoreSession = async () => {
      if (isRestoring.current) return;
      isRestoring.current = true;

      if (getTokenClaims()) {
        authStore.set({ isAuthenticated: true, isInitializing: false });
        return;
      }

      logger.log("No cached claims, attempting silent refresh...");
      try {
        const res = await auth.refresh();
        if (res.success) {
          authStore.set({ isAuthenticated: true, isInitializing: false });
          logger.log("Session restored.");
        } else {
          authStore.reset();
          logger.warn("No active session.");
        }
      } catch {
        authStore.reset();
        logger.warn("Could not restore session (offline?).");
      } finally {
        authStore.set({ isInitializing: false });
      }
    };

    restoreSession();
  }, [auth, isClient]);

  const contextValue = useMemo(() => ({
    auth,
    isAuthenticated,
    isInitializing,
    isClient
  }), [auth, isAuthenticated, isInitializing, isClient]);

  if (waitSession && isInitializing) return fallback ?? null;

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>");
  return ctx;
}
