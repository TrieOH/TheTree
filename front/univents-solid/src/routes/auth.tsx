import { createFileRoute } from "@tanstack/solid-router";
import { ModernAuth, useAuth } from "@trieoh/identityx-sdk-ts-solid";
import {
  readAuthReturnTo,
  storeAuthReturnTo,
} from "@/features/auths/lib/auth-path";
import z from "zod";
import Logo from "@/shared/ui/Logo";
import { toast } from "@/shared/ui/toast";
import { useSessionActions } from "@/features/auths/hooks/use-session-actions";
import { requireGuest } from "@/features/auths/lib/route-guard";

const authSearchSchema = z.object({
  redirect: z.string().optional().catch(""),
});

export const Route = createFileRoute("/auth")({
  head: () => ({ meta: [{ title: "Entrar - Univents" }] }),
  validateSearch: (search) => authSearchSchema.parse(search),
  beforeLoad: requireGuest,
  component: AuthPage,
});

function AuthPage() {
  const search = Route.useSearch();
  const { auth } = useAuth();
  const { completeLogin } = useSessionActions();

  const handleLoginSuccess = async (message?: string) => {
    const redirectTo = search().redirect;
    const destination =
      redirectTo || readAuthReturnTo(localStorage) || "/profile";
    if (redirectTo) storeAuthReturnTo(localStorage, redirectTo);
    if (
      await completeLogin(
        destination,
        message || "Login realizado com sucesso!",
      )
    ) {
      if (!auth.profile()?.verified_at) {
        toast({
          type: "error",
          message: "Seu e-mail ainda não foi verificado.",
        });
      }
    }
  };

  const handleSignUpSuccess = async (message?: string) => {
    toast.success(message || "Conta criada com sucesso!");
  };

  const handleFailure = async (message: string, trace?: string[]) => {
    toast.error(trace?.length ? `${message}\n${trace.join("\n")}` : message);
  };

  return (
    <div class="relative [&>main]:py-16 [&>main]:pt-52 md:[&>main]:pt-56">
      <div class="absolute left-1/2 top-16 md:top-20 z-20 w-32 md:w-40 -translate-x-1/2">
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
