import { createSignal } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button, FormError, FormInput } from "./Shared";
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
  const submit = async (e: SubmitEvent) => {
    e.preventDefault();
    if (!email()) return setError("E-mail é obrigatório");
    setLoading(true);
    try {
      const result = await auth.sendForgotPassword(email());
      if (result.success) await props.onSuccess?.(result.message);
      else await props.onFailed?.(result.message, result.trace);
    } finally {
      setLoading(false);
    }
  };
  return (
    <form onSubmit={submit} class="w-full flex flex-col gap-4">
      <FormInput
        label="E-mail"
        type="email"
        value={email()}
        onInput={(e) => setEmail(e.currentTarget.value)}
      />
      <FormError message={error()} />
      <Button type="submit" disabled={loading()}>
        {loading() ? "Enviando..." : "Enviar link"}
      </Button>
      {props.signInRedirect && (
        <button
          type="button"
          onClick={props.signInRedirect}
          class="text-sm text-primary hover:underline"
        >
          Voltar ao login
        </button>
      )}
    </form>
  );
}
