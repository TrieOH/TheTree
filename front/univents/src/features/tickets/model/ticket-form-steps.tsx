import type { StepConfig } from "@/widgets/multi-step-form/model/types";
import type { TicketCreateInputI } from ".";

export function createTicketFormSteps(): StepConfig<TicketCreateInputI>[] {
  return [
    {
      id: "identidade",
      label: "Identidade",
      fields: [
        {
          kind: "text",
          name: "name",
          label: "Nome",
          placeholder: "Nome do ticket",
        },
        {
          kind: "custom",
          name: "description",
          optional: true,
          render: ({ form }) => (
            <label className="block space-y-2">
              <span className="block text-sm font-medium text-foreground">
                Descrição
              </span>
              <textarea
                {...form.register("description")}
                rows={4}
                placeholder="Descreva o ticket"
                className="min-h-28 w-full rounded-xl border border-border/60 bg-background px-3 py-2.5 text-sm text-foreground outline-none transition-colors placeholder:text-muted-foreground/70 focus:border-primary focus:ring-2 focus:ring-primary/15"
              />
            </label>
          ),
        },
      ],
    },
    {
      id: "comercial",
      label: "Comercial",
      fields: [
        {
          kind: "money",
          name: "price_cents",
          label: "Preço (em centavos)",
          currency: "BRL",
          valueType: "number",
          maxCents: 99999999999,
          // placeholder: 'Ex: 5000 para R$50,00',
          // inputType: 'number',
        },
        {
          kind: "text",
          name: "access_level",
          label: "Nível de acesso",
          placeholder: "0",
          inputType: "number",
        },
        {
          kind: "text",
          name: "max_quantity",
          label: "Quantidade máxima",
          placeholder: "Deixe em branco para ilimitado",
          inputType: "number",
          optional: true,
        },
      ],
    },
  ];
}
