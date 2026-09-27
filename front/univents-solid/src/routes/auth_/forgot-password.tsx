import { createFileRoute } from "@tanstack/solid-router";
import { ModernForgotPassword } from "@trieoh/identityx-sdk-ts-solid";
import { AuthActionPage } from "@/features/auths/ui/AuthActionPage";

export const Route = createFileRoute("/auth_/forgot-password")({
  component: ForgotPasswordPage,
});

function ForgotPasswordPage() {
  const navigate = Route.useNavigate();
  return (
    <AuthActionPage
      title="Recuperar senha"
      description="Enviaremos um link seguro para seu e-mail."
    >
      <ModernForgotPassword
        signInRedirect={() =>
          void navigate({ to: "/auth", search: { redirect: "" } })
        }
      />
    </AuthActionPage>
  );
}
