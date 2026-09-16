import { createEffect, createSignal, untrack } from "solid-js";
import { Button, Dialog, Field, Input, Textarea } from "@trieoh/ui-solid";
import type { VariantCreateOutputI, VariantI } from "../model";

export interface ManageVariantValues {
  vendor_code: string;
  name: string;
  description: string;
  price: string;
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
  price: "0.00",
  stock: "",
};

const valuesOf = (variant: VariantI | null): ManageVariantValues =>
  variant
    ? {
        vendor_code: variant.vendor_code,
        name: variant.name,
        description: variant.description ?? "",
        price: (variant.price / 100).toFixed(2),
        stock: variant.stock != null ? String(variant.stock) : "",
      }
    : { ...emptyValues };

function validate(values: ManageVariantValues): Partial<Record<keyof ManageVariantValues, string>> {
  const errors: Partial<Record<keyof ManageVariantValues, string>> = {};

  if (values.vendor_code.trim().length < 2) {
    errors.vendor_code = "O código deve ter pelo menos 2 caracteres.";
  }

  if (values.name.trim().length < 2) {
    errors.name = "O nome deve ter pelo menos 2 caracteres.";
  }

  const numPrice = parseFloat(values.price.replace(",", "."));
  if (Number.isNaN(numPrice) || numPrice < 0) {
    errors.price = "Informe um preço válido (>= 0).";
  }

  if (values.stock.trim() !== "") {
    const qty = Number(values.stock);
    if (!Number.isInteger(qty) || qty <= 0) {
      errors.stock = "Quantidade em estoque deve ser um número inteiro maior que 0.";
    }
  }

  return errors;
}

export function ManageVariantDialog(props: ManageVariantDialogProps) {
  const [values, setValues] = createSignal(untrack(() => valuesOf(props.variant)));
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<Partial<Record<keyof ManageVariantValues, string>>>({});

  createEffect(
    () => props.variant,
    (v) => {
      setValues(valuesOf(v));
      setErrors({});
    },
  );

  const update = <K extends keyof ManageVariantValues>(key: K, value: ManageVariantValues[K]) => {
    setValues((prev) => ({ ...prev, [key]: value }));
    setErrors((prev) => ({ ...prev, [key]: undefined }));
  };

  const isEditing = () => Boolean(props.variant);

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    const current = values();
    const validationErrors = validate(current);

    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }

    const numPrice = parseFloat(current.price.replace(",", "."));
    const priceCents = Math.round(numPrice * 100);
    const stockQty = current.stock.trim() !== "" ? parseInt(current.stock, 10) : null;

    setSubmitting(true);
    try {
      const ok = await props.onSubmit({
        vendor_code: current.vendor_code.trim(),
        name: current.name.trim(),
        description: current.description.trim() || null,
        price: priceCents,
        stock: stockQty,
        gallery_urls: props.variant?.gallery_urls ?? [],
      });
      if (ok) props.onOpenChange(false);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={isEditing() ? "Editar variação" : "Nova variação"}
      description={
        isEditing()
          ? "Atualize as informações de preço e estoque da variação."
          : "Cadastre um novo tamanho, cor ou modelo para este produto."
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
            form="manage-variant-form"
            disabled={submitting()}
          >
            {submitting()
              ? "Salvando..."
              : isEditing()
                ? "Salvar alterações"
                : "Criar variação"}
          </Button>
        </div>
      }
    >
      <form id="manage-variant-form" onSubmit={handleSubmit} class="space-y-4">
        <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <Field
            label="Código da variação"
            hint="Identificador único (SKU), ex: CAMISETA-G"
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
                placeholder="Ex: CAMISETA-G"
                disabled={submitting()}
                required
              />
            )}
          </Field>

          <Field
            label="Nome da variação"
            hint="Ex: Tamanho G - Branca"
            required
            error={errors().name}
          >
            {(ids) => (
              <Input
                id={ids.id}
                aria-describedby={ids.describedBy}
                aria-invalid={errors().name ? "true" : undefined}
                value={values().name}
                onInput={(e) => update("name", e.currentTarget.value)}
                placeholder="Ex: Tamanho G - Branca"
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
            error={errors().price}
          >
            {(ids) => (
              <Input
                id={ids.id}
                aria-describedby={ids.describedBy}
                aria-invalid={errors().price ? "true" : undefined}
                type="text"
                inputmode="decimal"
                value={values().price}
                onInput={(e) => update("price", e.currentTarget.value)}
                placeholder="0.00"
                disabled={submitting()}
                required
              />
            )}
          </Field>

          <Field
            label="Estoque"
            hint="Deixe vazio para ilimitado"
            error={errors().stock}
          >
            {(ids) => (
              <Input
                id={ids.id}
                aria-describedby={ids.describedBy}
                aria-invalid={errors().stock ? "true" : undefined}
                type="number"
                min="1"
                step="1"
                value={values().stock}
                onInput={(e) => update("stock", e.currentTarget.value)}
                placeholder="Ilimitado"
                disabled={submitting()}
              />
            )}
          </Field>
        </div>

        <Field label="Descrição (opcional)">
          {(ids) => (
            <Textarea
              id={ids.id}
              rows={3}
              value={values().description}
              onInput={(e) => update("description", e.currentTarget.value)}
              placeholder="Especificações, medidas, material..."
              disabled={submitting()}
            />
          )}
        </Field>
      </form>
    </Dialog>
  );
}
