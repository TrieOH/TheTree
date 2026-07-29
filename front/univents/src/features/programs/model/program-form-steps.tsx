import type { StepConfig } from "@/widgets/multi-step-form/model/types";
import type { ProgramCreateInput } from ".";

export function createProgramFormSteps(): StepConfig<ProgramCreateInput>[] {
  return [
    {
      id: "identidade",
      label: "Identidade",
      fields: [
        {
          kind: "combobox",
          name: "kind",
          label: "Tipo",
          placeholder: "Selecione o tipo",
          options: [
            { label: "Atividade", value: "activity" },
            { label: "Checkpoint", value: "checkpoint" },
          ],
        },
        {
          kind: "text",
          name: "name",
          label: "Nome",
          placeholder: "Nome do programa",
        },
        {
          kind: "custom",
          name: "description",
          optional: true,
          render: ({ form }) => (
            <label className="block space-y-2">
              <span className="block text-sm font-medium">Descrição</span>
              <textarea
                {...form.register("description")}
                rows={4}
                className="w-full rounded-xl border border-border/60 bg-background px-3 py-2.5 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/15"
              />
            </label>
          ),
        },
      ],
    },
    {
      id: "acesso",
      label: "Acesso",
      fields: [
        {
          kind: "text",
          name: "min_access_level",
          label: "Nível mínimo de acesso",
          inputType: "number",
          optional: true,
        },
        {
          kind: "money",
          name: "price",
          label: "Preço",
          currency: "BRL",
          valueType: "number",
          maxCents: 99999999999,
        },
        { kind: "toggle", name: "staff_only", label: "Somente equipe" },
      ],
    },
  ];
}
