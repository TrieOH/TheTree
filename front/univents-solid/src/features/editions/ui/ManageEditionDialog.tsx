import type { JSX } from "@solidjs/web";
import { Show, createEffect, createSignal, untrack } from "solid-js";
import {
  MultiStepDialog,
  type MultiStepField,
  type MultiStepItem,
} from "@trieoh/ui-solid";
import { toast } from "@/shared/ui/toast";
import type { EditionI } from "../model";

export interface ManageEditionValues {
  name: string;
  slug: string;
  starts_at: string;
  ends_at: string;
  location_name?: string | null;
  location_description?: string | null;
  tagline?: string | null;
  description?: string | null;
  contact_email?: string | null;
  registration_opens_at?: string | null;
  logo_url?: string | null;
  banner_url?: string | null;
}

export interface ManageEditionDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  edition: EditionI | null;
  onSubmit: (values: ManageEditionValues) => Promise<boolean>;
}

type FormInternalValues = {
  name: string;
  slug: string;
  starts_at: string;
  ends_at: string;
  location_name: string;
  location_description: string;
  tagline: string;
  description: string;
  contact_email: string;
  registration_opens_at: string;
  logo_url: string | null;
  banner_url: string | null;
};

function toSlug(value: string): string {
  return value
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

function toInputDate(iso: string | null | undefined): string {
  if (!iso) return "";
  try {
    const d = new Date(iso);
    if (isNaN(d.getTime())) return "";
    return d.toISOString().slice(0, 16);
  } catch {
    return "";
  }
}

function formatDisplayDate(iso: string | null | undefined): string {
  if (!iso) return "Não informado";
  try {
    const d = new Date(iso);
    if (isNaN(d.getTime())) return "Não informado";
    return new Intl.DateTimeFormat("pt-BR", {
      dateStyle: "medium",
      timeStyle: "short",
    }).format(d);
  } catch {
    return iso;
  }
}

const emptyValues: FormInternalValues = {
  name: "",
  slug: "",
  starts_at: "",
  ends_at: "",
  location_name: "",
  location_description: "",
  tagline: "",
  description: "",
  contact_email: "",
  registration_opens_at: "",
  logo_url: null,
  banner_url: null,
};

const valuesOf = (edition: EditionI | null): FormInternalValues =>
  edition
    ? {
      name: edition.name,
      slug: edition.slug,
      starts_at: toInputDate(edition.starts_at),
      ends_at: toInputDate(edition.ends_at),
      location_name: edition.location_name ?? "",
      location_description: edition.location_description ?? "",
      tagline: edition.tagline ?? "",
      description: edition.description ?? "",
      contact_email: edition.contact_email ?? "",
      registration_opens_at: toInputDate(edition.registration_opens_at),
      logo_url: edition.logo_url ?? null,
      banner_url: edition.banner_url ?? null,
    }
    : { ...emptyValues };

const EDITION_FIELDS: MultiStepField<FormInternalValues>[] = [
  {
    name: "name",
    label: "Nome da Edição",
    placeholder: "Ex: Edição 2026",
    required: true,
  },
  {
    name: "slug",
    label: "Slug",
    placeholder: "edicao-2026",
    hint: "Identificador na URL pública da edição.",
    required: true,
  },
];

const SCHEDULE_FIELDS: MultiStepField<FormInternalValues>[] = [
  {
    kind: "datetime-local",
    name: "starts_at",
    label: "Início",
    layout: "half",
    required: true,
  },
  {
    kind: "datetime-local",
    name: "ends_at",
    label: "Término",
    layout: "half",
    required: true,
  },
  {
    name: "location_name",
    label: "Nome do Local (opcional)",
    placeholder: "Ex: Centro de Convenções",
    hint: "Espaço ou endereço onde ocorrerá a edição.",
  },
];

export function ManageEditionDialog(props: ManageEditionDialogProps): JSX.Element {
  const [values, setValues] = createSignal<FormInternalValues>(
    untrack(() => valuesOf(props.edition)),
  );
  const [currentStep, setCurrentStep] = createSignal(0);
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<Partial<Record<keyof FormInternalValues, string>>>({});

  // Reset or re-sync values when dialog opens or edition changes
  createEffect(
    () => (props.open ? valuesOf(props.edition) : null),
    (computedValues) => {
      if (computedValues) {
        setValues(computedValues);
        setCurrentStep(0);
        setErrors({});
      }
    },
  );

  const update = <K extends keyof FormInternalValues>(key: K, value: FormInternalValues[K]) => {
    const current = values();
    const next: FormInternalValues = { ...current, [key]: value };

    if (key === "name") {
      if (!current.slug || current.slug === toSlug(current.name)) {
        next.slug = toSlug(String(value));
      }
    }

    setValues(next);

    // Clear errors when fields become valid
    setErrors((prev) => {
      const updated = { ...prev };
      delete updated[key];

      if (next.name.trim().length >= 2) {
        delete updated.name;
      }
      if (next.slug.trim().length >= 2) {
        delete updated.slug;
      }
      if (next.starts_at) {
        delete updated.starts_at;
      }
      if (
        next.ends_at &&
        (!next.starts_at || new Date(next.ends_at).getTime() > new Date(next.starts_at).getTime())
      ) {
        delete updated.ends_at;
      }

      return updated;
    });
  };

  const validateStep = (stepIndex: number): boolean => {
    const v = values();
    const stepErrors: Partial<Record<keyof FormInternalValues, string>> = {};

    if (stepIndex === 0) {
      if (v.name.trim().length < 2) {
        stepErrors.name = "Informe ao menos 2 caracteres.";
      }
      if (v.slug.trim().length < 2) {
        stepErrors.slug = "Informe ao menos 2 caracteres.";
      }
    } else if (stepIndex === 1) {
      if (!v.starts_at) {
        stepErrors.starts_at = "Informe a data de início.";
      }
      if (!v.ends_at) {
        stepErrors.ends_at = "Informe a data de término.";
      } else if (v.starts_at && new Date(v.ends_at).getTime() <= new Date(v.starts_at).getTime()) {
        stepErrors.ends_at = "O término deve ser posterior ao início.";
      }
    }

    setErrors((prev) => ({ ...prev, ...stepErrors }));
    return Object.keys(stepErrors).length === 0;
  };

  const submitCurrent = async (): Promise<boolean> => {
    const isStep0Valid = validateStep(0);
    if (!isStep0Valid) {
      setCurrentStep(0);
      return false;
    }

    const isStep1Valid = validateStep(1);
    if (!isStep1Valid) {
      setCurrentStep(1);
      return false;
    }

    const raw = values();
    const current: ManageEditionValues = {
      name: raw.name.trim(),
      slug: raw.slug.trim(),
      starts_at: raw.starts_at ? new Date(raw.starts_at).toISOString() : "",
      ends_at: raw.ends_at ? new Date(raw.ends_at).toISOString() : "",
      location_name: raw.location_name?.trim() || null,
      location_description: raw.location_description?.trim() || props.edition?.location_description || null,
      tagline: raw.tagline?.trim() || props.edition?.tagline || null,
      description: raw.description?.trim() || props.edition?.description || null,
      contact_email: raw.contact_email?.trim() || props.edition?.contact_email || null,
      registration_opens_at: raw.registration_opens_at
        ? new Date(raw.registration_opens_at).toISOString()
        : props.edition?.registration_opens_at ?? null,
      logo_url: raw.logo_url ?? props.edition?.logo_url ?? null,
      banner_url: raw.banner_url ?? props.edition?.banner_url ?? null,
    };

    setSubmitting(true);
    try {
      const ok = await props.onSubmit(current);
      if (ok) {
        toast.success(
          props.edition
            ? "Edição atualizada com sucesso!"
            : "Edição criada com sucesso!",
        );
        props.onOpenChange(false);
        return true;
      } else {
        toast.error("Não foi possível salvar a edição. Verifique os dados.");
        return false;
      }
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : "Erro ao salvar a edição.";
      toast.error(message);
      return false;
    } finally {
      setSubmitting(false);
    }
  };

  const steps: MultiStepItem<FormInternalValues>[] = [
    {
      id: "edicao",
      title: "Edição",
      description: "Nome e identificador",
      fields: EDITION_FIELDS,
    },
    {
      id: "cronograma",
      title: "Cronograma",
      description: "Datas e local do evento",
      fields: SCHEDULE_FIELDS,
    },
    {
      id: "resumo",
      title: "Resumo",
      description: "Revisão e confirmação",
      summary: {
        badge: () => (props.edition ? "Pronto para atualizar" : "Pronto para criar"),
        items: [
          { label: "Nome da edição", value: () => values().name },
          {
            label: "Slug",
            value: () => (values().slug ? `/${values().slug}` : ""),
            mono: true,
          },
          {
            label: "Início",
            value: () => formatDisplayDate(values().starts_at),
          },
          {
            label: "Término",
            value: () => formatDisplayDate(values().ends_at),
          },
          {
            label: "Local",
            value: () => values().location_name || "Não informado",
            fullWidth: true,
          },
        ],
        extra: () => (
          <Show when={values().logo_url || values().banner_url}>
            <div class="overflow-hidden rounded-xl border border-border/70 bg-card shadow-xs">
              <div class="relative h-20 w-full bg-muted overflow-hidden">
                <Show
                  when={values().banner_url}
                  fallback={
                    <div class="size-full bg-linear-to-r from-muted via-muted/60 to-muted/40" />
                  }
                >
                  {(banner) => (
                    <img
                      src={banner()}
                      alt="Banner da edição"
                      class="size-full object-cover"
                    />
                  )}
                </Show>
              </div>

              <div class="relative flex items-center gap-3 px-3.5 pb-3 pt-2 bg-card">
                <div class="relative -mt-7 size-12 shrink-0 overflow-hidden rounded-xl border-2 border-card bg-card shadow-sm">
                  <Show
                    when={values().logo_url}
                    fallback={
                      <div class="flex size-full items-center justify-center bg-muted text-xs font-bold text-muted-foreground/60">
                        ED
                      </div>
                    }
                  >
                    {(logo) => (
                      <img
                        src={logo()}
                        alt="Logo da edição"
                        class="size-full object-cover"
                      />
                    )}
                  </Show>
                </div>

                <div class="min-w-0 flex-1">
                  <span class="block text-xs font-semibold text-foreground truncate">
                    {values().name || "Identidade visual"}
                  </span>
                  <span class="block text-[11px] text-muted-foreground truncate">
                    {values().logo_url && values().banner_url
                      ? "Logo e banner vinculados"
                      : values().logo_url
                        ? "Logo vinculado"
                        : "Banner vinculado"}
                  </span>
                </div>
              </div>
            </div>
          </Show>
        ),
      },
    },
  ];

  return (
    <MultiStepDialog<FormInternalValues>
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={props.edition ? "Editar edição" : "Nova edição"}
      description={
        props.edition
          ? "Atualize as informações da edição em etapas organizadas."
          : "Cadastre uma nova edição para o evento em etapas organizadas."
      }
      currentStep={currentStep()}
      onStepChange={setCurrentStep}
      steps={steps}
      values={values()}
      errors={errors()}
      onChange={(key, val) => update(key, val as FormInternalValues[typeof key])}
      formId="manage-edition-form"
      loading={submitting()}
      submitLabel={props.edition ? "Salvar alterações" : "Criar edição"}
      nextLabel="Continuar"
      onBeforeNext={(step) => validateStep(step)}
      onFormSubmit={submitCurrent}
      onSubmit={submitCurrent}
    />
  );
}
