import type { JSX } from "@solidjs/web";
import { Show, createSignal } from "solid-js";
import EraserIcon from "~icons/lucide/eraser";
import ImageIcon from "~icons/lucide/image";
import Loader2Icon from "~icons/lucide/loader-2";
import PenLineIcon from "~icons/lucide/pen-line";
import UploadIcon from "~icons/lucide/upload";

import { Button, Dialog, Input, Label, cn } from "@trieoh/ui-solid";
import { uploadFile } from "@/features/storage/api/index";
import { toast } from "@/shared/ui/toast";
import { useCreateSignatureMutation } from "../api/mutations";
import { SignatureCanvas, type SignatureCanvasRef } from "./SignatureCanvas";

const Eraser = EraserIcon as unknown as (props: { class?: string }) => JSX.Element;
const ImageLucide = ImageIcon as unknown as (props: { class?: string }) => JSX.Element;
const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;
const PenLine = PenLineIcon as unknown as (props: { class?: string }) => JSX.Element;
const Upload = UploadIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface CreateSignatureModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  eventId: string;
  editionId: string;
  onSuccess?: () => void;
}

type Mode = "draw" | "upload";

export function CreateSignatureModal(
  props: CreateSignatureModalProps,
): JSX.Element {
  const [mode, setMode] = createSignal<Mode>("draw");
  const [signatoryName, setSignatoryName] = createSignal("");
  const [signatoryTitle, setSignatoryTitle] = createSignal("");
  const [signatoryEmail, setSignatoryEmail] = createSignal("");
  const [importedFile, setImportedFile] = createSignal<File | null>(null);
  const [previewUrl, setPreviewUrl] = createSignal<string | null>(null);
  const [isSubmitting, setIsSubmitting] = createSignal(false);
  const [hasDrawn, setHasDrawn] = createSignal(false);

  let canvasRef: SignatureCanvasRef | undefined;

  const createSignatureMutation = useCreateSignatureMutation();

  const resetForm = () => {
    setSignatoryName("");
    setSignatoryTitle("");
    setSignatoryEmail("");
    setImportedFile(null);
    if (previewUrl()) {
      URL.revokeObjectURL(previewUrl()!);
      setPreviewUrl(null);
    }
    canvasRef?.clear();
    setHasDrawn(false);
  };

  const handleFileChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (file) {
      setImportedFile(file);
      if (previewUrl()) {
        URL.revokeObjectURL(previewUrl()!);
      }
      setPreviewUrl(URL.createObjectURL(file));
    }
  };

  const handleSubmit = async () => {
    const name = signatoryName().trim();
    if (!name || name.length < 2) {
      toast.error("Informe o nome do signatário (mínimo 2 caracteres).");
      return;
    }

    setIsSubmitting(true);
    try {
      let imageUrl = "";

      if (mode() === "draw") {
        if (!canvasRef || canvasRef.isEmpty()) {
          toast.error("Por favor, desenhe uma assinatura no quadro.");
          setIsSubmitting(false);
          return;
        }

        const blob = await canvasRef.toBlob("image/png");
        if (!blob) {
          throw new Error("Não foi possível capturar a assinatura do canvas.");
        }

        const file = new File(
          [blob],
          `signature-${Date.now()}.png`,
          { type: "image/png" },
        );

        imageUrl = await uploadFile(
          file,
          `events/${props.eventId}/editions/${props.editionId}/signatures`,
        );
      } else {
        const file = importedFile();
        if (!file) {
          toast.error("Selecione um arquivo de imagem para a assinatura.");
          setIsSubmitting(false);
          return;
        }

        imageUrl = await uploadFile(
          file,
          `events/${props.eventId}/editions/${props.editionId}/signatures`,
        );
      }

      await createSignatureMutation.mutateAsync({
        editionId: props.editionId,
        data: {
          signatory_name: name,
          signatory_title: signatoryTitle().trim() || undefined,
          signatory_email: signatoryEmail().trim() || undefined,
          image_url: imageUrl,
        },
      });

      toast.success("Assinatura criada com sucesso!");
      resetForm();
      props.onOpenChange(false);
      props.onSuccess?.();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Erro ao salvar assinatura.";
      toast.error(msg);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={(open) => {
        if (!open) resetForm();
        props.onOpenChange(open);
      }}
      title="Adicionar assinatura"
      description="Crie uma nova assinatura para autenticar certificados desta edição."
      footer={
        <div class="flex w-full items-center justify-end gap-2 pt-3 border-t border-border">
          <Button
            variant="outline"
            disabled={isSubmitting()}
            onClick={() => props.onOpenChange(false)}
          >
            Cancelar
          </Button>
          <Button
            disabled={isSubmitting()}
            onClick={handleSubmit}
            class="min-w-28"
          >
            <Show when={isSubmitting()} fallback="Salvar assinatura">
              <span class="inline-flex items-center gap-2">
                <Loader2 class="size-4 animate-spin" />
                <span>Salvando...</span>
              </span>
            </Show>
          </Button>
        </div>
      }
    >
      <div class="space-y-4 py-2 text-left">
        {/* Fields */}
        <div class="grid gap-3 sm:grid-cols-2">
          <div class="space-y-1.5 sm:col-span-2">
            <Label for="sig-name" class="text-xs font-medium">
              Nome do signatário <span class="text-destructive">*</span>
            </Label>
            <Input
              id="sig-name"
              placeholder="Ex.: Prof. Dra. Maria Souza"
              value={signatoryName()}
              onInput={(e) => setSignatoryName((e.target as HTMLInputElement).value)}
              disabled={isSubmitting()}
            />
          </div>

          <div class="space-y-1.5">
            <Label for="sig-title" class="text-xs font-medium">
              Cargo / Função
            </Label>
            <Input
              id="sig-title"
              placeholder="Ex.: Coordenador do Evento"
              value={signatoryTitle()}
              onInput={(e) => setSignatoryTitle((e.target as HTMLInputElement).value)}
              disabled={isSubmitting()}
            />
          </div>

          <div class="space-y-1.5">
            <Label for="sig-email" class="text-xs font-medium">
              E-mail do signatário
            </Label>
            <Input
              id="sig-email"
              type="email"
              placeholder="nome@instituicao.org"
              value={signatoryEmail()}
              onInput={(e) => setSignatoryEmail((e.target as HTMLInputElement).value)}
              disabled={isSubmitting()}
            />
          </div>
        </div>

        {/* Mode Selector */}
        <div class="space-y-1.5 pt-1">
          <Label class="text-xs font-medium">Origem da assinatura</Label>
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              onClick={() => setMode("draw")}
              class={cn(
                "inline-flex items-center justify-center gap-2 rounded-lg border px-3 py-2 text-xs font-medium transition-colors cursor-pointer",
                mode() === "draw"
                  ? "border-primary bg-primary/10 text-primary font-semibold"
                  : "border-border bg-card text-muted-foreground hover:text-foreground",
              )}
            >
              <PenLine class="size-3.5" />
              <span>Desenhar no quadro</span>
            </button>
            <button
              type="button"
              onClick={() => setMode("upload")}
              class={cn(
                "inline-flex items-center justify-center gap-2 rounded-lg border px-3 py-2 text-xs font-medium transition-colors cursor-pointer",
                mode() === "upload"
                  ? "border-primary bg-primary/10 text-primary font-semibold"
                  : "border-border bg-card text-muted-foreground hover:text-foreground",
              )}
            >
              <Upload class="size-3.5" />
              <span>Importar imagem</span>
            </button>
          </div>
        </div>

        {/* Canvas or File Upload area */}
        <Show
          when={mode() === "draw"}
          fallback={
            <div class="space-y-2">
              <label
                for="sig-file-upload"
                class="flex flex-col items-center justify-center gap-2 rounded-lg border-2 border-dashed border-border/80 bg-muted/20 p-6 text-center transition-colors hover:border-primary/60 hover:bg-muted/40 cursor-pointer"
              >
                <Show
                  when={previewUrl()}
                  fallback={
                    <>
                      <div class="flex size-10 items-center justify-center rounded-full bg-muted text-muted-foreground">
                        <ImageLucide class="size-5" />
                      </div>
                      <div>
                        <p class="text-xs font-medium text-foreground">
                          Clique para escolher ou arraste uma imagem
                        </p>
                        <p class="mt-0.5 text-[11px] text-muted-foreground">
                          PNG com fundo transparente, WebP ou JPEG
                        </p>
                      </div>
                    </>
                  }
                >
                  <div class="relative max-h-32 w-full overflow-hidden rounded bg-white p-2">
                    <img
                      src={previewUrl()!}
                      alt="Prévia da assinatura"
                      class="mx-auto max-h-28 object-contain"
                    />
                  </div>
                  <span class="text-[11px] text-muted-foreground hover:underline">
                    Clique para trocar de imagem
                  </span>
                </Show>
                <input
                  id="sig-file-upload"
                  type="file"
                  accept="image/png,image/jpeg,image/webp"
                  class="sr-only"
                  onChange={handleFileChange}
                />
              </label>
            </div>
          }
        >
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <span class="text-[11px] text-muted-foreground">
                Assine com o mouse, trackpad ou caneta stylus:
              </span>
              <button
                type="button"
                onClick={() => {
                  canvasRef?.clear();
                  setHasDrawn(false);
                }}
                class="inline-flex items-center gap-1 text-[11px] font-medium text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
              >
                <Eraser class="size-3" />
                <span>Limpar</span>
              </button>
            </div>

            <div class="overflow-hidden rounded-lg border border-border bg-white shadow-xs">
              <SignatureCanvas
                ref={(r) => {
                  canvasRef = r;
                }}
                onStroke={() => setHasDrawn(true)}
              />
            </div>
          </div>
        </Show>
      </div>
    </Dialog>
  );
}
