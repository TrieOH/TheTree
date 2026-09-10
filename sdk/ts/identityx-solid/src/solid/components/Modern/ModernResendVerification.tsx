import { createSignal } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button } from "./Shared";
export interface ModernResendVerificationProps {
  email: string;
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
}
export function ModernResendVerification(props: ModernResendVerificationProps) {
  const { auth } = useAuth();
  const [loading, setLoading] = createSignal(false);
  const [sent, setSent] = createSignal(false);
  const resend = async () => {
    setLoading(true);
    try {
      const result = await auth.resendVerifyEmail(props.email);
      if (result.success) {
        setSent(true);
        await props.onSuccess?.(result.message);
      } else await props.onFailed?.(result.message, result.trace);
    } finally {
      setLoading(false);
    }
  };
  return (
    <div class="flex flex-col gap-4">
      <Button onClick={() => void resend()} disabled={loading()}>
        {sent()
          ? "E-mail reenviado"
          : loading()
            ? "Enviando..."
            : "Reenviar verificação"}
      </Button>
    </div>
  );
}
