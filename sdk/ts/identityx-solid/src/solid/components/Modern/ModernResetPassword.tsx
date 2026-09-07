import { createSignal } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button, FormError, FormInput } from "./Shared";
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
  const submit = async (e: SubmitEvent) => {
    e.preventDefault();
    if (password().length < 8 || password() !== confirm())
      return setError("A senha deve ter pelo menos 8 caracteres e coincidir.");
    setLoading(true);
    try {
      const result = await auth.resetPassword(props.token, password());
      if (result.success) await props.onSuccess?.(result.message);
      else await props.onFailed?.(result.message, result.trace);
    } finally {
      setLoading(false);
    }
  };
  return (
    <form onSubmit={submit} class="w-full flex flex-col gap-4">
      <FormInput
        label="Nova senha"
        type="password"
        value={password()}
        onInput={(e) => setPassword(e.currentTarget.value)}
      />
      <FormInput
        label="Confirmar senha"
        type="password"
        value={confirm()}
        onInput={(e) => setConfirm(e.currentTarget.value)}
      />
      <FormError message={error()} />
      <Button type="submit" disabled={loading()}>
        {loading() ? "Salvando..." : "Redefinir senha"}
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
