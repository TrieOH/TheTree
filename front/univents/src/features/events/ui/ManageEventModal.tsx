import { useEffect, useMemo, useRef } from "react";
import { useMultiStepForm } from "@/widgets/multi-step-form/hooks/use-multi-step-form";
import { eventCreateSchema } from "../model";
import type { EventCreateInputI, EventCreateSubmitI, EventI } from "../model";
import { MultiStepFormModal } from "@/widgets/multi-step-form/ui/multi-step-form-modal";
import { createEventFormSteps } from "../model/event-form-steps";
import { useImageFieldTracking } from "@/widgets/multi-step-form/hooks/use-image-field-tracking";
import { useSyncEventMediaMutation } from "../api/mutations";

export interface ManageEventModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** Omit for create mode. Pass the record being edited for edit mode. */
  event?: EventI;
  onCreate?: (values: EventCreateSubmitI) => Promise<EventI | null | boolean> | EventI | null | boolean;
  onUpdate?: (id: string, values: EventCreateSubmitI) => Promise<EventI | null | boolean> | EventI | null | boolean;
}

const emptyDefaultValues: EventCreateInputI = {
  name: "",
  slug: "",
  contact_email: "",
};

const MEDIA_FIELD_KEYS = new Set(["logo_url", "banner_url", "gallery_urls"] as const);

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
    gallery_urls: event.gallery_urls ?? [],
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
    .map((word) => word.charAt(0).toUpperCase())
    .join("");
}

function hasDirtyFieldsOutsideMedia(dirtyFields: unknown): boolean {
  if (!dirtyFields || dirtyFields === false) return false;
  if (dirtyFields === true) return true;
  if (Array.isArray(dirtyFields)) return dirtyFields.some((entry) => hasDirtyFieldsOutsideMedia(entry));
  if (typeof dirtyFields !== "object") return false;

  return Object.entries(dirtyFields as Record<string, unknown>).some(([key, value]) => {
    if (MEDIA_FIELD_KEYS.has(key as (typeof MEDIA_FIELD_KEYS extends Set<infer T> ? T : never))) {
      return false;
    }

    return hasDirtyFieldsOutsideMedia(value);
  });
}

export function ManageEventModal({ open, onOpenChange, event, onCreate, onUpdate }: ManageEventModalProps) {
  const isEditing = Boolean(event);
  const editValues = useMemo(() => (event ? toFormValues(event) : undefined), [event]);
  const autoAcronymRef = useRef("");
  const autoSlugRef = useRef("");

  const { track, getChanges, reset: resetImageTracking } = useImageFieldTracking();
  const steps = useMemo(() => createEventFormSteps(track), [track]);
  const syncEventMediaMutation = useSyncEventMediaMutation();

  const controller = useMultiStepForm({
    schema: eventCreateSchema,
    steps: steps,
    defaultValues: emptyDefaultValues,
    resetOnSuccessValues: isEditing ? undefined : emptyDefaultValues,
    values: editValues,
    // Blocks the submit button until something actually changed only meaningful while editing.
    requireDirtyToSubmit: isEditing,
    onSubmit: async (values): Promise<boolean> => {
      const { logo_url: _logoUrl, banner_url: _bannerUrl, gallery_urls: _galleryUrls, ...submitValues } = values;
      const payload: EventCreateSubmitI = submitValues;
      const hasNonMediaChanges = hasDirtyFieldsOutsideMedia(controller.form.formState.dirtyFields);

      let persistedEvent = event;

      if (!isEditing || hasNonMediaChanges) {
        const result = event
          ? await onUpdate?.(event.id, payload)
          : await onCreate?.(payload);

        if (!result) return false;

        persistedEvent = typeof result === "object" ? result : event;
      }

      if (!persistedEvent) return false;

      const logoChanges = getChanges("logo_url");
      const bannerChanges = getChanges("banner_url");
      const galleryChanges = getChanges("gallery_urls");

      const mediaChangesWereMade =
        logoChanges.added.length > 0 ||
        logoChanges.removed.length > 0 ||
        bannerChanges.added.length > 0 ||
        bannerChanges.removed.length > 0 ||
        galleryChanges.added.length > 0 ||
        galleryChanges.removed.length > 0;

      if (mediaChangesWereMade) {
        const syncedEvent = await syncEventMediaMutation.mutateAsync({
          eventId: persistedEvent.id,
          values,
          logoChanges,
          bannerChanges,
          galleryChanges,
        });

        if (!syncedEvent) return false;
      }

      resetImageTracking();
      return true;
    },
    onSubmitSuccess: () => onOpenChange(false),
  });

  const name = controller.form.watch("name");

  useEffect(() => {
    if (isEditing) return;

    const trimmedName = name.trim();
    if (!trimmedName) {
      const currentAcronym = controller.form.getValues("acronym");
      const currentSlug = controller.form.getValues("slug");

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
    const currentAcronym = controller.form.getValues("acronym");
    const currentSlug = controller.form.getValues("slug");

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
