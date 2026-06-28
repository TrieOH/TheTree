import { useState } from "react";
import { useAuth } from "../../AuthProvider";
import BasicSubmitButton from "../Form/BasicSubmitButton";
import CardAvatar from "../Form/CardAvatar";

export interface ResendVerifyEmailProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
}

export function ResendVerifyEmail({
  onSuccess,
  onFailed,
}: ResendVerifyEmailProps) {
  const [loading, setLoading] = useState(false);
  const [sent, setSent] = useState(false);
  const { auth } = useAuth();

  const handleResend = async () => {
    setLoading(true);
    const res = await auth.resendVerifyEmail();
    if (res.success) {
      setSent(true);
      if (onSuccess) await onSuccess(res.message);
    } else if (onFailed) {
      await onFailed(res.message, res.trace);
    }
    setLoading(false);
  };

  return (
    <div className="font-sans flex flex-col w-full max-w-120 min-w-60 max-h-[max(75dvh,32rem)] overflow-hidden gap-5 items-center bg-background text-foreground p-[1.25rem_1.5rem] shadow-[0_0.25rem_0.25rem_0_rgba(0,0,0,0.25)] transition-transform duration-150 ease-in-out rounded-lg">
      <CardAvatar
        mainText="Reenviar Verificação"
        subText="Não recebeu o e-mail de verificação? Clique no botão abaixo para reenviar."
      />

      <BasicSubmitButton
        label={sent ? "E-mail Reenviado" : "Reenviar E-mail de Verificação"}
        onSubmit={handleResend}
        loading={loading}
      />

      {sent && (
        <p className="text-sm font-semibold text-primary">
          Um novo link de verificação foi enviado para seu e-mail.
        </p>
      )}
    </div>
  );
}
