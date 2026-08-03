import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { useEffect, useRef } from "react";
import { toast } from "sonner";
import { z } from "zod";

const providerSchema = z.enum(["google", "github"]);
const callbackSearchSchema = z.object({
  code: z.string().optional(),
  state: z.string().optional(),
});

export const Route = createFileRoute("/auth/$provider/callback")({
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
  const router = useRouter();
  const { provider } = Route.useParams();
  const { code, state } = Route.useSearch();
  const { auth } = useAuth();
  const called = useRef(false);

  useEffect(() => {
    if (called.current) return;
    called.current = true;

    const completeLogin = async () => {
      if (!code || !state) {
        toast.error("Callback OAuth inválido");
        await navigate({ to: "/auth", replace: true });
        return;
      }

      try {
        const response = await auth.completeProviderLogin(provider, code, state);
        if (!response.success) {
          throw new Error(response.message || "Não foi possível concluir o login");
        }

        const currentAuth = router.options.context.auth;
        if (currentAuth) {
          router.update({
            context: {
              ...router.options.context,
              auth: { ...currentAuth, isAuthenticated: true },
            },
          });
        }
        router.options.context.queryClient.invalidateQueries();
        toast.success("Login realizado com sucesso");
        await navigate({ to: "/profile", replace: true });
      } catch (error) {
        toast.error(
          error instanceof Error ? error.message : "Falha na autenticação OAuth",
        );
        await navigate({ to: "/auth", replace: true });
      }
    };

    void completeLogin();
  }, [auth, code, navigate, provider, router, state]);

  return (
    <main className="flex min-h-dvh items-center justify-center bg-background">
      <div className="text-center">
        <div className="mx-auto mb-4 size-12 animate-spin rounded-full border-4 border-secondary border-t-primary" />
        <h1 className="text-lg font-semibold">Entrando…</h1>
        <p className="text-sm text-muted-foreground">
          Aguarde enquanto concluímos sua autenticação.
        </p>
      </div>
    </main>
  );
}
