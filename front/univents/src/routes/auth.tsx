import { createFileRoute } from "@tanstack/react-router";
import { ModernAuth, useAuth } from "@trieoh/identityx-sdk-ts/react";
import { useEffect } from "react";
import { toast } from "sonner";
import z from "zod";
import { useSessionActions } from "@/features/auths/hooks/use-session-actions";
import {
  readAuthReturnTo,
  storeAuthReturnTo,
} from "@/features/auths/lib/auth-path";
import { requireGuest } from "@/features/auths/lib/route-guard";
import { Logo } from "@/shared/ui/logo";

const authSearchSchema = z.object({
  redirect: z.string().optional().catch(""),
});

export const Route = createFileRoute("/auth")({
  validateSearch: (search) => authSearchSchema.parse(search),
  beforeLoad: requireGuest,
  component: AuthPage,
});

function AuthPage() {
  const search = Route.useSearch();
  const { auth: sessionAuth } = useAuth();
  const { completeLogin } = useSessionActions();

  useEffect(() => {
    if (search.redirect) {
      storeAuthReturnTo(localStorage, search.redirect);
    }
  }, [search.redirect]);

  const handleLoginSuccess = async (message?: string) => {
    const destination =
      search.redirect || readAuthReturnTo(localStorage) || "/profile";
    const completed = await completeLogin(destination, message);
    if (!completed) return;
    if (!sessionAuth.profile()?.verified_at) {
      toast.warning("Seu e-mail ainda não foi verificado", {
        description: "Verifique sua conta para liberar todos os recursos.",
      });
    }
  };

  const handleSignUpSuccess = async (message?: string) => {
    toast.success(message ?? "Conta criada com sucesso");
  };

  const handleFailure = async (message: string, trace?: string[]) => {
    const description = trace?.join("\n").replaceAll("trace: ", "");
    toast.error(message || "Não foi possível concluir a autenticação", {
      description,
    });
  };

  return (
    <div className="relative [&>main]:py-16 [&>main]:pt-52 md:[&>main]:pt-56">
      <div className="absolute left-1/2 top-16 md:top-20 z-20 w-32 md:w-40 -translate-x-1/2">
        <Logo
          variant="complete"
          priority
          imgClassName="h-auto max-h-16 md:max-h-20"
        />
      </div>

      <ModernAuth
        initialView="signin"
        onLoginSuccess={handleLoginSuccess}
        onSignUpSuccess={handleSignUpSuccess}
        onFailed={handleFailure}
      />
    </div>
  );
}
