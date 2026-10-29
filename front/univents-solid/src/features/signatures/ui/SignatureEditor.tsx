import type { JSX } from "@solidjs/web";
import { useNavigate } from "@tanstack/solid-router";
import { Show, createSignal } from "solid-js";
import ArrowLeftIcon from "~icons/lucide/arrow-left";
import EraserIcon from "~icons/lucide/eraser";
import ImageIcon from "~icons/lucide/image";
import Loader2Icon from "~icons/lucide/loader-2";
import PenLineIcon from "~icons/lucide/pen-line";
import UploadIcon from "~icons/lucide/upload";

import { Button, Input, Label, cn } from "@trieoh/ui-solid";
import { uploadFile } from "@/features/storage/api/index";
import { toast } from "@/shared/ui/toast";
import { useCreateSignatureMutation } from "../api/mutations";
import { SignatureCanvas, type SignatureCanvasRef } from "./SignatureCanvas";

const ArrowLeft = ArrowLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const Eraser = EraserIcon as unknown as (props: { class?: string }) => JSX.Element;
const ImageLucide = ImageIcon as unknown as (props: { class?: string }) => JSX.Element;
const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;
const PenLine = PenLineIcon as unknown as (props: { class?: string }) => JSX.Element;
const Upload = UploadIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface SignatureEditorProps {
  eventId: string;
  editionId: string;
}

type Mode = "draw" | "upload";

