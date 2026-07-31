import type { StepConfig } from "@/widgets/multi-step-form/model/types";
import type { OccurrenceCreateInput } from ".";

export function createOccurrenceFormSteps(): StepConfig<OccurrenceCreateInput>[] {
  return [
    {
      id: "agenda",
      label: "Agenda",
      fields: [
        { kind: "datetime", name: "starts_at", label: "Início" },
        { kind: "datetime", name: "ends_at", label: "Término" },
        {
          kind: "text",
          name: "max_capacity",
          label: "Capacidade máxima",
          inputType: "number",
          optional: true,
        },
      ],
    },
  ];
}
