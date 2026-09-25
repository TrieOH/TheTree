import type { JSX } from "@solidjs/web";
import { Show, createEffect, createMemo, createSignal, untrack } from "solid-js";
import PackageIcon from "~icons/lucide/package";
import ShieldCheckIcon from "~icons/lucide/shield-check";

import {
  MultiStepDialog,
  type MultiStepItem,
  cn,
} from "@trieoh/ui-solid";
import { formatPrice } from "@/shared/lib/money";
import type {
  CreateInitialProductOutputI,
  ProductI,
  ProductPatchOutputI,
} from "../model";

const Package = PackageIcon as unknown as (props: { class?: string }) => JSX.Element;
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

export function ManageProductDialog(props: ManageProductDialogProps): JSX.Element {
  const [currentStep, setCurrentStep] = createSignal(0);
  let currentValues: ManageProductValues = untrack(() => valuesOf(props.product));
  const [values, setValues] = createSignal<ManageProductValues>(
    untrack(() => currentValues),
  );
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<
    Partial<Record<keyof ManageProductValues, string>>
  >({});

  const isEditing = createMemo(() => Boolean(props.product));

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

  const resetForm = () => {
    currentValues = valuesOf(props.product);
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
    const nextErrors: Partial<Record<keyof ManageProductValues, string>> = {};

    if (stepIndex === 0) {
      if (current.vendor_code.trim().length < 2) {
        nextErrors.vendor_code = "O código do produto deve ter pelo menos 2 caracteres.";
      }
    }

    if (!isEditing() && stepIndex === 1) {
      if (current.variant_vendor_code.trim().length < 2) {
        nextErrors.variant_vendor_code = "O código da variação deve ter pelo menos 2 caracteres.";
      }

      if (current.variant_name.trim().length < 2) {
        nextErrors.variant_name = "O nome da variação deve ter pelo menos 2 caracteres.";
      }

      if (
        typeof current.variant_price_cents !== "number" ||
        isNaN(current.variant_price_cents) ||
        current.variant_price_cents < 0
      ) {
        nextErrors.variant_price_cents = "Informe um preço válido (>= 0).";
      }

      if (current.variant_stock.trim() !== "") {
        const qty = Number(current.variant_stock);
        if (!Number.isInteger(qty) || qty <= 0) {
          nextErrors.variant_stock = "Quantidade em estoque deve ser um número inteiro maior que 0.";
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
    setSubmitting(true);
    try {
      if (isEditing()) {
        const ok = await props.onUpdate?.({
          vendor_code: current.vendor_code.trim(),
          requires_registration: current.requires_registration,
        });
        if (ok) {
          resetForm();
          props.onOpenChange(false);
        }
        return Boolean(ok);
      }

      const stockQty =
        current.variant_stock.trim() !== ""
          ? parseInt(current.variant_stock, 10)
          : null;

      const ok = await props.onCreate?.({
        vendor_code: current.vendor_code.trim(),
        requires_registration: current.requires_registration,
        variant_vendor_code: current.variant_vendor_code.trim(),
        name: current.variant_name.trim(),
        description: current.variant_description.trim() || null,
        price: current.variant_price_cents,
        stock: stockQty,
      });

      if (ok) {
        resetForm();
        props.onOpenChange(false);
      }
      return Boolean(ok);
    } finally {
      setSubmitting(false);
    }
  };

  const steps = createMemo<MultiStepItem<ManageProductValues>[]>(() => {
    if (isEditing()) {
      return [
        {
          id: "produto",
          title: "Identificação",
          description: "Código SKU e configurações do produto",
          fields: [
            {
              name: "vendor_code",
              label: "Código do produto (SKU geral)",
              required: true,
              placeholder: "Ex: CAMISETA-OFICIAL",
              hint: "Identificador único do produto no evento.",
            },
            {
              name: "requires_registration",
              kind: "custom",
              label: "Requisitos de compra",
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
                        Exigir cadastro do comprador
                      </span>
                      <span class="text-[11px] text-muted-foreground block">
                        Solicita dados de participante e formulário no momento da compra
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
          description: "Revise os dados antes de salvar o produto",
          summary: {
            title: "Dados do produto",
            badge: () => "Pronto para salvar",
            items: [
              {
                label: "Código do produto (SKU)",
                value: () => values().vendor_code || "Sem código",
              },
              {
                label: "Exige cadastro",
                value: () =>
                  values().requires_registration
                    ? "Sim (Exige formulário de participante)"
                    : "Não (Venda direta)",
              },
            ],
            extra: () => (
              <div class="relative overflow-hidden rounded-xl border border-border/80 bg-card p-4 shadow-xs space-y-3">
                <div class="flex items-center justify-between pb-2 border-b border-border/60">
                  <span class="text-xs font-semibold text-foreground flex items-center gap-1.5">
                    <Package class="size-4 text-primary" />
                    {values().vendor_code || "SKU-PRODUTO"}
                  </span>
                  <span
                    class={cn(
                      "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[11px] font-medium",
                      values().requires_registration
                        ? "bg-primary/10 text-primary"
                        : "bg-muted text-muted-foreground",
                    )}
                  >
                    <ShieldCheck class="size-3" />
                    {values().requires_registration ? "Exige cadastro" : "Livre"}
                  </span>
                </div>
                <p class="text-xs text-muted-foreground">
                  As variações e fotos deste produto podem ser gerenciadas na tela de variações.
                </p>
              </div>
            ),
          },
        },
      ];
    }

    return [
      {
        id: "produto",
        title: "Identificação",
        description: "Código geral do produto no catálogo",
        fields: [
          {
            name: "vendor_code",
            label: "Código do produto (SKU geral)",
            required: true,
            placeholder: "Ex: CAMISETA-OFICIAL",
            hint: "Identificador único do produto no evento.",
          },
          {
            name: "requires_registration",
            kind: "custom",
            label: "Requisitos de compra",
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
                      Exigir cadastro do comprador
                    </span>
                    <span class="text-[11px] text-muted-foreground block">
                      Solicita dados de participante e formulário no momento da compra
                    </span>
                  </div>
                </label>
              </div>
            ),
          },
        ],
      },
      {
        id: "variacao",
        title: "Primeira Variação",
        description: "Defina o modelo, preço e estoque inicial",
        fields: [
          {
            name: "variant_vendor_code",
            label: "Código da variação (SKU)",
            required: true,
            layout: "half",
            placeholder: "Ex: CAM-PRETA-M",
            hint: "SKU específico deste modelo.",
          },
          {
            name: "variant_name",
            label: "Nome da variação",
            required: true,
            layout: "half",
            placeholder: "Ex: Tamanho M - Preta",
            hint: "Como o comprador verá este item.",
          },
          {
            name: "variant_price_cents",
            kind: "money",
            currency: "BRL",
            label: "Preço",
            layout: "half",
            hint: "Valor unitário de venda.",
          },
          {
            name: "variant_stock",
            kind: "number",
            min: "1",
            step: "1",
            label: "Estoque inicial",
            layout: "half",
            placeholder: "Ilimitado",
            hint: "Deixe em branco para estoque ilimitado.",
          },
          {
            name: "variant_description",
            kind: "textarea",
            rows: 2,
            label: "Descrição da variação (opcional)",
            layout: "full",
            placeholder: "Detalhes sobre tecido, medidas, cor ou especificações...",
            hint: "Informações adicionais para os participantes.",
          },
        ],
      },
      {
        id: "resumo",
        title: "Resumo",
        description: "Revise os dados antes de cadastrar o produto",
        summary: {
          title: "Dados do produto e variação",
          badge: () => "Pronto para cadastrar",
          items: [
            {
              label: "Código do produto",
              value: () => values().vendor_code || "Sem código",
            },
            {
              label: "Exige cadastro",
              value: () =>
                values().requires_registration ? "Sim" : "Não",
            },
            {
              label: "Primeira variação",
              value: () => values().variant_name || "Sem nome",
            },
            {
              label: "SKU da variação",
              value: () => values().variant_vendor_code || "Sem SKU",
            },
            {
              label: "Preço unitário",
              value: () => formatPrice(values().variant_price_cents),
            },
            {
              label: "Estoque",
              value: () =>
                values().variant_stock.trim() !== ""
                  ? `${values().variant_stock} unidades`
                  : "Ilimitado",
            },
          ],
          extra: () => (
            <div class="relative overflow-hidden rounded-xl border border-border/80 bg-card p-4 shadow-xs space-y-3">
              <div class="flex items-center justify-between pb-2 border-b border-border/60">
                <span class="text-xs font-semibold text-foreground flex items-center gap-1.5">
                  <Package class="size-4 text-primary" />
                  {values().vendor_code || "SKU-PRODUTO"}
                </span>
                <span
                  class={cn(
                    "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[11px] font-medium",
                    values().requires_registration
                      ? "bg-primary/10 text-primary"
                      : "bg-muted text-muted-foreground",
                  )}
                >
                  <ShieldCheck class="size-3" />
                  {values().requires_registration ? "Exige cadastro" : "Livre"}
                </span>
              </div>

              <div class="flex items-center justify-between rounded-lg border border-border/60 bg-muted/20 p-3">
                <div class="space-y-0.5 min-w-0 flex-1">
                  <span class="text-xs font-semibold text-foreground block truncate">
                    {values().variant_name || "Nome da Variação"}
                  </span>
                  <span class="text-[11px] font-mono text-muted-foreground block truncate">
                    {values().variant_vendor_code || "SKU-VAR"}
                  </span>
                </div>
                <div class="text-right shrink-0">
                  <span class="text-xs font-bold text-foreground block">
                    {formatPrice(values().variant_price_cents)}
                  </span>
                  <span class="text-[11px] text-muted-foreground block">
                    <Show
                      when={values().variant_stock.trim() !== ""}
                      fallback="Estoque ilimitado"
                    >
                      {values().variant_stock} un.
                    </Show>
                  </span>
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
      title={isEditing() ? "Editar produto" : "Novo produto"}
      description={
        isEditing()
          ? "Atualize o código de referência e configurações de registro."
          : "Cadastre o produto no catálogo com sua primeira variação de estoque e preço."
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
