import type { JSX } from "@solidjs/web";
import { Button, Dialog, Field, Input, Label } from "@trieoh/ui-solid";
import { For, Show, createEffect, createMemo, createSignal, untrack } from "solid-js";
import { DEFAULT_BADGE_TEMPLATE } from "../default-template";
import type { BadgeElement, BadgeTemplate, BadgeTemplateCreate } from "../model";
import { BadgePreview } from "./BadgePreview";

export interface ManageBadgeTemplateValues {
  name: string;
  ticket_type_id: string | null;
  origin: "staff" | null;
  backgroundColor: string;
  textColor: string;
}

export interface ManageBadgeTemplateDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  template: BadgeTemplate | null;
  tickets: Array<{ id: string; name: string }>;
  onSubmit: (values: BadgeTemplateCreate) => Promise<boolean>;
}

const PRESET_COLORS = [
  { label: "Branco", bg: "#ffffff", text: "#0f172a" },
  { label: "Escuro", bg: "#0f172a", text: "#f8fafc" },
  { label: "Marinho", bg: "#1e293b", text: "#f1f5f9" },
  { label: "Índigo", bg: "#312e81", text: "#e0e7ff" },
  { label: "Esmeralda", bg: "#064e3b", text: "#ecfdf5" },
  { label: "Vinho", bg: "#881337", text: "#ffe4e6" },
];

const valuesOf = (template: BadgeTemplate | null): ManageBadgeTemplateValues => {
  if (!template) {
    return {
      name: "",
      ticket_type_id: null,
      origin: null,
      backgroundColor: "#ffffff",
      textColor: "#0f172a",
    };
  }
  const design = template.design_data ?? DEFAULT_BADGE_TEMPLATE.design_data;
  const nameEl = design.elements?.find((el): el is Extract<BadgeElement, { type: "text" }> => el.id === "default-name" && el.type === "text");
  return {
    name: template.name,
    ticket_type_id: template.ticket_type_id,
    origin: template.origin ?? null,
    backgroundColor: design.backgroundColor ?? "#ffffff",
    textColor: nameEl?.paragraphs?.[0]?.runs?.[0]?.color ?? "#0f172a",
  };
};

