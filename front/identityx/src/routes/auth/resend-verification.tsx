import { createFileRoute } from "@tanstack/react-router";
import { ModernResendVerification } from "@trieoh/identityx-sdk-ts/react";
import { AuthActionPage } from "@/features/auth/ui/auth-action-page";

export const Route = createFileRoute("/auth/resend-verification")({
  component: ResendVerificationPage,
});

function ResendVerificationPage() {
  return (
    <AuthActionPage
      title="Reenviar verificação"
      description="Informe seu e-mail para receber um novo link de verificação."
    >
      <ModernResendVerification />
    </AuthActionPage>
  );
}
