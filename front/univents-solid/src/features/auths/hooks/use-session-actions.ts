import { useQueryClient, type RouterSession } from '@trieoh/front-core-solid';
import { useNavigate, useRouter } from "@tanstack/solid-router";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { toast } from '@/shared/ui/toast';
import { clearAuthReturnTo } from "../lib/auth-path";

export function useSessionActions() {
  const { auth } = useAuth();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const router = useRouter();

  const updateAuthState = (isAuthenticated: boolean) => {
    const context = router.options.context as { session?: RouterSession };
    const current = context.session;
    if (!current) return false;

    router.update({
      context: {
        ...context,
        session: { ...current, isAuthenticated },
      },
    });
    return true;
  };

  const completeLogin = async (destination: string, message?: string) => {
    if (!updateAuthState(true)) {
      toast.error("Não foi possível inicializar a autenticação");
      return false;
    }

    clearAuthReturnTo(localStorage);
    queryClient.clear();
    await navigate({ to: destination, replace: true });
    toast.success(message ?? "Login realizado com sucesso");

    return true;
  };

  const logoutTo = async (destination: string) => {
    try {
      const response = await auth.logout();
      if (!response.success) {
        throw new Error(response.message || "Não foi possível sair");
      }

      clearAuthReturnTo(localStorage);
      updateAuthState(false);
      queryClient.clear();
      await navigate({ to: destination, replace: true });
      toast.success("Sessão encerrada");

      return true;
    } catch (error) {
      toast.error(
        error instanceof Error && error.message
          ? error.message
          : "Não foi possível encerrar a sessão",
      );
      return false;
    }
  };

  const logout = () => logoutTo("/");

  return { completeLogin, logout, logoutTo };
}
