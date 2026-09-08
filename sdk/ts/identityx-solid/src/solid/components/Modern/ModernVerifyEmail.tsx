import { createSignal, onSettled } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button } from "./Shared";
import { Loader2, CheckCircle2, XCircle, ArrowRight } from "./Shared/Icons";

export interface ModernVerifyEmailProps {
  token: string;
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
}

export function ModernVerifyEmail(props: ModernVerifyEmailProps) {
  const { auth } = useAuth();
  const [status, setStatus] = createSignal<"loading" | "success" | "error">("loading");
  const [message, setMessage] = createSignal("Verificando seu e-mail...");

  const verify = async () => {
    setStatus("loading");
    setMessage("Verificando seu e-mail...");
    try {
      const result = await auth.verifyEmail(props.token);
      if (result.success) {
        setStatus("success");
        setMessage(result.message || "E-mail verificado com sucesso!");
        await props.onSuccess?.(result.message);
      } else {
        setStatus("error");
        setMessage(result.message || "Falha ao verificar e-mail.");
        await props.onFailed?.(result.message, result.trace);
      }
    } catch {
      setStatus("error");
      setMessage("Ocorreu um erro inesperado.");
    }
  };

  onSettled(() => {
    if (props.token) {
      void verify();
    } else {
      setStatus("error");
      setMessage("Token de verificação inválido ou ausente.");
    }
  });

  return (
    <div class="w-full max-w-md mx-auto">
      <div class="bg-card p-8 rounded-2xl border shadow-sm text-center flex flex-col items-center gap-6">
        {status() === "loading" && (
          <div class="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center text-primary">
            <Loader2 class="w-8 h-8 animate-spin" />
          </div>
        )}

        {status() === "success" && (
          <div class="w-16 h-16 rounded-full bg-emerald-500/10 flex items-center justify-center text-emerald-600">
            <CheckCircle2 class="w-8 h-8" />
          </div>
        )}

        {status() === "error" && (
          <div class="w-16 h-16 rounded-full bg-destructive/10 flex items-center justify-center text-destructive">
            <XCircle class="w-8 h-8" />
          </div>
        )}

        <div class="space-y-2">
          <h1 class="font-heading text-2xl font-bold tracking-tight">
            {status() === "loading" ? "Verificando..."
              : status() === "success" ? "Tudo certo!"
                : "Ops, algo deu errado"}
          </h1>
          <p class="text-muted-foreground text-sm max-w-62.5 mx-auto">
            {message()}
          </p>
        </div>

        {status() === "error" && (
          <Button
            onClick={verify}
            class="w-full mt-2 flex items-center justify-center gap-2"
          >
            Tentar novamente
            <ArrowRight class="w-4 h-4" />
          </Button>
        )}
      </div>
    </div>
  );
}
