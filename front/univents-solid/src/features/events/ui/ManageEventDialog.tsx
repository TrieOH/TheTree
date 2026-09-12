import { createSignal, untrack } from "solid-js";

import { Button, Dialog, Field, Input, Textarea } from "@trieoh/ui-solid";

import type { EventI } from "../model";

export interface ManageEventValues {
  full_name: string;
  slug: string;
  acronym: string;
  description: string;
  contact_email: string;
}

export interface ManageEventDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** `null` creates; an event edits it. */
  event: EventI | null;
  onSubmit: (values: ManageEventValues) => Promise<boolean>;
}

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

const emptyValues: ManageEventValues = {
  full_name: "",
  slug: "",
  acronym: "",
  description: "",
  contact_email: "",
};

const valuesOf = (event: EventI | null): ManageEventValues =>
  event
    ? {
      full_name: event.full_name,
      slug: event.slug,
      acronym: event.acronym ?? "",
      description: event.description ?? "",
      contact_email: event.contact_email ?? "",
    }
    : { ...emptyValues };

function validate(values: ManageEventValues): Partial<Record<keyof ManageEventValues, string>> {
  const errors: Partial<Record<keyof ManageEventValues, string>> = {};

  if (values.full_name.trim().length < 2) {
    errors.full_name = "Informe ao menos 2 caracteres.";
  }
  if (values.slug.trim().length < 2) {
    errors.slug = "Informe ao menos 2 caracteres.";
  }
  if (values.contact_email && !values.contact_email.includes("@")) {
    errors.contact_email = "E-mail inválido.";
  }

  return errors;
}

/**
 * Create/edit form. Single step on purpose — the React app opens a multi-step
 * wizard here, and porting that widget is its own task; the fields, the slug
 * rules and the payload are the same.
 */
export function ManageEventDialog(props: ManageEventDialogProps) {
  // Seeded once, on mount: the caller remounts this component when it points at
  // another event (`<Show keyed>`), so there is no re-seed effect to fight the
  // user's typing. `untrack` documents that this is a snapshot, not a
  // subscription.
  const [values, setValues] = createSignal(untrack(() => valuesOf(props.event)));
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<Partial<Record<keyof ManageEventValues, string>>>({});

  /**
   * Derivation is decided by comparing with the current values instead of a
   * "touched" signal: the moment the field stops matching what the name would
   * produce, it is the user's and we leave it alone. One signal, one write.
   */
  const update = <K extends keyof ManageEventValues>(key: K, value: ManageEventValues[K]) => {
    const current = values();
    const next: ManageEventValues = { ...current, [key]: value };

    if (key === "full_name") {
      if (!current.slug || current.slug === toSlug(current.full_name)) {
        next.slug = toSlug(String(value));
      }
      if (!current.acronym || current.acronym === toAcronym(current.full_name)) {
        next.acronym = toAcronym(String(value));
      }
    }

    setValues(next);
  };

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();

    // Normalise at the form boundary: nothing leaves here with stray spaces.
    const raw = values();
    const current: ManageEventValues = {
      full_name: raw.full_name.trim(),
      slug: raw.slug.trim(),
      acronym: raw.acronym.trim(),
      description: raw.description.trim(),
      contact_email: raw.contact_email.trim(),
    };

    const found = validate(current);
    setErrors(found);
    if (Object.keys(found).length > 0) return;

    setSubmitting(true);
    try {
      await props.onSubmit(current);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={props.event ? "Editar evento" : "Criar novo evento"}
      description={
        props.event
          ? "Alterações entram em vigor imediatamente."
          : "O evento nasce como rascunho e pode ser publicado depois."
      }
      footer={
        <>
          <Button variant="ghost" onClick={() => props.onOpenChange(false)}>
            Cancelar
          </Button>
          <Button type="submit" form="manage-event-form" disabled={submitting()}>
            {props.event ? "Salvar alterações" : "Criar evento"}
          </Button>
        </>
      }
    >
      <form id="manage-event-form" class="flex flex-col gap-4" onSubmit={submit}>
        <Field label="Nome" required error={errors().full_name}>
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              aria-invalid={errors().full_name ? "true" : undefined}
              value={values().full_name}
              autocomplete="off"
              onInput={(event) => update("full_name", event.currentTarget.value)}
            />
          )}
        </Field>

        <div class="grid gap-4 sm:grid-cols-2">
          <Field label="Slug" required error={errors().slug} hint="Usado na URL pública.">
            {(ids) => (
              <Input
                id={ids.id}
                aria-describedby={ids.describedBy}
                aria-invalid={errors().slug ? "true" : undefined}
                value={values().slug}
                autocomplete="off"
                onInput={(event) => update("slug", event.currentTarget.value)}
              />
            )}
          </Field>

          <Field label="Sigla" hint="Gerada a partir do nome.">
            {(ids) => (
              <Input
                id={ids.id}
                aria-describedby={ids.describedBy}
                value={values().acronym}
                autocomplete="off"
                onInput={(event) => update("acronym", event.currentTarget.value)}
              />
            )}
          </Field>
        </div>

        <Field label="E-mail de contato" error={errors().contact_email}>
          {(ids) => (
            <Input
              id={ids.id}
              type="email"
              aria-describedby={ids.describedBy}
              aria-invalid={errors().contact_email ? "true" : undefined}
              value={values().contact_email}
              onInput={(event) => update("contact_email", event.currentTarget.value)}
            />
          )}
        </Field>

        <Field label="Descrição">
          {(ids) => (
            <Textarea
              id={ids.id}
              rows={3}
              aria-describedby={ids.describedBy}
              value={values().description}
              onInput={(event) => update("description", event.currentTarget.value)}
            />
          )}
        </Field>
      </form>
    </Dialog>
  );
}
