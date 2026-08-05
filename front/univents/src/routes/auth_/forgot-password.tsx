import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { ModernForgotPassword } from "@trieoh/identityx-sdk-ts/react";
import { AuthActionPage } from "@/features/auths/ui/auth-action-page";

export const Route = createFileRoute("/auth_/forgot-password")({
  component: ForgotPasswordPage,
});

function ForgotPasswordPage() {
  const navigate = useNavigate();
  return (
    <AuthActionPage
      title="Recuperar senha"
      description="Enviaremos um link seguro para o seu e-mail."
    >
      <ModernForgotPassword
        signInRedirect={() =>
          navigate({ to: "/auth", search: { redirect: "" } })
        }
      />
    </AuthActionPage>
  );
}
