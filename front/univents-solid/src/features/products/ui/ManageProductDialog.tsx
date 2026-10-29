import type { JSX } from "@solidjs/web";
import { Show, createEffect, createMemo, createSignal, untrack } from "solid-js";
import ShieldCheckIcon from "~icons/lucide/shield-check";

import {
  MultiStepDialog,
  type MultiStepItem,
  cn,
} from "@trieoh/ui-solid";
import { formatPrice } from "@/shared/lib/money";
import { getFieldErrors } from "@/shared/lib/effect-schema";
import {
  productStep1Schema,
  productVariantStepSchema,
  type CreateInitialProductOutputI,
  type ProductI,
  type ProductPatchOutputI,
} from "../model";

const ShieldCheck = ShieldCheckIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface ManageProductValues {
  vendor_code: string;
  requires_registration: boolean;
  // Campos da primeira variação (ao criar)
  variant_vendor_code: string;
  variant_name: string;
  variant_description: string;
  variant_price_cents: number;
  variant_stock: string;
}

export interface ManageProductDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  product: ProductI | null;
  onCreate?: (values: CreateInitialProductOutputI) => Promise<boolean>;
  onUpdate?: (values: ProductPatchOutputI) => Promise<boolean>;
}

const emptyValues: ManageProductValues = {
  vendor_code: "",
  requires_registration: false,
  variant_vendor_code: "",
  variant_name: "",
  variant_description: "",
  variant_price_cents: 0,
  variant_stock: "",
};

const valuesOf = (product: ProductI | null): ManageProductValues =>
  product
    ? {
        ...emptyValues,
        vendor_code: product.vendor_code,
        requires_registration: product.requires_registration,
      }
    : { ...emptyValues };

