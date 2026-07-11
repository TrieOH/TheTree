import { useEffect, useMemo, useRef } from "react";
import { useMultiStepForm } from "@/widgets/multi-step-form/hooks/use-multi-step-form";
import { eventCreateSchema, type EventCreateInputI, type EventCreateOutputI, type EventI } from "../model";
import { MultiStepFormModal } from "@/widgets/multi-step-form/ui/multi-step-form-modal";
import { eventFormSteps } from "../model/event-form-steps";

export interface ManageEventModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** Omit for create mode. Pass the record being edited for edit mode. */
  event?: EventI;
  onCreate?: (values: EventCreateOutputI) => Promise<boolean> | boolean;
  onUpdate?: (id: string, values: EventCreateOutputI) => Promise<boolean> | boolean;
}

const emptyDefaultValues: EventCreateInputI = {
  name: "",
  slug: "",
  contact_email: "",
};

/** Maps an API record (which uses `null` for "empty") to the form's
 * input shape (which uses `""`, since that's what a text input holds). */
function toFormValues(event: EventI): EventCreateInputI {
  return {
    name: event.name,
    slug: event.slug,
    acronym: event.acronym ?? "",
    tagline: event.tagline ?? "",
    logo_url: event.logo_url ?? "",
    banner_url: event.banner_url ?? "",
    contact_email: event.contact_email,
    social_links: {
      twitter: event.social_links?.twitter ?? "",
      instagram: event.social_links?.instagram ?? "",
      linkedin: event.social_links?.linkedin ?? "",
      website: event.social_links?.website ?? "",
    },
  };
}

function toSlug(value: string): string {
  return value
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

function toAcronym(value: string): string {
  return value
    .trim()
    .split(/\s+/)
    .filter((word) => word.length > 2)
    .map((word) => word[0]?.toUpperCase() ?? "")
    .join("");
}

export function ManageEventModal({ open, onOpenChange, event, onCreate, onUpdate }: ManageEventModalProps) {
  const isEditing = Boolean(event);
  const editValues = useMemo(() => (event ? toFormValues(event) : undefined), [event]);
  const autoAcronymRef = useRef("");
  const autoSlugRef = useRef("");

  const controller = useMultiStepForm({
    schema: eventCreateSchema,
    steps: eventFormSteps,
    defaultValues: emptyDefaultValues,
    resetOnSuccessValues: isEditing ? undefined : emptyDefaultValues,
    // Only pass `values` in edit mode - RHF resets/re-syncs the form to
    // this object whenever its reference changes (e.g. a different
    // event gets picked, or fresh data comes back from a refetch).
    values: editValues,
    // Blocks the submit button until something actually changed only meaningful while editing.
    requireDirtyToSubmit: isEditing,
    onSubmit: async (values): Promise<boolean> => {
      const result = event
        ? await onUpdate?.(event.id, values)
        : await onCreate?.(values);

      return result !== false;
    },
    onSubmitSuccess: () => onOpenChange(false),
  });

  const name = controller.form.watch("name");

  useEffect(() => {
    if (isEditing) return;

    const trimmedName = name.trim();
    if (!trimmedName) {
      const currentAcronym = controller.form.getValues("acronym") ?? "";
      const currentSlug = controller.form.getValues("slug") ?? "";

      if (currentAcronym === "" || currentAcronym === autoAcronymRef.current) {
        controller.form.setValue("acronym", "", {
          shouldDirty: false,
          shouldTouch: false,
          shouldValidate: false,
        });
      }

      if (currentSlug === "" || currentSlug === autoSlugRef.current) {
        controller.form.setValue("slug", "", {
          shouldDirty: false,
          shouldTouch: false,
          shouldValidate: false,
        });
      }

      autoAcronymRef.current = "";
      autoSlugRef.current = "";
      return;
    }

    const nextAcronym = toAcronym(trimmedName);
    const nextSlug = toSlug(trimmedName);
    const currentAcronym = controller.form.getValues("acronym") ?? "";
    const currentSlug = controller.form.getValues("slug") ?? "";

    const acronymIsAuto = currentAcronym === "" || currentAcronym === autoAcronymRef.current;
    const slugIsAuto = currentSlug === "" || currentSlug === autoSlugRef.current;

    if (acronymIsAuto && currentAcronym !== nextAcronym) {
      controller.form.setValue("acronym", nextAcronym, {
        shouldDirty: false,
        shouldTouch: false,
        shouldValidate: true,
      });
    }

    if (slugIsAuto && currentSlug !== nextSlug) {
      controller.form.setValue("slug", nextSlug, {
        shouldDirty: false,
        shouldTouch: false,
        shouldValidate: true,
      });
    }

    autoAcronymRef.current = nextAcronym;
    autoSlugRef.current = nextSlug;
  }, [controller.form, isEditing, name]);

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEditing ? "Editar Evento" : "Criar Novo Evento"}
      controller={controller}
      submitLabel={isEditing ? "Salvar alterações" : "Criar Evento"}
    />
  );
}
