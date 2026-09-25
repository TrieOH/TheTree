import type { JSX } from "@solidjs/web";
import { Show, createEffect, createMemo, createSignal, untrack } from "solid-js";
import BoxesIcon from "~icons/lucide/boxes";
import CoinsIcon from "~icons/lucide/coins";
import InfinityIcon from "~icons/lucide/infinity";
import LayersIcon from "~icons/lucide/layers";

import {
  MultiStepDialog,
  type MultiStepItem,
} from "@trieoh/ui-solid";
import { formatPrice } from "@/shared/lib/money";
import type { VariantCreateOutputI, VariantI } from "../model";

const Boxes = BoxesIcon as unknown as (props: { class?: string }) => JSX.Element;
const Coins = CoinsIcon as unknown as (props: { class?: string }) => JSX.Element;
const InfinityLucide = InfinityIcon as unknown as (props: { class?: string }) => JSX.Element;
const Layers = LayersIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface ManageVariantValues {
  vendor_code: string;
  name: string;
  description: string;
  price_cents: number;
  stock: string;
}

export interface ManageVariantDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  variant: VariantI | null;
  onSubmit: (values: VariantCreateOutputI) => Promise<boolean>;
}

const emptyValues: ManageVariantValues = {
  vendor_code: "",
  name: "",
  description: "",
  price_cents: 0,
  stock: "",
};

const valuesOf = (variant: VariantI | null): ManageVariantValues =>
  variant
    ? {
        vendor_code: variant.vendor_code,
        name: variant.name,
        description: variant.description ?? "",
        price_cents: variant.price,
        stock: variant.stock != null ? String(variant.stock) : "",
      }
    : { ...emptyValues };

