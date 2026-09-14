import { createEffect, createSignal, untrack } from "solid-js";
import { Button, Dialog, Field, Input } from "@trieoh/ui-solid";
import type { EditionI } from "../model";

export interface ManageEditionValues {
  name: string;
  slug: string;
  starts_at: string;
  ends_at: string;
  location_name?: string;
}

export interface ManageEditionDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  edition: EditionI | null;
  onSubmit: (values: ManageEditionValues) => Promise<boolean>;
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

function toInputDate(iso: string | undefined): string {
  if (!iso) return "";
  try {
    const d = new Date(iso);
    if (isNaN(d.getTime())) return "";
    return d.toISOString().slice(0, 16);
  } catch {
    return "";
  }
}

const emptyValues: ManageEditionValues = {
  name: "",
  slug: "",
  starts_at: "",
  ends_at: "",
  location_name: "",
};

const valuesOf = (edition: EditionI | null): ManageEditionValues =>
  edition
    ? {
      name: edition.name,
      slug: edition.slug,
      starts_at: toInputDate(edition.starts_at),
      ends_at: toInputDate(edition.ends_at),
      location_name: edition.location_name ?? "",
    }
    : { ...emptyValues };

function validate(values: ManageEditionValues): Partial<Record<keyof ManageEditionValues, string>> {
  const errors: Partial<Record<keyof ManageEditionValues, string>> = {};

  if (values.name.trim().length < 2) {
    errors.name = "Informe ao menos 2 caracteres.";
  }
  if (values.slug.trim().length < 2) {
    errors.slug = "Informe ao menos 2 caracteres.";
  }
  if (!values.starts_at) {
    errors.starts_at = "Informe a data de início.";
  }
  if (!values.ends_at) {
    errors.ends_at = "Informe a data de término.";
  } else if (values.starts_at && new Date(values.ends_at).getTime() <= new Date(values.starts_at).getTime()) {
    errors.ends_at = "O término deve ser posterior ao início.";
  }

  return errors;
}

export function ManageEditionDialog(props: ManageEditionDialogProps) {
  const [values, setValues] = createSignal(untrack(() => valuesOf(props.edition)));
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<Partial<Record<keyof ManageEditionValues, string>>>({});

  createEffect(
    () => props.edition,
    (ed) => {
      setValues(valuesOf(ed));
      setErrors({});
    },
  );

  const update = <K extends keyof ManageEditionValues>(key: K, value: ManageEditionValues[K]) => {
    const current = values();
    const next: ManageEditionValues = { ...current, [key]: value };

    if (key === "name") {
      if (!current.slug || current.slug === toSlug(current.name)) {
        next.slug = toSlug(String(value));
      }
    }

    setValues(next);
  };

  const submit = async (e: SubmitEvent) => {
    e.preventDefault();

    const raw = values();
    const current: ManageEditionValues = {
      name: raw.name.trim(),
      slug: raw.slug.trim(),
      starts_at: raw.starts_at ? new Date(raw.starts_at).toISOString() : "",
      ends_at: raw.ends_at ? new Date(raw.ends_at).toISOString() : "",
      location_name: raw.location_name?.trim() || undefined,
    };

    const found = validate(raw);
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
      title={props.edition ? "Editar edição" : "Criar nova edição"}
      description={
        props.edition
          ? "Atualize os detalhes da edição do evento."
          : "Defina o nome, identificador e as datas da nova edição."
      }
      footer={
        <>
          <Button variant="ghost" onClick={() => props.onOpenChange(false)}>
            Cancelar
          </Button>
          <Button type="submit" form="manage-edition-form" disabled={submitting()}>
            {props.edition ? "Salvar alterações" : "Criar edição"}
          </Button>
        </>
      }
    >
      <form id="manage-edition-form" class="flex flex-col gap-4" onSubmit={submit}>
        <Field label="Nome da Edição" required error={errors().name}>
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              aria-invalid={errors().name ? "true" : undefined}
              value={values().name}
              placeholder="Ex: Edição 2026"
              autocomplete="off"
              onInput={(e) => update("name", e.currentTarget.value)}
            />
          )}
        </Field>

        <Field label="Slug" required error={errors().slug} hint="Usado na URL pública da edição.">
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              aria-invalid={errors().slug ? "true" : undefined}
              value={values().slug}
              placeholder="edicao-2026"
              autocomplete="off"
              onInput={(e) => update("slug", e.currentTarget.value)}
            />
          )}
        </Field>

        <div class="grid gap-4 sm:grid-cols-2">
          <Field label="Início" required error={errors().starts_at}>
            {(ids) => (
              <Input
                id={ids.id}
                type="datetime-local"
                aria-describedby={ids.describedBy}
                aria-invalid={errors().starts_at ? "true" : undefined}
                value={values().starts_at}
                onInput={(e) => update("starts_at", e.currentTarget.value)}
              />
            )}
          </Field>

          <Field label="Término" required error={errors().ends_at}>
            {(ids) => (
              <Input
                id={ids.id}
                type="datetime-local"
                aria-describedby={ids.describedBy}
                aria-invalid={errors().ends_at ? "true" : undefined}
                value={values().ends_at}
                onInput={(e) => update("ends_at", e.currentTarget.value)}
              />
            )}
          </Field>
        </div>

        <Field label="Nome do Local (opcional)">
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              value={values().location_name ?? ""}
              placeholder="Ex: Centro de Convenções"
              autocomplete="off"
              onInput={(e) => update("location_name", e.currentTarget.value)}
            />
          )}
        </Field>
      </form>
    </Dialog>
  );
}