export function ManageBadgeTemplateDialog(
  props: ManageBadgeTemplateDialogProps,
): JSX.Element {
  const [values, setValues] = createSignal<ManageBadgeTemplateValues>(
    untrack(() => valuesOf(props.template)),
  );
  const [submitting, setSubmitting] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);

  createEffect(
    () => props.template,
    (t) => {
      setValues(valuesOf(t));
      setError(null);
    },
  );

  const isEditing = () => Boolean(props.template);

  // Derived preview badge
  const previewBadge = createMemo<BadgeTemplate>(() => {
    const v = values();
    const base = props.template?.design_data
      ? JSON.parse(JSON.stringify(props.template.design_data))
      : JSON.parse(JSON.stringify(DEFAULT_BADGE_TEMPLATE.design_data));

    base.backgroundColor = v.backgroundColor;

    // Update text runs colors if matching defaults
    if (base.elements) {
      for (const el of base.elements) {
        if (el.type === "text" && el.paragraphs) {
          for (const p of el.paragraphs) {
            for (const r of p.runs) {
              if (el.id === "default-name") {
                r.color = v.textColor;
              } else if (el.id === "default-event" || el.id === "default-ticket") {
                r.color = v.backgroundColor === "#ffffff" ? "#64748b" : "#94a3b8";
              }
            }
          }
        }
      }
    }

    return {
      id: props.template?.id ?? "preview-id",
      edition_id: props.template?.edition_id ?? "preview-edition",
      name: v.name || "Pré-visualização do Crachá",
      ticket_type_id: v.ticket_type_id,
      origin: v.origin,
      design_data: base,
      created_at: new Date().toISOString(),
      updated_at: null,
    };
  });

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    const v = values();
    if (v.name.trim().length < 3) {
      setError("O nome deve ter pelo menos 3 caracteres.");
      return;
    }

    setSubmitting(true);
    try {
      const design = previewBadge().design_data;
      const success = await props.onSubmit({
        name: v.name.trim(),
        ticket_type_id: v.origin === "staff" ? null : v.ticket_type_id,
        origin: v.origin,
        design_data: design,
      });
      if (success) {
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
      title={isEditing() ? "Editar template de crachá" : "Novo template de crachá"}
      description="Configure as informações básicas e o visual do modelo de crachá."
      class="sm:max-w-xl"
    >
      <form onSubmit={handleSubmit} class="space-y-4 pt-2">
        <Field label="Nome do template" required error={error() ?? undefined}>
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              type="text"
              placeholder="Ex: Crachá Participante VIP"
              value={values().name}
              onInput={(e) => {
                const val = e.currentTarget.value;
                setValues((prev) => ({ ...prev, name: val }));
                setError(null);
              }}
              required
              class={error() ? "border-destructive" : ""}
            />
          )}
        </Field>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {/* Tipo de Destinatário */}
          <div class="space-y-1.5">
            <Label>Destinatário</Label>
            <div class="flex gap-2">
              <button
                type="button"
                onClick={() =>
                  setValues((prev) => ({ ...prev, origin: null }))
                }
                class={`flex-1 rounded-lg border px-3 py-2 text-xs font-medium transition-colors ${
                  values().origin !== "staff"
                    ? "border-primary bg-primary/10 text-primary"
                    : "border-border bg-card text-muted-foreground hover:bg-muted"
                }`}
              >
                Participante
              </button>
              <button
                type="button"
                onClick={() =>
                  setValues((prev) => ({
                    ...prev,
                    origin: "staff",
                    ticket_type_id: null,
                  }))
                }
                class={`flex-1 rounded-lg border px-3 py-2 text-xs font-medium transition-colors ${
                  values().origin === "staff"
                    ? "border-primary bg-primary/10 text-primary"
                    : "border-border bg-card text-muted-foreground hover:bg-muted"
                }`}
              >
                Staff (Equipe)
              </button>
            </div>
          </div>

          {/* Ingresso vinculado (se participante) */}
          <div class="space-y-1.5">
            <Label for="badge-ticket">Ingresso vinculado</Label>
            <select
              id="badge-ticket"
              disabled={values().origin === "staff"}
              value={values().ticket_type_id ?? ""}
              onChange={(e) => {
                const val = e.currentTarget.value;
                setValues((prev) => ({
                  ...prev,
                  ticket_type_id: val ? val : null,
                }));
              }}
              class="h-9 w-full rounded-md border border-border bg-background px-3 text-xs text-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="">Padrão da edição (todos)</option>
              <For each={props.tickets}>
                {(ticket) => (
                  <option value={ticket.id}>{ticket.name}</option>
                )}
              </For>
            </select>
          </div>
        </div>

        {/* Color Presets */}
        <div class="space-y-1.5">
          <Label>Cor de fundo do crachá</Label>
          <div class="flex flex-wrap items-center gap-2">
            <For each={PRESET_COLORS}>
              {(preset) => (
                <button
                  type="button"
                  title={preset.label}
                  onClick={() =>
                    setValues((prev) => ({
                      ...prev,
                      backgroundColor: preset.bg,
                      textColor: preset.text,
                    }))
                  }
                  class={`size-7 rounded-full border shadow-xs transition-transform hover:scale-110 ${
                    values().backgroundColor.toLowerCase() ===
                    preset.bg.toLowerCase()
                      ? "ring-2 ring-primary ring-offset-2"
                      : "border-border"
                  }`}
                  style={{ "background-color": preset.bg }}
                />
              )}
            </For>
            <div class="ml-2 flex items-center gap-2">
              <input
                type="color"
                aria-label="Cor customizada"
                value={values().backgroundColor}
                onInput={(e) => {
                  const val = e.currentTarget.value;
                  setValues((prev) => ({ ...prev, backgroundColor: val }));
                }}
                class="size-7 cursor-pointer rounded-md border border-border bg-transparent p-0.5"
              />
              <span class="text-xs text-muted-foreground uppercase">
                {values().backgroundColor}
              </span>
            </div>
          </div>
        </div>

        {/* Live Preview Container */}
        <div class="space-y-1.5">
          <Label>Pré-visualização</Label>
          <div class="flex h-44 w-full items-center justify-center overflow-hidden rounded-xl border border-border bg-muted/40 p-3">
            <BadgePreview
              badge={previewBadge()}
              contain
              showVariables
              eventName="Nome do Evento"
              editionName="Edição 2026"
              ticketName={
                values().origin === "staff"
                  ? "Staff / Organização"
                  : props.tickets.find((t) => t.id === values().ticket_type_id)
                      ?.name ?? "Padrão da Edição"
              }
              participantName="Nome do Participante"
              location="Local do Evento"
              class="max-h-full max-w-full"
            />
          </div>
        </div>

        {/* Action Buttons */}
        <div class="flex w-full items-center justify-end gap-2 pt-2">
          <Button
            type="button"
            variant="outline"
            disabled={submitting()}
            onClick={() => props.onOpenChange(false)}
          >
            Cancelar
          </Button>
          <Button type="submit" disabled={submitting()}>
            <Show when={submitting()} fallback={isEditing() ? "Salvar alterações" : "Criar template"}>
              Salvando...
            </Show>
          </Button>
        </div>
      </form>
    </Dialog>
  );
}