export function SignatureEditor(props: SignatureEditorProps): JSX.Element {
  const navigate = useNavigate();
  const [mode, setMode] = createSignal<Mode>("draw");
  const [signatoryName, setSignatoryName] = createSignal("");
  const [signatoryTitle, setSignatoryTitle] = createSignal("");
  const [signatoryEmail, setSignatoryEmail] = createSignal("");
  const [importedFile, setImportedFile] = createSignal<File | null>(null);
  const [previewUrl, setPreviewUrl] = createSignal<string | null>(null);
  const [isSubmitting, setIsSubmitting] = createSignal(false);

  let canvasRef: SignatureCanvasRef | undefined;

  const createSignatureMutation = useCreateSignatureMutation();

  const handleFileChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (file) {
      setImportedFile(file);
      const existing = previewUrl();
      if (existing) {
        URL.revokeObjectURL(existing);
      }
      setPreviewUrl(URL.createObjectURL(file));
    }
  };

  const handleSave = async () => {
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
          throw new Error("Não foi possível capturar a assinatura do quadro.");
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
          toast.error("Selecione um arquivo de imagem da assinatura.");
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
      navigate({
        to: "/admin/events/$eventId/editions/$editionId/signatures",
        params: { eventId: props.eventId, editionId: props.editionId },
      });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Erro ao salvar assinatura.";
      toast.error(msg);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div class="mx-auto max-w-5xl space-y-6">
      {/* Top Header */}
      <div class="flex items-center gap-3">
        <button
          type="button"
          onClick={() =>
            navigate({
              to: "/admin/events/$eventId/editions/$editionId/signatures",
              params: { eventId: props.eventId, editionId: props.editionId },
            })
          }
          class="flex size-9 items-center justify-center rounded-lg border border-border/80 bg-background text-muted-foreground transition-colors hover:border-border hover:text-foreground cursor-pointer"
          title="Voltar para a lista de assinaturas"
        >
          <ArrowLeft class="size-4" />
        </button>
        <div>
          <h1 class="text-xl font-bold tracking-tight text-foreground sm:text-2xl">
            Nova Assinatura Digital
          </h1>
          <p class="text-xs text-muted-foreground">
            Cadastre uma assinatura para uso automático na emissão de crachás e certificados.
          </p>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-6 lg:grid-cols-12">
        {/* Left Column: Signatory Information */}
        <div class="space-y-5 lg:col-span-5">
          <div class="rounded-xl border border-border/70 bg-card p-5 shadow-xs space-y-4">
            <h2 class="text-sm font-semibold text-foreground">
              Dados do Signatário
            </h2>

            <div class="space-y-1.5">
              <Label for="signatory-name">
                Nome completo <span class="text-destructive">*</span>
              </Label>
              <Input
                id="signatory-name"
                type="text"
                placeholder="Ex: Prof. Dr. Carlos Eduardo Mendes"
                value={signatoryName()}
                onInput={(e) => setSignatoryName(e.currentTarget.value)}
                disabled={isSubmitting()}
              />
              <p class="text-[11px] text-muted-foreground">
                Nome como deve constar nos documentos emitidos.
              </p>
            </div>

            <div class="space-y-1.5">
              <Label for="signatory-title">Cargo / Função (opcional)</Label>
              <Input
                id="signatory-title"
                type="text"
                placeholder="Ex: Coordenador Geral, Diretor, Reitor"
                value={signatoryTitle()}
                onInput={(e) => setSignatoryTitle(e.currentTarget.value)}
                disabled={isSubmitting()}
              />
              <p class="text-[11px] text-muted-foreground">
                Título institucional que confere autoridade à assinatura.
              </p>
            </div>

            <div class="space-y-1.5">
              <Label for="signatory-email">E-mail (opcional)</Label>
              <Input
                id="signatory-email"
                type="email"
                placeholder="carlos.mendes@instituicao.org"
                value={signatoryEmail()}
                onInput={(e) => setSignatoryEmail(e.currentTarget.value)}
                disabled={isSubmitting()}
              />
              <p class="text-[11px] text-muted-foreground">
                Para contato interno e auditoria de emissão.
              </p>
            </div>
          </div>

          <div class="rounded-xl border border-border/70 bg-card p-5 shadow-xs space-y-3">
            <h3 class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Como funciona
            </h3>
            <ul class="space-y-2 text-xs text-muted-foreground">
              <li class="flex items-start gap-2">
                <span class="mt-0.5 flex size-4 shrink-0 items-center justify-center rounded-full bg-primary/10 text-[10px] font-bold text-primary">
                  1
                </span>
                <span>
                  Cadastre a assinatura desenhando na tela ou enviando uma imagem escaneada.
                </span>
              </li>
              <li class="flex items-start gap-2">
                <span class="mt-0.5 flex size-4 shrink-0 items-center justify-center rounded-full bg-primary/10 text-[10px] font-bold text-primary">
                  2
                </span>
                <span>
                  No editor de crachá e certificado, adicione o bloco de assinatura e vincule a este registro.
                </span>
              </li>
              <li class="flex items-start gap-2">
                <span class="mt-0.5 flex size-4 shrink-0 items-center justify-center rounded-full bg-primary/10 text-[10px] font-bold text-primary">
                  3
                </span>
                <span>
                  Durante a emissão em lote, cada documento é renderizado com a assinatura vinculada.
                </span>
              </li>
            </ul>
          </div>

          <div class="flex items-center gap-3">
            <Button
              type="button"
              variant="outline"
              onClick={() =>
                navigate({
                  to: "/admin/events/$eventId/editions/$editionId/signatures",
                  params: { eventId: props.eventId, editionId: props.editionId },
                })
              }
              disabled={isSubmitting()}
              class="flex-1"
            >
              Cancelar
            </Button>
            <Button
              type="button"
              variant="default"
              onClick={handleSave}
              disabled={isSubmitting()}
              class="flex-1 gap-2"
            >
              <Show
                when={isSubmitting()}
                fallback={<span>Salvar Assinatura</span>}
              >
                <Loader2 class="size-4 animate-spin" />
                <span>Salvando...</span>
              </Show>
            </Button>
          </div>
        </div>

        {/* Right Column: Signature input (Canvas or Upload) */}
        <div class="space-y-4 lg:col-span-7">
          <div class="flex items-center justify-between border-b border-border pb-3">
            <div>
              <h2 class="text-sm font-semibold text-foreground">
                Traço da Assinatura
              </h2>
              <p class="text-xs text-muted-foreground">
                Escolha o método mais conveniente para registrar o traço.
              </p>
            </div>

            <div class="flex items-center rounded-lg border border-border/80 bg-muted/40 p-0.5">
              <button
                type="button"
                onClick={() => setMode("draw")}
                class={cn(
                  "flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium transition-all cursor-pointer",
                  mode() === "draw"
                    ? "bg-card text-foreground shadow-xs font-semibold"
                    : "text-muted-foreground hover:text-foreground",
                )}
              >
                <PenLine class="size-3.5" />
                <span>Desenhar</span>
              </button>
              <button
                type="button"
                onClick={() => setMode("upload")}
                class={cn(
                  "flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium transition-all cursor-pointer",
                  mode() === "upload"
                    ? "bg-card text-foreground shadow-xs font-semibold"
                    : "text-muted-foreground hover:text-foreground",
                )}
              >
                <Upload class="size-3.5" />
                <span>Importar arquivo</span>
              </button>
            </div>
          </div>

          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-foreground">
                {mode() === "draw" ? "Desenhe no quadro" : "Envie a imagem"}
              </p>
              <p class="text-[11px] text-muted-foreground">
                {mode() === "draw"
                  ? "Assine no quadro abaixo usando o mouse ou tela de toque."
                  : "Carregue um arquivo PNG transparente com a assinatura escaneada."}
              </p>
            </div>

            <Show when={mode() === "draw"}>
              <button
                type="button"
                onClick={() => canvasRef?.clear()}
                class="inline-flex items-center gap-1.5 rounded-lg border border-border/60 bg-background/80 px-2.5 py-1 text-xs font-medium text-muted-foreground transition-colors hover:border-border hover:text-foreground cursor-pointer"
              >
                <Eraser class="size-3.5" />
                <span>Limpar</span>
              </button>
            </Show>
          </div>

          <Show
            when={mode() === "draw"}
            fallback={
              <label
                for="editor-file-upload"
                class="flex min-h-64 flex-col items-center justify-center gap-3 rounded-lg border-2 border-dashed border-border/80 bg-muted/10 p-8 text-center transition-colors hover:border-primary/60 hover:bg-muted/20 cursor-pointer"
              >
                <Show
                  when={previewUrl()}
                  fallback={
                    <>
                      <div class="flex size-12 items-center justify-center rounded-full bg-muted text-muted-foreground">
                        <ImageLucide class="size-6" />
                      </div>
                      <div>
                        <p class="text-sm font-medium text-foreground">
                          Clique ou arraste a imagem da assinatura aqui
                        </p>
                        <p class="mt-1 text-xs text-muted-foreground">
                          PNG com transparência, WebP ou JPEG
                        </p>
                      </div>
                    </>
                  }
                >
                  <div class="relative max-h-48 w-full overflow-hidden rounded bg-white p-4">
                    <img
                      src={previewUrl() ?? ""}
                      alt="Prévia da assinatura"
                      class="mx-auto max-h-40 object-contain"
                    />
                  </div>
                  <span class="text-xs text-muted-foreground hover:underline">
                    Clique para selecionar outra imagem
                  </span>
                </Show>
                <input
                  id="editor-file-upload"
                  type="file"
                  accept="image/png,image/jpeg,image/webp"
                  class="sr-only"
                  onChange={handleFileChange}
                />
              </label>
            }
          >
            <div class="overflow-hidden rounded-lg border border-border bg-white shadow-xs">
              <SignatureCanvas
                ref={(r) => {
                  canvasRef = r;
                }}
              />
            </div>
          </Show>
        </div>
      </div>
    </div>
  );
}
