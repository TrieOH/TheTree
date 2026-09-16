import { Show, createEffect, createSignal, untrack } from "solid-js";
import { Button, Dialog, Field, Input, Textarea } from "@trieoh/ui-solid";
import type {
  CreateInitialProductOutputI,
  ProductI,
  ProductPatchOutputI,
} from "../model";

export interface ManageProductValues {
  vendor_code: string;
  requires_registration: boolean;
  // Initial variant fields (only when creating)
  variant_vendor_code: string;
  variant_name: string;
  variant_description: string;
  variant_price: string;
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
  variant_price: "0.00",
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

function validate(
  values: ManageProductValues,
  isEditing: boolean,
): Partial<Record<keyof ManageProductValues, string>> {
  const errors: Partial<Record<keyof ManageProductValues, string>> = {};

  if (values.vendor_code.trim().length < 2) {
    errors.vendor_code = "O código do produto deve ter pelo menos 2 caracteres.";
  }

  if (!isEditing) {
    if (values.variant_vendor_code.trim().length < 2) {
      errors.variant_vendor_code = "O código da variação deve ter pelo menos 2 caracteres.";
    }

    if (values.variant_name.trim().length < 2) {
      errors.variant_name = "O nome da variação deve ter pelo menos 2 caracteres.";
    }

    const numPrice = parseFloat(values.variant_price.replace(",", "."));
    if (Number.isNaN(numPrice) || numPrice < 0) {
      errors.variant_price = "Informe um preço válido (>= 0).";
    }

    if (values.variant_stock.trim() !== "") {
      const qty = Number(values.variant_stock);
      if (!Number.isInteger(qty) || qty <= 0) {
        errors.variant_stock = "Quantidade em estoque deve ser um número inteiro maior que 0.";
      }
    }
  }

  return errors;
}

export function ManageProductDialog(props: ManageProductDialogProps) {
  const [values, setValues] = createSignal(untrack(() => valuesOf(props.product)));
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<Partial<Record<keyof ManageProductValues, string>>>({});

  createEffect(
    () => props.product,
    (p) => {
      setValues(valuesOf(p));
      setErrors({});
    },
  );

  const update = <K extends keyof ManageProductValues>(key: K, value: ManageProductValues[K]) => {
    setValues((prev) => ({ ...prev, [key]: value }));
    setErrors((prev) => ({ ...prev, [key]: undefined }));
  };

  const isEditing = () => Boolean(props.product);

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    const current = values();
    const validationErrors = validate(current, isEditing());

    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }

