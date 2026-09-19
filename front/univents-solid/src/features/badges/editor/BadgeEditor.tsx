import type { JSX } from "@solidjs/web";
import { Link, useNavigate } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { Button, Input, Label, cn, buttonVariants } from "@trieoh/ui-solid";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
} from "solid-js";

import ArrowDownIcon from "~icons/lucide/arrow-down";
import ArrowLeftIcon from "~icons/lucide/arrow-left";
import ArrowUpIcon from "~icons/lucide/arrow-up";
import BadgeCheckIcon from "~icons/lucide/badge-check";
import ImageIcon from "~icons/lucide/image";
import Loader2Icon from "~icons/lucide/loader-2";
import MonitorIcon from "~icons/lucide/monitor";
import SaveIcon from "~icons/lucide/save";
import Trash2Icon from "~icons/lucide/trash-2";
import TypeIcon from "~icons/lucide/type";

import { badgeTemplateQueryOptions } from "@/features/badges/api";
import {
  useCreateBadgeTemplateMutation,
  useUpdateBadgeTemplateMutation,
} from "@/features/badges/api/mutations";
import { DEFAULT_BADGE_TEMPLATE } from "@/features/badges/default-template";
import {
  type BadgeDesign,
  type BadgeElement,
  type BadgeTemplateCreate,
  MIN_BADGE_CANVAS_SIZE_MM,
  MIN_BADGE_CANVAS_SIZE_PX,
  badgeMmToPx,
  badgePxToMm,
  badgeTemplateCreateSchema,
} from "@/features/badges/model";
import { resizeBadgeDesign } from "@/features/badges/model/resize-design";
import {
  DEFAULT_EDITOR_FONT,
  DEFAULT_EDITOR_TEXT_COLOR,
  RichTextToolbar,
  ToolbarCombobox,
  type RichTextController,
  type TextElementAdapter,
  type TextSelectionStyles,
} from "@/features/editor";
import { allTicketsQueryOptions } from "@/features/tickets/api";
import { toast } from "@/shared/ui/toast";

import { BadgeCanvas } from "./BadgeCanvas";
import { uploadBadgeAssets } from "./upload-assets";

const ArrowDown = ArrowDownIcon as unknown as (props: { class?: string }) => JSX.Element;
const ArrowLeft = ArrowLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const ArrowUp = ArrowUpIcon as unknown as (props: { class?: string }) => JSX.Element;
const BadgeCheck = BadgeCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const ImageLucide = ImageIcon as unknown as (props: { class?: string }) => JSX.Element;
const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Monitor = MonitorIcon as unknown as (props: { class?: string }) => JSX.Element;
const Save = SaveIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash2 = Trash2Icon as unknown as (props: { class?: string }) => JSX.Element;
const TypeLucide = TypeIcon as unknown as (props: { class?: string }) => JSX.Element;

export const VARIABLES: readonly [
  variable: string,
  label: string,
  description: string,
][] = [
  [
    "{{participant_name}}",
    "Nome civil do participante",
    "Nome civil informado no perfil do participante",
  ],
  ["{{event_name}}", "Nome do evento", "Nome do evento"],
  ["{{edition_name}}", "Nome da edição", "Nome da edição do evento"],
  ["{{ticket_name}}", "Tipo de ingresso", "Ingresso associado ao participante"],
  ["{{location}}", "Local", "Local informado na edição"],
];

const DEFAULT_PREVIEW_VALUES: Record<string, string> = {
  participant_name: "Maria da Silva",
  event_name: "Tech Summit 2026",
  edition_name: "Edição 2026",
  ticket_name: "Ingresso VIP",
  location: "Centro de Convenções",
};

const BADGE_SIZE_PRESETS = [
  {
    value: "portrait",
    label: "Vertical (54 × 85 mm)",
    width: 204,
    height: 321,
  },
  {
    value: "landscape",
    label: "Horizontal (85 × 54 mm)",
    width: 321,
    height: 204,
  },
  {
    value: "square",
    label: "Quadrado (60 × 60 mm)",
    width: 227,
    height: 227,
  },
] as const;

function uid() {
  return typeof crypto !== "undefined" && crypto.randomUUID
    ? crypto.randomUUID()
    : Math.random().toString(36).slice(2, 11);
}

function layerLabel(element: BadgeElement) {
  if (element.type === "image") return "Imagem";
  if (element.type === "qr") return "QR Code";
  const text = element.paragraphs
    .flatMap((paragraph) => paragraph.runs.map((run) => run.text))
    .join("")
    .trim();
  return text || "Texto vazio";
}

