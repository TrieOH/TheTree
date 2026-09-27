import type { JSX } from "@solidjs/web";
import { Show, createSignal, untrack } from "solid-js";
import CheckCircle2Icon from "~icons/lucide/check-circle-2";
import EraserIcon from "~icons/lucide/eraser";
import ImageIcon from "~icons/lucide/image";
import PenLineIcon from "~icons/lucide/pen-line";
import UploadIcon from "~icons/lucide/upload";

import {
  MultiStepDialog,
  type MultiStepItem,
  cn,
} from "@trieoh/ui-solid";
import { uploadFile } from "@/features/storage/api/index";
import { toast } from "@/shared/ui/toast";
import { useCreateSignatureMutation } from "../api/mutations";
import { SignatureCanvas, type SignatureCanvasRef } from "./SignatureCanvas";

const CheckCircle2 = CheckCircle2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Eraser = EraserIcon as unknown as (props: { class?: string }) => JSX.Element;
const ImageLucide = ImageIcon as unknown as (props: { class?: string }) => JSX.Element;
const PenLine = PenLineIcon as unknown as (props: { class?: string }) => JSX.Element;
const Upload = UploadIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface CreateSignatureModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  eventId: string;
  editionId: string;
  onSuccess?: () => void;
}

export type SignatureOriginMode = "draw" | "upload";

export interface CreateSignatureValues {
  signatory_name: string;
  signatory_title: string;
  signatory_email: string;
  mode: SignatureOriginMode;
}

const emptyValues: CreateSignatureValues = {
  signatory_name: "",
  signatory_title: "",
  signatory_email: "",
  mode: "draw",
};

