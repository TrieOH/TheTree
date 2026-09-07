import { createSignal, onSettled } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button } from "./Shared";
export interface ModernVerifyEmailProps {
  token: string;
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
}
export function ModernVerifyEmail(props: ModernVerifyEmailProps) {
  const { auth } = useAuth();
  const [message, setMessage] = createSignal("Verificando...");
  const [loading, setLoading] = createSignal(true);
  const verify = async () => {
    try {
      const result = await auth.verifyEmail(props.token);
      if (result.success) {
        setMessage(result.message || "E-mail verificado com sucesso!");
        await props.onSuccess?.(result.message);
      } else {
        setMessage(result.message);
        await props.onFailed?.(result.message, result.trace);
      }
    } catch {
      setMessage("Ocorreu um erro inesperado.");
    } finally {
      setLoading(false);
    }
  };
  onSettled(() => void verify());
  return (
    <div class="w-full flex flex-col items-center gap-4 text-center">
      <p>{message()}</p>
      {loading() && (
        <span class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
      )}{" "}
      {!loading() && (
        <Button
          onClick={() => {
            setLoading(true);
            void verify();
          }}
        >
          Tentar novamente
        </Button>
      )}
    </div>
  );
}
