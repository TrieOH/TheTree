import type { JSX } from "@solidjs/web";
import { Show, createEffect, createMemo, createSignal, untrack } from "solid-js";
import BoxesIcon from "~icons/lucide/boxes";

import {
  MultiStepDialog,
  type MultiStepItem,
} from "@trieoh/ui-solid";
import { formatPrice } from "@/shared/lib/money";
import { getFieldErrors } from "@/shared/lib/effect-schema";
import {
  variantStep1Schema,
  variantStep2Schema,
  type VariantCreateOutputI,
  type VariantI,
} from "../model";

const Boxes = BoxesIcon as unknown as (props: { class?: string }) => JSX.Element;

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

export function ManageVariantDialog(props: ManageVariantDialogProps) {
  const [currentStep, setCurrentStep] = createSignal(0);
  let currentValues: ManageVariantValues = untrack(() => valuesOf(props.variant));
  const [values, setValues] = createSignal<ManageVariantValues>(
    untrack(() => currentValues),
  );
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<
    Partial<Record<keyof ManageVariantValues, string>>
  >({});

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

  const isEditing = () => Boolean(props.variant);

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
    let nextErrors: Partial<Record<keyof ManageVariantValues, string>> = {};

    if (stepIndex === 0) {
      nextErrors = getFieldErrors<ManageVariantValues>(
        variantStep1Schema.safeParse(current),
      );
    } else if (stepIndex === 1) {
      nextErrors = getFieldErrors<ManageVariantValues>(
        variantStep2Schema.safeParse(current),
      );
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
        description: current.description?.trim() || null,
        price: current.price_cents,
        stock: stockQty,
        gallery_urls: [],
      });
      if (ok) {
        props.onOpenChange(false);
        return true;
      }
      return false;
    } finally {
      setSubmitting(false);
    }
  };

  const steps = createMemo<MultiStepItem<ManageVariantValues>[]>(() => [
    {
      id: "identificacao",
      title: "Variação",
      description: "Identificação e nome desta opção",
      fields: [
        {
          name: "vendor_code",
          label: "Código da variação (SKU)",
          required: true,
          placeholder: "Ex: CAMISA-AZUL-G",
          hint: "Identificador único desta opção no seu inventário.",
        },
        {
          name: "name",
          label: "Nome da variação",
          required: true,
          placeholder: "Ex: Azul / Tamanho G, Edição Especial",
          hint: "Nome visível para os compradores.",
        },
        {
          name: "description",
          label: "Descrição (opcional)",
          kind: "textarea",
          rows: 3,
          placeholder: "Detalhes adicionais sobre esta opção específica (opcional)...",
          hint: "Informações extras que diferenciam esta opção das demais.",
        },
      ],
    },
    {
      id: "valores",
      title: "Preço e Estoque",
      description: "Valor de venda e controle de estoque",
      fields: [
        {
          name: "price_cents",
          label: "Preço unitário",
          kind: "money",
          currency: "BRL",
          locale: "pt-BR",
          layout: "half",
          required: true,
          hint: "R$ 0,00 para item gratuito.",
        },
        {
          name: "stock",
          label: "Estoque disponível",
          kind: "number",
          step: "1",
          min: "1",
          layout: "half",
          placeholder: "Ilimitado",
          hint: "Deixe vazio para estoque ilimitado.",
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
          props.variant ? "Pronto para salvar alterações" : "Pronto para cadastrar",
        items: [
          {
            label: "Nome",
            value: () => values().name,
          },
          {
            label: "Código (SKU)",
            value: () => values().vendor_code,
            mono: true,
          },
          {
            label: "Preço",
            value: () => formatPrice(values().price_cents ?? 0),
          },
          {
            label: "Estoque",
            value: () =>
              values().stock ? `${values().stock} unidades` : "Ilimitado",
          },
          {
            label: "Descrição",
            value: () => values().description || "Sem descrição",
            fullWidth: true,
          },
        ],
        extra: () => (
          <div class="relative overflow-hidden rounded-xl border border-border/80 bg-card p-4 shadow-xs">
            <div class="flex items-start justify-between gap-3">
              <div class="space-y-1 min-w-0">
                <p class="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                  Prévia da variação
                </p>
                <h4 class="text-base font-semibold text-foreground truncate">
                  {values().name || "Nome da variação"}
                </h4>
                <p class="text-xs font-mono text-muted-foreground">
                  {values().vendor_code || "SKU-EXEMPLO"}
                </p>
              </div>
              <span class="inline-flex shrink-0 items-center rounded-full px-2.5 py-0.5 text-xs font-semibold bg-primary/10 text-primary">
                {formatPrice(values().price_cents ?? 0)}
              </span>
            </div>

            <Show when={values().description}>
              <p class="mt-2.5 text-xs text-muted-foreground line-clamp-2">
                {values().description}
              </p>
            </Show>

            <div class="mt-4 flex flex-wrap items-center gap-2 pt-3 border-t border-border/60 text-xs text-muted-foreground">
              <span class="inline-flex items-center gap-1 rounded-md bg-muted px-2 py-0.5 font-medium">
                <Boxes class="size-3" />
                {values().stock
                  ? `${values().stock} disponíveis`
                  : "Estoque ilimitado"}
              </span>
            </div>
          </div>
        ),
      },
    },
  ]);

  return (
    <MultiStepDialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={isEditing() ? "Editar variação" : "Nova variação"}
      description={
        isEditing()
          ? "Atualize as informações desta variação."
          : "Cadastre uma nova opção ou combinação deste produto."
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
