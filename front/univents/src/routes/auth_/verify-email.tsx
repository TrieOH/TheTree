import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { AlertCircle, MailWarning, RefreshCw } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { z } from "zod";
import { safeInternalReturnTo } from "@/features/auths/lib/auth-path";
import { AuthActionPage } from "@/features/auths/ui/auth-action-page";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";

export const Route = createFileRoute("/auth_/verify-email")({
  validateSearch: (search) =>
    z
      .object({
        token: z.string().catch(""),
        returnTo: z.string().optional(),
      })
      .parse(search),
  component: VerifyEmailPage,
});

function VerifyEmailPage() {
  const { token, returnTo } = Route.useSearch();
  const navigate = useNavigate();
  const { auth } = useAuth();
  const [status, setStatus] = useState<
    "required" | "loading" | "error" | "success"
  >(token ? "loading" : "required");
  const [email, setEmail] = useState(auth.profile()?.email ?? "");
  const [resending, setResending] = useState(false);
  const [message, setMessage] = useState("");
  const isValidEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim());

  useEffect(() => {
    if (returnTo) localStorage.setItem("univents:verify-return", returnTo);
    if (!token) return;

    void auth
      .verifyEmail(token)
      .then(async (response) => {
        if (response.success) {
          const refreshed = await auth.refresh();
          if (refreshed.success) {
            const destination = safeInternalReturnTo(
              returnTo,
              localStorage.getItem("univents:verify-return"),
            );
            localStorage.removeItem("univents:verify-return");
            toast.success("E-mail verificado com sucesso.");
            await navigate({ to: destination, replace: true });
            return;
          }
          setStatus("success");
          return;
        }
        setStatus("error");
        setMessage(
          "Este link expirou ou não é válido. Solicite um novo link para verificar seu e-mail.",
        );
      })
      .catch(() => {
        setStatus("error");
        setMessage(
          "Não foi possível verificar este link. Solicite um novo para continuar.",
        );
        toast.error("Não foi possível verificar o link de e-mail.");
      });
  }, [auth, navigate, returnTo, token]);

  async function continueAfterVerification() {
    const response = await auth.refresh();
    if (!response.success || !auth.profile()?.verified_at) {
      toast.info("A verificação ainda não foi confirmada.");
      return;
    }
    const destination = safeInternalReturnTo(
      returnTo,
      localStorage.getItem("univents:verify-return"),
    );
    localStorage.removeItem("univents:verify-return");
    await navigate({ to: destination, replace: true });
  }

  async function resendVerification() {
    if (!isValidEmail) {
      toast.error("Digite um endereço de e-mail válido.");
      return;
    }

    setResending(true);
    try {
      const response = await auth.resendVerifyEmail(email.trim());
      if (!response.success) {
        toast.error(response.message || "Não foi possível enviar o novo link.");
        return;
      }
      setMessage("Enviamos um novo link de verificação para seu e-mail.");
      toast.success("Novo link de verificação enviado.");
    } catch {
      toast.error("Não foi possível enviar o novo link.");
    } finally {
      setResending(false);
    }
  }

  return (
    <AuthActionPage
      title="Verificação de e-mail"
      description="Confirme seu endereço de e-mail para acessar sua conta."
    >
      {status === "loading" ? (
        <p className="text-center text-sm text-muted-foreground">
          Confirmando seu e-mail…
        </p>
      ) : status === "success" ? (
        <div className="space-y-4 text-center">
          <p className="text-sm text-emerald-600">
            E-mail verificado com sucesso.
          </p>
          <Button
            className="h-11 w-full"
            onClick={() =>
              navigate({
                to: "/auth",
                search: {
                  redirect: safeInternalReturnTo(
                    returnTo,
                    localStorage.getItem("univents:verify-return"),
                  ),
                },
                replace: true,
              })
            }
          >
            Ir para o login
          </Button>
        </div>
      ) : status === "required" ? (
        <div className="space-y-5 text-center">
          <div className="rounded-2xl border border-amber-500/20 bg-amber-500/5 p-5">
            <MailWarning className="mx-auto mb-3 size-8 text-amber-600" />
            <p className="font-semibold">Verifique seu e-mail para continuar</p>
            <p className="mt-2 text-sm text-muted-foreground">
              Enviamos um link para {email || "o e-mail da sua conta"}.
            </p>
          </div>
          <Button
            className="h-11 w-full"
            disabled={resending || !isValidEmail}
            onClick={() => void resendVerification()}
          >
            <RefreshCw className={resending ? "animate-spin" : undefined} />
            {resending ? "Enviando…" : "Reenviar e-mail"}
          </Button>
          <Button
            variant="outline"
            className="h-11 w-full"
            onClick={() => void continueAfterVerification()}
          >
            Já verifiquei
          </Button>
        </div>
      ) : (
        <div className="space-y-5">
          <div className="rounded-2xl border border-destructive/20 bg-destructive/5 p-4 text-center">
            <AlertCircle className="mx-auto mb-2 size-7 text-destructive" />
            <p className="text-sm font-semibold text-foreground">
              Link de verificação indisponível
            </p>
            <p
              role="alert"
              className="mt-1 text-sm leading-relaxed text-muted-foreground"
            >
              {message ||
                (token
                  ? "Link de verificação inválido ou incompleto."
                  : "Informe seu e-mail para receber um novo link de verificação.")}
            </p>
          </div>
          <div className="space-y-3">
            <div>
              <label
                htmlFor="verification-email"
                className="text-sm font-medium"
              >
                E-mail da sua conta
              </label>
              <p className="mt-1 text-xs text-muted-foreground">
                Enviaremos um novo link para este endereço.
              </p>
            </div>
            <Input
              id="verification-email"
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              placeholder="voce@exemplo.com"
              aria-invalid={email.length > 0 && !isValidEmail}
            />
            {email.length > 0 && !isValidEmail && (
              <p className="text-xs text-destructive">
                Informe um endereço de e-mail válido.
              </p>
            )}
          </div>
          <Button
            className="h-11 w-full"
            disabled={resending || !isValidEmail}
            onClick={() => void resendVerification()}
          >
            <RefreshCw className={resending ? "animate-spin" : undefined} />
            {resending ? "Enviando…" : "Enviar novo link"}
          </Button>
        </div>
      )}
    </AuthActionPage>
  );
}
