import { useState, useEffect, useRef } from 'react';
import { motion } from "motion/react";
import { Loader2, CheckCircle2, XCircle, ArrowRight } from 'lucide-react';
import { useAuth } from "../../AuthProvider";
import { Button } from './Shared/Button';

export interface ModernVerifyEmailProps {
  token: string;
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  signInRedirect?: () => void;
}

type VerifyStatus = "verifying" | "success" | "error" | "already_verified";

export function ModernVerifyEmail({
  token,
  onSuccess,
  onFailed,
  signInRedirect,
}: ModernVerifyEmailProps) {
  const [status, setStatus] = useState<VerifyStatus>("verifying");
  const [message, setMessage] = useState("Verificando seu e-mail...");
  const [isLoading, setIsLoading] = useState(true);
  const { auth } = useAuth();
  const hasAttempted = useRef(false);

  const handleVerify = async () => {
    if (hasAttempted.current) return;
    hasAttempted.current = true;

    setIsLoading(true);
    setStatus("verifying");

    try {
      const profile = auth.profile();
      if (profile?.is_verified) {
        setStatus("already_verified");
        setMessage("Seu e-mail já está verificado.");
        setIsLoading(false);
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
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (token && !hasAttempted.current) {
      void handleVerify();
    }
  }, [token]);

  return (
    <div className="w-full max-w-md z-10 flex flex-col antialiased selection:bg-primary/10 selection:text-primary">
      <div className="text-center mb-6">
        <motion.div
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          className="flex justify-center mb-4"
        >
          {status === "verifying" && <Loader2 className="w-12 h-12 text-primary animate-spin" />}
          {(status === "success" || status === "already_verified") && <CheckCircle2 className="w-12 h-12 text-green-500" />}
          {status === "error" && <XCircle className="w-12 h-12 text-destructive" />}
        </motion.div>

        <motion.h1
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="font-heading text-3xl font-bold tracking-tight mb-2"
        >
          {status === "verifying" ? "Verificando..." :
            status === "error" ? "Falha na verificação" : "E-mail Verificado!"}
        </motion.h1>
        <p className="text-muted-foreground text-sm mb-6">
          {message}
        </p>

        {status === "error" && (
          <Button
            onClick={() => {
              hasAttempted.current = false;
              void handleVerify();
            }}
            disabled={isLoading}
            className="w-full flex items-center justify-center gap-2"
          >
            Tentar novamente
          </Button>
        )}

        {(status === "success" || status === "already_verified") && (
          <div className="space-y-4">
            <p className="text-xs text-muted-foreground">
              Você já pode continuar navegando em sua conta.
            </p>
            {signInRedirect && (
              <Button
                onClick={signInRedirect}
                className="w-full flex items-center justify-center gap-2"
              >
                Ir para o login
                <ArrowRight className="w-4 h-4" />
              </Button>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
