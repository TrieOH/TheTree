import type { StepConfig } from "@/widgets/multi-step-form/model/types";
import { MemberRoleField } from "../ui/field/MemberRoleField";
import type { EventMemberCreateInput } from "./member";

export function createEventMemberFormSteps(): StepConfig<EventMemberCreateInput>[] {
  return [
    {
      id: "membro",
      label: "Membro",
      fields: [
        {
          kind: "text",
          name: "email",
          label: "E-mail",
          placeholder: "membro@exemplo.com",
          inputType: "email",
          description:
            "O usuário precisa possuir uma conta vinculada a este e-mail.",
        },
        {
          kind: "custom",
          name: "role",
          render: ({ form }) => <MemberRoleField form={form} />,
        },
      ],
    },
  ];
}
