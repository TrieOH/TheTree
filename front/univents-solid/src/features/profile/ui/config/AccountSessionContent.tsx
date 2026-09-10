import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import {
  Show,
  createEffect,
  createSignal,
} from "solid-js";
import type { JSX } from "@solidjs/web";

import CheckCircle2Icon from "~icons/lucide/circle-check";
import FingerprintIcon from "~icons/lucide/fingerprint";
import MailIcon from "~icons/lucide/mail";
import MailWarningIcon from "~icons/lucide/mail-warning";

import { toast } from "@/shared/ui/toast";

type IconComponent = () => JSX.Element;

const CheckCircle2 = CheckCircle2Icon as unknown as IconComponent;

const Fingerprint = FingerprintIcon as unknown as IconComponent;

const Mail = MailIcon as unknown as IconComponent;

const MailWarning = MailWarningIcon as unknown as IconComponent;

type MaybeAccessor<T> = | T | (() => T);

function readReactive<T>(
  value: MaybeAccessor<T>,
): T {
  return typeof value === "function"
    ? (value as () => T)()
    : value;
}

export function AccountSessionContent() {
  const {
    auth,
    isInitializing,
    isAuthenticated,
  } = useAuth();

  const [resending, setResending] =
    createSignal(false);

  const profile = () => auth.profile();

  const initializing = () => readReactive(isInitializing);

  const authenticated = () => readReactive(isAuthenticated);
  let refreshedEmail: string | undefined;

  createEffect(
    () => ({
      initializing: initializing(),
      authenticated: authenticated(),
      email: profile()?.email,
      verifiedAt: profile()?.verified_at,
    }),

    ({
      initializing,
      authenticated,
      email,
      verifiedAt,
    }) => {
      if (
        initializing ||
        !authenticated ||
        !email ||
        verifiedAt
      ) {
        return;
      }

      if (refreshedEmail === email) return;

      refreshedEmail = email;

      void auth.refresh().catch((error) => {
        console.error("Erro ao atualizar sessão:", error);
      });
    },
  );

  const resendVerification = async (email: string) => {
    if (resending()) return;

    setResending(true);

    try {
      const response = await auth.resendVerifyEmail(email);

      if (!response.success) {
        throw new Error(
          response.message ??
          "Não foi possível enviar o e-mail",
        );
      }

      toast.success("E-mail de verificação enviado");
    } catch (error) {
      toast.error(
        error instanceof Error
          ? error.message
          : "Não foi possível enviar o e-mail",
      );
    } finally {
      setResending(false);
    }
  };

  return (
    <Show
      when={!initializing()}
      fallback={<AccountSkeleton />}
    >
      <Show
        when={authenticated() && profile()}
        fallback={
          <div class="py-3 text-sm text-muted-foreground">
            Nenhuma sessão ativa encontrada.
          </div>
        }
      >
        {(currentProfile) => {
          const email = () => currentProfile().email;

          const id = () => currentProfile().id;

          const verifiedAt = () => currentProfile().verified_at;

          return (
            <div class="space-y-2">

              <div class="rounded-lg border border-border/50 bg-card p-3 shadow-sm">
                <div class="flex min-w-0 flex-nowrap items-center gap-3 max-[424px]:flex-wrap">
                  <div class="flex size-9 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary [&>svg]:size-4">
                    <Fingerprint />
                  </div>

                  <div class="min-w-0 flex-1">
                    <p class="m-0! text-sm font-medium">
                      Sua conta
                    </p>

                    <p class="mt-0.5 text-xs text-muted-foreground">
                      Perfil autenticado
                    </p>
                  </div>

                  <Show
                    when={verifiedAt()}
                    fallback={
                      <span class="ml-auto inline-flex shrink-0 items-center gap-1.5 whitespace-nowrap text-xs font-medium text-amber-600 max-[424px]:ml-0 max-[424px]:basis-full [&>svg]:size-4">
                        <MailWarning />
                        Não verificado
                      </span>
                    }
                  >
                    <span class="ml-auto inline-flex shrink-0 items-center gap-1.5 whitespace-nowrap text-xs font-medium text-emerald-600 max-[424px]:ml-0 max-[424px]:basis-full [&>svg]:size-4">
                      <CheckCircle2 />
                      Verificado
                    </span>
                  </Show>
                </div>

                <div class="mt-3 divide-y divide-border/40 border-t border-border/40">

                  <div class="flex min-w-0 flex-wrap items-center gap-3 py-3">
                    <span class="flex size-8 shrink-0 items-center justify-center rounded-md bg-muted text-muted-foreground [&>svg]:size-4">
                      <Mail />
                    </span>

                    <div class="min-w-0 flex-1">
                      <p class="m-0! text-sm font-medium">
                        E-mail
                      </p>

                      <p class="mt-0.5 truncate text-xs text-muted-foreground">
                        {email() || "Não disponível"}
                      </p>
                    </div>

                    <Show
                      when={!verifiedAt() && email()}
                    >
                      {(currentEmail) => (
                        <button
                          type="button"
                          disabled={resending()}
                          class="h-9 w-full shrink-0 rounded-md border border-border bg-background px-3 text-sm font-medium transition-colors hover:bg-muted disabled:cursor-not-allowed disabled:opacity-50 sm:w-auto"
                          onClick={() => resendVerification(currentEmail())}
                        >
                          {resending() ? "Enviando…" : "Reenviar verificação"}
                        </button>
                      )}
                    </Show>
                  </div>

                  <div class="flex min-w-0 items-center gap-3 py-3">
                    <span class="flex size-8 shrink-0 items-center justify-center rounded-md bg-muted text-muted-foreground [&>svg]:size-4">
                      <Fingerprint />
                    </span>

                    <div class="min-w-0 flex-1">
                      <p class="m-0! text-sm font-medium">ID da conta</p>
                      <p
                        class="mt-0.5 truncate font-mono text-xs text-muted-foreground"
                        title={id()}
                      >
                        {id()}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          );
        }}
      </Show>
    </Show>
  );
}

function AccountSkeleton() {
  return (
    <div class="space-y-2">
      <div class="rounded-lg border border-border/50 bg-card p-3">
        <div class="flex items-center gap-3">
          <div class="size-9 animate-pulse rounded-md bg-muted" />

          <div class="flex-1 space-y-2">
            <div class="h-4 w-24 animate-pulse rounded bg-muted" />
            <div class="h-3 w-32 animate-pulse rounded bg-muted" />
          </div>

          <div class="h-4 w-20 animate-pulse rounded bg-muted" />
        </div>

        <div class="mt-3 space-y-3 border-t border-border/40 pt-3">
          <SkeletonRow />

          <SkeletonRow />
        </div>
      </div>
    </div>
  );
}

function SkeletonRow() {
  return (
    <div class="flex items-center gap-3">
      <div class="size-8 shrink-0 animate-pulse rounded-md bg-muted" />

      <div class="flex-1 space-y-2">
        <div class="h-3.5 w-20 animate-pulse rounded bg-muted" />
        <div class="h-3 w-2/3 animate-pulse rounded bg-muted" />
      </div>
    </div>
  );
}