function readImage(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(String(reader.result));
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

export interface BadgeEditorProps {
  eventId: string;
  editionId: string;
  templateId?: string;
  duplicate?: boolean;
}

export function BadgeEditor(props: BadgeEditorProps): JSX.Element {
  const navigate = useNavigate();
  const createMutation = useCreateBadgeTemplateMutation();
  const updateMutation = useUpdateBadgeTemplateMutation();

  const templateQuery = useQuery(() => ({
    ...badgeTemplateQueryOptions(props.templateId ?? ""),
    enabled: Boolean(props.templateId),
  }));

  const ticketsQuery = useQuery(() => allTicketsQueryOptions(props.editionId));

  const [draft, setDraft] = createSignal<BadgeTemplateCreate>(
    structuredClone(DEFAULT_BADGE_TEMPLATE),
    { ownedWrite: true } as any,
  );
  const [selectedId, setSelectedId] = createSignal<string | null>(null, {
    ownedWrite: true,
  } as any);
  const [uploading, setUploading] = createSignal(false, {
    ownedWrite: true,
  } as any);
  const [previewValues, setPreviewValues] = createSignal<Record<string, string>>(
    DEFAULT_PREVIEW_VALUES,
    { ownedWrite: true } as any,
  );

  const [textController, setTextController] =
    createSignal<RichTextController | null>(null, { ownedWrite: true } as any);
  const [textSelectionStyles, setTextSelectionStyles] =
    createSignal<TextSelectionStyles | null>(null, { ownedWrite: true } as any);

  createEffect(
    () => ({
      data: templateQuery().data,
      duplicate: props.duplicate,
    }),
    ({ data, duplicate }) => {
      if (data) {
        const copy = structuredClone(data);
        setDraft({
          ...copy,
          name: duplicate ? `${copy.name} (cópia)` : copy.name,
        });
      }
    },
  );

  const selected = createMemo(() => {
    return selectedId()
      ? (draft().design_data.elements.find((item) => item.id === selectedId()) ?? null)
      : null;
  });

  const updateDesign = (changes: Partial<BadgeDesign>) => {
    setDraft((curr) => ({
      ...curr,
      design_data: { ...curr.design_data, ...changes },
    }));
  };

  const updateElement = (id: string, changes: Partial<BadgeElement>) => {
    setDraft((curr) => ({
      ...curr,
      design_data: {
        ...curr.design_data,
        elements: curr.design_data.elements.map((item) =>
          item.id === id ? ({ ...item, ...changes } as BadgeElement) : item,
        ),
      },
    }));
  };

  const resizeCanvas = (canvas: { width: number; height: number }) => {
    if (
      canvas.width < MIN_BADGE_CANVAS_SIZE_PX ||
      canvas.height < MIN_BADGE_CANVAS_SIZE_PX
    ) {
      return;
    }
    updateDesign(resizeBadgeDesign(draft().design_data, canvas));
  };

  const textAdapter: TextElementAdapter = {
    updateParagraphs: (id, paragraphs) => updateElement(id, { paragraphs }),
    setController: setTextController,
    setSelectionStyles: setTextSelectionStyles,
    stopEditing: () => undefined,
  };

  const addElement = (element: BadgeElement) => {
    updateDesign({ elements: [...draft().design_data.elements, element] });
    setSelectedId(element.id);
  };

  const moveElement = (id: string, direction: -1 | 1) => {
    const elements = [...draft().design_data.elements];
    const index = elements.findIndex((item) => item.id === id);
    const target = index + direction;
    if (index < 0 || target < 0 || target >= elements.length) return;
    const temp = elements[index]!;
    elements[index] = elements[target]!;
    elements[target] = temp;
    updateDesign({ elements });
  };

  const deleteElement = (id: string) => {
    updateDesign({
      elements: draft().design_data.elements.filter((item) => item.id !== id),
    });
    setSelectedId(null);
  };

  // Fallback controller when a text element is selected but not inline-editing
  const fallbackController = createMemo<RichTextController | null>(() => {
    const initialEl = selected();
    if (!initialEl || initialEl.type !== "text") return null;
    const elementId = initialEl.id;

    return {
      elementId,
      commit: () => {},
      toggleBold: () => {
        const item = selected();
        if (!item || item.type !== "text") return;
        const anyBold = item.paragraphs.some((p) => p.runs.some((r) => r.bold));
        updateElement(item.id, {
          paragraphs: item.paragraphs.map((p) => ({
            ...p,
            runs: p.runs.map((r) => ({ ...r, bold: !anyBold })),
          })),
        });
      },
      toggleItalic: () => {
        const item = selected();
        if (!item || item.type !== "text") return;
        const anyItalic = item.paragraphs.some((p) => p.runs.some((r) => r.italic));
        updateElement(item.id, {
          paragraphs: item.paragraphs.map((p) => ({
            ...p,
            runs: p.runs.map((r) => ({ ...r, italic: !anyItalic })),
          })),
        });
      },
      toggleUnderline: () => {
        const item = selected();
        if (!item || item.type !== "text") return;
        const anyUnderline = item.paragraphs.some((p) =>
          p.runs.some((r) => r.underline),
        );
        updateElement(item.id, {
          paragraphs: item.paragraphs.map((p) => ({
            ...p,
            runs: p.runs.map((r) => ({ ...r, underline: !anyUnderline })),
          })),
        });
      },
      setAlign: (align) => {
        const item = selected();
        if (!item || item.type !== "text") return;
        updateElement(item.id, {
          paragraphs: item.paragraphs.map((p) => ({ ...p, align })),
        });
      },
      setLineHeight: (lineHeight) => {
        const item = selected();
        if (!item || item.type !== "text") return;
        updateElement(item.id, {
          paragraphs: item.paragraphs.map((p) => ({ ...p, lineHeight })),
        });
      },
      setColor: (color) => {
        const item = selected();
        if (!item || item.type !== "text") return;
        updateElement(item.id, {
          paragraphs: item.paragraphs.map((p) => ({
            ...p,
            runs: p.runs.map((r) => ({ ...r, color })),
          })),
        });
      },
      setFontSize: (fontSize) => {
        const item = selected();
        if (!item || item.type !== "text") return;
        updateElement(item.id, {
          paragraphs: item.paragraphs.map((p) => ({
            ...p,
            runs: p.runs.map((r) => ({ ...r, fontSize })),
          })),
        });
      },
      setFontFamily: (fontFamily) => {
        const item = selected();
        if (!item || item.type !== "text") return;
        updateElement(item.id, {
          paragraphs: item.paragraphs.map((p) => ({
            ...p,
            runs: p.runs.map((r) => ({ ...r, fontFamily })),
          })),
        });
      },
      insertText: (text) => {
        const item = selected();
        if (!item || item.type !== "text") return;
        const lastIdx = item.paragraphs.length - 1;
        updateElement(item.id, {
          paragraphs: item.paragraphs.map((p, idx) =>
            idx === lastIdx
              ? {
                  ...p,
                  runs: [
                    ...p.runs,
                    {
                      text: ` ${text} `,
                      fontSize: p.runs[0]?.fontSize ?? 18,
                      fontFamily: p.runs[0]?.fontFamily ?? DEFAULT_EDITOR_FONT,
                      color: p.runs[0]?.color ?? DEFAULT_EDITOR_TEXT_COLOR,
                      bold: p.runs[0]?.bold ?? false,
                      italic: p.runs[0]?.italic ?? false,
                      underline: p.runs[0]?.underline ?? false,
                    },
                  ],
                }
              : p,
          ),
        });
      },
    };
  });

  const fallbackStyles = createMemo<TextSelectionStyles | null>(() => {
    const el = selected();
    if (!el || el.type !== "text") return null;
    const firstRun = el.paragraphs[0]?.runs[0];
    return {
      bold: Boolean(firstRun?.bold),
      italic: Boolean(firstRun?.italic),
      underline: Boolean(firstRun?.underline),
      align: el.paragraphs[0]?.align ?? "left",
      lineHeight: el.paragraphs[0]?.lineHeight ?? 1.25,
      color: firstRun?.color ?? DEFAULT_EDITOR_TEXT_COLOR,
      fontSize: firstRun?.fontSize ?? 18,
      fontFamily: firstRun?.fontFamily ?? DEFAULT_EDITOR_FONT,
    };
  });

  const effectiveController = () => textController() ?? fallbackController();
  const effectiveStyles = () => textSelectionStyles() ?? fallbackStyles();

  const isSaving = () =>
    createMutation.result().status === "pending" ||
    updateMutation.result().status === "pending" ||
    uploading();

  const save = async () => {
    const parsed = badgeTemplateCreateSchema.safeParse(draft());
    if (!parsed.success) {
      toast.error(parsed.error.issues[0]?.message ?? "Template inválido");
      return;
    }

    setUploading(true);
    try {
      const data = await uploadBadgeAssets(
        parsed.data,
        props.eventId,
        props.editionId,
      );
      if (props.templateId && !props.duplicate) {
        await updateMutation.mutateAsync({
          templateId: props.templateId,
          editionId: props.editionId,
          data,
        });
      } else {
        await createMutation.mutateAsync({
          editionId: props.editionId,
          data,
        });
      }

      toast.success("Template de crachá salvo com sucesso");
      void navigate({
        to: "/admin/events/$eventId/editions/$editionId/badges",
        params: { eventId: props.eventId, editionId: props.editionId },
      });
    } catch {
      toast.error("Não foi possível salvar o template de crachá");
    } finally {
      setUploading(false);
    }
  };

  const tickets = () => ticketsQuery().data ?? [];

  return (
    <div class="h-dvh min-h-0 bg-background text-foreground">
      {/* Mobile unsupported notice */}
      <div class="flex h-full flex-col items-center justify-center gap-3 px-6 text-center lg:hidden!">
        <Monitor class="size-10 text-muted-foreground" />
        <p class="max-w-xs text-sm text-muted-foreground">
          O editor de crachás foi desenvolvido para telas maiores. Abra esta
          página em um computador para editar o template.
        </p>
      </div>

      {/* Desktop Workspace */}
      <div class="hidden h-full min-h-0 flex-col lg:flex">
        {/* Header */}
        <header class="flex h-16 shrink-0 items-center justify-between border-b border-muted bg-card px-4 shadow-xs">
          <div class="flex items-center gap-3">
            <Link
              to="/admin/events/$eventId/editions/$editionId/badges"
              params={{ eventId: props.eventId, editionId: props.editionId }}
              class="inline-flex size-8 items-center justify-center rounded-lg border border-muted"
            >
              <ArrowLeft class="size-4" />
            </Link>
            <strong class="text-sm">Editor de crachás</strong>
          </div>
          <Button onClick={() => void save()} disabled={isSaving()}>
            <Show when={isSaving()} fallback={<Save class="size-4" />}>
              <Loader2 class="size-4 animate-spin" />
            </Show>
            Salvar
          </Button>
        </header>

        {/* Workspace Body */}
        <div class="flex min-h-0 flex-1 overflow-x-auto">
          {/* Left Sidebar */}
          <aside class="w-72 shrink-0 overflow-y-auto border-r border-muted bg-card p-4">
            <div class="space-y-2">
              <Label for="badge-name">Nome do template</Label>
              <Input
                id="badge-name"
                value={draft().name}
                onInput={(event) =>
                  setDraft((curr) => ({
                    ...curr,
                    name: event.currentTarget.value,
                  }))
                }
              />
            </div>

            <div class="mt-5 space-y-2">
              <Label for="badge-ticket">Ingresso associado</Label>
              <ToolbarCombobox
                value={draft().ticket_type_id ?? ""}
                options={[
                  { value: "", label: "Padrão da edição" },
                  ...tickets().map((ticket) => ({
                    value: ticket.id,
                    label: ticket.name,
                  })),
                ]}
                placeholder="Selecione o ingresso"
                searchPlaceholder="Buscar ingresso..."
                class="w-full"
                triggerClass="h-9"
                disabled={Boolean(props.templateId) && !props.duplicate}
                onChange={(value) =>
                  setDraft((curr) => ({
                    ...curr,
                    ticket_type_id: value || null,
                  }))
                }
              />
              <Show when={props.templateId && !props.duplicate}>
                <p class="text-xs text-muted-foreground">
                  O ingresso associado não pode ser alterado após a criação.
                </p>
              </Show>
            </div>

            <div class="mt-5 space-y-2">
              <Label for="badge-origin">Origem do crachá</Label>
              <ToolbarCombobox
                value={draft().origin ?? ""}
                options={[
                  { value: "", label: "Participante / padrão" },
                  { value: "staff", label: "Equipe (staff)" },
                ]}
                placeholder="Selecione a origem"
                class="w-full"
                triggerClass="h-9"
                disabled={Boolean(props.templateId) && !props.duplicate}
                onChange={(value) =>
                  setDraft((curr) => ({
                    ...curr,
                    origin: value === "staff" ? "staff" : null,
                  }))
                }
              />
              <Show when={props.templateId && !props.duplicate}>
                <p class="text-xs text-muted-foreground">
                  A origem não pode ser alterada após a criação.
                </p>
              </Show>
            </div>

            <div class="mt-6">
              <Label>Adicionar</Label>
              <div class="mt-2 grid grid-cols-2 gap-2">
                <Button
                  variant="outline"
                  class="h-16 flex-col gap-1"
                  onClick={() =>
                    addElement({
                      id: uid(),
                      type: "text",
                      x: draft().design_data.canvas.width * 0.15,
                      y: draft().design_data.canvas.height * 0.4,
                      width: draft().design_data.canvas.width * 0.7,
                      height: Math.max(
                        32,
                        draft().design_data.canvas.height * 0.15,
                      ),
                      paragraphs: [
                        {
                          align: "center",
                          lineHeight: 1.25,
                          runs: [
                            {
                              text: "Novo texto",
                              fontSize: 18,
                              fontFamily: "Inter, sans-serif",
                              color: "#0f172a",
                              bold: false,
                              italic: false,
                              underline: false,
                            },
                          ],
                        },
                      ],
                    })
                  }
                >
                  <TypeLucide class="size-4" />
                  Texto
                </Button>
                <label class="inline-flex h-16 cursor-pointer flex-col items-center justify-center gap-1 rounded-lg border border-border bg-background text-sm font-medium hover:bg-muted hover:text-foreground">
                  <ImageLucide class="size-4" />
                  Imagem
                  <input
                    type="file"
                    accept="image/*"
                    class="hidden"
                    onChange={async (event) => {
                      const file = event.target.files?.[0];
                      if (file) {
                        addElement({
                          id: uid(),
                          type: "image",
                          x: draft().design_data.canvas.width * 0.2,
                          y: draft().design_data.canvas.height * 0.2,
                          width: draft().design_data.canvas.width * 0.6,
                          height: draft().design_data.canvas.height * 0.6,
                          src: await readImage(file),
                          fit: "contain",
                          radius: 0,
                          opacity: 1,
                        });
                      }
                      event.target.value = "";
                    }}
                  />
                </label>
              </div>
            </div>

            <div class="mt-6 space-y-2">
              <Label>Fundo</Label>
              <div class="flex gap-2">
                <input
                  type="color"
                  class="h-9 w-14 cursor-pointer rounded-md border border-input bg-background p-1 disabled:opacity-40"
                  value={
                    draft().design_data.backgroundColor === "transparent"
                      ? "#ffffff"
                      : draft().design_data.backgroundColor
                  }
                  disabled={draft().design_data.backgroundColor === "transparent"}
                  onInput={(event) =>
                    updateDesign({ backgroundColor: event.currentTarget.value })
                  }
                />
                <label class="inline-flex h-9 flex-1 cursor-pointer items-center justify-center rounded-md border border-muted text-sm hover:bg-muted">
                  Imagem de fundo
                  <input
                    type="file"
                    accept="image/*"
                    class="hidden"
                    onChange={async (event) => {
                      const file = event.target.files?.[0];
                      if (file) {
                        updateDesign({ background: await readImage(file) });
                      }
                    }}
                  />
                </label>
              </div>
              <label class="flex items-center gap-2 text-sm text-muted-foreground cursor-pointer">
                <input
                  type="checkbox"
                  checked={draft().design_data.backgroundColor === "transparent"}
                  onChange={(event) =>
                    updateDesign({
                      backgroundColor: event.currentTarget.checked
                        ? "transparent"
                        : "#ffffff",
                    })
                  }
                />
                Fundo transparente
              </label>
              <Show when={draft().design_data.background}>
                <Button
                  variant="ghost"
                  class="w-full"
                  onClick={() => updateDesign({ background: null })}
                >
                  Remover imagem
                </Button>
              </Show>
            </div>

            <div class="mt-6 space-y-2">
              <Label>Tamanho (mm)</Label>
              <ToolbarCombobox
                value={
                  BADGE_SIZE_PRESETS.find(
                    (preset) =>
                      preset.width === draft().design_data.canvas.width &&
                      preset.height === draft().design_data.canvas.height,
                  )?.value
                }
                options={BADGE_SIZE_PRESETS}
                placeholder="Predefinições"
                class="w-full"
                triggerClass="h-10"
                onChange={(value) => {
                  const preset = BADGE_SIZE_PRESETS.find(
                    (item) => item.value === value,
                  );
                  if (preset) {
                    resizeCanvas({
                      width: preset.width,
                      height: preset.height,
                    });
                  }
                }}
              />
              <div class="grid grid-cols-2 gap-2">
                <Input
                  type="number"
                  min={MIN_BADGE_CANVAS_SIZE_MM}
                  step={0.1}
                  value={badgePxToMm(draft().design_data.canvas.width)}
                  onChange={(e) =>
                    resizeCanvas({
                      ...draft().design_data.canvas,
                      width: badgeMmToPx(Number(e.currentTarget.value)),
                    })
                  }
                />
                <Input
                  type="number"
                  min={MIN_BADGE_CANVAS_SIZE_MM}
                  step={0.1}
                  value={badgePxToMm(draft().design_data.canvas.height)}
                  onChange={(e) =>
                    resizeCanvas({
                      ...draft().design_data.canvas,
                      height: badgeMmToPx(Number(e.currentTarget.value)),
                    })
                  }
                />
              </div>
            </div>

            <div class="mt-6 space-y-2">
              <Label>Valores de pré-visualização</Label>
              <For each={VARIABLES}>
                {([variable, label]) => (
                  <PreviewVariableField
                    label={label}
                    variable={variable}
                    previewValues={previewValues}
                    setPreviewValues={setPreviewValues}
                  />
                )}
              </For>
              <p class="text-xs text-muted-foreground">
                Usados somente no editor; não são salvos no template.
              </p>
            </div>

            <div class="mt-6 space-y-2">
              <Label>Camadas</Label>
              <ul class="space-y-1">
                <For
                  each={[...draft().design_data.elements]
                    .reverse()
                    .map((element, reversedIndex) => ({
                      element,
                      index: draft().design_data.elements.length - 1 - reversedIndex,
                    }))}
                >
                  {(item) => (
                    <LayerRowItem
                      element={item.element}
                      index={item.index}
                      totalElements={draft().design_data.elements.length}
                      selectedId={selectedId}
                      onSelect={setSelectedId}
                      onMove={moveElement}
                      onDelete={item.element.type !== "qr" ? deleteElement : undefined}
                    />
                  )}
                </For>
              </ul>
            </div>
          </aside>

          {/* Canvas Center Area */}
          <div class="flex min-w-90 flex-1 flex-col">
            <RichTextToolbar
              controller={effectiveController()}
              selectionStyles={effectiveStyles()}
            />
            <BadgeCanvas
              design={draft().design_data}
              selectedId={selectedId()}
              onSelect={setSelectedId}
              onChangeElement={updateElement}
              textAdapter={textAdapter}
              previewValues={previewValues()}
              onDeleteElement={(id) => deleteElement(id)}
            />
          </div>

          {/* Right Sidebar: Properties */}
          <aside class="w-80 shrink-0 overflow-y-auto border-l border-muted bg-card p-4">
            <h2 class="font-semibold">Propriedades</h2>
            <Show
              when={selected()}
              fallback={
                <p class="mt-4 text-sm text-muted-foreground">
                  Selecione um elemento no crachá para personalizá-lo.
                </p>
              }
            >
              {(element) => (
                <div class="mt-4 space-y-4">
                  {/* Common: Coordinates and Dimensions */}
                  <div class="grid grid-cols-2 gap-2">
                    <div class="space-y-1">
                      <Label>Posição X (px)</Label>
                      <Input
                        type="number"
                        value={Math.round(element().x)}
                        onChange={(e) =>
                          updateElement(element().id, {
                            x: Number(e.currentTarget.value),
                          })
                        }
                      />
                    </div>
                    <div class="space-y-1">
                      <Label>Posição Y (px)</Label>
                      <Input
                        type="number"
                        value={Math.round(element().y)}
                        onChange={(e) =>
                          updateElement(element().id, {
                            y: Number(e.currentTarget.value),
                          })
                        }
                      />
                    </div>
                    <div class="space-y-1">
                      <Label>Largura (px)</Label>
                      <Input
                        type="number"
                        min={10}
                        value={Math.round(element().width)}
                        onChange={(e) =>
                          updateElement(element().id, {
                            width: Number(e.currentTarget.value),
                          })
                        }
                      />
                    </div>
                    <div class="space-y-1">
                      <Label>Altura (px)</Label>
                      <Input
                        type="number"
                        min={10}
                        value={Math.round(element().height)}
                        onChange={(e) =>
                          updateElement(element().id, {
                            height: Number(e.currentTarget.value),
                          })
                        }
                      />
                    </div>
                  </div>

                  {/* QR Specific Properties */}
                  <Show when={element().type === "qr"}>
                    {(() => {
                      const qr = () =>
                        element() as Extract<BadgeElement, { type: "qr" }>;

                      return (
                        <div class="space-y-3">
                          <div class="space-y-1">
                            <Label>Cor do QR</Label>
                            <input
                              type="color"
                              class="h-9 w-full cursor-pointer rounded-md border border-input bg-background p-1"
                              value={qr().foreground}
                              onInput={(e) =>
                                updateElement(element().id, {
                                  foreground: e.currentTarget.value,
                                })
                              }
                            />
                          </div>
                          <div class="space-y-1">
                            <Label>Cor de fundo do QR</Label>
                            <input
                              type="color"
                              class="h-9 w-full cursor-pointer rounded-md border border-input bg-background p-1"
                              value={qr().background}
                              onInput={(e) =>
                                updateElement(element().id, {
                                  background: e.currentTarget.value,
                                })
                              }
                            />
                          </div>
                          <div class="space-y-1">
                            <Label>Estilo dos pontos</Label>
                            <ToolbarCombobox
                              value={qr().style}
                              options={[
                                { value: "square", label: "Quadrado" },
                                { value: "dots", label: "Pontos" },
                                { value: "rounded", label: "Arredondado" },
                              ]}
                              placeholder="Estilo"
                              class="w-full"
                              triggerClass="h-9"
                              onChange={(value) =>
                                updateElement(element().id, {
                                  style: value as any,
                                })
                              }
                            />
                          </div>
                        </div>
                      );
                    })()}
                  </Show>

                  {/* Image Specific Properties */}
                  <Show when={element().type === "image"}>
                    {(() => {
                      const img = () =>
                        element() as Extract<BadgeElement, { type: "image" }>;

                      return (
                        <div class="space-y-3">
                          <div class="space-y-1">
                            <Label>Ajuste (fit)</Label>
                            <ToolbarCombobox
                              value={img().fit}
                              options={[
                                { value: "contain", label: "Conter (contain)" },
                                { value: "cover", label: "Cobrir (cover)" },
                                { value: "fill", label: "Preencher (fill)" },
                              ]}
                              placeholder="Ajuste da imagem"
                              class="w-full"
                              triggerClass="h-9"
                              onChange={(value) =>
                                updateElement(element().id, {
                                  fit: value as any,
                                })
                              }
                            />
                          </div>
                          <div class="space-y-1">
                            <Label>Opacidade ({Math.round(img().opacity * 100)}%)</Label>
                            <input
                              type="range"
                              min={0}
                              max={1}
                              step={0.05}
                              class="w-full"
                              value={img().opacity}
                              onInput={(e) =>
                                updateElement(element().id, {
                                  opacity: Number(e.currentTarget.value),
                                })
                              }
                            />
                          </div>
                          <div class="space-y-1">
                            <Label>Borda arredondada ({img().radius}px)</Label>
                            <input
                              type="range"
                              min={0}
                              max={100}
                              step={1}
                              class="w-full"
                              value={img().radius}
                              onInput={(e) =>
                                updateElement(element().id, {
                                  radius: Number(e.currentTarget.value),
                                })
                              }
                            />
                          </div>
                        </div>
                      );
                    })()}
                  </Show>

                  {/* Dynamic Information Tokens for Text */}
                  <Show when={element().type === "text"}>
                    <div class="border-t border-border" />
                    <div class="space-y-2.5">
                      <div>
                        <h3 class="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
                          Informações dinâmicas
                        </h3>
                        <p class="mt-1 text-xs leading-relaxed text-muted-foreground">
                          Clique para inserir no texto selecionado.
                        </p>
                      </div>
                      <div class="space-y-1.5">
                        <For each={VARIABLES}>
                          {([token, label, description]) => (
                            <VariableInsertButton
                              token={token}
                              label={label}
                              description={description}
                              effectiveController={effectiveController}
                            />
                          )}
                        </For>
                      </div>
                    </div>
                  </Show>

                  {/* Layer ordering Frente / Atrás */}
                  <div class="border-t border-border" />
                  <div class="grid grid-cols-2 gap-2">
                    <Button
                      variant="outline"
                      disabled={
                        draft().design_data.elements.indexOf(element()) ===
                        draft().design_data.elements.length - 1
                      }
                      onClick={() => {
                        const items = [...draft().design_data.elements];
                        const i = items.indexOf(element());
                        [items[i], items[i + 1]] = [items[i + 1], items[i]];
                        updateDesign({ elements: items });
                      }}
                    >
                      <ArrowUp class="size-4" />
                      Frente
                    </Button>
                    <Button
                      variant="outline"
                      disabled={
                        draft().design_data.elements.indexOf(element()) === 0
                      }
                      onClick={() => {
                        const items = [...draft().design_data.elements];
                        const i = items.indexOf(element());
                        [items[i], items[i - 1]] = [items[i - 1], items[i]];
                        updateDesign({ elements: items });
                      }}
                    >
                      <ArrowDown class="size-4" />
                      Atrás
                    </Button>
                  </div>

                  {/* Delete element button */}
                  <Show when={element().type !== "qr"}>
                    <Button
                      variant="destructive"
                      class="w-full"
                      onClick={() => deleteElement(element().id)}
                    >
                      <Trash2 class="size-4" />
                      Excluir elemento
                    </Button>
                  </Show>
                </div>
              )}
            </Show>
          </aside>
        </div>
      </div>
    </div>
  );
}

function PreviewVariableField(props: {
  label: string;
  variable: string;
  previewValues: () => Record<string, string>;
  setPreviewValues: (
    updater: (current: Record<string, string>) => Record<string, string>,
  ) => void;
}): JSX.Element {
  const variableKey = () => props.variable.slice(2, -2);
  const value = () => props.previewValues()[variableKey()] ?? "";

  const handleInput = (event: { currentTarget: HTMLInputElement }) => {
    const staticKey = variableKey();
    const nextValue = event.currentTarget.value;
    props.setPreviewValues((current) => ({
      ...current,
      [staticKey]: nextValue,
    }));
  };

  return (
    <Input
      aria-label={props.label}
      placeholder={props.label}
      value={value()}
      onInput={handleInput}
    />
  );
}

function LayerRowItem(props: {
  element: BadgeElement;
  index: number;
  totalElements: number;
  selectedId: () => string | null;
  onSelect: (id: string) => void;
  onMove: (id: string, delta: -1 | 1) => void;
  onDelete?: (id: string) => void;
}): JSX.Element {
  const isSelected = createMemo(() => props.selectedId() === props.element.id);

  return (
    <li
      class={cn(
        "flex cursor-pointer items-center gap-2 rounded-md border px-2 py-1.5 text-sm",
        isSelected()
          ? "border-ring bg-muted"
          : "border-transparent hover:bg-muted/60",
      )}
      onClick={() => props.onSelect(props.element.id)}
    >
      <Show when={props.element.type === "text"}>
        <TypeLucide class="size-3.5 shrink-0 text-muted-foreground" />
      </Show>
      <Show when={props.element.type === "image"}>
        <ImageLucide class="size-3.5 shrink-0 text-muted-foreground" />
      </Show>
      <Show when={props.element.type === "qr"}>
        <BadgeCheck class="size-3.5 shrink-0 text-muted-foreground" />
      </Show>

      <span class="min-w-0 flex-1 truncate">{layerLabel(props.element)}</span>

      <button
        type="button"
        aria-label="Trazer para frente"
        disabled={props.index === props.totalElements - 1}
        class="rounded p-0.5 disabled:opacity-30 cursor-pointer"
        onClick={(event) => {
          event.stopPropagation();
          props.onMove(props.element.id, 1);
        }}
      >
        <ArrowUp class="size-3.5" />
      </button>

      <button
        type="button"
        aria-label="Enviar para trás"
        disabled={props.index === 0}
        class="rounded p-0.5 disabled:opacity-30 cursor-pointer"
        onClick={(event) => {
          event.stopPropagation();
          props.onMove(props.element.id, -1);
        }}
      >
        <ArrowDown class="size-3.5" />
      </button>

      <Show when={props.element.type !== "qr" && props.onDelete}>
        <button
          type="button"
          aria-label="Excluir camada"
          class="rounded p-0.5 hover:bg-destructive hover:text-destructive-foreground cursor-pointer"
          onClick={(event) => {
            event.stopPropagation();
            props.onDelete?.(props.element.id);
          }}
        >
          <Trash2 class="size-3.5" />
        </button>
      </Show>
    </li>
  );
}

function VariableInsertButton(props: {
  token: string;
  label: string;
  description: string;
  effectiveController: () => RichTextController | null;
}): JSX.Element {
  const disabled = createMemo(() => !props.effectiveController());

  return (
    <button
      type="button"
      class={cn(
        buttonVariants({ variant: "outline" }),
        "h-auto w-full justify-start whitespace-normal px-3 py-2 text-left text-xs cursor-pointer disabled:pointer-events-none disabled:opacity-50",
      )}
      disabled={disabled()}
      title={props.description}
      onPointerDown={(e) => e.preventDefault()}
      onClick={() => props.effectiveController()?.insertText(props.token)}
    >
      <span class="min-w-0">
        <span class="block truncate">{props.label}</span>
        <span class="mt-0.5 block text-[11px] font-normal leading-tight text-muted-foreground">
          {props.description}
        </span>
      </span>
    </button>
  );
}