export function ManageVariantDialog(props: ManageVariantDialogProps): JSX.Element {
  const [currentStep, setCurrentStep] = createSignal(0);
  let currentValues: ManageVariantValues = untrack(() => valuesOf(props.variant));
  const [values, setValues] = createSignal<ManageVariantValues>(
    untrack(() => currentValues),
  );
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<
    Partial<Record<keyof ManageVariantValues, string>>
  >({});

  const isEditing = createMemo(() => Boolean(props.variant));

  let isMounted = false;
  createEffect(
    () => (props.variant ? props.variant.id : null),
    () => {
      if (!isMounted) {
        isMounted = true;
        return;
      }
      currentValues = valuesOf(props.variant);
      setValues(currentValues);
      setErrors({});
      setCurrentStep(0);
    },
  );

  const resetForm = () => {
    currentValues = valuesOf(props.variant);
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
    key: keyof ManageVariantValues & string,
    value: unknown,
  ) => {
    currentValues = { ...currentValues, [key]: value };
    setValues(currentValues);
    setErrors((prev) => {
      if (!prev[key as keyof ManageVariantValues]) return prev;
      const nextErrors = { ...prev };
      delete nextErrors[key as keyof ManageVariantValues];
      return nextErrors;
    });
  };

  const validateStep = (stepIndex: number): boolean => {
    const current = currentValues;
    const nextErrors: Partial<Record<keyof ManageVariantValues, string>> = {};

    if (stepIndex === 0) {
      if (current.vendor_code.trim().length < 2) {
        nextErrors.vendor_code = "O código deve ter pelo menos 2 caracteres.";
      }

      if (current.name.trim().length < 2) {
        nextErrors.name = "O nome deve ter pelo menos 2 caracteres.";
      }
    }

    if (stepIndex === 1) {
      if (
        typeof current.price_cents !== "number" ||
        isNaN(current.price_cents) ||
        current.price_cents < 0
      ) {
        nextErrors.price_cents = "Informe um preço válido (>= 0).";
      }

      if (current.stock.trim() !== "") {
        const qty = Number(current.stock);
        if (!Number.isInteger(qty) || qty <= 0) {
          nextErrors.stock = "Quantidade em estoque deve ser um número inteiro maior que 0.";
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
    const stockQty =
      current.stock.trim() !== "" ? parseInt(current.stock, 10) : null;

    setSubmitting(true);
    try {
      const ok = await props.onSubmit({
        vendor_code: current.vendor_code.trim(),
        name: current.name.trim(),
        description: current.description.trim() || null,
        price: current.price_cents,
        stock: stockQty,
        gallery_urls: props.variant?.gallery_urls ?? [],
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

  const steps = createMemo<MultiStepItem<ManageVariantValues>[]>(() => [
    {
      id: "identificacao",
      title: "Identificação",
      description: "Código SKU e nome da variação",
      fields: [
        {
          name: "vendor_code",
          label: "Código da variação (SKU)",
          required: true,
          layout: "half",
          placeholder: "Ex: CAMISETA-G",
          hint: "Identificador único (SKU) desta opção.",
        },
        {
          name: "name",
          label: "Nome da variação",
          required: true,
          layout: "half",
          placeholder: "Ex: Tamanho G - Branca",
          hint: "Como o comprador verá este item.",
        },
        {
          name: "description",
          kind: "textarea",
          rows: 2,
          label: "Descrição (opcional)",
          layout: "full",
          placeholder: "Detalhes sobre tecido, medidas, cor ou especificações...",
          hint: "Informações adicionais para os participantes.",
        },
      ],
    },
    {
      id: "valores",
      title: "Preço e Estoque",
      description: "Defina o valor de venda e a disponibilidade",
      fields: [
        {
          name: "price_cents",
          kind: "money",
          currency: "BRL",
          label: "Preço unitário",
          layout: "half",
          hint: "Valor unitário de venda deste modelo.",
        },
        {
          name: "stock",
          kind: "number",
          min: "1",
          step: "1",
          label: "Estoque disponível",
          layout: "half",
          placeholder: "Ilimitado",
          hint: "Deixe em branco para estoque ilimitado.",
        },
      ],
    },
    {
      id: "resumo",
      title: "Resumo",
      description: "Revise os dados antes de salvar a variação",
      summary: {
        title: "Dados da variação",
        badge: () =>
          isEditing() ? "Pronto para salvar alterações" : "Pronto para cadastrar",
        items: [
          {
            label: "Nome",
            value: () => values().name || "Sem nome",
          },
          {
            label: "Código SKU",
            value: () => values().vendor_code || "Sem código",
          },
          {
            label: "Preço unitário",
            value: () => formatPrice(values().price_cents),
          },
          {
            label: "Estoque",
            value: () =>
              values().stock.trim() !== ""
                ? `${values().stock} unidades`
                : "Ilimitado",
          },
          {
            label: "Descrição",
            value: () => values().description || "Sem descrição",
          },
        ],
        extra: () => (
          <div class="relative overflow-hidden rounded-xl border border-border/80 bg-card p-4 shadow-xs space-y-3">
            <div class="flex items-center justify-between pb-2 border-b border-border/60">
              <span class="text-xs font-semibold text-foreground flex items-center gap-1.5">
                <Layers class="size-4 text-primary" />
                {values().name || "Nova Variação"}
              </span>
              <span class="font-mono text-xs text-muted-foreground">
                {values().vendor_code || "SKU-VAR"}
              </span>
            </div>

            <div class="flex items-center justify-between rounded-lg border border-border/60 bg-muted/20 p-3">
              <div class="flex items-center gap-2">
                <Coins class="size-4 text-emerald-600 dark:text-emerald-400" />
                <span class="text-sm font-bold text-foreground">
                  {formatPrice(values().price_cents)}
                </span>
              </div>

              <div class="flex items-center gap-1.5 text-xs">
                <Show
                  when={values().stock.trim() !== ""}
                  fallback={
                    <span class="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 text-[11px] font-medium text-emerald-600 dark:text-emerald-400">
                      <InfinityLucide class="size-3" />
                      Estoque ilimitado
                    </span>
                  }
                >
                  <span class="inline-flex items-center gap-1 rounded-full bg-primary/10 px-2 py-0.5 text-[11px] font-medium text-primary">
                    <Boxes class="size-3" />
                    {values().stock} un.
                  </span>
                </Show>
              </div>
            </div>

            <Show when={values().description}>
              <p class="text-xs text-muted-foreground leading-relaxed">
                {values().description}
              </p>
            </Show>
          </div>
        ),
      },
    },
  ]);

  return (
    <MultiStepDialog
      open={props.open}
      onOpenChange={handleOpenChange}
      title={isEditing() ? "Editar variação" : "Nova variação"}
      description={
        isEditing()
          ? "Atualize as informações de preço e estoque da variação."
          : "Cadastre um novo tamanho, cor ou modelo para este produto."
      }
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
            : "Criar variação"
      }
      onBeforeNext={validateStep}
      onSubmit={handleSubmit}
    />
  );
}
