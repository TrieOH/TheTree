import { useMemo } from "react";
import { useMultiStepForm } from "@/widgets/multi-step-form/hooks/use-multi-step-form";
import { MultiStepFormModal } from "@/widgets/multi-step-form/ui/multi-step-form-modal";
import {
  type EditionI,
  type EditionPatchInputI,
  type EditionPatchOutputI,
  editionPatchSchema,
} from "../model";
import { createEditionPatchFormSteps } from "../model/edition-form-steps";

interface EditEditionModalProps {
  open: boolean;
  edition: EditionI;
  onOpenChange: (open: boolean) => void;
  onUpdate: (
    values: EditionPatchOutputI,
  ) => Promise<EditionI | null | boolean> | EditionI | null | boolean;
}

export function EditEditionModal({
  open,
  edition,
  onOpenChange,
  onUpdate,
}: EditEditionModalProps) {
  const steps = useMemo(() => createEditionPatchFormSteps(), []);
  const defaultValues = useMemo<EditionPatchInputI>(
    () => ({
      name: edition.name,
      slug: edition.slug,
      tagline: edition.tagline ?? "",
      description: edition.description ?? "",
      registration_opens_at: edition.registration_opens_at ?? "",
      starts_at: edition.starts_at,
      ends_at: edition.ends_at,
      location_name: edition.location_name ?? "",
      location_description: edition.location_description ?? "",
      logo_url: edition.logo_url ?? "",
      banner_url: edition.banner_url ?? "",
      contact_email: edition.contact_email ?? "",
    }),
    [edition],
  );
  const controller = useMultiStepForm({
    schema: editionPatchSchema,
    steps,
    defaultValues,
    resetOnSuccessValues: defaultValues,
    onSubmit: async (values) => Boolean(await onUpdate(values)),
    onSubmitSuccess: () => onOpenChange(false),
  });

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Editar edição"
      controller={controller}
      submitLabel="Salvar alterações"
    />
  );
}
