import type { JSX } from "@solidjs/web";
import { Link } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { Button, buttonVariants, cn } from "@trieoh/ui-solid";
import { Show, createMemo, createSignal } from "solid-js";

import CheckCircle2Icon from "~icons/lucide/check-circle-2";
import FileWarningIcon from "~icons/lucide/file-warning";
import Loader2Icon from "~icons/lucide/loader-2";
import Trash2Icon from "~icons/lucide/trash-2";
import XCircleIcon from "~icons/lucide/x-circle";

import { toast } from "@/shared/ui/toast";
import { signatureQueryOptions } from "../api";
import { useRevokeSignatureMutation } from "../api/mutations";
import {
  parseJwtPayload,
  signatureRevocationTokenClaimsSchema,
  type SignatureRevocationTokenClaims,
} from "../model";

const CheckCircle2 = CheckCircle2Icon as unknown as (props: { class?: string }) => JSX.Element;
const FileWarning = FileWarningIcon as unknown as (props: { class?: string }) => JSX.Element;
const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Trash2 = Trash2Icon as unknown as (props: { class?: string }) => JSX.Element;
const XCircle = XCircleIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface RevokeSignatureViewProps {
  token: string;
}

export function RevokeSignatureView(
  props: RevokeSignatureViewProps,
): JSX.Element {
  const claims = createMemo(() =>
    parseJwtPayload<SignatureRevocationTokenClaims>(
      props.token,
      signatureRevocationTokenClaimsSchema,
    ),
  );

  const signatureId = () => claims()?.signature_id ?? "";

  const signatureQuery = useQuery(() => ({
    ...signatureQueryOptions(signatureId()),
    enabled: Boolean(signatureId()),
  }));

  const signature = () => signatureQuery().data;

  const [confirmed, setConfirmed] = createSignal(false);
  const [isRevoking, setIsRevoking] = createSignal(false);
  const [isRevokedSuccess, setIsRevokedSuccess] = createSignal(false);

  const revokeMutation = useRevokeSignatureMutation();

  const handleRevoke = async () => {
    if (!props.token) {
      toast.error("Token de autorização inválido ou ausente.");
      return;
    }

    if (!confirmed()) {
      toast.error("Você precisa confirmar que está ciente da revogação.");
      return;
    }

    setIsRevoking(true);
    try {
      await revokeMutation.mutateAsync({
        token: props.token,
        signatureId: signature()?.id ?? signatureId(),
        editionId: claims()?.edition_id ?? signature()?.edition_id,
      });

      setIsRevokedSuccess(true);
      toast.success("Assinatura revogada com sucesso.");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Erro ao revogar assinatura.";
      toast.error(msg);
    } finally {
      setIsRevoking(false);
    }
  };

  const formattedDate = () => {
    const d = signature()?.created_at;
    if (!d) return "";
    try {
      return new Date(d).toLocaleDateString("pt-BR", {
        day: "2-digit",
        month: "long",
        year: "numeric",
      });
    } catch {
      return "";
    }
  };

  return (
    <div class="min-h-screen bg-background text-foreground flex flex-col justify-start pt-8 sm:pt-12 pb-32 sm:pb-36 px-4 sm:px-6">
      <div class="mx-auto w-full max-w-2xl">
        {/* 1. Missing or Invalid Token State */}
        <Show when={!props.token || !claims()}>
          <div class="py-16 text-center max-w-md mx-auto space-y-4">
            <FileWarning class="size-10 text-muted-foreground mx-auto" />
            <div class="space-y-1.5">
              <h1 class="text-xl font-semibold tracking-tight text-foreground">
                Link Inválido ou Incompleto
              </h1>
              <p class="text-xs text-muted-foreground leading-relaxed">
                Não foi possível validar o token de revogação. Certifique-se de acessar o link completo enviado para o seu e-mail.
              </p>
            </div>
            <div class="pt-2">
              <Link
                to="/"
                search={{ as: "guest" }}
                class={cn(buttonVariants({ variant: "outline" }), "text-xs cursor-pointer")}
              >
                Voltar à Página Inicial
              </Link>
            </div>
          </div>
        </Show>

        {/* 2. Loading State */}
        <Show when={props.token && claims() && signatureQuery().isLoading}>
          <div class="py-24 text-center max-w-sm mx-auto space-y-3">
            <Loader2 class="mx-auto size-7 animate-spin text-muted-foreground" />
            <p class="text-xs text-muted-foreground">
              Localizando assinatura...
            </p>
          </div>
        </Show>

        {/* 3. Success State */}
        <Show when={isRevokedSuccess()}>
          <div class="py-16 max-w-md mx-auto text-center space-y-5">
            <CheckCircle2 class="size-12 text-foreground mx-auto" />
            <div class="space-y-2">
              <h1 class="text-2xl font-semibold tracking-tight text-foreground">
                Assinatura Revogada
              </h1>
              <p class="text-xs text-muted-foreground leading-relaxed">
                Esta assinatura foi desativada e não poderá mais ser vinculada a novos certificados deste evento.
              </p>
            </div>
            <div class="pt-2">
              <Link
                to="/"
                search={{ as: "guest" }}
                class={cn(buttonVariants({ variant: "outline" }), "text-xs cursor-pointer")}
              >
                Fechar
              </Link>
            </div>
          </div>
        </Show>

        {/* 4. Not Found / Already Deleted State */}
        <Show
          when={
            props.token &&
            claims() &&
            !signatureQuery().isLoading &&
            !isRevokedSuccess() &&
            (!signature() || signature()?.deleted_at)
          }
        >
          <div class="py-16 max-w-md mx-auto text-center space-y-4">
            <XCircle class="size-10 text-muted-foreground mx-auto" />
            <div class="space-y-1.5">
              <h1 class="text-xl font-semibold tracking-tight text-foreground">
                Assinatura Inexistente ou Já Revogada
              </h1>
              <p class="text-xs text-muted-foreground leading-relaxed">
                Esta assinatura já foi revogada anteriormente ou não foi localizada.
              </p>
            </div>
            <div class="pt-2">
              <Link
                to="/"
                search={{ as: "guest" }}
                class={cn(buttonVariants({ variant: "outline" }), "text-xs cursor-pointer")}
              >
                Voltar à Página Inicial
              </Link>
            </div>
          </div>
        </Show>

        {/* 5. Active Revocation Layout */}
        <Show
          when={
            props.token &&
            claims() &&
            !signatureQuery().isLoading &&
            !isRevokedSuccess() &&
            signature() &&
            !signature()?.deleted_at
          }
        >
          <div class="space-y-6">
            {/* Header with clean contextual metadata */}
            <div class="space-y-1.5 border-b border-border pb-5">
              <h1 class="text-xl sm:text-2xl font-semibold tracking-tight text-foreground">
                Revogação de Assinatura Digital
              </h1>
              <div class="flex flex-wrap items-center gap-x-2 gap-y-1 text-xs text-muted-foreground">
                <span>Signatário: <strong class="font-medium text-foreground">{signature()?.signatory_name}</strong></span>
                <Show when={signature()?.signatory_title}>
                  <span>•</span>
                  <span>{signature()?.signatory_title}</span>
                </Show>
                <Show when={signature()?.signatory_email}>
                  <span>•</span>
                  <span>{signature()?.signatory_email}</span>
                </Show>
                <Show when={formattedDate()}>
                  <span>•</span>
                  <span>Registrada em {formattedDate()}</span>
                </Show>
              </div>
            </div>

            {/* Signature Preview Canvas */}
            <div class="space-y-2">
              <span class="text-xs font-medium text-muted-foreground">
                Assinatura vinculada
              </span>
              <div
                class="flex items-center justify-center rounded-xl border border-border p-6 min-h-48 overflow-hidden shadow-xs"
                style={{
                  "background-image":
                    "linear-gradient(45deg, #f1f5f9 25%, transparent 25%), linear-gradient(-45deg, #f1f5f9 25%, transparent 25%), linear-gradient(45deg, transparent 75%, #f1f5f9 75%), linear-gradient(-45deg, transparent 75%, #f1f5f9 75%)",
                  "background-size": "16px 16px",
                  "background-position": "0 0, 0 8px, 8px -8px, -8px 0",
                  "background-color": "#ffffff",
                }}
              >
                <img
                  src={signature()?.image_url}
                  alt={`Assinatura de ${signature()?.signatory_name}`}
                  class="max-h-36 max-w-full object-contain filter drop-shadow-xs"
                />
              </div>
            </div>

            {/* Warning Note */}
            <div class="rounded-lg border border-border/80 bg-muted/20 p-4 space-y-1.5 text-xs text-muted-foreground">
              <p class="font-medium text-foreground">
                Impacto da revogação
              </p>
              <p class="leading-relaxed">
                A assinatura será desativada de forma permanente. Novos certificados emitidos pela organização deste evento não poderão mais utilizá-la. Certificados já emitidos continuarão arquivados normalmente.
              </p>
            </div>

            {/* Confirmation Checkbox */}
            <div class="flex items-center gap-3 rounded-lg border border-border/80 bg-muted/20 p-3.5">
              <input
                id="confirm-revoke-signature"
                type="checkbox"
                checked={confirmed()}
                onInput={(e) => {
                  const target = (e.target ?? e.currentTarget) as HTMLInputElement;
                  setConfirmed(target.checked);
                }}
                onChange={(e) => {
                  const target = (e.target ?? e.currentTarget) as HTMLInputElement;
                  setConfirmed(target.checked);
                }}
                class="size-4 rounded border-border text-destructive focus:ring-destructive cursor-pointer shrink-0"
              />
              <label
                for="confirm-revoke-signature"
                class="text-xs text-muted-foreground leading-normal select-none cursor-pointer"
              >
                Confirmo que desejo revogar definitivamente esta assinatura digital.
              </label>
            </div>

            {/* Actions Bar */}
            <div class="flex flex-col-reverse sm:flex-row items-stretch sm:items-center justify-between gap-3 pt-3 border-t border-border">
              <Link
                to="/"
                search={{ as: "guest" }}
                class="text-xs text-muted-foreground hover:text-foreground transition-colors cursor-pointer text-center py-2 sm:py-0"
              >
                Cancelar
              </Link>

              <Button
                type="button"
                variant="destructive"
                onClick={handleRevoke}
                disabled={isRevoking() || !confirmed()}
                class="inline-flex items-center justify-center gap-2 cursor-pointer text-xs font-semibold px-6 py-2.5 h-10 sm:h-9 w-full sm:w-auto"
              >
                <Show when={isRevoking()} fallback={<Trash2 class="size-3.5" />}>
                  <Loader2 class="size-3.5 animate-spin" />
                </Show>
                {isRevoking() ? "Revogando..." : "Revogar Assinatura"}
              </Button>
            </div>
          </div>
        </Show>
      </div>
    </div>
  );
}