    setSubmitting(true);
    try {
      if (isEditing()) {
        const ok = await props.onUpdate?.({
          vendor_code: current.vendor_code.trim(),
          requires_registration: current.requires_registration,
        });
        if (ok) props.onOpenChange(false);
      } else {
        const numPrice = parseFloat(current.variant_price.replace(",", "."));
        const priceCents = Math.round(numPrice * 100);
        const stockQty = current.variant_stock.trim() !== "" ? parseInt(current.variant_stock, 10) : null;

        const ok = await props.onCreate?.({
          vendor_code: current.vendor_code.trim(),
          requires_registration: current.requires_registration,
          variant_vendor_code: current.variant_vendor_code.trim(),
          name: current.variant_name.trim(),
          description: current.variant_description.trim() || null,
          price: priceCents,
          stock: stockQty,
        });
        if (ok) props.onOpenChange(false);
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={isEditing() ? "Editar produto" : "Novo produto"}
      description={
        isEditing()
          ? "Atualize o código de referência e configurações de registro."
          : "Cadastre o produto no catálogo com sua primeira variação de estoque e preço."
      }
      footer={
        <div class="flex w-full items-center justify-end gap-2 pt-2">
          <Button
            variant="outline"
            type="button"
            disabled={submitting()}
            onClick={() => props.onOpenChange(false)}
          >
            Cancelar
          </Button>
          <Button
            type="submit"
            form="manage-product-form"
            disabled={submitting()}
          >
            {submitting()
              ? "Salvando..."
              : isEditing()
                ? "Salvar alterações"
                : "Criar produto"}
          </Button>
        </div>
      }
    >
      <form id="manage-product-form" onSubmit={handleSubmit} class="space-y-4">
        {/* Identidade do Produto */}
        <div class="space-y-3">
          <Field
            label="Código do produto"
            hint="Identificador único (SKU geral) no evento, ex: CAMISETA-OFICIAL"
            required
            error={errors().vendor_code}
          >
            {(ids) => (
              <Input
                id={ids.id}
                aria-describedby={ids.describedBy}
                aria-invalid={errors().vendor_code ? "true" : undefined}
                value={values().vendor_code}
                onInput={(e) => update("vendor_code", e.currentTarget.value)}
                placeholder="Ex: CAMISETA-OFICIAL"
                disabled={submitting()}
                required
              />
            )}
          </Field>

          <label class="flex items-center justify-between gap-3 rounded-lg border border-border bg-muted/20 p-3 text-sm cursor-pointer hover:bg-muted/40 transition-colors">
            <div class="space-y-0.5">
              <span class="font-medium text-foreground">Exigir cadastro</span>
              <p class="text-xs text-muted-foreground">
                Exige que o comprador preencha dados de participante ao adquirir este item.
              </p>
            </div>
            <input
              type="checkbox"
              checked={values().requires_registration}
              onChange={(e) => update("requires_registration", e.currentTarget.checked)}
              disabled={submitting()}
              class="size-4 rounded border-border accent-primary cursor-pointer"
            />
          </label>
        </div>

        {/* Campos da primeira variação (apenas ao criar) */}
        <Show when={!isEditing()}>
          <div class="space-y-3 border-t border-border/60 pt-4">
            <div class="space-y-0.5">
              <h4 class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Primeira Variação
              </h4>
              <p class="text-xs text-muted-foreground">
                Defina o primeiro modelo ou tamanho deste produto.
              </p>
            </div>

            <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <Field
                label="Código da variação"
                hint="SKU específico, ex: CAMISETA-M"
                required
                error={errors().variant_vendor_code}
              >
                {(ids) => (
                  <Input
                    id={ids.id}
                    aria-describedby={ids.describedBy}
                    aria-invalid={errors().variant_vendor_code ? "true" : undefined}
                    value={values().variant_vendor_code}
                    onInput={(e) => update("variant_vendor_code", e.currentTarget.value)}
                    placeholder="Ex: CAMISETA-M"
                    disabled={submitting()}
                    required
                  />
                )}
              </Field>

              <Field
                label="Nome da variação"
                hint="Ex: Tamanho M - Preta"
                required
                error={errors().variant_name}
              >
                {(ids) => (
                  <Input
                    id={ids.id}
                    aria-describedby={ids.describedBy}
                    aria-invalid={errors().variant_name ? "true" : undefined}
                    value={values().variant_name}
                    onInput={(e) => update("variant_name", e.currentTarget.value)}
                    placeholder="Ex: Tamanho M - Preta"
                    disabled={submitting()}
                    required
                  />
                )}
              </Field>
            </div>

            <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <Field
                label="Preço (R$)"
                hint="Valor da unidade"
                required
                error={errors().variant_price}
              >
                {(ids) => (
                  <Input
                    id={ids.id}
                    aria-describedby={ids.describedBy}
                    aria-invalid={errors().variant_price ? "true" : undefined}
                    type="text"
                    inputmode="decimal"
                    value={values().variant_price}
                    onInput={(e) => update("variant_price", e.currentTarget.value)}
                    placeholder="0.00"
                    disabled={submitting()}
                    required
                  />
                )}
              </Field>

              <Field
                label="Estoque inicial"
                hint="Vazio = estoque ilimitado"
                error={errors().variant_stock}
              >
                {(ids) => (
                  <Input
                    id={ids.id}
                    aria-describedby={ids.describedBy}
                    aria-invalid={errors().variant_stock ? "true" : undefined}
                    type="number"
                    min="1"
                    step="1"
                    value={values().variant_stock}
                    onInput={(e) => update("variant_stock", e.currentTarget.value)}
                    placeholder="Ilimitado"
                    disabled={submitting()}
                  />
                )}
              </Field>
            </div>

            <Field label="Descrição da variação (opcional)">
              {(ids) => (
                <Textarea
                  id={ids.id}
                  rows={2}
                  value={values().variant_description}
                  onInput={(e) => update("variant_description", e.currentTarget.value)}
                  placeholder="Detalhes sobre tecido, medidas ou especificações..."
                  disabled={submitting()}
                />
              )}
            </Field>
          </div>
        </Show>
      </form>
    </Dialog>
  );
}
