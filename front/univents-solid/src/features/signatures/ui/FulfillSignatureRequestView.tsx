import type { JSX } from "@solidjs/web";
import { Link } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { Button, Dialog, Label, buttonVariants, cn } from "@trieoh/ui-solid";
import { Show, createMemo, createSignal, onCleanup } from "solid-js";

import CheckCircle2Icon from "~icons/lucide/check-circle-2";
import ClockIcon from "~icons/lucide/clock";
import EraserIcon from "~icons/lucide/eraser";
import FileWarningIcon from "~icons/lucide/file-warning";
import Loader2Icon from "~icons/lucide/loader-2";
import PenLineIcon from "~icons/lucide/pen-line";
import RefreshCwIcon from "~icons/lucide/refresh-cw";
import Trash2Icon from "~icons/lucide/trash-2";
import UploadCloudIcon from "~icons/lucide/upload-cloud";
import XCircleIcon from "~icons/lucide/x-circle";

import { uploadFile } from "@/features/storage/api/index";
import { toast } from "@/shared/ui/toast";
import { signatureRequestQueryOptions } from "../api";
import {
  useDenySignatureRequestMutation,
  useFulfillSignatureRequestMutation,
} from "../api/mutations";
import {
  parseJwtPayload,
  signatureRequestTokenClaimsSchema,
  type SignatureRequestTokenClaims,
} from "../model";
import { SignatureCanvas, type SignatureCanvasRef } from "./SignatureCanvas";

const CheckCircle2 = CheckCircle2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Clock = ClockIcon as unknown as (props: { class?: string }) => JSX.Element;
const Eraser = EraserIcon as unknown as (props: { class?: string }) => JSX.Element;
const FileWarning = FileWarningIcon as unknown as (props: { class?: string }) => JSX.Element;
const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;
const PenLine = PenLineIcon as unknown as (props: { class?: string }) => JSX.Element;
const RefreshCw = RefreshCwIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash2 = Trash2Icon as unknown as (props: { class?: string }) => JSX.Element;
const UploadCloud = UploadCloudIcon as unknown as (props: { class?: string }) => JSX.Element;
const XCircle = XCircleIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface FulfillSignatureRequestViewProps {
  token: string;
}

type InputMode = "draw" | "upload";

