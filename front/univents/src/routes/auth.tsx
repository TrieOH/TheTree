import { createFileRoute, useRouter } from "@tanstack/react-router";
import { ModernAuth } from "@trieoh/identityx-sdk-ts/react";
import { toast } from "sonner";
import z from "zod";
import { requireGuest } from "@/features/auths/lib/route-guard";

const authSearchSchema = z.object({
  redirect: z.string().optional().catch(""),
});

export const Route = createFileRoute("/auth")({
  validateSearch: (search) => authSearchSchema.parse(search),
  beforeLoad: requireGuest,
  component: AuthPage,
});

function AuthPage() {
  const navigate = Route.useNavigate();
  const router = useRouter();
  const search = Route.useSearch();

  const handleLoginSuccess = async (message?: string) => {
    const auth = router.options.context.auth;

    if (!auth) {
      toast.error("Não foi possível inicializar a autenticação");
      return;
    }

    router.update({
      context: {
        ...router.options.context,
        auth: { ...auth, isAuthenticated: true },
      },
    });

    const destination = search.redirect || "/profile";
    await navigate({ to: destination, replace: true });
    toast.success(message ?? "Login realizado com sucesso");
    router.options.context.queryClient.invalidateQueries();
  };

  const handleSignUpSuccess = async (message?: string) => {
    toast.success(message ?? "Conta criada com sucesso");
  };

  const handleFailure = async (message: string, trace?: string[]) => {
    const description = trace?.join("\n").replaceAll("trace: ", "");
    toast.error(message || "Não foi possível concluir a autenticação", {
      description,
    });
  };

  return (
    <div className="[&>main]:py-16">
      <ModernAuth
        initialView="signin"
        onLoginSuccess={handleLoginSuccess}
        onSignUpSuccess={handleSignUpSuccess}
        onFailed={handleFailure}
      />
    </div>
  );
}
