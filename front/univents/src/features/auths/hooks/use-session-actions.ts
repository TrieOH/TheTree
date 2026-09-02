import { useQueryClient } from "@tanstack/react-query";
import { useNavigate, useRouter } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { useCallback } from "react";
import { toast } from "sonner";
import { clearAuthReturnTo } from "../lib/auth-path";

export function useSessionActions() {
  const { auth } = useAuth();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const router = useRouter();

  const updateAuthState = useCallback(
    (isAuthenticated: boolean) => {
      const current = router.options.context.auth;
      if (!current) return false;
      router.update({
        context: {
          ...router.options.context,
          auth: { ...current, isAuthenticated },
        },
      });
      return true;
    },
    [router],
  );

  const completeLogin = useCallback(
    async (destination: string, message?: string) => {
      if (!updateAuthState(true)) {
        toast.error("Não foi possível inicializar a autenticação");
        return false;
      }
      clearAuthReturnTo(localStorage);
      queryClient.clear();
      await navigate({ to: destination, replace: true });
      toast.success(message ?? "Login realizado com sucesso");
      return true;
    },
    [navigate, queryClient, updateAuthState],
  );

  const logoutTo = useCallback(
    async (destination: string) => {
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
    },
    [auth, navigate, queryClient, updateAuthState],
  );

  const logout = useCallback(() => logoutTo("/"), [logoutTo]);

  return { completeLogin, logout, logoutTo };
}
