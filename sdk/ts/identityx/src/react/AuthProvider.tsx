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
import {
  createAuthService,
  type AuthCallbacks,
  type AuthService,
} from "../core/services";
import { getTokenClaims, isRefreshSessionExpired } from "../utils/token-utils";
import { validateProjectKey } from "../utils/env-validator";
import { configure } from "../core/env";
import { authStore } from "../store/auth-store";
import { logger, type DefaultFetchClientConfig } from "@trieoh/envoy-fetch-ts";
import type { AuthTokenClaims, TokenSubject } from "../types/token-types";

type AuthContextType = {
  auth: AuthService;
  isAuthenticated: boolean;
  isInitializing: boolean;
  isProjectMode?: boolean;
};

const AuthContext = createContext<AuthContextType | null>(null);

export interface AuthProviderAdapterContext {
  callbacks: AuthCallbacks;
  defaultAuth: AuthService;
  getProfile(): TokenSubject | null;
  setProfile(profile: TokenSubject | null): void;
  setAuthenticated(isAuthenticated: boolean): void;
}

export interface RestoredAuthSession {
  isAuthenticated: boolean;
  profile: TokenSubject | null;
}

/**
 * Replaces browser-side IdentityX calls with an application-owned transport.
 * This is intended for BFF/server-function integrations where tokens remain in
 * HttpOnly cookies or server-side session storage.
 */
export interface AuthProviderAdapter {
  createAuth(context: AuthProviderAdapterContext): AuthService;
  restoreSession(): Promise<boolean | RestoredAuthSession>;
}

export interface AuthProviderProps extends AuthCallbacks {
  children: React.ReactNode;
  baseURL?: string;
  projectId?: string;
  isProjectMode?: boolean;
  /** Component to show while initial auth check is in progress */
  fallback?: React.ReactNode;
  /** Whether to wait for the session restoration before rendering children. Defaults to true. */
  waitSession?: boolean;
  /** Extra config forwarded to the API client (e.g. timeout) */
  clientConfig?: Omit<DefaultFetchClientConfig, "adapter">;
  /** Override auth calls and session restoration, typically with server functions. */
  adapter?: AuthProviderAdapter;
}

export function AuthProvider({
  children,
  baseURL,
  projectId,
  isProjectMode = true,
  fallback,
  waitSession = true,
  clientConfig,
  adapter,
  onLogin,
  onSetup,
  onResetPassword,
  onRegister,
  onVerify,
  onRefresh,
}: AuthProviderProps) {
  const isRestoring = useRef(false);
  const profile = useRef<TokenSubject | null>(null);

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
    onRefresh?.();
  }, [onRefresh]);

  const apiInstance = useMemo(() => new Api(
    baseURL,
    undefined,
    onTokenRefreshed,
    clientConfig,
  ), [baseURL, onTokenRefreshed, clientConfig]);

  const setAuthenticated = useCallback((authenticated: boolean) => {
    authStore.set({
      isAuthenticated: authenticated,
      isInitializing: false,
    });
  }, []);

  const callbacks = useMemo<AuthCallbacks>(() => ({
    onLogin,
    onSetup,
    onResetPassword,
    onRegister,
    onVerify,
    onRefresh,
  }), [onLogin, onSetup, onResetPassword, onRegister, onVerify, onRefresh]);

  const defaultAuth = useMemo(
    () => createAuthService(apiInstance, callbacks),
    [apiInstance, callbacks],
  );

  const auth = useMemo(
    () => adapter
      ? adapter.createAuth({
          callbacks,
          defaultAuth,
          getProfile: () => profile.current,
          setProfile: (value) => { profile.current = value; },
          setAuthenticated,
        })
      : defaultAuth,
    [adapter, defaultAuth, callbacks, setAuthenticated],
  );

  useEffect(() => {
    if (isProjectMode) validateProjectKey();

    const restoreSession = async () => {
      if (isRestoring.current) return;
      isRestoring.current = true;

      if (adapter) {
        try {
          const restored = await adapter.restoreSession();
          if (typeof restored === "boolean") {
            setAuthenticated(restored);
          } else {
            profile.current = restored.profile;
            setAuthenticated(restored.isAuthenticated);
          }
        } catch {
          setAuthenticated(false);
          logger.warn("Could not restore server-managed session.");
        }
        return;
      }

      if (getTokenClaims()) {
        authStore.set({ isAuthenticated: true, isInitializing: false });
        return;
      }

      if (isRefreshSessionExpired()) {
        authStore.reset();
        authStore.set({ isInitializing: false });
        return;
      }

      logger.log("No cached claims, attempting silent refresh...");
      try {
        await apiInstance.interceptor.refreshToken();
        logger.log("Session restored.");
      } catch {
        authStore.reset();
        logger.warn("Could not restore session (offline?).");
      } finally {
        authStore.set({ isInitializing: false });
      }
    };

    restoreSession();
  }, [adapter, apiInstance, isProjectMode, setAuthenticated]);

  const contextValue = useMemo(() => ({
    auth,
    isAuthenticated,
    isInitializing,
    isProjectMode
  }), [auth, isAuthenticated, isInitializing, isProjectMode]);

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
