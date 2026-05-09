import { useEffect, useState, useRef, type MouseEvent } from "react";
import { useAuth } from "../../AuthProvider";
import BasicSubmitButton from "../Form/BasicSubmitButton";
import CardAvatar from "../Form/CardAvatar";
import { FiCheck, FiInfo } from "react-icons/fi";

export interface VerifyEmailProps {
  token: string;
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
}

type VerifyStatus = "verifying" | "success" | "error" | "already_verified";

/**
 * Hook to manage verification logic and allow mocking for Storybook
 * without exposing internal props in the public VerifyEmail component.
 */
export function useVerifyEmailLogic(token: string, onSuccess?: (m?: string) => void, onFailed?: (m: string, t?: string[]) => void) {
  const [loading, setLoading] = useState(true);
  const [status, setStatus] = useState<VerifyStatus>("verifying");
  const [message, setMessage] = useState("Verificando seu e-mail...");
  const { auth } = useAuth();
  const hasAttempted = useRef(false);

  const handleVerify = async () => {
    if (hasAttempted.current) return;
    hasAttempted.current = true;

    setLoading(true);
    setStatus("verifying");

    try {
      const profile = auth.profile();
      if (profile?.is_verified) {
        setStatus("already_verified");
        setMessage("Seu e-mail já está verificado.");
        setLoading(false);
        return;
      }

      const res = await auth.verifyEmail(token);

      if (res.success) {
        await auth.refresh();
        setStatus("success");
        setMessage(res.message || "E-mail verificado com sucesso!");
        if (onSuccess) await onSuccess(res.message);
      } else {
        setStatus("error");
        setMessage(res.message || "Falha na verificação do e-mail.");
        if (onFailed) await onFailed(res.message, res.trace);
      }
    } catch (error) {
      setStatus("error");
      setMessage("Ocorreu um erro inesperado durante a verificação.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (token && !hasAttempted.current) handleVerify();
  }, [token]);

  const onRetry = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    hasAttempted.current = false;
    handleVerify();
  };

  return {
    loading,
    status,
    message,
    onRetry,
    setStatus,
    setLoading,
    setMessage
  };
}

export function VerifyEmail({
  token,
  onSuccess,
  onFailed,
}: VerifyEmailProps) {
  const { loading, status, message, onRetry } = useVerifyEmailLogic(token, onSuccess, onFailed);

  const mainText = status === "verifying" ? "Verificando..." :
    status === "error" ? "Ops! Falha na verificação" : "Tudo pronto!";

  return (
    <div className="font-inter flex flex-col w-full h-full min-h-screen items-center justify-center bg-trieoh-neutral1 text-trieoh-neutral2 p-6">
      <div className="w-full max-w-[30rem] flex flex-col gap-8 items-center bg-trieoh-neutral1 text-trieoh-neutral2 p-[1.25rem_1.5rem] shadow-[0_0.25rem_0.25rem_0_rgba(0,0,0,0.25)] rounded-[0.25rem]">
        <CardAvatar
          mainText={mainText}
          subText={message}
        />

        {status === "verifying" && (
          <div className="w-full flex justify-center py-4">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-trieoh-secondary"></div>
          </div>
        )}

        {status === "error" && (
          <div className="w-full flex flex-col gap-6 items-center">
            <div className="w-full mt-2">
              <BasicSubmitButton
                label="Tentar novamente"
                onSubmit={onRetry}
                loading={loading}
              />
            </div>
          </div>
        )}

        {(status === "success" || status === "already_verified") && (
          <div className="w-full flex flex-col gap-4 items-center">
            <div className="w-16 h-16 flex items-center justify-center rounded-full bg-trieoh-secondary/10 text-trieoh-secondary">
              {status === "success" ? <FiCheck size={40} /> : <FiInfo size={40} />}
            </div>
            <p className="text-trieoh-sm font-semibold opacity-50 text-center">
              Você já pode fechar esta janela ou continuar navegando em sua conta.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