export function CreateSignatureModal(
  props: CreateSignatureModalProps,
): JSX.Element {
  const [currentStep, setCurrentStep] = createSignal(0);
  let currentValues: CreateSignatureValues = { ...emptyValues };
  const [values, setValues] = createSignal<CreateSignatureValues>(
    untrack(() => currentValues),
  );
  const [errors, setErrors] = createSignal<
    Partial<Record<keyof CreateSignatureValues | "signature", string>>
  >({});
  const [isSubmitting, setIsSubmitting] = createSignal(false);
  const [isDragging, setIsDragging] = createSignal(false);

  const [importedFile, setImportedFile] = createSignal<File | null>(null);
  const [previewUrl, setPreviewUrl] = createSignal<string | null>(null);
  const [drawingDataUrl, setDrawingDataUrl] = createSignal<string | null>(null);

  let canvasRef: SignatureCanvasRef | undefined;

  const createSignatureMutation = useCreateSignatureMutation();

  const resetForm = () => {
    currentValues = { ...emptyValues };
    setValues({ ...emptyValues });
    setErrors({});
    setCurrentStep(0);
    setImportedFile(null);
    setIsDragging(false);
    if (previewUrl()) {
      try {
        if (typeof URL.revokeObjectURL === "function") {
          URL.revokeObjectURL(previewUrl()!);
        }
      } catch { }
      setPreviewUrl(null);
    }
    setDrawingDataUrl(null);
    canvasRef?.clear();
  };

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      resetForm();
    }
    props.onOpenChange(open);
  };

  const handleChange = (
    key: (keyof CreateSignatureValues | "signature") & string,
    value: unknown,
  ) => {
    if (key in emptyValues) {
      currentValues = { ...currentValues, [key]: value };
      setValues(currentValues);
    }
    setErrors((prev) => {
      if (!prev[key as keyof typeof prev]) return prev;
      const nextErrors = { ...prev };
      delete nextErrors[key as keyof typeof prev];
      return nextErrors;
    });
  };

  const processFile = (file?: File) => {
    if (!file) return;

    if (!file.type.startsWith("image/")) {
      toast.error("Por favor, selecione um arquivo de imagem (PNG, JPG ou WebP).");
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      toast.error("A imagem deve ter no máximo 5MB.");
      return;
    }

    setImportedFile(file);
    if (previewUrl()) {
      try {
        if (typeof URL.revokeObjectURL === "function") {
          URL.revokeObjectURL(previewUrl()!);
        }
      } catch { }
    }
    try {
      if (typeof URL.createObjectURL === "function") {
        setPreviewUrl(URL.createObjectURL(file));
      } else {
        setPreviewUrl("blob:mock-file");
      }
    } catch {
      setPreviewUrl("blob:mock-file");
    }
    setErrors((prev) => {
      if (!prev.signature) return prev;
      const nextErrors = { ...prev };
      delete nextErrors.signature;
      return nextErrors;
    });
  };

  const handleFileChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    processFile(file);
  };

  const handleDrop = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    const file = e.dataTransfer?.files?.[0];
    processFile(file);
  };

  const validateStep = (stepIndex: number): boolean => {
    const current = currentValues;
    const nextErrors: Partial<
      Record<keyof CreateSignatureValues | "signature", string>
    > = {};

    if (stepIndex === 0) {
      const name = current.signatory_name.trim();
      if (name.length < 2) {
        nextErrors.signatory_name =
          "O nome do signatário deve ter pelo menos 2 caracteres.";
      }

      const email = current.signatory_email.trim();
      if (email.length > 0 && (!email.includes("@") || !email.includes("."))) {
        nextErrors.signatory_email = "Informe um e-mail válido.";
      }
    }

    if (stepIndex === 1) {
      if (current.mode === "draw") {
        if (!canvasRef || canvasRef.isEmpty()) {
          nextErrors.signature = "Por favor, desenhe uma assinatura no quadro.";
          toast.error("Por favor, desenhe uma assinatura no quadro.");
        } else {
          setDrawingDataUrl(canvasRef.toDataURL());
        }
      } else {
        if (!importedFile()) {
          nextErrors.signature =
            "Selecione um arquivo de imagem com a assinatura.";
          toast.error("Selecione um arquivo de imagem com a assinatura.");
        }
      }
    }

    if (Object.keys(nextErrors).length > 0) {
      setErrors(nextErrors);
      return false;
    }

    setErrors({});
    return true;
  };

  const handleSubmit = async (): Promise<boolean> => {
    const current = currentValues;
    const name = current.signatory_name.trim();

    setIsSubmitting(true);
    try {
      let imageUrl = "";

      if (current.mode === "draw") {
        if (!canvasRef || canvasRef.isEmpty()) {
          toast.error("Por favor, desenhe uma assinatura no quadro.");
          return false;
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
          `editions/${props.editionId}/signatures`,
        );
      } else {
        const file = importedFile();
        if (!file) {
          toast.error("Selecione um arquivo de imagem com a assinatura.");
          return false;
        }

        imageUrl = await uploadFile(
          file,
          `editions/${props.editionId}/signatures`,
        );
      }

      await createSignatureMutation.mutateAsync({
        editionId: props.editionId,
        data: {
          signatory_name: name,
          signatory_title: current.signatory_title.trim() || undefined,
          signatory_email: current.signatory_email.trim() || undefined,
          image_url: imageUrl,
        },
      });

      toast.success("Assinatura adicionada com sucesso!");
      resetForm();
      props.onOpenChange(false);
      props.onSuccess?.();
      return true;
    } catch (err: unknown) {
      const msg =
        err instanceof Error ? err.message : "Erro ao cadastrar assinatura.";
      toast.error(msg);
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  const activeSignaturePreview = () => {
    if (values().mode === "draw") {
      return drawingDataUrl() ?? (canvasRef && !canvasRef.isEmpty() ? canvasRef.toDataURL() : null);
    }
    return previewUrl();
  };

  const steps: MultiStepItem<CreateSignatureValues>[] = [
    {
      id: "signatario",
      title: "Signatário",
      description: "Identificação da autoridade que assinará os documentos",
      fields: [
        {
          name: "signatory_name",
          label: "Nome do signatário",
          required: true,
          placeholder: "Ex.: Prof. Dra. Carolina Mendes",
          hint: "Nome completo como deve constar no certificado.",
        },
        {
          name: "signatory_title",
          label: "Cargo / Função",
          layout: "half",
          placeholder: "Ex.: Coordenadora Científica, Reitora",
          hint: "Titulação ou cargo oficial.",
        },
        {
          name: "signatory_email",
          label: "E-mail de contato",
          kind: "email",
          layout: "half",
          placeholder: "carolina.mendes@evento.com",
          hint: "E-mail institucional (opcional).",
        },
      ],
    },
    {
      id: "assinatura",
      title: "Assinatura Digital",
      description: "Desenhe diretamente na tela ou importe uma imagem",
      render: () => (
        <div class="space-y-4">
          {/* Mode Selector */}
          <div class="space-y-1.5">
            <span class="text-xs font-medium text-foreground">
              Origem da assinatura
            </span>
            <div class="grid grid-cols-2 gap-2.5">
              <button
                type="button"
                onClick={() => handleChange("mode", "draw")}
                class={cn(
                  "flex items-center justify-center gap-2 rounded-xl border p-3 text-xs font-medium transition-all cursor-pointer",
                  values().mode === "draw"
                    ? "border-primary bg-primary/10 text-primary font-semibold ring-1 ring-primary/20 shadow-xs"
                    : "border-border bg-card text-muted-foreground hover:bg-muted/50 hover:text-foreground",
                )}
              >
                <PenLine class="size-4" />
                <span>Desenhar no quadro</span>
              </button>
              <button
                type="button"
                onClick={() => handleChange("mode", "upload")}
                class={cn(
                  "flex items-center justify-center gap-2 rounded-xl border p-3 text-xs font-medium transition-all cursor-pointer",
                  values().mode === "upload"
                    ? "border-primary bg-primary/10 text-primary font-semibold ring-1 ring-primary/20 shadow-xs"
                    : "border-border bg-card text-muted-foreground hover:bg-muted/50 hover:text-foreground",
                )}
              >
                <Upload class="size-4" />
                <span>Importar imagem</span>
              </button>
            </div>
          </div>

          {/* Draw or Upload Content */}
          <Show
            when={values().mode === "draw"}
            fallback={
              <div class="space-y-2">
                <label
                  for="sig-file-upload"
                  onDragOver={(e) => {
                    e.preventDefault();
                    setIsDragging(true);
                  }}
                  onDragLeave={(e) => {
                    e.preventDefault();
                    setIsDragging(false);
                  }}
                  onDrop={handleDrop}
                  class={cn(
                    "flex flex-col items-center justify-center gap-2 rounded-xl border-2 border-dashed bg-muted/10 p-6 text-center transition-all cursor-pointer hover:border-primary/60 hover:bg-muted/20",
                    isDragging() && "border-primary bg-primary/10",
                    errors().signature
                      ? "border-destructive/70 bg-destructive/5"
                      : "border-border/80",
                  )}
                >
                  <Show
                    when={previewUrl()}
                    fallback={
                      <>
                        <div class="flex size-11 items-center justify-center rounded-full bg-muted text-muted-foreground shadow-xs">
                          <ImageLucide class="size-5" />
                        </div>
                        <div class="space-y-1">
                          <p class="text-xs font-semibold text-foreground">
                            Clique para escolher ou arraste o arquivo
                          </p>
                          <p class="text-[11px] text-muted-foreground">
                            Formatos aceitos: PNG transparente, WebP ou JPEG (máx. 5MB)
                          </p>
                        </div>
                      </>
                    }
                  >
                    <div class="relative max-h-32 w-full overflow-hidden rounded-lg bg-card p-3 shadow-xs">
                      <img
                        src={previewUrl()!}
                        alt="Prévia da assinatura"
                        class="mx-auto max-h-24 object-contain"
                      />
                    </div>
                    <span class="text-xs font-medium text-primary hover:underline">
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
                <span class="text-xs text-muted-foreground">
                  Assine com o mouse, trackpad ou caneta:
                </span>
                <button
                  type="button"
                  onClick={() => {
                    canvasRef?.clear();
                    setDrawingDataUrl(null);
                  }}
                  class="inline-flex items-center gap-1.5 text-xs font-medium text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                >
                  <Eraser class="size-3.5" />
                  <span>Limpar quadro</span>
                </button>
              </div>

              <div
                class={cn(
                  "overflow-hidden rounded-xl border bg-white shadow-xs transition-colors",
                  errors().signature
                    ? "border-destructive ring-1 ring-destructive/30"
                    : "border-border",
                )}
              >
                <SignatureCanvas
                  ref={(r) => {
                    canvasRef = r;
                  }}
                  onStroke={() => {
                    if (errors().signature) {
                      setErrors((prev) => {
                        const next = { ...prev };
                        delete next.signature;
                        return next;
                      });
                    }
                  }}
                />
              </div>
            </div>
          </Show>

          <Show when={errors().signature}>
            <p class="text-xs font-medium text-destructive">
              {errors().signature}
            </p>
          </Show>
        </div>
      ),
    },
    {
      id: "resumo",
      title: "Resumo",
      description: "Revise os dados antes de salvar a assinatura digital",
      summary: {
        title: "Dados da assinatura",
        badge: () => "Pronto para cadastrar",
        items: [
          {
            label: "Nome do signatário",
            value: () => values().signatory_name || "Não informado",
          },
          {
            label: "Cargo / Função",
            value: () => values().signatory_title || "Não informado",
          },
          {
            label: "E-mail de contato",
            value: () => values().signatory_email || "Não informado",
          },
          {
            label: "Origem do traço",
            value: () =>
              values().mode === "draw"
                ? "Quadro digital interativo"
                : "Arquivo de imagem importado",
          },
        ],
        extra: () => (
          <div class="relative overflow-hidden rounded-xl border border-border/80 bg-card p-5 shadow-xs">
            <div class="flex items-center justify-between pb-3 border-b border-border/60">
              <span class="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                Prévia do Certificado
              </span>
              <span class="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2.5 py-0.5 text-xs font-semibold text-emerald-600 dark:text-emerald-400">
                <CheckCircle2 class="size-3.5" />
                <span>Autenticação Digital</span>
              </span>
            </div>

            <div class="mt-4 flex flex-col items-center justify-center rounded-lg border border-dashed border-border/70 bg-muted/15 p-6">
              <Show
                when={activeSignaturePreview()}
                fallback={
                  <div class="h-16 flex items-center justify-center text-xs text-muted-foreground italic">
                    Assinatura não capturada
                  </div>
                }
              >
                <img
                  src={activeSignaturePreview()!}
                  alt="Prévia da assinatura"
                  class="max-h-24 max-w-full object-contain filter contrast-125 dark:invert"
                />
              </Show>

              <div class="mt-2 w-52 border-b-2 border-foreground/40" />

              <h4 class="mt-2 text-sm font-bold text-foreground">
                {values().signatory_name || "Nome do Signatário"}
              </h4>
              <Show when={values().signatory_title}>
                <p class="text-xs text-muted-foreground italic mt-0.5">
                  {values().signatory_title}
                </p>
              </Show>
            </div>
          </div>
        ),
      },
    },
  ];

  return (
    <MultiStepDialog
      open={props.open}
      onOpenChange={handleOpenChange}
      title="Adicionar assinatura"
      description="Cadastre uma assinatura digital para autenticar certificados desta edição."
      steps={steps}
      currentStep={currentStep()}
      onStepChange={setCurrentStep}
      values={values()}
      errors={errors()}
      onChange={handleChange}
      loading={isSubmitting()}
      submitLabel={isSubmitting() ? "Salvando..." : "Salvar assinatura"}
      onBeforeNext={validateStep}
      onSubmit={handleSubmit}
    />
  );
}

export { CreateSignatureModal as ManageSignatureDialog };
