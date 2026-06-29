import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { useEffect, useRef } from "react";
import { toast } from "sonner";
import z from "zod";

const callbackSearchSchema = z.object({
  code: z.string().optional(),
});

const providerSchema = z.enum(["google", "github"]);

export const Route = createFileRoute("/auth/$provider/callback")({
  component: CallbackPage,
  validateSearch: callbackSearchSchema,
  params: {
    parse: (params) => ({
      provider: providerSchema.parse(params.provider),
    }),
  },
});

function CallbackPage() {
  const navigate = Route.useNavigate();
  const { provider } = Route.useParams();
  const { code } = Route.useSearch();
  const router = useRouter()
  const { auth } = useAuth();
  const calledRef = useRef(false);

  useEffect(() => {
    if (calledRef.current) return;
    calledRef.current = true;

    async function authenticate() {
      if (!code) {
        navigate({ to: "/auth" });
        return;
      }
      try {
        const res = await auth.completeProviderLogin(provider, code)
        if (res.success) {
          await navigate({ to: '/admin' })
          toast.success(res.message ?? "Login successful!")
          router.options.context.queryClient.invalidateQueries();
          return;
        }
        toast.error(res.message ?? "Auth Initialization Failed")
        navigate({ to: "/auth" });
      } catch {
        navigate({ to: "/auth" });
      }
    }

    authenticate();
  }, [provider, code, navigate]);

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <div className="mb-4 h-12 w-12 animate-spin rounded-full border-4 border-secondary border-t-primary mx-auto" />
        <h1 className="text-lg font-semibold text-primary">Entrando...</h1>
        <p className="text-sm text-muted-foreground">
          Aguarde enquanto concluímos sua autenticação.
        </p>
      </div>
    </div>
  );
}