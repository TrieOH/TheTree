import { createSignal } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button, FormError, FormInput } from "./Shared";
import { ArrowRight, Loader2 } from "./Shared/Icons";

export interface ModernSignUpProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  signInRedirect?: () => void;
}

export function ModernSignUp(props: ModernSignUpProps) {
  const { auth } = useAuth();
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [confirm, setConfirm] = createSignal("");
  const [emailError, setEmailError] = createSignal("");
  const [passwordError, setPasswordError] = createSignal("");
  const [confirmError, setConfirmError] = createSignal("");
  const [loading, setLoading] = createSignal(false);

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    setEmailError(email() ? "" : "E-mail inválido");
    setPasswordError(password().length >= 8 ? "" : "A senha deve ter pelo menos 8 caracteres");
    setConfirmError(password() === confirm() ? "" : "As senhas não coincidem");
    if (!email() || password().length < 8 || password() !== confirm()) {
      return;
    }
    setLoading(true);
    try {
      const result = await auth.register(email(), password());
      if (result.success) await props.onSuccess?.(result.message);
      else await props.onFailed?.(result.message, result.trace);
    } catch {
      await props.onFailed?.("Ocorreu um erro inesperado");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={submit} class="space-y-4">
      <div class="space-y-1">
        <FormInput
          label="E-mail"
          name="email"
          type="email"
          autocomplete="email"
          value={email()}
          error={!!emailError()}
          onInput={(e) => setEmail(e.currentTarget.value)}
        />
        <FormError id="email-error" message={emailError()} />
      </div>
      <div class="space-y-1">
        <FormInput
          label="Senha"
          name="password"
          type="password"
          autocomplete="new-password"
          value={password()}
          error={!!passwordError()}
          onInput={(e) => setPassword(e.currentTarget.value)}
        />
        <FormError id="password-error" message={passwordError()} />
      </div>
      <div class="space-y-1">
        <FormInput
          label="Confirmar senha"
          name="confirmPassword"
          type="password"
          autocomplete="new-password"
          value={confirm()}
          error={!!confirmError()}
          onInput={(e) => setConfirm(e.currentTarget.value)}
        />
        <FormError id="confirmPassword-error" message={confirmError()} />
      </div>

      <Button
        type="submit"
        disabled={loading()}
        class="w-full flex items-center justify-center gap-2"
      >
        {loading() ? (
          <Loader2 class="w-5 h-5 animate-spin" />
        ) : (
          <>
            Criar conta
            <ArrowRight class="w-4 h-4" />
          </>
        )}
      </Button>
    </form>
  );
}
