import type { ImageFieldChange, StepConfig } from "@/widgets/multi-step-form/model/types";
import { SocialLinksField } from "../ui/field/SocialLinkField";
import type { EventCreateInputI } from ".";

export function createEventFormSteps(
  track: (key: string) => (change: ImageFieldChange) => void,
): StepConfig<EventCreateInputI>[] {
  return [
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
        {
          kind: "image",
          name: "logo_url",
          label: "Logo do evento",
          hint: "PNG/JPEG (Máx. 2MB)",
          accept: "image/png,image/jpeg",
          maxSizeMB: 2,
          optional: true,
          layout: "half",
          onTrackingChange: track("logo_url"),
        },
        {
          kind: "image",
          name: "banner_url",
          label: "Banner principal",
          hint: "JPG/PNG (1920x1080)",
          accept: "image/jpeg,image/png",
          maxSizeMB: 5,
          optional: true,
          layout: "half",
          onTrackingChange: track("banner_url"),
        },
        {
          kind: "gallery",
          name: "gallery_urls",
          label: "Galeria de fotos",
          hint: "Até 10 fotos, JPG/PNG (Máx. 2MB cada)",
          accept: "image/jpeg,image/png",
          maxSizeMB: 2,
          maxItems: 10,
          optional: true,
          onTrackingChange: track("gallery_urls"),
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
          optional: true,
          name: "social-links-picker",
          render: ({ form }) => <SocialLinksField form={form} optional />,
        },
      ],
    },
  ];
}
