import { useMemo } from "react";
import { useMultiStepForm } from "@/widgets/multi-step-form/hooks/use-multi-step-form";
import { MultiStepFormModal } from "@/widgets/multi-step-form/ui/multi-step-form-modal";
import {
  type TicketCreateInputI,
  type TicketCreateOutputI,
  type TicketI,
  ticketCreateSchema,
} from "../model";
import { createTicketFormSteps } from "../model/ticket-form-steps";

export interface ManageTicketModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  ticket?: TicketI;
  onCreate?: (
    values: TicketCreateOutputI,
  ) => Promise<TicketI | null | boolean> | TicketI | null | boolean;
  onUpdate?: (
    id: string,
    values: TicketCreateOutputI,
  ) => Promise<TicketI | null | boolean> | TicketI | null | boolean;
}

const emptyDefaultValues: TicketCreateInputI = {
  name: "",
  description: "",
  access_level: 0,
  price_cents: 0,
  max_quantity: null,
};

function toFormValues(ticket: TicketI): TicketCreateInputI {
  return {
    name: ticket.name,
    description: ticket.description ?? "",
    access_level: ticket.access_level,
    price_cents: ticket.price_cents,
    max_quantity: ticket.max_quantity,
  };
}

export function ManageTicketModal({
  open,
  onOpenChange,
  ticket,
  onCreate,
  onUpdate,
}: ManageTicketModalProps) {
  const isEditing = Boolean(ticket);
  const editValues = useMemo(
    () => (ticket ? toFormValues(ticket) : undefined),
    [ticket],
  );
  const steps = useMemo(() => createTicketFormSteps(), []);

  const controller = useMultiStepForm({
    schema: ticketCreateSchema,
    steps,
    defaultValues: emptyDefaultValues,
    values: editValues,
    requireDirtyToSubmit: isEditing,
    onSubmit: async (values): Promise<boolean> => {
      const result = ticket
        ? await onUpdate?.(ticket.id, values)
        : await onCreate?.(values);

      return Boolean(result);
    },
    onSubmitSuccess: () => onOpenChange(false),
  });

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEditing ? "Editar ticket" : "Criar ticket"}
      controller={controller}
      submitLabel={isEditing ? "Salvar alterações" : "Criar ticket"}
    />
  );
}
