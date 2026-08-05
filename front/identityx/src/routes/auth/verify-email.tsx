import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { ModernVerifyEmail } from "@trieoh/identityx-sdk-ts/react";
import { z } from "zod";
import { AuthActionPage } from "@/features/auth/ui/auth-action-page";

export const Route = createFileRoute("/auth/verify-email")({
  validateSearch: (search) =>
    z.object({ token: z.string().catch("") }).parse(search),
  component: VerifyEmailPage,
});

function VerifyEmailPage() {
  const { token } = Route.useSearch();
  const navigate = useNavigate();

  return (
    <AuthActionPage>
      {token ? (
        <ModernVerifyEmail
          token={token}
          signInRedirect={() =>
            navigate({
              to: "/auth",
              search: { redirect: "" },
              replace: true,
            })
          }
        />
      ) : (
        <p role="alert" className="text-center text-sm text-destructive">
          Link de verificação inválido ou incompleto.
        </p>
      )}
    </AuthActionPage>
  );
}
