import {
  createComponent,
  createContext,
  createSignal,
  onCleanup,
  onSettled,
  useContext,
  type Accessor,
  type ParentProps,
} from "solid-js";

import { logger, type DefaultFetchClientConfig } from "@trieoh/envoy-fetch-ts";

import {
  Api,
  authStore,
  configure,
  createAuthService,
  getTokenClaims,
  isRefreshSessionExpired,
  validateProjectKey,
  type AuthCallbacks,
  type AuthService,
  type AuthTokenClaims,
  type TokenSubject,
} from "@trieoh/identityx-sdk-ts";

export interface AuthProviderProps extends ParentProps, AuthCallbacks {
  baseURL?: string;
  projectId?: string;
  isProjectMode?: boolean;
  waitSession?: boolean;
  fallback?: ParentProps["children"];
  clientConfig?: Omit<DefaultFetchClientConfig, "adapter">;
  adapter?: AuthProviderAdapter;
}

export interface AuthProviderAdapterContext {
  callbacks: AuthCallbacks;
  defaultAuth: AuthService;
  getProfile(): TokenSubject | null;
  setProfile(profile: TokenSubject | null): void;
  setAuthenticated(authenticated: boolean): void;
}

export interface RestoredAuthSession {
  isAuthenticated: boolean;
  profile: TokenSubject | null;
}

export interface AuthProviderAdapter {
  createAuth(context: AuthProviderAdapterContext): AuthService;
  restoreSession(): Promise<boolean | RestoredAuthSession>;
}

export interface AuthContextValue {
  auth: AuthService;
  isAuthenticated: Accessor<boolean>;
  isInitializing: Accessor<boolean>;
  isProjectMode?: boolean;
}

const AuthContext = createContext<AuthContextValue>();

export function AuthProvider(props: AuthProviderProps) {
  const [state, setState] = createSignal(authStore.getSnapshot());
  const [profile, setProfile] = createSignal<TokenSubject | null>(null);
  const unsubscribe = authStore.subscribe(() =>
    setState(authStore.getSnapshot()),
  );
  onCleanup(unsubscribe);

  configure({
    ...(props.projectId ? { PROJECT_ID: props.projectId } : {}),
    ...(props.baseURL ? { BASE_URL: props.baseURL } : {}),
  });

  const onTokenRefreshed = (claims: AuthTokenClaims) => {
    authStore.set({
      isAuthenticated: !!claims.access_data,
      isInitializing: false,
    });
    props.onRefresh?.();
  };
  const api = new Api(
    props.baseURL,
    undefined,
    onTokenRefreshed,
    props.clientConfig,
  );
  const defaultAuth = createAuthService(api, props);
  const auth = props.adapter
    ? props.adapter.createAuth({
      callbacks: props,
      defaultAuth,
      getProfile: profile,
      setProfile,
      setAuthenticated: (authenticated) =>
        authStore.set({
          isAuthenticated: authenticated,
          isInitializing: false,
        }),
    })
    : defaultAuth;

  onSettled(
    () =>
      void (async () => {
        if (props.isProjectMode !== false) validateProjectKey();
        if (props.adapter) {
          try {
            const restored = await props.adapter.restoreSession();
            if (typeof restored === "boolean")
              authStore.set({
                isAuthenticated: restored,
                isInitializing: false,
              });
            else {
              setProfile(restored.profile);
              authStore.set({
                isAuthenticated: restored.isAuthenticated,
                isInitializing: false,
              });
            }
          } catch {
            authStore.reset();
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
        try {
          await api.interceptor.refreshToken();
        } catch {
          authStore.reset();
          logger.warn("Could not restore session.");
        } finally {
          authStore.set({ isInitializing: false });
        }
      })(),
  );

  const value: AuthContextValue = {
    auth,
    isAuthenticated: () => state().isAuthenticated,
    isInitializing: () => state().isInitializing,
    isProjectMode: props.isProjectMode,
  };
  return createComponent(AuthContext, {
    value,
    get children() {
      return props.waitSession !== false && state().isInitializing
        ? props.fallback
        : props.children;
    },
  });
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used inside <AuthProvider>");
  return context;
}

export type {
  AuthCallbacks,
  AuthService,
  ActorType,
  TokenSubject,
} from "@trieoh/identityx-sdk-ts";
