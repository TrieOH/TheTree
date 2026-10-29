import type { JSX } from "@solidjs/web";
import { Show, createSignal } from "solid-js";
import CheckCircle2Icon from "~icons/lucide/check-circle-2";
import EraserIcon from "~icons/lucide/eraser";
import PenLineIcon from "~icons/lucide/pen-line";
import UploadIcon from "~icons/lucide/upload";

import {
  MultiStepDialog,
  type MultiStepItem,
  cn,
} from "@trieoh/ui-solid";
import { uploadFile } from "@/features/storage/api/index";
import { toast } from "@/shared/ui/toast";
import { getFieldErrors } from "@/shared/lib/effect-schema";
import { useCreateSignatureMutation } from "../api/mutations";
import { signatureStep1Schema } from "../model";
import { SignatureCanvas, type SignatureCanvasRef } from "./SignatureCanvas";

const CheckCircle2 = CheckCircle2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Eraser = EraserIcon as unknown as (props: { class?: string }) => JSX.Element;
const PenLine = PenLineIcon as unknown as (props: { class?: string }) => JSX.Element;
const Upload = UploadIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface CreateSignatureModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  eventId: string;
  editionId: string;
  onSuccess?: () => void;
}

export interface CreateSignatureValues {
  signatory_name: string;
  signatory_title: string;
  signatory_email: string;
  mode: "draw" | "upload";
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
  const createSignatureMutation = useCreateSignatureMutation();

  const [currentStep, setCurrentStep] = createSignal(0);
  let currentValues: CreateSignatureValues = { ...emptyValues };
  const [values, setValues] = createSignal<CreateSignatureValues>({
    ...emptyValues,
  });

