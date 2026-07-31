import { useMemo } from "react";
import { useMultiStepForm } from "@/widgets/multi-step-form/hooks/use-multi-step-form";
import { MultiStepFormModal } from "@/widgets/multi-step-form/ui/multi-step-form-modal";
import {
  type OccurrenceCreateOutput,
  type OccurrenceI,
  occurrenceSchema,
} from "../model";
import { createOccurrenceFormSteps } from "../model/occurrence-form-steps";

export function ManageOccurrenceModal({
  open,
  occurrence,
  onOpenChange,
  onSave,
}: {
  open: boolean;
  occurrence?: OccurrenceI;
  onOpenChange: (open: boolean) => void;
  onSave: (values: OccurrenceCreateOutput) => Promise<boolean>;
}) {
  const steps = useMemo(() => createOccurrenceFormSteps(), []);
  const values = useMemo(
    () =>
      occurrence
        ? {
            starts_at: occurrence.starts_at,
            ends_at: occurrence.ends_at,
            max_capacity: occurrence.max_capacity ?? undefined,
          }
        : undefined,
    [occurrence],
  );
  const controller = useMultiStepForm({
    schema: occurrenceSchema,
    steps,
    defaultValues: { starts_at: "", ends_at: "", max_capacity: undefined },
    values,
    requireDirtyToSubmit: Boolean(occurrence),
    onSubmit: onSave,
    onSubmitSuccess: () => onOpenChange(false),
  });
  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title={occurrence ? "Editar ocorrência" : "Nova ocorrência"}
      controller={controller}
      submitLabel={occurrence ? "Salvar alterações" : "Criar ocorrência"}
    />
  );
}
