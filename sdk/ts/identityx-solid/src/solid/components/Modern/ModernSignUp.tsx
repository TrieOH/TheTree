import { createSignal } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button, FormError, FormInput } from "./Shared";
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
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    if (!email() || password().length < 8 || password() !== confirm()) {
      setError("Confira o e-mail e as senhas informadas.");
      return;
    }
    setError("");
    setLoading(true);
    try {
      const result = await auth.register(email(), password());
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
      <FormInput
        label="Senha"
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
        {loading() ? "Criando..." : "Criar conta"}
      </Button>
      {props.signInRedirect && (
        <button
          type="button"
          onClick={props.signInRedirect}
          class="text-sm text-primary hover:underline"
        >
          Já possui uma conta? Entrar
        </button>
      )}
    </form>
  );
}
