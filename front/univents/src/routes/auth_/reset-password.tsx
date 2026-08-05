import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { ModernResetPassword } from "@trieoh/identityx-sdk-ts/react";
import { toast } from "sonner";
import { z } from "zod";
import { AuthActionPage } from "@/features/auths/ui/auth-action-page";

export const Route = createFileRoute("/auth_/reset-password")({
  validateSearch: (search) =>
    z.object({ token: z.string().catch("") }).parse(search),
  component: ResetPasswordPage,
});

function ResetPasswordPage() {
  const { token } = Route.useSearch();
  const navigate = useNavigate();
  const goToLogin = () =>
    navigate({ to: "/auth", search: { redirect: "" }, replace: true });

  return (
    <AuthActionPage>
      {token ? (
        <ModernResetPassword
          token={token}
          onSuccess={async (message) => {
            toast.success(message || "Senha redefinida com sucesso!");
            await goToLogin();
          }}
          signInRedirect={goToLogin}
        />
      ) : (
        <p role="alert" className="text-center text-sm text-destructive">
          Link de redefinição inválido ou incompleto.
        </p>
      )}
    </AuthActionPage>
  );
}
