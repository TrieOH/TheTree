import type { StepConfig } from "@/widgets/multi-step-form/model/types";
import { SocialLinksField } from "../ui/field/SocialLinkField";
import type { EventCreateInputI } from ".";

export const eventFormSteps: StepConfig<EventCreateInputI>[] = [
  {
    id: "identidade",
    label: "Identidade",
    fields: [
      { kind: "text", name: "name", label: "Nome do evento", placeholder: "Ex: Tech Summit 2026" },
      {
        kind: "text",
        name: "slug",
        label: "Slug",
        placeholder: "tech-summit-2026",
      },
      {
        kind: "text",
        name: "acronym",
        label: "Sigla",
        placeholder: "TS26",
        optional: true,
      },
      {
        kind: "text",
        name: "tagline",
        label: "Tagline",
        placeholder: "Uma frase curta sobre o evento",
        optional: true,
      },
    ],
  },
  {
    id: "midia",
    label: "Mídia",
    fields: [
      { kind: "text", name: "logo_url", label: "Logo", placeholder: "https://...", optional: true, inputType: "url" },
      {
        kind: "text",
        name: "banner_url",
        label: "Banner",
        placeholder: "https://...",
        optional: true,
        inputType: "url",
      },
    ],
  },
  {
    id: "conexao",
    label: "Conexão",
    fields: [
      {
        kind: "text",
        name: "contact_email",
        label: "E-mail de contato",
        placeholder: "contato@evento.com",
        inputType: "email",
      },
      {
        kind: "custom",
        name: "social-links-picker",
        render: ({ form }) => <SocialLinksField form={form} />,
      },
    ],
  },
];
