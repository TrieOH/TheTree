import { createSignal } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button, FormError, FormInput } from "./Shared";
import { ArrowRight, Loader2 } from "./Shared/Icons";

export interface ModernForgotPasswordProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  signInRedirect?: () => void;
}

export function ModernForgotPassword(props: ModernForgotPasswordProps) {
  const { auth } = useAuth();
  const [email, setEmail] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [successMessage, setSuccessMessage] = createSignal("");

  const submit = async (e: SubmitEvent) => {
    e.preventDefault();
    if (!email()) return setError("E-mail é obrigatório");
    setLoading(true);
    try {
      const result = await auth.sendForgotPassword(email());
      if (result.success) {
        setSuccessMessage("Link enviado para o seu e-mail.");
        await props.onSuccess?.(result.message);
      }
      else await props.onFailed?.(result.message, result.trace);
    } catch {
      await props.onFailed?.("Ocorreu um erro inesperado");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={submit} class="space-y-4">
      {successMessage() && (
        <div class="p-4 bg-primary/10 text-primary rounded-md text-sm text-center">
          {successMessage()}
        </div>
      )}

      <div class="space-y-1">
        <FormInput
          label="E-mail"
          type="email"
          autocomplete="email"
          value={email()}
          onInput={(e) => setEmail(e.currentTarget.value)}
        />
        <FormError message={error()} />
      </div>
      
      <Button
        type="submit"
        disabled={loading() || !!successMessage()}
        class="w-full flex items-center justify-center gap-2"
      >
        {loading() ? (
          <Loader2 class="w-5 h-5 animate-spin" />
        ) : (
          <>
            Enviar link de recuperação
            <ArrowRight class="w-4 h-4" />
          </>
        )}
      </Button>
    </form>
  );
}
