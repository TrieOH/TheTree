import { createQueryClient, TanStackQueryProvider } from "@trieoh/front-core";
import type { ReactNode } from "react";
import { toast } from "sonner";
import { isVerifiedEmailRequiredError } from "@/shared/lib/errors";

let context: { queryClient: ReturnType<typeof createQueryClient> } | undefined;

export function getContext() {
  if (context) return context;

  context = {
    queryClient: createQueryClient({
      onError: (error) => {
        if (isVerifiedEmailRequiredError(error)) {
          toast.warning("Você precisa verificar seu e-mail", {
            id: "verified-email-required",
            description: "Verifique sua conta para acessar este recurso.",
          });
        }
      },
    }),
  };

  return context;
}

export function Provider({ children }: { children: ReactNode }) {
  const { queryClient } = getContext();

  return (
    <TanStackQueryProvider queryClient={queryClient}>
      {children}
    </TanStackQueryProvider>
  );
}
