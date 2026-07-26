import type { StepConfig } from "@/widgets/multi-step-form/model/types";
import type { ActivityCreateInputI } from ".";

const difficultyOptions = [
  { label: "Sem pré-requisitos", value: "no_prerequisites" },
  { label: "Iniciante", value: "beginner" },
  { label: "Intermediário", value: "intermediate" },
  { label: "Avançado", value: "advanced" },
  { label: "Especialista", value: "expert" },
];

export function createActivityFormSteps(): StepConfig<ActivityCreateInputI>[] {
  return [
    {
      id: "identidade",
      label: "Identidade",
      fields: [
        {
          kind: "text",
          name: "title",
          label: "Título",
          placeholder: "Nome da atividade",
        },
        {
          kind: "text",
          name: "presenter_name",
          label: "Apresentador/Palestrante",
          placeholder: "Nome do responsável",
          optional: true,
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
                placeholder="Descreva a atividade"
                className="min-h-28 w-full rounded-xl border border-border/60 bg-background px-3 py-2.5 text-sm text-foreground outline-none transition-colors placeholder:text-muted-foreground/70 focus:border-primary focus:ring-2 focus:ring-primary/15"
              />
            </label>
          ),
        },
      ],
    },
    {
      id: "agenda",
      label: "Agenda",
      fields: [
        {
          kind: "text",
          name: "location",
          label: "Local",
          placeholder: "Sala, auditório, etc.",
        },
        {
          kind: "datetime",
          name: "starts_at",
          label: "Início",
        },
        {
          kind: "datetime",
          name: "ends_at",
          label: "Término",
        },
        {
          kind: "combobox",
          name: "difficulty",
          label: "Nível de dificuldade",
          placeholder: "Selecione a dificuldade",
          options: difficultyOptions,
        },
      ],
    },
    {
      id: "capacidade",
      label: "Capacidade",
      fields: [
        {
          kind: "text",
          name: "token_cost",
          label: "Custo em tokens",
          placeholder: "0",
          inputType: "number",
        },
        {
          kind: "toggle",
          name: "has_capacity",
          label: "Limitar vagas",
        },
        {
          kind: "text",
          name: "capacity",
          label: "Capacidade",
          placeholder: "Número de vagas",
          inputType: "number",
          optional: true,
          visibleIf: { type: "equals", field: "has_capacity", value: true },
        },
      ],
    },
  ];
}
