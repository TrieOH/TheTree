import { createFileRoute } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { z } from "zod";
import { AUTH_RETURN_TO_STORAGE_KEY } from "@/features/auths/lib/auth-path";

const providerSchema = z.enum(["google", "github"]);
const callbackSearchSchema = z.object({
  code: z.string().optional(),
  state: z.string().optional(),
});

export const Route = createFileRoute("/auth_/$provider/callback")({
  component: OAuthCallbackPage,
  validateSearch: callbackSearchSchema,
  params: {
    parse: (params) => ({
      provider: providerSchema.parse(params.provider),
    }),
  },
});

function OAuthCallbackPage() {
  const navigate = Route.useNavigate();
  const { provider } = Route.useParams();
  const { code, state } = Route.useSearch();
  const { auth } = useAuth();
  const called = useRef(false);
  const [result, setResult] = useState<string>();

  useEffect(() => {
    if (called.current) return;
    called.current = true;

    const completeLogin = async () => {
      if (!code || !state) {
        setResult("Callback OAuth inválido: code ou state ausente.");
        return;
      }

      try {
        const response = await auth.completeProviderLogin(
          provider,
          code,
          state,
        );
        if (!response.success) {
          throw new Error(
            response.message || "Não foi possível concluir o login",
          );
        }

        toast.success("Login realizado com sucesso");
        const destination =
          localStorage.getItem(AUTH_RETURN_TO_STORAGE_KEY) || "/profile";
        localStorage.removeItem(AUTH_RETURN_TO_STORAGE_KEY);
        window.location.assign(destination);
      } catch (error) {
        const message =
          error instanceof Error
            ? error.message
            : "Falha na autenticação OAuth";
        setResult(message);
        toast.error(message);
      }
    };

    void completeLogin();
  }, [auth, code, navigate, provider, state]);

  return (
    <main className="flex min-h-dvh items-center justify-center bg-background">
      <div className="max-w-md space-y-4 text-center">
        {!result && (
          <div className="mx-auto size-12 animate-spin rounded-full border-4 border-secondary border-t-primary" />
        )}
        <h1 className="text-lg font-semibold">
          {result ? "Resultado da autenticação" : "Entrando…"}
        </h1>
        <p className="text-sm text-muted-foreground">
          {result ?? "Aguarde enquanto concluímos sua autenticação."}
        </p>
        {result && (
          <button
            type="button"
            onClick={() =>
              navigate({
                to: "/profile",
                search: { tab: "about" },
                replace: true,
              })
            }
            className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"
          >
            Continuar
          </button>
        )}
      </div>
    </main>
  );
}
