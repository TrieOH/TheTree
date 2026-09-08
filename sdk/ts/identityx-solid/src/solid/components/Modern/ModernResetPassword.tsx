import { createSignal } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button, FormError, FormInput } from "./Shared";
import { ArrowRight, Loader2 } from "./Shared/Icons";

export interface ModernResetPasswordProps {
  token: string;
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  signInRedirect?: () => void;
}

export function ModernResetPassword(props: ModernResetPasswordProps) {
  const { auth } = useAuth();
  const [password, setPassword] = createSignal("");
  const [confirm, setConfirm] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [success, setSuccess] = createSignal(false);

  const submit = async (e: SubmitEvent) => {
    e.preventDefault();
    if (password().length < 8 || password() !== confirm()) {
      return setError("As senhas não coincidem ou são muito curtas.");
    }
    setError("");
    setLoading(true);
    try {
      const result = await auth.resetPassword(props.token, password());
      if (result.success) {
        setSuccess(true);
        await props.onSuccess?.(result.message);
      } else {
        await props.onFailed?.(result.message, result.trace);
      }
    } catch {
      await props.onFailed?.("Ocorreu um erro inesperado");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="w-full max-w-md mx-auto space-y-6">
      <div class="text-center space-y-2">
        <h1 class="font-heading text-3xl font-bold tracking-tight">
          Redefinir Senha
        </h1>
        <p class="text-muted-foreground text-sm">
          Crie uma nova senha para sua conta
        </p>
      </div>

      <form onSubmit={submit} class="space-y-4">
        {success() && (
          <div class="p-4 bg-emerald-500/10 text-emerald-600 rounded-md text-sm text-center">
            Senha redefinida com sucesso!
          </div>
        )}

        <div class="space-y-1">
          <FormInput
            label="Nova senha"
            name="newPassword"
            type="password"
            autocomplete="new-password"
            value={password()}
            onInput={(e) => setPassword(e.currentTarget.value)}
          />
        </div>
        <div class="space-y-1">
          <FormInput
            label="Confirmar senha"
            name="confirmPassword"
            type="password"
            autocomplete="new-password"
            value={confirm()}
            onInput={(e) => setConfirm(e.currentTarget.value)}
          />
          <FormError message={error()} />
        </div>

        <Button
          type="submit"
          disabled={loading() || success()}
          class="w-full flex items-center justify-center gap-2"
        >
          {loading() ? (
            <Loader2 class="w-5 h-5 animate-spin" />
          ) : (
            <>
              Redefinir senha
              <ArrowRight class="w-4 h-4" />
            </>
          )}
        </Button>
      </form>
      
      {props.signInRedirect && (
        <div class="text-center">
          <button
            type="button"
            onClick={props.signInRedirect}
            class="text-sm font-medium text-muted-foreground hover:text-primary transition-colors"
          >
            Voltar para o login
          </button>
        </div>
      )}
    </div>
  );
}
