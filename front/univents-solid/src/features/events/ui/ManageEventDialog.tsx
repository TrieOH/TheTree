import type { JSX } from "@solidjs/web";
import { Show, createEffect, createSignal, untrack } from "solid-js";
import {
  MultiStepDialog,
  type MultiStepField,
  type MultiStepItem,
} from "@trieoh/ui-solid";
import { toast } from "@/shared/ui/toast";
import type { EventI } from "../model";

export interface ManageEventValues {
  full_name: string;
  slug: string;
  acronym?: string | null;
  description?: string | null;
  contact_email?: string | null;
  logo_url?: string | null;
  banner_url?: string | null;
}

export interface ManageEventDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  event?: EventI | null;
  onSubmit: (values: ManageEventValues) => Promise<boolean>;
}

type FormInternalValues = {
  full_name: string;
  slug: string;
  acronym: string;
  description: string;
  contact_email: string;
  logo_url: string | null;
  banner_url: string | null;
};

/** Same slug rules as the React form, so the two apps generate identical URLs. */
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

const emptyValues: FormInternalValues = {
  full_name: "",
  slug: "",
  acronym: "",
  description: "",
  contact_email: "",
  logo_url: null,
  banner_url: null,
};

const valuesOf = (event: EventI | null | undefined): FormInternalValues =>
  event
    ? {
      full_name: event.full_name,
      slug: event.slug,
      acronym: event.acronym ?? "",
      description: event.description ?? "",
      contact_email: event.contact_email ?? "",
      logo_url: event.logo_url ?? null,
      banner_url: event.banner_url ?? null,
    }
    : { ...emptyValues };

const IDENTITY_FIELDS: MultiStepField<FormInternalValues>[] = [
  {
    name: "full_name",
    label: "Nome",
    placeholder: "Ex: Tech Summit 2026",
    required: true,
  },
  {
    name: "slug",
    label: "Slug",
    placeholder: "tech-summit-2026",
    hint: "Identificador na URL pública.",
    layout: "half",
    required: true,
  },
  {
    name: "acronym",
    label: "Sigla",
    placeholder: "TS26",
    hint: "Gerada a partir do nome.",
    layout: "half",
  },
  {
    kind: "textarea",
    name: "description",
    label: "Descrição",
    placeholder: "Escreva uma breve apresentação sobre o seu evento...",
    hint: "Conteúdo explicativo sobre o propósito e público do evento.",
    rows: 4,
  },
];

const CONTACT_FIELDS: MultiStepField<FormInternalValues>[] = [
  {
    kind: "email",
    name: "contact_email",
    label: "E-mail de contato",
    placeholder: "contato@evento.com",
    hint: "E-mail visível para dúvidas e suporte aos participantes.",
  },
];

