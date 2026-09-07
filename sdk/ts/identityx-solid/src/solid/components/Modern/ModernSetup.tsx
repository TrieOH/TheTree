import { createSignal } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button, FormError, FormInput } from "./Shared";
export interface ModernSetupProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
}
export function ModernSetup(props: ModernSetupProps) {
  const { auth } = useAuth();
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const submit = async (e: SubmitEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const result = await auth.setup(email(), password());
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
      <FormError message={error()} />
      <Button type="submit" disabled={loading()}>
        {loading() ? "Configurando..." : "Configurar"}
      </Button>
    </form>
  );
}
