import { createFileRoute } from "@tanstack/solid-router";
import {
  ModernResendVerification,
  ModernVerifyEmail,
  useAuth,
} from "@trieoh/identityx-sdk-ts-solid";
import { createEffect, createMemo, createSignal } from "solid-js";
import z from "zod";

import { AuthActionPage } from "@/features/auths/ui/AuthActionPage";
import {
  clearAuthReturnTo,
  readAuthReturnTo,
  safeInternalReturnTo,
  storeAuthReturnTo,
} from "@/features/auths/lib/auth-path";
import { toast } from "@/shared/ui/toast";

export const Route = createFileRoute("/auth_/verify")({
  validateSearch: z.object({
    token: z.string().catch(""),
    returnTo: z.string().optional(),
  }),
  component: VerifyEmailPage,
});

function VerifyEmailPage() {
  const search = Route.useSearch();
  const { auth } = useAuth();
  const token = createMemo(() => search().token);
  const [showHeader, setShowHeader] = createSignal(true);
  const destination = () =>
    safeInternalReturnTo(search().returnTo, readAuthReturnTo(localStorage));
  createEffect(
    () => search().returnTo,
    (returnTo) => {
      if (returnTo) storeAuthReturnTo(localStorage, returnTo);
    },
  );
  const handleSuccess = async () => {
    const refreshed = await auth.refresh();
    if (refreshed.success) {
      clearAuthReturnTo(localStorage);
      toast.success("E-mail verificado com sucesso.");
      window.location.assign(destination());
    }
  };

  return (
    <AuthActionPage
      title={showHeader() ? "Verificação de e-mail" : undefined}
      description={
        showHeader()
          ? "Confirme seu endereço de e-mail para acessar sua conta."
          : undefined
      }
    >
      {token() ? (
        <ModernVerifyEmail
          token={token()}
          onSuccess={handleSuccess}
          onFailed={() => {
            setShowHeader(false);
            return Promise.resolve();
          }}
        />
      ) : (
        <ModernResendVerification email={auth.profile()?.email ?? ""} />
      )}
    </AuthActionPage>
  );
}
