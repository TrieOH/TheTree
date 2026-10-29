import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { createEffect, createSignal } from "solid-js";
import Loader2Icon from "~icons/lucide/loader-2";

import {
  clearAuthReturnTo,
  readAuthReturnTo,
} from "@/features/auths/lib/auth-path";
import { toast } from "@/shared/ui/toast";
import { Schema, createSchema } from "@/shared/lib/effect-schema";

const Loader2 = Loader2Icon as unknown as (props: {
  class?: string;
}) => JSX.Element;

const providerSchema = Schema.Literals(["google", "github"]);
const decodeProvider = Schema.decodeUnknownSync(providerSchema);

const searchSchema = createSchema(
  Schema.Struct({
    code: Schema.optional(Schema.String),
    state: Schema.optional(Schema.String),
  }),
);

export const Route = createFileRoute("/auth_/$provider/callback")({
  params: {
    parse: (params) => ({
      provider: decodeProvider(params.provider),
    }),
  },
  validateSearch: searchSchema,
  component: OAuthCallbackPage,
});

function OAuthCallbackPage() {
  const params = Route.useParams();
  const search = Route.useSearch();
  const { auth } = useAuth();
  const [error, setError] = createSignal<string | null>(null);
  let called = false;

  createEffect(
    () => ({ provider: params().provider, search: search() }),
    ({ provider, search: currentSearch }) => {
      if (called) return;
      called = true;

      const code = currentSearch.code;
      const state = currentSearch.state;

      if (!code || !state) {
        setError("Código de autorização ou estado ausente.");
        return;
      }

      void auth
        .completeProviderLogin(provider, code, state)
        .then((response) => {
          if (!response.success) {
            throw new Error(
              response.message || "Não foi possível concluir o login",
            );
          }
          toast.success("Login realizado com sucesso");
          const destination = readAuthReturnTo(localStorage) || "/profile";
          clearAuthReturnTo(localStorage);
          window.location.assign(destination);
        })
        .catch((err: unknown) => {
          const message =
            err instanceof Error ? err.message : "Falha na autenticação OAuth.";
          setError(message);
          toast.error(message);
        });
    },
  );

  return (
    <div class="min-h-screen flex items-center justify-center p-4 bg-background">
      <div class="max-w-md w-full rounded-2xl border border-border bg-card p-8 text-center space-y-4 shadow-xl">
        {error() ? (
          <div class="space-y-4">
            <h2 class="text-lg font-semibold text-destructive">
              Erro na Autenticação
            </h2>
            <p class="text-sm text-muted-foreground">{error()}</p>
            <a
              href="/auth"
              class="inline-flex items-center justify-center px-4 py-2 rounded-lg bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors"
            >
              Voltar para o login
            </a>
          </div>
        ) : (
          <div class="space-y-4">
            <Loader2 class="size-8 animate-spin mx-auto text-primary" />
            <p class="text-sm text-muted-foreground">
              Finalizando login com {params().provider}...
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
