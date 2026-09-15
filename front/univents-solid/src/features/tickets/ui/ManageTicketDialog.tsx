import { createEffect, createSignal, untrack } from "solid-js";
import { Button, Dialog, Field, Input } from "@trieoh/ui-solid";
import type { TicketCreateOutputI, TicketI } from "../model";

export interface ManageTicketValues {
  name: string;
  description: string;
  price: string; // in BRL, e.g. "50.00" or "0"
  access_level: number;
  max_quantity: string; // empty string for unlimited
}

export interface ManageTicketDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  ticket: TicketI | null;
  onSubmit: (values: TicketCreateOutputI) => Promise<boolean>;
}

const emptyValues: ManageTicketValues = {
  name: "",
  description: "",
  price: "0",
  access_level: 0,
  max_quantity: "",
};

const valuesOf = (ticket: TicketI | null): ManageTicketValues =>
  ticket
    ? {
        name: ticket.name,
        description: ticket.description ?? "",
        price: (ticket.price_cents / 100).toFixed(2),
        access_level: ticket.access_level,
        max_quantity: ticket.max_quantity != null ? String(ticket.max_quantity) : "",
      }
    : { ...emptyValues };

function validate(values: ManageTicketValues): Partial<Record<keyof ManageTicketValues, string>> {
  const errors: Partial<Record<keyof ManageTicketValues, string>> = {};

  if (values.name.trim().length < 2) {
    errors.name = "O nome deve ter pelo menos 2 caracteres.";
  }

  const numPrice = parseFloat(values.price.replace(",", "."));
  if (isNaN(numPrice) || numPrice < 0) {
    errors.price = "Informe um preço válido (>= 0).";
  }

  if (values.access_level < 0 || !Number.isInteger(Number(values.access_level))) {
    errors.access_level = "O nível de acesso deve ser um número inteiro >= 0.";
  }

  if (values.max_quantity.trim() !== "") {
    const qty = Number(values.max_quantity);
    if (!Number.isInteger(qty) || qty <= 0) {
      errors.max_quantity = "Quantidade deve ser um número inteiro maior que 0.";
    }
  }

  return errors;
}

export function ManageTicketDialog(props: ManageTicketDialogProps) {
  const [values, setValues] = createSignal(untrack(() => valuesOf(props.ticket)));
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<Partial<Record<keyof ManageTicketValues, string>>>({});

  createEffect(
    () => props.ticket,
    (t) => {
      setValues(valuesOf(t));
      setErrors({});
    },
  );

  const update = <K extends keyof ManageTicketValues>(key: K, value: ManageTicketValues[K]) => {
    setValues((prev) => ({ ...prev, [key]: value }));
    setErrors((prev) => ({ ...prev, [key]: undefined }));
  };

  const isEditing = () => Boolean(props.ticket);

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
    const maxQty = current.max_quantity.trim() !== "" ? parseInt(current.max_quantity, 10) : null;

    setSubmitting(true);
    try {
      const ok = await props.onSubmit({
        name: current.name.trim(),
        description: current.description.trim() || null,
        price_cents: priceCents,
        access_level: Number(current.access_level),
        max_quantity: maxQty,
      });
      if (ok) {
        props.onOpenChange(false);
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={isEditing() ? "Editar ticket" : "Novo ticket"}
      description={
        isEditing()
          ? "Atualize as informações do tipo de ingresso."
          : "Crie um novo tipo de ingresso para venda nesta edição."
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
            form="manage-ticket-form"
            disabled={submitting()}
          >
            {submitting()
              ? "Salvando..."
              : isEditing()
                ? "Salvar alterações"
                : "Criar ticket"}
          </Button>
        </div>
      }
    >
      <form id="manage-ticket-form" onSubmit={handleSubmit} class="space-y-4">
        <Field label="Nome do ticket" required error={errors().name}>
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              aria-invalid={errors().name ? "true" : undefined}
              value={values().name}
              placeholder="Ex: Ingresso Geral, VIP, Meia-entrada"
              disabled={submitting()}
              onInput={(e) => update("name", e.currentTarget.value)}
            />
          )}
        </Field>

        <Field label="Descrição" error={errors().description}>
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              aria-invalid={errors().description ? "true" : undefined}
              value={values().description}
              placeholder="Ex: Acesso a todas as palestras e coffee break"
              disabled={submitting()}
              onInput={(e) => update("description", e.currentTarget.value)}
            />
          )}
        </Field>

        <div class="grid grid-cols-2 gap-3">
          <Field label="Preço (R$)" required error={errors().price}>
            {(ids) => (
              <Input
                id={ids.id}
                type="number"
                step="0.01"
                min="0"
                aria-describedby={ids.describedBy}
                aria-invalid={errors().price ? "true" : undefined}
                value={values().price}
                placeholder="0.00"
                disabled={submitting()}
                onInput={(e) => update("price", e.currentTarget.value)}
              />
            )}
          </Field>

          <Field label="Nível de acesso" required error={errors().access_level}>
            {(ids) => (
              <Input
                id={ids.id}
                type="number"
                step="1"
                min="0"
                aria-describedby={ids.describedBy}
                aria-invalid={errors().access_level ? "true" : undefined}
                value={values().access_level}
                placeholder="0"
                disabled={submitting()}
                onInput={(e) => update("access_level", parseInt(e.currentTarget.value || "0", 10))}
              />
            )}
          </Field>
        </div>

        <Field
          label="Quantidade máxima (vagas)"
          error={errors().max_quantity}
          hint="Deixe vazio para quantidade ilimitada."
        >
          {(ids) => (
            <Input
              id={ids.id}
              type="number"
              step="1"
              min="1"
              aria-describedby={ids.describedBy}
              aria-invalid={errors().max_quantity ? "true" : undefined}
              value={values().max_quantity}
              placeholder="Ilimitado"
              disabled={submitting()}
              onInput={(e) => update("max_quantity", e.currentTarget.value)}
            />
          )}
        </Field>
      </form>
    </Dialog>
  );
}
