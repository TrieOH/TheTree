import { createFileRoute } from "@tanstack/react-router";
import { XCircle } from "lucide-react";
import { useEffect, useState } from "react";
import { revokeSignatureFn } from "@/features/signatures/api";
import { SignatureRequestStatus } from "@/routes/signature-requests/ui/SignatureRequestStatus";
import { cn } from "@/shared/lib/utils";

export const Route = createFileRoute("/signatures/revoke")({
  validateSearch: (search: Record<string, unknown>) => ({
    token: typeof search.token === "string" ? search.token : "",
  }),
  component: RevokeSignaturePage,
});

function RevokeSignaturePage() {
  const { token } = Route.useSearch();
  const [state, setState] = useState<"loading" | "success" | "error">(
    token ? "loading" : "error",
  );

  useEffect(() => {
    if (!token) return;
    let active = true;
    void revokeSignatureFn(token)
      .then((response) => {
        if (active) setState(response.success ? "success" : "error");
      })
      .catch(() => {
        if (active) setState("error");
      });
    return () => {
      active = false;
    };
  }, [token]);

  if (state === "loading") {
    return (
      <SignatureRequestStatus
        title="Revogando assinatura"
        message="Aguarde enquanto processamos sua solicitação."
        loading
      />
    );
  }

  if (state === "error") {
    return (
      <SignatureRequestStatus
        title="Não foi possível revogar"
        message="O link é inválido, expirou ou a assinatura já foi revogada."
      />
    );
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-muted px-3 pb-24 sm:px-4">
      <section
        className={cn(
          "flex w-full max-w-sm flex-col items-center gap-4 rounded-2xl border border-border bg-card px-6 py-8 text-center shadow-sm",
          "sm:px-8 sm:py-10",
        )}
      >
        <span className="flex size-14 items-center justify-center rounded-2xl bg-destructive text-destructive-foreground">
          <XCircle className="size-7" />
        </span>
        <h1 className="text-xl font-semibold text-card-foreground">
          Assinatura revogada
        </h1>
        <p className="text-sm text-muted-foreground">
          A assinatura feita por convite foi revogada com sucesso.
        </p>
        <button
          type="button"
          onClick={() => window.close()}
          className="text-sm font-medium text-primary hover:underline"
        >
          Você pode fechar esta aba
        </button>
      </section>
    </main>
  );
}
