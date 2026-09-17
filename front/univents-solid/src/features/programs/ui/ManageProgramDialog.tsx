import type { JSX } from "@solidjs/web";
import { createEffect, createSignal, untrack } from "solid-js";

import { Button, Dialog, Field, Input } from "@trieoh/ui-solid";
import type { ProgramCreateInput, ProgramI } from "../model";

export interface ManageProgramDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  program?: ProgramI | null;
  onSubmit: (data: ProgramCreateInput) => Promise<boolean>;
}

interface FormValues {
  kind: "activity" | "checkpoint";
  name: string;
  description: string;
  min_access_level: string;
  staff_only: boolean;
  banner_url: string;
  price: string;
}

function valuesOf(program?: ProgramI | null): FormValues {
  if (!program) {
    return {
      kind: "activity",
      name: "",
      description: "",
      min_access_level: "0",
      staff_only: false,
      banner_url: "",
      price: "",
    };
  }

  return {
    kind: program.kind,
    name: program.name,
    description: program.description ?? "",
    min_access_level: String(program.min_access_level ?? 0),
    staff_only: Boolean(program.staff_only),
    banner_url: program.banner_url ?? "",
    price: program.price ? String(program.price / 100) : "",
  };
}

export function ManageProgramDialog(props: ManageProgramDialogProps): JSX.Element {
  const [values, setValues] = createSignal<FormValues>(untrack(() => valuesOf(props.program)));
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<Partial<Record<keyof FormValues, string>>>({});

  createEffect(
    () => props.program,
    (p) => {
      setValues(valuesOf(p));
      setErrors({});
    },
  );

  const update = <K extends keyof FormValues>(key: K, value: FormValues[K]) => {
    setValues((prev) => ({ ...prev, [key]: value }));
    setErrors((prev) => ({ ...prev, [key]: undefined }));
  };

  const isEditing = () => Boolean(props.program);

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    const current = values();
    const errs: Partial<Record<keyof FormValues, string>> = {};

    if (!current.name.trim() || current.name.trim().length < 2) {
      errs.name = "Nome deve ter pelo menos 2 caracteres.";
    }

    if (Object.keys(errs).length > 0) {
      setErrors(errs);
      return;
    }

    const priceNumber = current.price.trim() ? parseFloat(current.price.replace(",", ".")) : undefined;
    const priceCents = priceNumber !== undefined && !isNaN(priceNumber) ? Math.round(priceNumber * 100) : undefined;

    setSubmitting(true);
    try {
      const ok = await props.onSubmit({
        kind: current.kind,
        name: current.name.trim(),
        description: current.description.trim() || undefined,
        min_access_level: current.min_access_level ? parseInt(current.min_access_level, 10) : 0,
        staff_only: current.staff_only,
        banner_url: current.banner_url.trim() || null,
        price: priceCents,
      });

      if (ok) {
        props.onOpenChange(false);
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={isEditing() ? "Editar programa" : "Novo programa"}
      description="Configure os detalhes da atividade ou checkpoint na programação."
    >
      <form onSubmit={handleSubmit} class="space-y-4">
        {/* Kind selector */}
        <div class="space-y-1.5">
          <label class="text-xs font-medium text-foreground">Tipo de programa</label>
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              onClick={() => update("kind", "activity")}
              class={`flex items-center justify-center rounded-lg border p-2.5 text-xs font-medium transition-colors cursor-pointer ${
                values().kind === "activity"
                  ? "border-primary bg-primary/10 text-primary"
                  : "border-border bg-card text-muted-foreground hover:bg-muted"
              }`}
            >
              Atividade / Palestra
            </button>
            <button
              type="button"
              onClick={() => update("kind", "checkpoint")}
              class={`flex items-center justify-center rounded-lg border p-2.5 text-xs font-medium transition-colors cursor-pointer ${
                values().kind === "checkpoint"
                  ? "border-amber-500 bg-amber-500/10 text-amber-600 dark:text-amber-400"
                  : "border-border bg-card text-muted-foreground hover:bg-muted"
              }`}
            >
              Checkpoint de presença
            </button>
          </div>
        </div>

        {/* Name */}
        <Field label="Nome da programação" error={errors().name}>
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              placeholder="ex: Workshop de Rust, Credenciamento..."
              value={values().name}
              onInput={(e) => update("name", e.currentTarget.value)}
              disabled={submitting()}
            />
          )}
        </Field>

        {/* Description */}
        <Field label="Descrição (opcional)" error={errors().description}>
          {(ids) => (
            <textarea
              id={ids.id}
              aria-describedby={ids.describedBy}
              rows={2}
              placeholder="Detalhes, palestrantes, tópicos abordados..."
              value={values().description}
              onInput={(e) => update("description", e.currentTarget.value)}
              disabled={submitting()}
              class="w-full resize-none rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
            />
          )}
        </Field>

        {/* Access level & Staff only */}
        <div class="grid grid-cols-2 gap-3">
          <Field label="Nível de acesso mínimo" error={errors().min_access_level}>
            {(ids) => (
              <Input
                id={ids.id}
                aria-describedby={ids.describedBy}
                type="number"
                min="0"
                value={values().min_access_level}
                onInput={(e) => update("min_access_level", e.currentTarget.value)}
                disabled={submitting()}
              />
            )}
          </Field>

          <div class="flex flex-col justify-center space-y-1 pt-4">
            <label class="flex items-center gap-2 cursor-pointer select-none">
              <input
                type="checkbox"
                checked={values().staff_only}
                onChange={(e) => update("staff_only", e.currentTarget.checked)}
                class="size-4 rounded border-border text-primary focus:ring-primary"
              />
              <span class="text-xs font-medium text-foreground">Apenas para equipe</span>
            </label>
            <p class="text-[10px] text-muted-foreground">Visível somente para staff</p>
          </div>
        </div>

        {/* Banner URL */}
        <Field label="URL da imagem de capa (opcional)" error={errors().banner_url}>
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              type="url"
              placeholder="https://..."
              value={values().banner_url}
              onInput={(e) => update("banner_url", e.currentTarget.value)}
              disabled={submitting()}
            />
          )}
        </Field>

        <div class="flex justify-end gap-2 pt-2 border-t border-border">
          <Button
            type="button"
            variant="outline"
            disabled={submitting()}
            onClick={() => props.onOpenChange(false)}
          >
            Cancelar
          </Button>
          <Button type="submit" disabled={submitting()}>
            {submitting() ? "Salvando..." : isEditing() ? "Salvar alterações" : "Criar programa"}
          </Button>
        </div>
      </form>
    </Dialog>
  );
}