export function ManageEventDialog(props: ManageEventDialogProps) {
  const [values, setValues] = createSignal<FormInternalValues>(untrack(() => valuesOf(props.event)));
  const [currentStep, setCurrentStep] = createSignal(0);
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<Partial<Record<keyof FormInternalValues, string>>>({});

  // Reset or re-sync values only when the dialog opens (reading proxies inside compute function to avoid STRICT_READ_UNTRACKED)
  createEffect(
    () => (props.open ? valuesOf(props.event) : null),
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

    if (key === "full_name") {
      if (!current.slug || current.slug === toSlug(current.full_name)) {
        next.slug = toSlug(String(value));
      }
      if (!current.acronym || current.acronym === toAcronym(current.full_name)) {
        next.acronym = toAcronym(String(value));
      }
    }

    setValues(next);

    // Clear errors for fields that become valid
    setErrors((prev) => {
      const updated = { ...prev };
      delete updated[key];

      if (next.full_name.trim().length >= 2) {
        delete updated.full_name;
      }
      if (next.slug.trim().length >= 2) {
        delete updated.slug;
      }
      if (!next.contact_email || next.contact_email.includes("@")) {
        delete updated.contact_email;
      }

      return updated;
    });
  };

  const validateStep = (stepIndex: number): boolean => {
    const v = values();
    const stepErrors: Partial<Record<keyof FormInternalValues, string>> = {};

    if (stepIndex === 0) {
      if (v.full_name.trim().length < 2) {
        stepErrors.full_name = "Informe ao menos 2 caracteres.";
      }
      if (v.slug.trim().length < 2) {
        stepErrors.slug = "Informe ao menos 2 caracteres.";
      }
    } else if (stepIndex === 1) {
      if (v.contact_email && !v.contact_email.includes("@")) {
        stepErrors.contact_email = "E-mail inválido.";
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
    const current: ManageEventValues = {
      full_name: raw.full_name.trim(),
      slug: raw.slug.trim(),
      acronym: raw.acronym.trim() ? raw.acronym.trim() : null,
      description: raw.description.trim() ? raw.description.trim() : null,
      contact_email: raw.contact_email.trim() ? raw.contact_email.trim() : null,
      logo_url: raw.logo_url ?? props.event?.logo_url ?? null,
      banner_url: raw.banner_url ?? props.event?.banner_url ?? null,
    };

    setSubmitting(true);
    try {
      const ok = await props.onSubmit(current);
      if (ok) {
        toast.success(
          props.event
            ? "Evento atualizado com sucesso!"
            : "Evento criado com sucesso!",
        );
        props.onOpenChange(false);
        return true;
      } else {
        toast.error("Não foi possível salvar o evento. Verifique os dados.");
        return false;
      }
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : "Erro ao salvar o evento.";
      toast.error(message);
      return false;
    } finally {
      setSubmitting(false);
    }
  };

  const steps: MultiStepItem<FormInternalValues>[] = [
    {
      id: "identidade",
      title: "Identidade",
      description: "Nome, link e detalhes",
      fields: IDENTITY_FIELDS,
    },
    {
      id: "contato",
      title: "Contato",
      description: "E-mail de atendimento",
      fields: CONTACT_FIELDS,
    },
    {
      id: "resumo",
      title: "Resumo",
      description: "Revisão e confirmação",
      summary: {
        badge: "Pronto para salvar",
        items: [
          { label: "Nome do evento", value: () => values().full_name },
          { label: "Sigla", value: () => values().acronym },
          {
            label: "Link público",
            value: () => (values().slug ? `/events/${values().slug}` : ""),
            href: () => (values().slug ? `/events/${values().slug}` : undefined),
            mono: true,
          },
          {
            label: "E-mail de contato",
            value: () => values().contact_email || "Não informado",
          },
          {
            label: "Descrição",
            value: () => values().description,
            fullWidth: true,
          },
        ],
        extra: () => {
          const hasVisual = () => !!(values().logo_url || values().banner_url);
          return (
            <Show when={hasVisual()}>
              <div class="col-span-1 sm:col-span-2 overflow-hidden rounded-xl border border-border/70 bg-card shadow-xs">
                {/* Capa com o banner do evento ou gradiente harmonioso */}
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
                        alt="Banner do evento"
                        class="size-full object-cover"
                      />
                    )}
                  </Show>
                </div>

                {/* Faixa inferior com o logo sobreposto no banner e identificação */}
                <div class="relative flex items-center gap-3 px-3.5 pb-3 pt-2 bg-card">
                  <div class="relative -mt-7 size-12 shrink-0 overflow-hidden rounded-xl border-2 border-card bg-card shadow-sm">
                    <Show
                      when={values().logo_url}
                      fallback={
                        <div class="flex size-full items-center justify-center bg-muted text-xs font-bold text-muted-foreground/60">
                          {values().acronym || "EV"}
                        </div>
                      }
                    >
                      {(logo) => (
                        <img
                          src={logo()}
                          alt="Logo do evento"
                          class="size-full object-cover"
                        />
                      )}
                    </Show>
                  </div>

                  <div class="min-w-0 flex-1">
                    <span class="block text-xs font-semibold text-foreground truncate">
                      {values().full_name || "Identidade visual"}
                    </span>
                    <span class="block text-[11px] text-muted-foreground truncate">
                      {values().logo_url && values().banner_url
                        ? "Logo e banner configurados"
                        : values().logo_url
                          ? "Logo configurado"
                          : "Banner configurado"}
                    </span>
                  </div>
                </div>
              </div>
            </Show>
          );
        },
      },
    },
  ];

  return (
    <MultiStepDialog<FormInternalValues>
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={props.event ? "Editar evento" : "Novo evento"}
      description={
        props.event
          ? "Atualize as informações do evento em etapas organizadas."
          : "Cadastre um novo evento na plataforma em etapas organizadas."
      }
      currentStep={currentStep()}
      onStepChange={setCurrentStep}
      steps={steps}
      values={values()}
      errors={errors()}
      onChange={(key, val) => update(key, val as FormInternalValues[typeof key])}
      formId="manage-event-form"
      loading={submitting()}
      submitLabel={props.event ? "Atualizar evento" : "Criar evento"}
      nextLabel="Continuar"
      onBeforeNext={(step) => validateStep(step)}
      onFormSubmit={submitCurrent}
      onSubmit={submitCurrent}
    />
  );
}
