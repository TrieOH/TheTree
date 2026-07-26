import type { StepConfig } from "@/widgets/multi-step-form/model/types";
import type { EditionCreateInputI, EditionPatchInputI } from ".";

export function createEditionFormSteps(): StepConfig<EditionCreateInputI>[] {
  return [
    {
      id: "edicao",
      label: "Edição",
      fields: [
        {
          kind: "text",
          name: "name",
          label: "Nome da edição",
          placeholder: "Ex: Edição 2026",
        },
        {
          kind: "text",
          name: "slug",
          label: "Slug",
          placeholder: "edicao-2026",
        },
      ],
    },
    {
      id: "cronograma",
      label: "Cronograma",
      fields: [
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
      ],
    },
  ];
}

export function createEditionPatchFormSteps(): StepConfig<EditionPatchInputI>[] {
  return [
    {
      id: "edicao",
      label: "Edição",
      fields: [
        {
          kind: "text",
          name: "name",
          label: "Nome da edição",
          placeholder: "Ex: Edição 2026",
        },
        {
          kind: "text",
          name: "slug",
          label: "Slug",
          placeholder: "edicao-2026",
        },
      ],
    },
    {
      id: "cronograma",
      label: "Cronograma",
      fields: [
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
      ],
    },
    {
      id: "detalhes",
      label: "Detalhes",
      fields: [
        {
          kind: "text",
          name: "tagline",
          label: "Tagline",
          optional: true,
        },
        {
          kind: "text",
          name: "description",
          label: "Descrição",
          optional: true,
        },
        {
          kind: "datetime",
          name: "registration_opens_at",
          label: "Abertura das inscrições",
          optional: true,
        },
        {
          kind: "text",
          name: "contact_email",
          label: "E-mail de contato",
          inputType: "email",
          optional: true,
        },
      ],
    },
    {
      id: "local",
      label: "Local",
      fields: [
        {
          kind: "text",
          name: "location_name",
          label: "Nome do local",
          optional: true,
        },
        {
          kind: "text",
          name: "location_description",
          label: "Descrição do local",
          optional: true,
        },
      ],
    },
  ];
}
