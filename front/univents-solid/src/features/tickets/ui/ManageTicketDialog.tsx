import { createEffect, createSignal, Show, untrack } from "solid-js";
import {
  MultiStepDialog,
  type MultiStepItem,
} from "@trieoh/ui-solid";
import { formatPrice } from "@/shared/lib/money";
import type { TicketCreateOutputI, TicketI } from "../model";

export interface ManageTicketValues {
  name: string;
  description: string;
  price_cents: number;
  access_level: number | "";
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
  price_cents: 0,
  access_level: "",
  max_quantity: "",
};

const valuesOf = (ticket: TicketI | null): ManageTicketValues =>
  ticket
    ? {
        name: ticket.name,
        description: ticket.description ?? "",
        price_cents: ticket.price_cents,
        access_level: ticket.access_level,
        max_quantity: ticket.max_quantity != null ? String(ticket.max_quantity) : "",
      }
    : { ...emptyValues };

export function ManageTicketDialog(props: ManageTicketDialogProps) {
  const [currentStep, setCurrentStep] = createSignal(0);
  let currentValues: ManageTicketValues = untrack(() => valuesOf(props.ticket));
  const [values, setValues] = createSignal<ManageTicketValues>(
    untrack(() => currentValues),
  );
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<
    Partial<Record<keyof ManageTicketValues, string>>
  >({});

  let isMounted = false;
  createEffect(
    () => (props.ticket ? props.ticket.id : null),
    () => {
      if (!isMounted) {
        isMounted = true;
        return;
      }
      currentValues = valuesOf(props.ticket);
      setValues(currentValues);
      setErrors({});
      setCurrentStep(0);
    },
  );

  const isEditing = () => Boolean(props.ticket);

  const handleChange = (
    key: keyof ManageTicketValues & string,
    value: unknown,
  ) => {
    currentValues = { ...currentValues, [key]: value };
    setValues(currentValues);
    setErrors((prev) => {
      if (!prev[key as keyof ManageTicketValues]) return prev;
      const nextErrors = { ...prev };
      delete nextErrors[key as keyof ManageTicketValues];
      return nextErrors;
    });
  };

  const validateStep = (stepIndex: number): boolean => {
    const current = currentValues;
    const nextErrors: Partial<Record<keyof ManageTicketValues, string>> = {};

    if (stepIndex === 0) {
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
        nextErrors.price_cents = "Informe um preço válido.";
      }

      const rawLevel = String(current.access_level ?? "").trim();
      if (rawLevel === "") {
        nextErrors.access_level = "O nível de acesso é obrigatório.";
      } else {
        const accessLevel = Number(rawLevel);
        if (
          isNaN(accessLevel) ||
          accessLevel < 0 ||
          !Number.isInteger(accessLevel)
        ) {
          nextErrors.access_level =
            "O nível de acesso deve ser um número inteiro >= 0.";
        }
      }

      if (String(current.max_quantity ?? "").trim() !== "") {
        const qty = Number(current.max_quantity);
        if (!Number.isInteger(qty) || qty <= 0) {
          nextErrors.max_quantity =
            "Quantidade deve ser um número inteiro maior que 0.";
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
    const rawLevel = String(current.access_level ?? "").trim();
    const accessLevel = rawLevel === "" ? 0 : parseInt(rawLevel, 10);
    const maxQty =
      String(current.max_quantity ?? "").trim() !== ""
        ? parseInt(String(current.max_quantity), 10)
        : null;

    setSubmitting(true);
    try {
      const ok = await props.onSubmit({
        name: current.name.trim(),
        description: current.description?.trim() || null,
        price_cents: current.price_cents,
        access_level: accessLevel,
        max_quantity: maxQty,
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

  const steps: MultiStepItem<ManageTicketValues>[] = [
    {
      id: "ingresso",
      title: "Ingresso",
      description: "Informações principais do tipo de ingresso",
      fields: [
        {
          name: "name",
          label: "Nome do ticket",
          required: true,
          placeholder: "Ex: Ingresso Geral, VIP, Meia-entrada",
        },
        {
          name: "description",
          label: "Descrição",
          kind: "textarea",
          rows: 3,
          placeholder: "Ex: Acesso a todas as palestras, workshops e coffee break",
          hint: "Detalhes do que está incluso neste ingresso (opcional).",
        },
      ],
    },
    {
      id: "valores",
      title: "Valores e Vagas",
      description: "Preço, limite de capacidade e nível de acesso",
      fields: [
        {
          name: "price_cents",
          label: "Preço",
          kind: "money",
          currency: "BRL",
          locale: "pt-BR",
          layout: "half",
          required: true,
          hint: "R$ 0,00 para ingresso gratuito.",
        },
        {
          name: "max_quantity",
          label: "Quantidade máxima (vagas)",
          kind: "number",
          step: "1",
          min: "1",
          layout: "half",
          placeholder: "Ilimitado",
          hint: "Deixe vazio para ilimitado.",
        },
        {
          name: "access_level",
          label: "Nível de acesso",
          kind: "number",
          step: "1",
          min: "0",
          layout: "half",
          required: true,
          placeholder: "0",
          hint: "Nível numérico de permissão.",
        },
      ],
    },
    {
      id: "resumo",
      title: "Resumo",
      description: "Revise os dados antes de salvar o ticket",
      summary: {
        title: "Dados do ticket",
        badge: () => (props.ticket ? "Pronto para atualizar" : "Pronto para criar"),
        items: [
          {
            label: "Nome do ticket",
            value: () => values().name,
          },
          {
            label: "Valor",
            value: () => formatPrice(values().price_cents ?? 0),
          },
          {
            label: "Capacidade",
            value: () =>
              values().max_quantity ? `${values().max_quantity} vagas` : "Ilimitada",
          },
          {
            label: "Nível de acesso",
            value: () =>
              values().access_level !== ""
                ? `Nível ${values().access_level}`
                : "Não definido",
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
                  Prévia do ingresso
                </p>
                <h4 class="text-base font-semibold text-foreground truncate">
                  {values().name || "Nome do ticket"}
                </h4>
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
                Acesso nível{" "}
                {values().access_level !== "" ? values().access_level : 0}
              </span>
              <span class="inline-flex items-center gap-1 rounded-md bg-muted px-2 py-0.5 font-medium">
                {values().max_quantity
                  ? `${values().max_quantity} vagas`
                  : "Vagas ilimitadas"}
              </span>
            </div>
          </div>
        ),
      },
    },
  ];

  return (
    <MultiStepDialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={isEditing() ? "Editar ticket" : "Novo ticket"}
      description={
        isEditing()
          ? "Atualize as informações e regras deste tipo de ingresso."
          : "Crie um novo tipo de ingresso para venda nesta edição."
      }
      steps={steps}
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
            : "Criar ticket"
      }
      onBeforeNext={validateStep}
      onSubmit={handleSubmit}
    />
  );
}
