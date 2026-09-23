import { createFileRoute } from "@tanstack/solid-router";
import type { JSX } from "@solidjs/web";
import { ModernResetPassword } from "@trieoh/identityx-sdk-ts-solid";
import AlertCircleIcon from "~icons/lucide/alert-circle";
import { Button } from "@trieoh/ui-solid";
import { AuthActionPage } from "@/features/auths/ui/AuthActionPage";
import { toast } from "@/shared/ui/toast";
import z from "zod";

const AlertCircle = AlertCircleIcon as unknown as (props: {
  class?: string;
}) => JSX.Element;

export const Route = createFileRoute("/auth_/reset")({
  validateSearch: z.object({ token: z.string().catch("") }),
  component: ResetPasswordPage,
});

function ResetPasswordPage() {
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const goToLogin = () =>
    void navigate({ to: "/auth", search: { redirect: "" }, replace: true });
  const handleSuccess = async (message?: string) => {
    toast.success(message || "Senha redefinida com sucesso!");
    goToLogin();
  };

  return (
    <AuthActionPage
      title={search().token ? "Redefinir Senha" : undefined}
      description={search().token ? "Crie uma nova senha para sua conta." : undefined}
    >
      {search().token ? (
        <ModernResetPassword
          token={search().token}
          onSuccess={handleSuccess}
          signInRedirect={goToLogin}
        />
      ) : (
        <div class="space-y-5">
          <div class="rounded-2xl border border-destructive/20 bg-destructive/5 p-5 text-center">
            <AlertCircle class="mx-auto mb-3 size-8 text-destructive" />
            <p class="text-sm font-semibold text-foreground">
              Link de redefinição indisponível
            </p>
            <p role="alert" class="mt-1 text-sm leading-relaxed text-muted-foreground">
              O link é inválido ou está incompleto. Solicite um novo link para
              redefinir sua senha.
            </p>
          </div>
          <Button
            class="h-11 w-full"
            onClick={goToLogin}
          >
            Voltar para o login
          </Button>
        </div>
      )}
    </AuthActionPage>
  );
}
