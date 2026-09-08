import { createSignal, onSettled } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { FormError, FormInput, OAuthDivider, Button } from "./Shared";
import { OAuthProviderButton } from "./Shared/OAuthProviderButton";
import type { OAuthProviderI } from "@trieoh/identityx-sdk-ts";
import { ArrowRight, Loader2 } from "./Shared/Icons";

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
  const [oauthLoading, setOauthLoading] = createSignal<OAuthProviderI | null>(null);
  const [providers, setProviders] = createSignal<OAuthProviderI[]>(props.providers ?? []);

  onSettled(() => {
    if (props.providers) return;
    let cancelled = false;
    auth.getOAuthProviders().then((res) => {
      if (!cancelled && res.success) {
        setProviders(
          res.data
            .map(({ provider }) => provider)
            .filter((provider): provider is OAuthProviderI => provider === "google" || provider === "github"),
        );
      }
    }).catch(() => {});
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
    } catch {
      await props.onFailed?.("Ocorreu um erro inesperado");
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
    <form onSubmit={submit} class="space-y-4">
      <div class="space-y-1">
        <FormInput
          label="E-mail"
          name="email"
          type="email"
          value={email()}
          autocomplete="email"
          error={!!emailError()}
          onInput={(event) => setEmail(event.currentTarget.value)}
        />
        <FormError message={emailError()} />
      </div>

      <div class="space-y-1">
        <FormInput
          label="Senha"
          name="password"
          type="password"
          value={password()}
          autocomplete="current-password"
          error={!!passwordError()}
          onInput={(event) => setPassword(event.currentTarget.value)}
        />
        <FormError message={passwordError()} />
      </div>

      {props.forgotPasswordRedirect && (
        <div class="flex justify-end">
          <button
            type="button"
            onClick={props.forgotPasswordRedirect}
            class="text-xs text-primary hover:underline font-medium"
          >
            Esqueceu a senha?
          </button>
        </div>
      )}

      <Button
        type="submit"
        disabled={loading()}
        class="w-full flex items-center justify-center gap-2"
      >
        {loading() ? (
          <Loader2 class="w-5 h-5 animate-spin" />
        ) : (
          <>
            Entrar
            <ArrowRight class="w-4 h-4" />
          </>
        )}
      </Button>

      {providers().length > 0 && (
        <>
          <OAuthDivider />
          <div class="flex flex-col gap-2">
            {providers().map((provider) => (
              <OAuthProviderButton
                provider={provider}
                onClick={() => loginWithProvider(provider)}
                isLoading={oauthLoading() === provider}
              />
            ))}
          </div>
        </>
      )}
    </form>
  );
}