  const [isDragging, setIsDragging] = createSignal(false);
  const [importedFile, setImportedFile] = createSignal<File | null>(null);
  const [previewUrl, setPreviewUrl] = createSignal<string | null>(null);
  const [drawingDataUrl, setDrawingDataUrl] = createSignal<string | null>(null);
  const [isSubmitting, setIsSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<
    Partial<Record<keyof CreateSignatureValues | "signature", string>>
  >({});

  let canvasRef: SignatureCanvasRef | undefined;

  const resetForm = () => {
    currentValues = { ...emptyValues };
    setValues({ ...emptyValues });
    setImportedFile(null);
    const existing = previewUrl();
    if (existing) {
      try {
        if (typeof URL.revokeObjectURL === "function") {
          URL.revokeObjectURL(existing);
        }
      } catch {}
    }
    setPreviewUrl(null);
    setDrawingDataUrl(null);
    setErrors({});
    setCurrentStep(0);
  };

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      resetForm();
    }
    props.onOpenChange(open);
  };

  const handleChange = (
    key: keyof CreateSignatureValues & string,
    value: unknown,
  ) => {
    currentValues = { ...currentValues, [key]: value };
    setValues(currentValues);
    setErrors((prev) => {
      if (!prev[key as keyof CreateSignatureValues]) return prev;
      const nextErrors = { ...prev };
      delete nextErrors[key as keyof CreateSignatureValues];
      return nextErrors;
    });
  };

  const handleModeChange = (mode: "draw" | "upload") => {
    handleChange("mode", mode);
    setErrors((prev) => {
      const nextErrors = { ...prev };
      delete nextErrors.signature;
      return nextErrors;
    });
  };

  const processFile = (file?: File) => {
    if (!file) return;

    if (!file.type.startsWith("image/")) {
      toast.error("Por favor, selecione apenas arquivos de imagem.");
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      toast.error("O arquivo deve ter no máximo 5MB.");
      return;
    }

    setImportedFile(file);
    const existing = previewUrl();
    if (existing) {
      try {
        if (typeof URL.revokeObjectURL === "function") {
          URL.revokeObjectURL(existing);
        }
      } catch {}
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
    let nextErrors: Partial<
      Record<keyof CreateSignatureValues | "signature", string>
    > = {};

    if (stepIndex === 0) {
      nextErrors = getFieldErrors<CreateSignatureValues>(
        signatureStep1Schema.safeParse(current),
      );
    } else if (stepIndex === 1) {
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
          setErrors({
            signature: "Por favor, desenhe uma assinatura no quadro.",
          });
          toast.error("Por favor, desenhe uma assinatura no quadro.");
          return false;
        }

        const blob = await canvasRef.toBlob();
        if (!blob) {
          toast.error("Erro ao gerar imagem da assinatura.");
          return false;
        }

        const file = new File([blob], `signature_${Date.now()}.png`, {
          type: "image/png",
        });
        imageUrl = await uploadFile(file, "signatures");
      } else {
        const file = importedFile();
        if (!file) {
          setErrors({
            signature: "Selecione um arquivo de imagem com a assinatura.",
          });
          toast.error("Selecione um arquivo de imagem com a assinatura.");
          return false;
        }

        imageUrl = await uploadFile(file, "signatures");
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

      toast.success("Assinatura cadastrada com sucesso!");
      resetForm();
      props.onOpenChange(false);
      props.onSuccess?.();
      return true;
    } catch {
      toast.error("Não foi possível cadastrar a assinatura.");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  const steps: MultiStepItem<CreateSignatureValues>[] = [
    {
      id: "signatario",
      title: "Signatário",
      description: "Identificação do responsável pela assinatura",
      fields: [
        {
          name: "signatory_name",
          label: "Nome do signatário",
          required: true,
          placeholder: "Ex: Prof. Dr. Carlos Eduardo Mendes",
          hint: "Nome completo como deve constar nos certificados e crachás.",
        },
        {
          name: "signatory_title",
          label: "Cargo / Função",
          required: false,
          placeholder: "Ex: Coordenador Geral, Reitor, Diretor Científico",
          hint: "Título ou atribuição institucional que confere autoridade ao signatário.",
        },
        {
          name: "signatory_email",
          label: "E-mail de contato (opcional)",
          kind: "email",
          required: false,
          placeholder: "carlos.mendes@instituicao.org",
          hint: "E-mail para registro interno e eventuais notificações.",
        },
      ],
    },
    {
      id: "assinatura",
      title: "Assinatura Digital",
      description: "Desenhe diretamente na tela ou importe uma imagem",
      fields: [
        {
          name: "mode",
          label: "Origem da assinatura",
          kind: "custom",
          render: (fieldProps) => (
            <div class="space-y-4">
              <div class="flex items-center gap-2 p-1 bg-muted/50 rounded-xl border border-border/60">
                <button
                  type="button"
                  onClick={() => handleModeChange("draw")}
                  class={cn(
                    "flex-1 flex items-center justify-center gap-2 py-2 px-3 rounded-lg text-xs font-semibold transition-all cursor-pointer",
                    fieldProps.value === "draw"
                      ? "bg-card text-foreground shadow-xs border border-border/80"
                      : "text-muted-foreground hover:text-foreground",
                  )}
                >
                  <PenLine class="size-4 text-primary" />
                  Desenhar na tela
                </button>
                <button
                  type="button"
                  onClick={() => handleModeChange("upload")}
                  class={cn(
                    "flex-1 flex items-center justify-center gap-2 py-2 px-3 rounded-lg text-xs font-semibold transition-all cursor-pointer",
                    fieldProps.value === "upload"
                      ? "bg-card text-foreground shadow-xs border border-border/80"
                      : "text-muted-foreground hover:text-foreground",
                  )}
                >
                  <Upload class="size-4 text-primary" />
                  Importar imagem
                </button>
              </div>

              <Show when={fieldProps.value === "draw"}>
                <div class="space-y-2">
                  <div class="flex items-center justify-between">
                    <label class="text-xs font-medium text-muted-foreground">
                      Quadro de assinatura
                    </label>
                    <button
                      type="button"
                      onClick={() => canvasRef?.clear()}
                      class="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors cursor-pointer px-2 py-1 rounded-md hover:bg-muted"
                    >
                      <Eraser class="size-3.5" />
                      Limpar quadro
                    </button>
                  </div>

                  <div
                    class={cn(
                      "rounded-xl border overflow-hidden transition-colors bg-white",
                      errors().signature
                        ? "border-destructive ring-1 ring-destructive/40"
                        : "border-border/80",
                    )}
                  >
                    <SignatureCanvas
                      ref={(el) => (canvasRef = el)}
                      class="w-full"
                    />
                  </div>

                  <div class="flex items-center justify-between text-[11px] text-muted-foreground px-0.5">
                    <span>
                      Assine no quadro branco com o mouse ou caneta óptica
                    </span>
                    <span class="text-primary font-medium">Fundo transparente</span>
                  </div>

                  <Show when={errors().signature}>
                    <p class="text-xs text-destructive font-medium">
                      {errors().signature}
                    </p>
                  </Show>
                </div>
              </Show>

              <Show when={fieldProps.value === "upload"}>
                <div class="space-y-2">
                  <label
                    for="sig-file-upload"
                    onDragOver={(e) => {
                      e.preventDefault();
                      setIsDragging(true);
                    }}
                    onDragLeave={() => setIsDragging(false)}
                    onDrop={handleDrop}
                    class={cn(
                      "relative flex flex-col items-center justify-center p-6 rounded-xl border-2 border-dashed transition-all cursor-pointer",
                      isDragging()
                        ? "border-primary bg-primary/5 scale-[0.99]"
                        : previewUrl()
                          ? "border-border/70 bg-card hover:border-primary/50"
                          : "border-border/70 bg-muted/20 hover:bg-muted/40 hover:border-border",
                      errors().signature && !previewUrl()
                        ? "border-destructive bg-destructive/5"
                        : "",
                    )}
                  >
                    <input
                      id="sig-file-upload"
                      type="file"
                      accept="image/png,image/jpeg,image/webp,image/svg+xml"
                      onChange={handleFileChange}
                      class="sr-only"
                    />

                    <Show
                      when={previewUrl()}
                      fallback={
                        <div class="flex flex-col items-center gap-2 text-center py-4">
                          <div class="flex size-11 items-center justify-center rounded-xl bg-primary/10 text-primary">
                            <Upload class="size-5" />
                          </div>
                          <div>
                            <p class="text-xs font-semibold text-foreground">
                              Arraste a imagem da assinatura ou clique para buscar
                            </p>
                            <p class="text-[11px] text-muted-foreground mt-0.5">
                              PNG com fundo transparente recomendado (máx. 5MB)
                            </p>
                          </div>
                        </div>
                      }
                    >
                      <div class="flex flex-col items-center gap-3 w-full py-2">
                        <div class="relative max-h-36 max-w-full p-4 rounded-lg bg-checkerboard border border-border/40 flex items-center justify-center overflow-hidden">
                          <img
                            src={previewUrl() ?? ""}
                            alt="Prévia da assinatura"
                            class="max-h-28 object-contain"
                          />
                        </div>
                        <div class="flex items-center gap-1.5 text-xs text-muted-foreground">
                          <CheckCircle2 class="size-3.5 text-emerald-500" />
                          <span class="font-medium text-foreground">
                            {importedFile()?.name}
                          </span>
                          <span>•</span>
                          <span class="text-xs font-medium text-primary hover:underline">
                            Clique para trocar de imagem
                          </span>
                        </div>
                      </div>
                    </Show>
                  </label>

                  <Show when={errors().signature}>
                    <p class="text-xs text-destructive font-medium">
                      {errors().signature}
                    </p>
                  </Show>
                </div>
              </Show>
            </div>
          ),
        },
      ],
    },
    {
      id: "resumo",
      title: "Resumo",
      description: "Revise os dados antes de salvar a assinatura",
      summary: {
        title: "Dados da assinatura",
        badge: () => "Pronto para cadastrar",
        items: [
          {
            label: "Signatário",
            value: () => values().signatory_name,
          },
          {
            label: "Cargo / Função",
            value: () => values().signatory_title || "Não informado",
          },
          {
            label: "E-mail",
            value: () => values().signatory_email || "Não informado",
          },
          {
            label: "Tipo de assinatura",
            value: () =>
              values().mode === "draw"
                ? "Desenho feito na tela"
                : "Arquivo de imagem importado",
          },
        ],
        extra: () => (
          <div class="relative overflow-hidden rounded-xl border border-border/80 bg-card p-4 shadow-xs">
            <p class="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-2">
              Prévia do Registro
            </p>
            <div class="flex items-center justify-between gap-4">
              <div class="space-y-1 min-w-0">
                <h4 class="text-sm font-semibold text-foreground truncate">
                  {values().signatory_name || "Nome do signatário"}
                </h4>
                <p class="text-xs text-muted-foreground">
                  {values().signatory_title || "Cargo ou Função"}
                </p>
                <Show when={values().signatory_email}>
                  <p class="text-xs text-muted-foreground font-mono">
                    {values().signatory_email}
                  </p>
                </Show>
              </div>

              <div class="shrink-0 flex items-center justify-center h-16 w-32 rounded-lg border border-border/60 bg-muted/30 p-2">
                <Show
                  when={
                    values().mode === "draw" ? drawingDataUrl() : previewUrl()
                  }
                  fallback={
                    <span class="text-[11px] text-muted-foreground italic text-center">
                      Sem assinatura
                    </span>
                  }
                >
                  <img
                    src={
                      (values().mode === "draw"
                        ? drawingDataUrl()
                        : previewUrl()) ?? ""
                    }
                    alt="Assinatura"
                    class="max-h-full max-w-full object-contain filter dark:invert"
                  />
                </Show>
              </div>
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
      description="Cadastre uma assinatura para uso automático na emissão de crachás e certificados."
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
