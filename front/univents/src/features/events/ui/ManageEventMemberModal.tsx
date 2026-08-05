import { useMultiStepForm } from "@/widgets/multi-step-form/hooks/use-multi-step-form";
import { MultiStepFormModal } from "@/widgets/multi-step-form/ui/multi-step-form-modal";
import type { EventMemberWithEmailI } from "../api/members";
import {
  type EventMemberCreateInput,
  type EventMemberCreateOutput,
  eventMemberCreateSchema,
} from "../model/member";
import { createEventMemberFormSteps } from "../model/member-form-steps";

interface ManageEventMemberModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onCreate: (
    values: EventMemberCreateOutput,
  ) => boolean | Promise<boolean | EventMemberWithEmailI>;
}

const defaultValues: EventMemberCreateInput = {
  email: "",
  role: "staff",
};

export function ManageEventMemberModal({
  open,
  onOpenChange,
  onCreate,
}: ManageEventMemberModalProps) {
  const controller = useMultiStepForm({
    schema: eventMemberCreateSchema,
    steps: createEventMemberFormSteps(),
    defaultValues,
    resetOnSuccessValues: defaultValues,
    onSubmit: async (values): Promise<boolean> => {
      return Boolean(await onCreate(values));
    },
    onSubmitSuccess: () => onOpenChange(false),
  });

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Adicionar membro"
      controller={controller}
      submitLabel="Adicionar membro"
    />
  );
}