export function FulfillSignatureRequestView(
  props: FulfillSignatureRequestViewProps,
): JSX.Element {
  const claims = createMemo(() =>
    parseJwtPayload<SignatureRequestTokenClaims>(
      props.token,
      signatureRequestTokenClaimsSchema,
    ),
  );

  const requestId = () => claims()?.request_id ?? "";

  const requestQuery = useQuery(() => ({
    ...signatureRequestQueryOptions(requestId()),
    enabled: Boolean(requestId()),
  }));

  const request = () => requestQuery().data;

  // Local interaction states
  const [mode, setMode] = createSignal<InputMode>("draw");
  const [importedFile, setImportedFile] = createSignal<File | null>(null);
  const [previewUrl, setPreviewUrl] = createSignal<string>("");
  const [isDragging, setIsDragging] = createSignal(false);
  const [agreed, setAgreed] = createSignal(false);
  const [isSubmitting, setIsSubmitting] = createSignal(false);

  // Status transitions
  const [isFulfilledSuccess, setIsFulfilledSuccess] = createSignal(false);
  const [isDeniedSuccess, setIsDeniedSuccess] = createSignal(false);

  // Deny modal state
  const [denyOpen, setDenyOpen] = createSignal(false);
  const [denyReason, setDenyReason] = createSignal("");
  const [isDenying, setIsDenying] = createSignal(false);

  let canvasRef: SignatureCanvasRef | undefined;
  let fileInputRef: HTMLInputElement | undefined;
  let dragDepth = 0;

  const fulfillMutation = useFulfillSignatureRequestMutation();
  const denyMutation = useDenySignatureRequestMutation();

  const isExpired = createMemo(() => {
    const exp = claims()?.exp;
    if (exp && exp * 1000 < Date.now()) return true;
    if (request()?.status === "expired") return true;
    if (request()?.expires_at) {
      return new Date(request()!.expires_at).getTime() < Date.now();
    }
    return false;
  });

  const processFile = (file?: File) => {
    if (!file) return;

    if (!file.type.startsWith("image/")) {
      toast.error("Por favor, selecione um arquivo de imagem (PNG, JPG ou WebP).");
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      toast.error("O arquivo deve ter no máximo 5MB.");
      return;
    }

    setImportedFile(file);
    const reader = new FileReader();
    reader.onload = () => {
      setPreviewUrl(reader.result as string);
    };
    reader.readAsDataURL(file);
  };

  const handleFileInput = (e: Event) => {
    const target = e.target as HTMLInputElement;
    processFile(target.files?.[0]);
  };

  // Drag and drop handlers
  const handleDragEnter = (e: DragEvent) => {
    e.preventDefault();
    dragDepth += 1;
    setIsDragging(true);
  };

  const handleDragOver = (e: DragEvent) => {
    e.preventDefault();
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = "copy";
    }
  };

  const handleDragLeave = (e: DragEvent) => {
    e.preventDefault();
    dragDepth -= 1;
    if (dragDepth <= 0) {
      dragDepth = 0;
      setIsDragging(false);
    }
  };

  const handleDrop = (e: DragEvent) => {
    e.preventDefault();
    dragDepth = 0;
    setIsDragging(false);
    const file = e.dataTransfer?.files?.[0];
    processFile(file);
  };

  // Clipboard paste support for images
  const handlePaste = (e: ClipboardEvent) => {
    if (mode() !== "upload") return;
    const items = e.clipboardData?.items;
    if (!items) return;

    for (let i = 0; i < items.length; i++) {
      if (items[i].type.startsWith("image/")) {
        const file = items[i].getAsFile();
        if (file) {
          processFile(file);
          toast.success("Imagem colada da área de transferência.");
          break;
        }
      }
    }
  };

  if (typeof window !== "undefined") {
    window.addEventListener("paste", handlePaste);
    onCleanup(() => {
      window.removeEventListener("paste", handlePaste);
    });
  }

  const handleFulfill = async () => {
    if (!props.token) {
      toast.error("Token de autorização inválido ou ausente.");
      return;
    }

    if (!agreed()) {
      toast.error("Confirme a autorização de uso para continuar.");
      return;
    }

    setIsSubmitting(true);
    try {
      let imageUrl = "";

      if (mode() === "draw") {
        if (!canvasRef || canvasRef.isEmpty()) {
          toast.error("Por favor, desenhe sua assinatura antes de confirmar.");
          setIsSubmitting(false);
          return;
        }

        const blob = await canvasRef.toBlob("image/png");
        if (!blob) {
          throw new Error("Não foi possível capturar a assinatura desenhada.");
        }

        const file = new File(
          [blob],
          `signature-${Date.now()}.png`,
          { type: "image/png" },
        );

        const editionId = claims()?.edition_id ?? request()?.edition_id ?? "unknown";
        imageUrl = await uploadFile(file, `editions/${editionId}/signatures`);
      } else {
        const file = importedFile();
        if (!file) {
          toast.error("Selecione ou envie o arquivo da sua assinatura.");
          setIsSubmitting(false);
          return;
        }

        const editionId = claims()?.edition_id ?? request()?.edition_id ?? "unknown";
        imageUrl = await uploadFile(file, `editions/${editionId}/signatures`);
      }

      await fulfillMutation.mutateAsync({
        token: props.token,
        imageUrl,
        requestId: request()?.id,
        editionId: claims()?.edition_id ?? request()?.edition_id,
      });

      setIsFulfilledSuccess(true);
      toast.success("Assinatura vinculada com sucesso!");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Erro ao enviar assinatura.";
      toast.error(msg);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeny = async () => {
    if (!props.token) return;

    setIsDenying(true);
    try {
      await denyMutation.mutateAsync({
        token: props.token,
        reason: denyReason().trim() || undefined,
        requestId: request()?.id,
        editionId: claims()?.edition_id ?? request()?.edition_id,
      });

      setDenyOpen(false);
      setIsDeniedSuccess(true);
      toast.info("Solicitação de assinatura recusada.");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Erro ao recusar solicitação.";
      toast.error(msg);
    } finally {
      setIsDenying(false);
    }
  };

  const formattedExpiry = createMemo(() => {
    const raw = request()?.expires_at;
    if (!raw) return null;
    try {
      return new Date(raw).toLocaleDateString("pt-BR", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      });
    } catch {
      return null;
    }
  });

  return (
    <div class="min-h-screen bg-background text-foreground flex flex-col justify-start pt-8 sm:pt-12 pb-32 sm:pb-36 px-4 sm:px-6">
      <div class="mx-auto w-full max-w-3xl">
        {/* 1. Missing or Invalid Token State */}
        <Show when={!props.token || !claims()}>
          <div class="py-16 text-center max-w-md mx-auto space-y-4">
            <FileWarning class="size-10 text-muted-foreground mx-auto" />
            <div class="space-y-1.5">
              <h1 class="text-xl font-semibold tracking-tight text-foreground">
                Link Inválido ou Incompleto
              </h1>
              <p class="text-xs text-muted-foreground leading-relaxed">
                Não foi possível identificar o token de autenticação desta solicitação. Use o link original enviado ao seu e-mail.
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
        <Show when={props.token && claims() && requestQuery().isLoading}>
          <div class="py-24 text-center max-w-sm mx-auto space-y-3">
            <Loader2 class="mx-auto size-7 animate-spin text-muted-foreground" />
            <p class="text-xs text-muted-foreground">
              Carregando dados da assinatura...
            </p>
          </div>
        </Show>

        {/* 3. Success State */}
        <Show when={isFulfilledSuccess() || request()?.status === "completed"}>
          <div class="py-16 max-w-md mx-auto text-center space-y-5">
            <CheckCircle2 class="size-12 text-foreground mx-auto" />
            <div class="space-y-2">
              <h1 class="text-2xl font-semibold tracking-tight text-foreground">
                Assinatura Confirmada
              </h1>
              <p class="text-xs text-muted-foreground leading-relaxed">
                Sua assinatura foi registrada e será utilizada nos certificados oficiais emitidos para este evento.
              </p>
            </div>

            <div class="rounded-lg border border-border bg-muted/20 px-4 py-3 text-left text-xs space-y-1.5 text-muted-foreground">
              <div class="flex justify-between">
                <span>Signatário</span>
                <span class="font-medium text-foreground">{request()?.signatory_name ?? "Confirmado"}</span>
              </div>
              <Show when={request()?.signatory_title}>
                <div class="flex justify-between">
                  <span>Cargo / Título</span>
                  <span class="text-foreground">{request()?.signatory_title}</span>
                </div>
              </Show>
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

        {/* 4. Denied / Cancelled State */}
        <Show when={isDeniedSuccess() || request()?.status === "cancelled"}>
          <div class="py-16 max-w-md mx-auto text-center space-y-4">
            <XCircle class="size-10 text-muted-foreground mx-auto" />
            <div class="space-y-1.5">
              <h1 class="text-xl font-semibold tracking-tight text-foreground">
                Solicitação Encerrada
              </h1>
              <p class="text-xs text-muted-foreground leading-relaxed">
                Esta solicitação foi cancelada ou recusada e não pode mais receber assinaturas.
              </p>
            </div>
            <Show when={request()?.status_reason}>
              <div class="rounded-md border border-border bg-muted/30 px-3 py-2 text-xs text-muted-foreground italic">
                Motivo: {request()?.status_reason}
              </div>
            </Show>
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

        {/* 5. Expired State */}
        <Show
          when={
            !isFulfilledSuccess() &&
            !isDeniedSuccess() &&
            request()?.status !== "completed" &&
            request()?.status !== "cancelled" &&
            isExpired()
          }
        >
          <div class="py-16 max-w-md mx-auto text-center space-y-4">
            <Clock class="size-10 text-muted-foreground mx-auto" />
            <div class="space-y-1.5">
              <h1 class="text-xl font-semibold tracking-tight text-foreground">
                Solicitação Expirada
              </h1>
              <p class="text-xs text-muted-foreground leading-relaxed">
                O prazo para envio desta assinatura terminou. Caso necessário, solicite um novo envio aos organizadores.
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

        {/* 6. Active Signing Experience */}
        <Show
          when={
            props.token &&
            claims() &&
            !requestQuery().isLoading &&
            !isFulfilledSuccess() &&
            !isDeniedSuccess() &&
            request()?.status !== "completed" &&
            request()?.status !== "cancelled" &&
            !isExpired()
          }
        >
          <div class="space-y-6">
            {/* Header with clean contextual metadata */}
            <div class="space-y-1.5 border-b border-border pb-5">
              <h1 class="text-xl sm:text-2xl font-semibold tracking-tight text-foreground">
                Assinatura Oficial de Certificados
              </h1>
              <div class="flex flex-wrap items-center gap-x-2 gap-y-1 text-xs text-muted-foreground">
                <span>Signatário: <strong class="font-medium text-foreground">{request()?.signatory_name ?? "Convidado"}</strong></span>
                <Show when={request()?.signatory_title}>
                  <span>•</span>
                  <span>{request()?.signatory_title}</span>
                </Show>
                <Show when={request()?.signatory_email}>
                  <span>•</span>
                  <span>{request()?.signatory_email}</span>
                </Show>
                <Show when={formattedExpiry()}>
                  <span>•</span>
                  <span>Válido até {formattedExpiry()}</span>
                </Show>
              </div>
            </div>

            {/* Input Selection & Signature Area */}
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <div class="inline-flex rounded-lg border border-border bg-muted/40 p-0.5 text-xs font-medium">
                  <button
                    type="button"
                    onClick={() => setMode("draw")}
                    class={cn(
                      "px-3 py-1.5 rounded-md transition-colors cursor-pointer",
                      mode() === "draw"
                        ? "bg-background text-foreground shadow-xs font-semibold"
                        : "text-muted-foreground hover:text-foreground",
                    )}
                  >
                    Desenhar
                  </button>
                  <button
                    type="button"
                    onClick={() => setMode("upload")}
                    class={cn(
                      "px-3 py-1.5 rounded-md transition-colors cursor-pointer",
                      mode() === "upload"
                        ? "bg-background text-foreground shadow-xs font-semibold"
                        : "text-muted-foreground hover:text-foreground",
                    )}
                  >
                    Upload de arquivo
                  </button>
                </div>

                <Show when={mode() === "draw"}>
                  <button
                    type="button"
                    onClick={() => canvasRef?.clear()}
                    class="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                  >
                    <Eraser class="size-3.5" />
                    <span>Limpar</span>
                  </button>
                </Show>
              </div>

              {/* Mode: Draw */}
              <Show when={mode() === "draw"}>
                <div class="overflow-hidden rounded-xl border border-border bg-white shadow-xs">
                  <SignatureCanvas
                    ref={(ref) => {
                      canvasRef = ref;
                    }}
                    strokeColor="#0f172a"
                    lineWidth={3.2}
                    showBaseline={true}
                    class="w-full"
                  />
                </div>
              </Show>

              {/* Mode: Upload with full Drag & Drop and Preview */}
              <Show when={mode() === "upload"}>
                <div>
                  <Show
                    when={previewUrl()}
                    fallback={
                      <div
                        onDragEnter={handleDragEnter}
                        onDragOver={handleDragOver}
                        onDragLeave={handleDragLeave}
                        onDrop={handleDrop}
                        onClick={() => fileInputRef?.click()}
                        class={cn(
                          "flex flex-col items-center justify-center rounded-xl border-2 border-dashed p-10 text-center transition-colors cursor-pointer",
                          isDragging()
                            ? "border-primary bg-primary/5 text-primary"
                            : "border-border/80 hover:border-foreground/40 bg-muted/10 hover:bg-muted/20",
                        )}
                      >
                        <UploadCloud
                          class={cn(
                            "size-9 mb-3 transition-transform",
                            isDragging() ? "text-primary scale-110" : "text-muted-foreground",
                          )}
                        />
                        <p class="text-xs font-medium text-foreground">
                          {isDragging()
                            ? "Solte o arquivo de imagem aqui"
                            : "Arraste e solte o arquivo aqui, ou clique para selecionar"}
                        </p>
                        <p class="text-[11px] text-muted-foreground mt-1">
                          PNG, JPG ou WebP até 5MB. Recomendamos fundo transparente.
                        </p>
                        <p class="text-[10px] text-muted-foreground/70 mt-2">
                          Dica: você também pode colar uma imagem da área de transferência (Ctrl+V)
                        </p>
                      </div>
                    }
                  >
                    {/* Upload Preview with transparency checkerboard */}
                    <div class="rounded-xl border border-border bg-card p-4 space-y-3">
                      <div
                        class="flex items-center justify-center rounded-lg border border-border/60 p-6 min-h-48 overflow-hidden"
                        style={{
                          "background-image":
                            "linear-gradient(45deg, #f1f5f9 25%, transparent 25%), linear-gradient(-45deg, #f1f5f9 25%, transparent 25%), linear-gradient(45deg, transparent 75%, #f1f5f9 75%), linear-gradient(-45deg, transparent 75%, #f1f5f9 75%)",
                          "background-size": "16px 16px",
                          "background-position": "0 0, 0 8px, 8px -8px, -8px 0",
                          "background-color": "#ffffff",
                        }}
                      >
                        <img
                          src={previewUrl()}
                          alt="Pré-visualização da assinatura"
                          class="max-h-36 max-w-full object-contain filter drop-shadow-xs"
                        />
                      </div>

                      <div class="flex items-center justify-between text-xs text-muted-foreground pt-1">
                        <span class="truncate max-w-xs">{importedFile()?.name}</span>
                        <div class="flex items-center gap-3">
                          <button
                            type="button"
                            onClick={() => fileInputRef?.click()}
                            class="inline-flex items-center gap-1 hover:text-foreground cursor-pointer"
                          >
                            <RefreshCw class="size-3" />
                            <span>Trocar</span>
                          </button>
                          <span>•</span>
                          <button
                            type="button"
                            onClick={() => {
                              setImportedFile(null);
                              setPreviewUrl("");
                            }}
                            class="inline-flex items-center gap-1 text-destructive hover:underline cursor-pointer"
                          >
                            <Trash2 class="size-3" />
                            <span>Remover</span>
                          </button>
                        </div>
                      </div>
                    </div>
                  </Show>

                  <input
                    ref={(el) => {
                      fileInputRef = el;
                    }}
                    type="file"
                    accept="image/png,image/jpeg,image/webp"
                    onChange={handleFileInput}
                    class="hidden"
                  />
                </div>
              </Show>
            </div>

            {/* Legal Agreement */}
            <div class="pt-2">
              <div class="flex items-center gap-3 rounded-lg border border-border/80 bg-muted/20 p-3.5">
                <input
                  id="confirm-fulfill-agreement"
                  type="checkbox"
                  checked={agreed()}
                  onInput={(e) => {
                    const target = (e.target ?? e.currentTarget) as HTMLInputElement;
                    setAgreed(target.checked);
                  }}
                  onChange={(e) => {
                    const target = (e.target ?? e.currentTarget) as HTMLInputElement;
                    setAgreed(target.checked);
                  }}
                  class="size-4 rounded border-border text-primary focus:ring-primary cursor-pointer shrink-0"
                />
                <label
                  for="confirm-fulfill-agreement"
                  class="text-xs text-muted-foreground leading-normal select-none cursor-pointer"
                >
                  Declaro que esta assinatura é de minha autoria e autorizo sua utilização nos certificados digitais deste evento.
                </label>
              </div>
            </div>

            {/* Actions Bar */}
            <div class="flex flex-col-reverse sm:flex-row items-stretch sm:items-center justify-between gap-3 pt-3 border-t border-border">
              <button
                type="button"
                onClick={() => setDenyOpen(true)}
                disabled={isSubmitting()}
                class="text-xs text-muted-foreground hover:text-destructive transition-colors cursor-pointer text-center py-2 sm:py-0"
              >
                Recusar solicitação
              </button>

              <Button
                type="button"
                variant="default"
                onClick={handleFulfill}
                disabled={isSubmitting() || !agreed()}
                class="inline-flex items-center justify-center gap-2 cursor-pointer text-xs font-semibold px-6 py-2.5 h-10 sm:h-9 w-full sm:w-auto"
              >
                <Show when={isSubmitting()} fallback={<PenLine class="size-3.5" />}>
                  <Loader2 class="size-3.5 animate-spin" />
                </Show>
                {isSubmitting() ? "Enviando..." : "Confirmar e Assinar"}
              </Button>
            </div>
          </div>
        </Show>
      </div>

      {/* Deny Request Modal */}
      <Dialog
        open={denyOpen()}
        onOpenChange={setDenyOpen}
        title="Recusar Solicitação"
        description="Esta solicitação de assinatura será encerrada e os organizadores do evento serão notificados."
      >
        <div class="space-y-4 py-2">
          <div class="space-y-1.5">
            <Label class="text-xs font-medium text-foreground">
              Motivo da recusa (opcional):
            </Label>
            <textarea
              rows={3}
              value={denyReason()}
              onInput={(e) => setDenyReason(e.currentTarget.value)}
              placeholder="Ex: Não reconheço a solicitação ou os dados informados estão incorretos..."
              class="w-full rounded-md border border-border bg-background px-3 py-2 text-xs text-foreground placeholder:text-muted-foreground focus:outline-hidden focus:ring-2 focus:ring-primary"
            />
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => setDenyOpen(false)}
              disabled={isDenying()}
              class="cursor-pointer text-xs"
            >
              Cancelar
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={handleDeny}
              disabled={isDenying()}
              class="cursor-pointer text-xs font-medium inline-flex items-center gap-1.5"
            >
              <Show when={isDenying()}>
                <Loader2 class="size-3.5 animate-spin" />
              </Show>
              Recusar
            </Button>
          </div>
        </div>
      </Dialog>
    </div>
  );
}
