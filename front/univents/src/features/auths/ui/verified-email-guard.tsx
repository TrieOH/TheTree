import { useLocation, useRouter } from "@tanstack/react-router";
import { useAuthActions } from "@trieoh/front-core";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { LogOut, MailWarning, RefreshCw } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Button } from "@/shared/ui/shadcn/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/shadcn/dialog";
import { requiresVerifiedEmail } from "../lib/auth-path";

export function VerifiedEmailGuard({
  children,
}: {
  children: React.ReactNode;
}) {
  const { auth, isAuthenticated } = useAuth();
  const { handleLogout } = useAuthActions();
  const location = useLocation();
  const router = useRouter();
  const [resending, setResending] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const profile = auth.profile();
  const blocked =
    isAuthenticated &&
    !profile?.verified_at &&
    requiresVerifiedEmail(location.pathname);

  useEffect(() => {
    if (blocked) {
      localStorage.setItem("univents:verify-return", location.href);
    }
  }, [blocked, location.href]);

  async function resend() {
    if (!profile?.email) return;
    setResending(true);
    try {
      const response = await auth.resendVerifyEmail(profile.email);
      if (!response.success) throw new Error(response.message);
      toast.success("E-mail de verificação reenviado.");
    } catch (error) {
      toast.error(
        error instanceof Error && error.message
          ? error.message
          : "Não foi possível reenviar o e-mail.",
      );
    } finally {
      setResending(false);
    }
  }

  async function refresh() {
    setRefreshing(true);
    try {
      const response = await auth.refresh();
      if (!response.success || !auth.profile()?.verified_at) {
        toast.info("A verificação ainda não foi confirmada.");
        return;
      }
      localStorage.removeItem("univents:verify-return");
      toast.success("E-mail verificado. Acesso liberado.");
      await router.invalidate();
    } finally {
      setRefreshing(false);
    }
  }

  return (
    <>
      {children}
      <Dialog open={blocked} onOpenChange={() => undefined}>
        <DialogContent
          showCloseButton={false}
          overlayClassName="bg-background/35 backdrop-blur-md"
          className="sm:max-w-md"
        >
          <DialogHeader className="items-center text-center">
            <div className="mb-2 flex size-14 items-center justify-center rounded-full bg-amber-500/10 text-amber-600">
              <MailWarning className="size-7" />
            </div>
            <DialogTitle className="text-xl">
              Verifique seu e-mail para continuar
            </DialogTitle>
            <DialogDescription>
              Enviamos um link para {profile?.email ?? "o e-mail da sua conta"}.
              Depois de confirmar, volte aqui e libere o acesso.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="sm:flex-col">
            <Button
              className="w-full"
              disabled={refreshing}
              onClick={() => void refresh()}
            >
              <RefreshCw className={refreshing ? "animate-spin" : undefined} />
              {refreshing ? "Verificando…" : "Já verifiquei"}
            </Button>
            <Button
              variant="outline"
              className="w-full"
              disabled={resending || !profile?.email}
              onClick={() => void resend()}
            >
              {resending ? "Enviando…" : "Reenviar e-mail"}
            </Button>
            <Button
              variant="ghost"
              className="w-full text-muted-foreground"
              onClick={() => void handleLogout()}
            >
              <LogOut /> Sair da conta
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
