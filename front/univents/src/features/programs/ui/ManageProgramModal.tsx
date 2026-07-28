import { useMemo } from "react";
import { useMultiStepForm } from "@/widgets/multi-step-form/hooks/use-multi-step-form";
import { MultiStepFormModal } from "@/widgets/multi-step-form/ui/multi-step-form-modal";
import {
  type ProgramCreateInput,
  type ProgramCreateOutput,
  type ProgramI,
  programSchema,
} from "../model";
import { createProgramFormSteps } from "../model/program-form-steps";

interface Props {
  open: boolean;
  program?: ProgramI;
  onOpenChange: (open: boolean) => void;
  onSave: (values: ProgramCreateOutput) => Promise<boolean>;
}

const defaults: ProgramCreateInput = {
  kind: "activity",
  name: "",
  description: "",
  min_access_level: undefined,
  staff_only: false,
  price: undefined,
};

export function ManageProgramModal({
  open,
  program,
  onOpenChange,
  onSave,
}: Props) {
  const steps = useMemo(() => createProgramFormSteps(), []);
  const values = useMemo(
    () =>
      program
        ? { ...program, description: program.description ?? "" }
        : undefined,
    [program],
  );
  const controller = useMultiStepForm({
    schema: programSchema,
    steps,
    defaultValues: defaults,
    values,
    requireDirtyToSubmit: Boolean(program),
    onSubmit: onSave as (values: ProgramCreateInput) => Promise<boolean>,
    onSubmitSuccess: () => onOpenChange(false),
  });

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title={program ? "Editar programa" : "Novo programa"}
      controller={controller}
      submitLabel={program ? "Salvar alterações" : "Criar programa"}
    />
  );
}