export function ManageProductDialog(props: ManageProductDialogProps) {
  const [currentStep, setCurrentStep] = createSignal(0);
  let currentValues: ManageProductValues = untrack(() => valuesOf(props.product));
  const [values, setValues] = createSignal<ManageProductValues>(
    untrack(() => currentValues),
  );
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<
    Partial<Record<keyof ManageProductValues, string>>
  >({});

  let isMounted = false;
  createEffect(
    () => (props.product ? props.product.id : null),
    () => {
      if (!isMounted) {
        isMounted = true;
        return;
      }
      currentValues = valuesOf(props.product);
      setValues(currentValues);
      setErrors({});
      setCurrentStep(0);
    },
  );

  const isEditing = () => Boolean(props.product);

  const handleChange = (
    key: keyof ManageProductValues & string,
    value: unknown,
  ) => {
    currentValues = { ...currentValues, [key]: value };
    setValues(currentValues);
    setErrors((prev) => {
      if (!prev[key as keyof ManageProductValues]) return prev;
      const nextErrors = { ...prev };
      delete nextErrors[key as keyof ManageProductValues];
      return nextErrors;
    });
  };

  const validateStep = (stepIndex: number): boolean => {
    const current = currentValues;
    let nextErrors: Partial<Record<keyof ManageProductValues, string>> = {};

    if (stepIndex === 0) {
      nextErrors = getFieldErrors<ManageProductValues>(
        productStep1Schema.safeParse(current),
      );
    } else if (!isEditing() && stepIndex === 1) {
      nextErrors = getFieldErrors<ManageProductValues>(
        productVariantStepSchema.safeParse(current),
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
    setSubmitting(true);
    try {
      if (isEditing()) {
        const ok = await props.onUpdate?.({
          vendor_code: current.vendor_code.trim(),
          requires_registration: current.requires_registration,
        });
        if (ok) {
          props.onOpenChange(false);
          return true;
        }
        return false;
      }

      const stockRaw = current.variant_stock.trim();
      const stock = stockRaw === "" ? null : parseInt(stockRaw, 10);

      const ok = await props.onCreate?.({
        vendor_code: current.vendor_code.trim(),
        requires_registration: current.requires_registration,
        variant_vendor_code: current.variant_vendor_code.trim(),
        name: current.variant_name.trim(),
        description: current.variant_description.trim() || null,
        price: current.variant_price_cents,
        stock,
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

  const steps = createMemo<MultiStepItem<ManageProductValues>[]>(() => {
    const baseStep: MultiStepItem<ManageProductValues> = {
      id: "identificacao",
      title: "Produto",
      description: "Identificação do produto no catálogo",
      fields: [
        {
          name: "vendor_code",
          label: "Código do produto (SKU)",
          required: true,
          placeholder: "Ex: PROD-CAMISA-01",
          hint: "Identificador único do produto no seu inventário.",
        },
        {
          name: "requires_registration",
          label: "Exigir cadastro do comprador",
          kind: "custom",
          render: (ctx) => {
            const checked = () => Boolean(ctx.value);
            return (
              <div
                class={cn(
                  "flex items-start gap-3 rounded-xl border p-3.5 transition-colors cursor-pointer select-none",
                  checked()
                    ? "border-primary/50 bg-primary/5"
                    : "border-border hover:bg-muted/30",
                )}
                onClick={() => ctx.onChange(!checked())}
              >
                <div class="mt-0.5 shrink-0 text-muted-foreground">
                  <ShieldCheck class="size-5" />
                </div>
                <div class="flex-1 space-y-0.5">
                  <label
                    for={ctx.ids.id}
                    class="text-xs font-semibold text-foreground cursor-pointer block"
                  >
                    Exigir cadastro do comprador
                  </label>
                  <p class="text-xs text-muted-foreground">
                    Quando ativado, o comprador precisará preencher os dados de cadastro ao comprar este produto.
                  </p>
                </div>
                <input
                  id={ctx.ids.id}
                  type="checkbox"
                  checked={checked()}
                  onChange={(e) => ctx.onChange(e.currentTarget.checked)}
                  class="sr-only"
                />
                <div
                  class={cn(
                    "size-4 rounded border transition-colors flex items-center justify-center mt-0.5 shrink-0",
                    checked()
                      ? "border-primary bg-primary text-primary-foreground"
                      : "border-input bg-background",
                  )}
                  aria-hidden="true"
                >
                  <Show when={checked()}>
                    <svg
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="3"
                      class="size-3"
                    >
                      <path d="M20 6 9 17l-5-5" stroke-linecap="round" stroke-linejoin="round" />
                    </svg>
                  </Show>
                </div>
              </div>
            );
          },
        },
      ],
    };

    if (isEditing()) {
      return [
        baseStep,
        {
          id: "resumo",
          title: "Resumo",
          description: "Revise os dados antes de salvar o produto",
          summary: {
            title: "Dados do produto",
            badge: () => "Pronto para salvar",
            items: [
              {
                label: "Código (SKU)",
                value: () => values().vendor_code,
                mono: true,
              },
              {
                label: "Cadastro exigido",
                value: () =>
                  values().requires_registration ? "Sim" : "Não",
              },
            ],
          },
        },
      ];
    }

    return [
      baseStep,
      {
        id: "variacao",
        title: "Primeira Variação",
        description: "Configure a variação inicial deste produto",
        fields: [
          {
            name: "variant_vendor_code",
            label: "Código da variação (SKU)",
            required: true,
            layout: "half",
            placeholder: "Ex: CAMISA-AZUL-G",
            hint: "Identificador único desta opção.",
          },
          {
            name: "variant_name",
            label: "Nome da variação",
            required: true,
            layout: "half",
            placeholder: "Ex: Azul / Tamanho G",
          },
          {
            name: "variant_price_cents",
            label: "Preço",
            kind: "money",
            currency: "BRL",
            locale: "pt-BR",
            layout: "half",
            required: true,
            hint: "R$ 0,00 para item gratuito.",
          },
          {
            name: "variant_stock",
            label: "Estoque inicial",
            kind: "number",
            step: "1",
            min: "1",
            layout: "half",
            placeholder: "Ilimitado",
            hint: "Deixe vazio para estoque ilimitado.",
          },
          {
            name: "variant_description",
            label: "Descrição da variação",
            kind: "textarea",
            rows: 2,
            placeholder: "Detalhes adicionais sobre esta opção específica (opcional)...",
          },
        ],
      },
      {
        id: "resumo",
        title: "Resumo",
        description: "Revise os dados antes de salvar o produto",
        summary: {
          title: "Dados do produto e variação",
          badge: () => "Pronto para cadastrar",
          items: [
            {
              label: "Código do produto",
              value: () => values().vendor_code,
              mono: true,
            },
            {
              label: "Cadastro exigido",
              value: () =>
                values().requires_registration ? "Sim" : "Não",
            },
            {
              label: "Variação inicial",
              value: () => values().variant_name,
            },
            {
              label: "SKU da variação",
              value: () => values().variant_vendor_code,
              mono: true,
            },
            {
              label: "Preço",
              value: () =>
                formatPrice(values().variant_price_cents ?? 0),
            },
            {
              label: "Estoque",
              value: () =>
                values().variant_stock
                  ? `${values().variant_stock} unidades`
                  : "Ilimitado",
            },
          ],
        },
      },
    ];
  });

  return (
    <MultiStepDialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={isEditing() ? "Editar produto" : "Novo produto"}
      description={
        isEditing()
          ? "Atualize as informações do produto."
          : "Cadastre um produto com sua primeira opção de venda."
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
            : "Criar produto"
      }
      onBeforeNext={validateStep}
      onSubmit={handleSubmit}
    />
  );
}
