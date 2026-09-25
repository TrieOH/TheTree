import type { JSX } from "@solidjs/web";
import { Show, createEffect, createMemo, createSignal, untrack } from "solid-js";
import CalendarDaysIcon from "~icons/lucide/calendar-days";
import FlagIcon from "~icons/lucide/flag";
import LockIcon from "~icons/lucide/lock";
import ShieldAlertIcon from "~icons/lucide/shield-alert";
import SparklesIcon from "~icons/lucide/sparkles";

import {
  MultiStepDialog,
  type MultiStepItem,
  cn,
} from "@trieoh/ui-solid";
import { formatPrice } from "@/shared/lib/money";
import type { ProgramCreateInput, ProgramI } from "../model";

const CalendarDays = CalendarDaysIcon as unknown as (props: { class?: string }) => JSX.Element;
const Flag = FlagIcon as unknown as (props: { class?: string }) => JSX.Element;
const Lock = LockIcon as unknown as (props: { class?: string }) => JSX.Element;
const ShieldAlert = ShieldAlertIcon as unknown as (props: { class?: string }) => JSX.Element;
const Sparkles = SparklesIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface ManageProgramDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  program?: ProgramI | null;
  onSubmit: (data: ProgramCreateInput) => Promise<boolean>;
}

export interface ManageProgramValues {
  kind: "activity" | "checkpoint";
  name: string;
  description: string;
  banner_url: string;
  price_cents: number;
  min_access_level: number | "";
  staff_only: boolean;
}

const emptyValues: ManageProgramValues = {
  kind: "activity",
  name: "",
  description: "",
  banner_url: "",
  price_cents: 0,
  min_access_level: 0,
  staff_only: false,
};

function valuesOf(program?: ProgramI | null): ManageProgramValues {
  if (!program) {
    return { ...emptyValues };
  }

  return {
    kind: program.kind,
    name: program.name,
    description: program.description ?? "",
    banner_url: program.banner_url ?? "",
    price_cents: typeof program.price === "number" ? program.price : 0,
    min_access_level: program.min_access_level ?? 0,
    staff_only: Boolean(program.staff_only),
  };
}

