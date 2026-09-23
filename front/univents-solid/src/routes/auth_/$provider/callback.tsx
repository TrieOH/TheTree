import { createFileRoute } from "@tanstack/solid-router";
import type { JSX } from "@solidjs/web";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { createEffect, createSignal } from "solid-js";
import Loader2Icon from "~icons/lucide/loader-2";
import z from "zod";

import {
  clearAuthReturnTo,
  readAuthReturnTo,
} from "@/features/auths/lib/auth-path";
import { toast } from "@/shared/ui/toast";

const Loader2 = Loader2Icon as unknown as (props: {
  class?: string;
}) => JSX.Element;

const providerSchema = z.enum(["google", "github"]);

export const Route = createFileRoute("/auth_/$provider/callback")({
  params: {
    parse: (params) => ({
      provider: providerSchema.parse(params.provider),
    }),
  },
  validateSearch: z.object({
    code: z.string().optional(),
    state: z.string().optional(),
  }),
  component: OAuthCallbackPage,
});

function OAuthCallbackPage() {
  const params = Route.useParams();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { auth } = useAuth();
  const [result, setResult] = createSignal<string>();
  let called = false;

  createEffect(
    () => ({ search: search(), provider: params().provider }),
    ({ search: currentSearch, provider: currentProvider }) => {
      const { code, state } = currentSearch;
      if (called) return;
      called = true;
      if (!code || !state) {
        setResult("Callback OAuth inválido: code ou state ausente.");
        return;
      }
      void auth
        .completeProviderLogin(currentProvider, code, state)
        .then((response) => {
          if (!response.success)
            throw new Error(
              response.message || "Não foi possível concluir o login",
            );
          toast.success("Login realizado com sucesso");
          const destination = readAuthReturnTo(localStorage) || "/profile";
          clearAuthReturnTo(localStorage);
          window.location.assign(destination);
        })
        .catch((error) => {
          const message =
            error instanceof Error
              ? error.message
              : "Falha na autenticação OAuth";
          setResult(message);
          toast.error(message);
        });
    },
  );

  return (
    <main class="flex min-h-dvh items-center justify-center bg-background">
      <div class="max-w-md space-y-4 text-center">
        {!result() && (
          <Loader2 class="mx-auto size-12 animate-spin text-primary" />
        )}
        <h1 class="text-lg font-semibold">
          {result() ? "Resultado da autenticação" : "Entrando…"}
        </h1>
        <p class="text-sm text-muted-foreground">
          {result() || "Aguarde enquanto concluímos sua autenticação."}
        </p>
        {result() && (
          <button
            type="button"
            class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"
            onClick={() =>
              void navigate({
                to: "/profile",
                search: { tab: "about" },
                replace: true,
              })
            }
          >
            Continuar
          </button>
        )}
      </div>
    </main>
  );
}
