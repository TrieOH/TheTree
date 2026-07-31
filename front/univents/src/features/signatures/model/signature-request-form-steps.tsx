import type { StepConfig } from "@/widgets/multi-step-form/model/types";
import type { SignatureRequestCreateInputI } from ".";

export function createSignatureRequestFormSteps(): StepConfig<SignatureRequestCreateInputI>[] {
  return [
    {
      id: "signatario",
      label: "Signatário",
      fields: [
        {
          kind: "text",
          name: "signatory_name",
          label: "Nome do signatário",
          placeholder: "Nome completo",
        },
        {
          kind: "text",
          name: "signatory_title",
          label: "Cargo",
          placeholder: "Cargo ou função",
          optional: true,
        },
        {
          kind: "text",
          name: "signatory_email",
          label: "E-mail",
          placeholder: "nome@exemplo.com",
          inputType: "email",
        },
        {
          kind: "text",
          name: "signatory_user_id",
          label: "ID do usuário",
          placeholder: "Opcional",
          optional: true,
        },
      ],
    },
    {
      id: "validade",
      label: "Validade",
      fields: [
        {
          kind: "text",
          name: "expires_in_days",
          label: "Validade do convite (dias)",
          inputType: "number",
          description: "Defina entre 1 e 365 dias.",
        },
      ],
    },
  ];
}
