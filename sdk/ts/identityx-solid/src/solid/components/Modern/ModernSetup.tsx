import { createSignal } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button, FormError, FormInput } from "./Shared";
import { ArrowRight, Loader2, Sparkles } from "./Shared/Icons";

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
    } catch {
      await props.onFailed?.("Ocorreu um erro inesperado");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="w-full max-w-md mx-auto space-y-6">
      <div class="text-center space-y-2 mb-8">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-primary/10 text-primary mb-4">
          <Sparkles class="w-6 h-6" />
        </div>
        <h1 class="font-heading text-3xl font-bold tracking-tight">
          Configuração Inicial
        </h1>
        <p class="text-muted-foreground text-sm">
          Crie a conta de administrador do sistema
        </p>
      </div>

      <form onSubmit={submit} class="space-y-4 bg-card p-6 rounded-2xl border shadow-sm">
        <div class="space-y-1">
          <FormInput
            label="E-mail do Administrador"
            type="email"
            autocomplete="email"
            value={email()}
            onInput={(e) => setEmail(e.currentTarget.value)}
          />
        </div>
        <div class="space-y-1">
          <FormInput
            label="Senha (Mín. 8 caracteres)"
            type="password"
            value={password()}
            onInput={(e) => setPassword(e.currentTarget.value)}
          />
          <FormError message={error()} />
        </div>

        <Button
          type="submit"
          disabled={loading()}
          class="w-full flex items-center justify-center gap-2 mt-2"
        >
          {loading() ? (
            <Loader2 class="w-5 h-5 animate-spin" />
          ) : (
            <>
              Concluir configuração
              <ArrowRight class="w-4 h-4" />
            </>
          )}
        </Button>
      </form>
    </div>
  );
}
