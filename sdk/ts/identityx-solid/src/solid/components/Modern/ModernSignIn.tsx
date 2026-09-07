import { createSignal, onSettled } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { FormError, FormInput, OAuthDivider, Button } from "./Shared";
import type { OAuthProviderI } from "@trieoh/identityx-sdk-ts";

export interface ModernSignInProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  signUpRedirect?: () => void;
  forgotPasswordRedirect?: () => void;
  providers?: OAuthProviderI[];
}

export function ModernSignIn(props: ModernSignInProps) {
  const { auth } = useAuth();
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [emailError, setEmailError] = createSignal("");
  const [passwordError, setPasswordError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [oauthLoading, setOauthLoading] = createSignal<OAuthProviderI | null>(
    null,
  );
  const [providers, setProviders] = createSignal<OAuthProviderI[]>(
    props.providers ?? [],
  );

  onSettled(() => {
    if (props.providers) return;
    void auth
      .getOAuthProviders()
      .then((result) => {
        if (result.success) {
          setProviders(
            result.data
              .map(({ provider }) => provider)
              .filter(
                (provider): provider is OAuthProviderI =>
                  provider === "google" || provider === "github",
              ),
          );
        }
      })
      .catch(() => undefined);
  });

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    setEmailError(email() ? "" : "E-mail é obrigatório");
    setPasswordError(password() ? "" : "A senha é obrigatória");
    if (!email() || !password()) return;
    setLoading(true);
    try {
      const result = await auth.login(email(), password());
      if (result.success) await props.onSuccess?.(result.message);
      else await props.onFailed?.(result.message, result.trace);
    } finally {
      setLoading(false);
    }
  };

  const loginWithProvider = async (provider: OAuthProviderI) => {
    setOauthLoading(provider);
    try {
      const result = await auth.loginWithProvider(provider);
      if (result.success && result.data?.url)
        window.location.href = result.data.url;
      else await props.onFailed?.(result.message || "Falha na autenticação");
    } catch {
      await props.onFailed?.("Ocorreu um erro inesperado");
    } finally {
      setOauthLoading(null);
    }
  };

  return (
    <form onSubmit={submit} class="w-full flex flex-col gap-4">
      <FormInput
        label="E-mail"
        name="email"
        type="email"
        value={email()}
        autocomplete="email"
        onInput={(event) => setEmail(event.currentTarget.value)}
      />
      <FormError message={emailError()} />
      <FormInput
        label="Senha"
        name="password"
        type="password"
        value={password()}
        autocomplete="current-password"
        onInput={(event) => setPassword(event.currentTarget.value)}
      />
      <FormError message={passwordError()} />
      {props.forgotPasswordRedirect && (
        <button
          type="button"
          onClick={props.forgotPasswordRedirect}
          class="self-end text-xs text-primary hover:underline"
        >
          Esqueceu sua senha?
        </button>
      )}
      <Button type="submit" disabled={loading()} class="h-12 px-4">
        {loading() ? "Entrando..." : "Entrar"}
      </Button>
      {providers().length > 0 && (
        <>
          <OAuthDivider />
          {providers().map((provider) => (
            <Button
              type="button"
              variant="outline"
              disabled={oauthLoading() !== null}
              onClick={() => void loginWithProvider(provider)}
              class="h-12"
            >
              {oauthLoading() === provider
                ? "Conectando..."
                : `Continuar com ${provider === "github" ? "GitHub" : "Google"}`}
            </Button>
          ))}
        </>
      )}
      {props.signUpRedirect && (
        <p class="text-center text-sm text-muted-foreground">
          Ainda não tem uma conta?{" "}
          <button
            type="button"
            onClick={props.signUpRedirect}
            class="text-primary font-semibold hover:underline"
          >
            Cadastre-se
          </button>
        </p>
      )}
    </form>
  );
}