export function ManageProgramDialog(props: ManageProgramDialogProps): JSX.Element {
  const [currentStep, setCurrentStep] = createSignal(0);
  let currentValues: ManageProgramValues = valuesOf(props.program);
  const [values, setValues] = createSignal<ManageProgramValues>(
    untrack(() => currentValues),
  );
  // Dedicated signal for kind so steps does not recompute when typing into text/number fields
  const [activeKind, setActiveKind] = createSignal<"activity" | "checkpoint">(
    untrack(() => currentValues.kind),
  );
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<
    Partial<Record<keyof ManageProgramValues, string>>
  >({});

  let isMounted = false;
  createEffect(
    () => (props.program ? props.program.id : null),
    () => {
      if (!isMounted) {
        isMounted = true;
        return;
      }
      currentValues = valuesOf(props.program);
      setActiveKind(currentValues.kind);
      setValues(currentValues);
      setErrors({});
      setCurrentStep(0);
    },
  );

  const isEditing = createMemo(() => Boolean(props.program));

  const resetForm = () => {
    currentValues = valuesOf(props.program);
    setActiveKind(currentValues.kind);
    setValues(currentValues);
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
    key: keyof ManageProgramValues & string,
    value: unknown,
  ) => {
    currentValues = { ...currentValues, [key]: value };
    setValues(currentValues);
    setErrors((prev) => {
      if (!prev[key as keyof ManageProgramValues]) return prev;
      const nextErrors = { ...prev };
      delete nextErrors[key as keyof ManageProgramValues];
      return nextErrors;
    });
  };

  const handleKindSelect = (newKind: "activity" | "checkpoint") => {
    setActiveKind(newKind);
    handleChange("kind", newKind);
  };

  const validateStep = (stepIndex: number): boolean => {
    const current = currentValues;
    const nextErrors: Partial<Record<keyof ManageProgramValues, string>> = {};

    if (stepIndex === 0) {
      const name = current.name.trim();
      if (!name || name.length < 2) {
        nextErrors.name = "Nome deve ter pelo menos 2 caracteres.";
      }
    }

    if (stepIndex === 1) {
      if (current.kind === "activity") {
        if (
          typeof current.price_cents !== "number" ||
          isNaN(current.price_cents) ||
          current.price_cents < 0
        ) {
          nextErrors.price_cents = "Informe um valor válido em centavos.";
        }
      }

      const rawLevel = String(current.min_access_level ?? "").trim();
      if (rawLevel === "") {
        nextErrors.min_access_level = "O nível de acesso é obrigatório.";
      } else {
        const accessLevel = Number(rawLevel);
        if (
          isNaN(accessLevel) ||
          accessLevel < 0 ||
          !Number.isInteger(accessLevel)
        ) {
          nextErrors.min_access_level =
            "O nível de acesso deve ser um número inteiro >= 0.";
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
    const priceCents =
      current.kind === "activity"
        ? current.price_cents > 0
          ? current.price_cents
          : undefined
        : undefined;

    const accessLevel =
      current.min_access_level === "" ? 0 : Number(current.min_access_level);

    setSubmitting(true);
    try {
      const ok = await props.onSubmit({
        kind: current.kind,
        name: current.name.trim(),
        description: current.description.trim() || undefined,
        min_access_level: accessLevel,
        staff_only: current.staff_only,
        banner_url: current.banner_url.trim() || null,
        price: priceCents,
      });

      if (ok) {
        resetForm();
        props.onOpenChange(false);
      }
      return ok;
    } finally {
      setSubmitting(false);
    }
  };

  const steps = createMemo<MultiStepItem<ManageProgramValues>[]>(() => {
    const isActivity = activeKind() === "activity";

    return [
      {
        id: "identificacao",
        title: "Identificação",
        description: "Tipo, nome e descrição da atividade ou checkpoint",
        fields: [
          {
            name: "kind",
            kind: "custom",
            label: "Tipo de programação",
            layout: "full",
            render: (fieldProps) => (
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                <button
                  type="button"
                  onClick={() => handleKindSelect("activity")}
                  class={cn(
                    "flex items-start gap-3 rounded-xl border p-3.5 text-left text-xs transition-all cursor-pointer",
                    activeKind() === "activity"
                      ? "border-primary bg-primary/10 text-primary font-semibold ring-1 ring-primary/20 shadow-xs"
                      : "border-border bg-card text-muted-foreground hover:bg-muted/50 hover:text-foreground",
                  )}
                >
                  <div class="flex size-8 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                    <CalendarDays class="size-4" />
                  </div>
                  <div class="space-y-0.5">
                    <span class="font-medium text-foreground block">
                      Atividade / Palestra
                    </span>
                    <span class="text-[11px] text-muted-foreground font-normal leading-relaxed block">
                      Palestras, workshops ou mesas redondas com controle de presença e certificados.
                    </span>
                  </div>
                </button>

                <button
                  type="button"
                  onClick={() => handleKindSelect("checkpoint")}
                  class={cn(
                    "flex items-start gap-3 rounded-xl border p-3.5 text-left text-xs transition-all cursor-pointer",
                    activeKind() === "checkpoint"
                      ? "border-amber-500 bg-amber-500/10 text-amber-700 dark:text-amber-400 font-semibold ring-1 ring-amber-500/20 shadow-xs"
                      : "border-border bg-card text-muted-foreground hover:bg-muted/50 hover:text-foreground",
                  )}
                >
                  <div class="flex size-8 shrink-0 items-center justify-center rounded-lg bg-amber-500/15 text-amber-600 dark:text-amber-400">
                    <Flag class="size-4" />
                  </div>
                  <div class="space-y-0.5">
                    <span class="font-medium text-foreground block">
                      Checkpoint de Presença
                    </span>
                    <span class="text-[11px] text-muted-foreground font-normal leading-relaxed block">
                      Pontos de credenciamento, portarias de acesso ou validação rápida de crachás.
                    </span>
                  </div>
                </button>
              </div>
            ),
          },
          {
            name: "name",
            label: "Nome da programação",
            required: true,
            placeholder: "ex.: Workshop de SolidJS, Credenciamento Principal...",
            hint: "Título que aparecerá na agenda e no cronograma do evento.",
          },
          {
            name: "description",
            kind: "textarea",
            label: "Descrição (opcional)",
            rows: 2,
            placeholder: "Detalhes, palestrantes, tópicos abordados...",
            hint: "Informações adicionais para os participantes.",
          },
        ],
      },
      {
        id: "acesso",
        title: "Acesso e Valores",
        description: "Regras de entrada, público alvo e valores adicionais",
        fields: [
          ...(isActivity
            ? [
                {
                  name: "price_cents" as const,
                  kind: "money" as const,
                  label: "Valor adicional",
                  currency: "BRL",
                  layout: "half" as const,
                  hint: "Deixe R$ 0,00 caso esteja incluso no ingresso.",
                },
              ]
            : []),
          {
            name: "min_access_level",
            kind: "number",
            label: "Nível de acesso mínimo",
            min: "0",
            step: "1",
            layout: isActivity ? "half" : "full",
            required: true,
            placeholder: "0",
            hint: "Nível mínimo do ingresso para acessar (0 = público geral).",
          },
          {
            name: "staff_only",
            kind: "custom",
            label: "Restrição de acesso",
            layout: "full",
            render: (fieldProps) => (
              <div class="w-full pt-1">
                <label class="flex items-center gap-3 rounded-xl border border-border/70 bg-card p-3.5 transition-colors hover:bg-muted/40 cursor-pointer select-none">
                  <input
                    type="checkbox"
                    id={fieldProps.ids?.id}
                    checked={Boolean(fieldProps.value)}
                    onChange={(e) => fieldProps.onChange(e.currentTarget.checked)}
                    class="size-4 rounded border-border text-primary focus:ring-primary"
                  />
                  <div class="space-y-0.5">
                    <span class="text-xs font-semibold text-foreground block">
                      Apenas para equipe (Staff only)
                    </span>
                    <span class="text-[11px] text-muted-foreground block">
                      Visível e gerenciável somente por membros autorizados da equipe
                    </span>
                  </div>
                </label>
              </div>
            ),
          },
        ],
      },
      {
        id: "resumo",
        title: "Resumo",
        description: "Revise os dados antes de salvar a programação",
        summary: {
          title: "Dados da programação",
          badge: () =>
            isEditing() ? "Pronto para salvar alterações" : "Pronto para cadastrar",
          items: [
            {
              label: "Tipo",
              value: () =>
                values().kind === "activity"
                  ? "Atividade / Palestra"
                  : "Checkpoint de Presença",
            },
            {
              label: "Nome",
              value: () => values().name || "Sem nome",
            },
            {
              label: "Descrição",
              value: () => values().description || "Sem descrição",
            },
            {
              label: "Valor adicional",
              value: () =>
                values().kind === "activity"
                  ? values().price_cents > 0
                    ? formatPrice(values().price_cents)
                    : "Gratuito / Incluso no evento"
                  : "Não aplicável (Checkpoint)",
            },
            {
              label: "Nível de acesso",
              value: () =>
                values().min_access_level === 0
                  ? "Nível 0 (Livre para todos os ingressos)"
                  : `Nível ${values().min_access_level}`,
            },
            {
              label: "Restrição",
              value: () =>
                values().staff_only
                  ? "Exclusivo para equipe (Staff only)"
                  : "Aberto aos participantes",
            },
          ],
          extra: () => (
            <div class="relative overflow-hidden rounded-xl border border-border/80 bg-card p-5 shadow-xs">
              <div class="flex items-center justify-between pb-3 border-b border-border/60">
                <span class="text-xs font-medium text-muted-foreground uppercase tracking-wider flex items-center gap-1.5">
                  <Sparkles class="size-3.5" />
                  Prévia da Programação
                </span>
                <span
                  class={cn(
                    "inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold",
                    values().kind === "activity"
                      ? "bg-primary/10 text-primary"
                      : "bg-amber-500/10 text-amber-700 dark:text-amber-400",
                  )}
                >
                  <Show
                    when={values().kind === "activity"}
                    fallback={<Flag class="size-3.5" />}
                  >
                    <CalendarDays class="size-3.5" />
                  </Show>
                  <span>
                    {values().kind === "activity" ? "Atividade" : "Checkpoint"}
                  </span>
                </span>
              </div>

              <div class="mt-4 flex flex-col gap-3 rounded-lg border border-border/60 bg-muted/20 p-4">
                <Show when={values().banner_url}>
                  <div class="relative max-h-32 w-full overflow-hidden rounded-lg bg-card">
                    <img
                      src={values().banner_url}
                      alt={values().name}
                      class="h-28 w-full object-cover"
                      onError={(e) => {
                        (e.currentTarget as HTMLElement).style.display = "none";
                      }}
                    />
                  </div>
                </Show>

                <div>
                  <h4 class="text-sm font-semibold text-foreground">
                    {values().name || "Nome da Programação"}
                  </h4>
                  <Show when={values().description}>
                    <p class="mt-1 text-xs text-muted-foreground line-clamp-2 leading-relaxed">
                      {values().description}
                    </p>
                  </Show>
                </div>

                <div class="flex flex-wrap items-center gap-2 pt-1 border-t border-border/40 text-[11px]">
                  <span class="inline-flex items-center gap-1 font-medium text-muted-foreground">
                    <Lock class="size-3" />
                    Nível {values().min_access_level === "" ? 0 : values().min_access_level}
                  </span>

                  <Show when={values().kind === "activity"}>
                    <span class="text-muted-foreground">•</span>
                    <span class="font-medium text-foreground">
                      {values().price_cents > 0
                        ? formatPrice(values().price_cents)
                        : "Incluso no ingresso"}
                    </span>
                  </Show>

                  <Show when={values().staff_only}>
                    <span class="text-muted-foreground">•</span>
                    <span class="inline-flex items-center gap-1 text-amber-600 dark:text-amber-400 font-semibold">
                      <ShieldAlert class="size-3" />
                      Staff Only
                    </span>
                  </Show>
                </div>
              </div>
            </div>
          ),
        },
      },
    ];
  });

  return (
    <MultiStepDialog
      open={props.open}
      onOpenChange={handleOpenChange}
      title={isEditing() ? "Editar programação" : "Nova programação"}
      description="Configure os detalhes da atividade ou checkpoint na programação do evento."
      steps={steps()}
      currentStep={currentStep()}
      onStepChange={setCurrentStep}
      values={values()}
      errors={errors()}
      onChange={handleChange}
      loading={submitting()}
      submitLabel={
        submitting()
          ? "Salvando..."
          : isEditing()
            ? "Salvar alterações"
            : "Criar programação"
      }
      onBeforeNext={validateStep}
      onSubmit={handleSubmit}
    />
  );
}
