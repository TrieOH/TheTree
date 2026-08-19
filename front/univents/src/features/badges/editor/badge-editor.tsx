import { useQuery } from "@tanstack/react-query";
import { Link, useNavigate } from "@tanstack/react-router";
import {
  ArrowDown,
  ArrowLeft,
  ArrowUp,
  BadgeCheck,
  ImageIcon,
  Loader2,
  Monitor,
  Save,
  Trash2,
  Type,
} from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { badgeTemplateQueryOptions } from "@/features/badges/api";
import { allTicketsQueryOptions } from "@/features/tickets/api";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { Label } from "@/shared/ui/shadcn/label";
import type {
  CertificateRichTextController,
  CertificateTextSelectionStyles,
} from "../../certifications/editor/store";
import { RichTextToolbar } from "../../certifications/editor/ui/certificate-text-toolbar";
import { ToolbarCombobox } from "../../certifications/editor/ui/toolbar-combobox";
import {
  useCreateBadgeTemplateMutation,
  useUpdateBadgeTemplateMutation,
} from "../api/mutations";
import { DEFAULT_BADGE_TEMPLATE } from "../default-template";
import type { BadgeElement, BadgeTemplateCreate } from "../model";
import { badgeTemplateCreateSchema } from "../model";
import { BadgeCanvas } from "./badge-canvas";
import { uploadBadgeAssets } from "./upload-assets";

const VARIABLES = [
  [
    "{{participant_name}}",
    "Nome completo do participante",
    "Nome legal informado no perfil do participante",
  ],
  ["{{event_name}}", "Nome do evento"],
  ["{{edition_name}}", "Nome da edição"],
  ["{{ticket_name}}", "Tipo de ingresso"],
  ["{{location}}", "Local", "Local informado na edição"],
] as const;
const VARIABLE_OPTIONS = VARIABLES.map(([value, label, description]) => ({
  value,
  label,
  description,
}));

const uid = () => crypto.randomUUID();
const PX_PER_MM = 96 / 25.4;
const pxToMm = (value: number) => Number((value / PX_PER_MM).toFixed(1));
const mmToPx = (value: number) => Math.round(value * PX_PER_MM);

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

export function BadgeEditor({
  eventId,
  editionId,
  templateId,
  duplicate = false,
}: {
  eventId: string;
  editionId: string;
  templateId?: string;
  duplicate?: boolean;
}) {
  const navigate = useNavigate();
  const createMutation = useCreateBadgeTemplateMutation();
  const updateMutation = useUpdateBadgeTemplateMutation();
  const templateQuery = useQuery({
    ...badgeTemplateQueryOptions(templateId ?? ""),
    enabled: Boolean(templateId),
  });
  const { data: tickets = [] } = useQuery(allTicketsQueryOptions(editionId));
  const [draft, setDraft] = useState<BadgeTemplateCreate>(() =>
    structuredClone(DEFAULT_BADGE_TEMPLATE),
  );
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);
  const [textController, setTextController] =
    useState<CertificateRichTextController | null>(null);
  const [textSelectionStyles, setTextSelectionStyles] =
    useState<CertificateTextSelectionStyles | null>(null);
  const selected =
    draft.design_data.elements.find((item) => item.id === selectedId) ?? null;

  useEffect(() => {
    if (templateQuery.data) {
      const copy = structuredClone(templateQuery.data);
      setDraft({
        ...copy,
        name: duplicate ? `${copy.name} (cópia)` : copy.name,
      });
    }
  }, [duplicate, templateQuery.data]);

  const updateDesign = useCallback(
    (changes: Partial<BadgeTemplateCreate["design_data"]>) =>
      setDraft((current) => ({
        ...current,
        design_data: { ...current.design_data, ...changes },
      })),
    [],
  );
  const updateElement = useCallback(
    (id: string, changes: Partial<BadgeElement>) =>
      setDraft((current) => ({
        ...current,
        design_data: {
          ...current.design_data,
          elements: current.design_data.elements.map((item) =>
            item.id === id ? ({ ...item, ...changes } as BadgeElement) : item,
          ),
        },
      })),
    [],
  );
  const resizeCanvas = (canvas: { width: number; height: number }) => {
    const previous = draft.design_data.canvas;
    if (canvas.width < 200 || canvas.height < 200) return;
    const scaleX = canvas.width / previous.width;
    const scaleY = canvas.height / previous.height;
    const fontScale = Math.sqrt(scaleX * scaleY);
    updateDesign({
      canvas,
      elements: draft.design_data.elements.map((element) => ({
        ...element,
        x: element.x * scaleX,
        y: element.y * scaleY,
        width: element.width * scaleX,
        height: element.height * scaleY,
        ...(element.type === "text"
          ? {
              paragraphs: element.paragraphs.map((paragraph) => ({
                ...paragraph,
                runs: paragraph.runs.map((run) => ({
                  ...run,
                  fontSize: Math.max(6, Math.round(run.fontSize * fontScale)),
                })),
              })),
            }
          : {}),
      })) as BadgeElement[],
    });
  };
  const textAdapter = useMemo(
    () => ({
      updateParagraphs: (
        id: string,
        paragraphs: Extract<BadgeElement, { type: "text" }>["paragraphs"],
      ) => updateElement(id, { paragraphs }),
      setController: setTextController,
      setSelectionStyles: setTextSelectionStyles,
      stopEditing: () => undefined,
    }),
    [updateElement],
  );
  const addElement = (element: BadgeElement) => {
    updateDesign({ elements: [...draft.design_data.elements, element] });
    setSelectedId(element.id);
  };
  const moveElement = (id: string, direction: -1 | 1) => {
    const elements = [...draft.design_data.elements];
    const index = elements.findIndex((element) => element.id === id);
    const target = index + direction;
    if (index < 0 || target < 0 || target >= elements.length) return;
    [elements[index], elements[target]] = [elements[target], elements[index]];
    updateDesign({ elements });
  };
  const deleteElement = (id: string) => {
    updateDesign({
      elements: draft.design_data.elements.filter((item) => item.id !== id),
    });
    setSelectedId(null);
  };

  async function save() {
    const parsed = badgeTemplateCreateSchema.safeParse(draft);
    if (!parsed.success)
      return toast.error(
        parsed.error.issues[0]?.message ?? "Template inválido",
      );
    setUploading(true);
    try {
      const data = await uploadBadgeAssets(parsed.data, eventId, editionId);
      const onSuccess = () =>
        void navigate({
          to: "/admin/events/$eventId/editions/$editionId/badges",
          params: { eventId, editionId },
        });
      if (templateId && !duplicate)
        updateMutation.mutate({ templateId, data, editionId }, { onSuccess });
      else createMutation.mutate({ editionId, data }, { onSuccess });
    } catch {
      toast.error("Não foi possível enviar as imagens do crachá");
    } finally {
      setUploading(false);
    }
  }

  return (
    <div className="h-dvh min-h-0 bg-background text-foreground">
      <div className="flex h-full flex-col items-center justify-center gap-3 px-6 text-center lg:hidden!">
        <Monitor className="size-10 text-muted-foreground" />
        <p className="max-w-xs text-sm text-muted-foreground">
          O editor de crachás foi desenvolvido para telas maiores. Abra esta
          página em um computador para editar o template.
        </p>
      </div>
      <div className="hidden h-full min-h-0 flex-col lg:flex">
        <header className="flex h-16 shrink-0 items-center justify-between border-b border-muted bg-card px-4 shadow-sm">
          <div className="flex items-center gap-3">
            <Link
              to="/admin/events/$eventId/editions/$editionId/badges"
              params={{ eventId, editionId }}
              className="inline-flex size-8 items-center justify-center rounded-lg border border-muted"
            >
              <ArrowLeft className="size-4" />
            </Link>
            <strong className="text-sm">Editor de crachás</strong>
          </div>
          <Button
            onClick={() => void save()}
            disabled={
              createMutation.isPending || updateMutation.isPending || uploading
            }
          >
            {createMutation.isPending ||
            updateMutation.isPending ||
            uploading ? (
              <Loader2 className="size-4 animate-spin" />
            ) : (
              <Save className="size-4" />
            )}
            Salvar
          </Button>
        </header>
        <div className="flex min-h-0 flex-1 overflow-x-auto">
          <aside className="w-72 shrink-0 overflow-y-auto border-r border-muted bg-card p-4">
            <div className="space-y-2">
              <Label htmlFor="badge-name">Nome do template</Label>
              <Input
                id="badge-name"
                value={draft.name}
                onChange={(event) =>
                  setDraft({ ...draft, name: event.target.value })
                }
              />
            </div>
            <div className="mt-5 space-y-2">
              <Label htmlFor="badge-ticket">Ingresso associado</Label>
              <ToolbarCombobox
                value={draft.ticket_type_id ?? ""}
                options={[
                  { value: "", label: "Padrão da edição" },
                  ...tickets.map((ticket) => ({
                    value: ticket.id,
                    label: ticket.name,
                  })),
                ]}
                placeholder="Selecione o ingresso"
                searchPlaceholder="Buscar ingresso..."
                className="w-full"
                triggerClassName="h-9"
                disabled={Boolean(templateId) && !duplicate}
                onChange={(value) =>
                  setDraft({
                    ...draft,
                    ticket_type_id: value || null,
                  })
                }
              />
              {templateId && !duplicate && (
                <p className="text-xs text-muted-foreground">
                  O ingresso associado não pode ser alterado após a criação.
                </p>
              )}
            </div>
            <div className="mt-5 space-y-2">
              <Label htmlFor="badge-origin">Origem do crachá</Label>
              <ToolbarCombobox
                value={draft.origin ?? ""}
                options={[
                  { value: "", label: "Participante / padrão" },
                  { value: "staff", label: "Equipe (staff)" },
                ]}
                placeholder="Selecione a origem"
                className="w-full"
                triggerClassName="h-9"
                disabled={Boolean(templateId) && !duplicate}
                onChange={(value) =>
                  setDraft({
                    ...draft,
                    origin: value === "staff" ? "staff" : null,
                  })
                }
              />
              {templateId && !duplicate && (
                <p className="text-xs text-muted-foreground">
                  A origem não pode ser alterada após a criação.
                </p>
              )}
            </div>
            <div className="mt-6">
              <Label>Adicionar</Label>
              <div className="mt-2 grid grid-cols-2 gap-2">
                <Button
                  variant="outline"
                  className="h-16 flex-col gap-1"
                  onClick={() =>
                    addElement({
                      id: uid(),
                      type: "text",
                      x: draft.design_data.canvas.width * 0.15,
                      y: draft.design_data.canvas.height * 0.4,
                      width: draft.design_data.canvas.width * 0.7,
                      height: Math.max(
                        32,
                        draft.design_data.canvas.height * 0.15,
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
                  <Type className="size-4" />
                  Texto
                </Button>
                <label className="inline-flex h-16 cursor-pointer flex-col items-center justify-center gap-1 rounded-lg border border-border bg-background text-sm font-medium hover:bg-muted hover:text-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50">
                  <ImageIcon className="size-4" />
                  Imagem
                  <input
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={async (event) => {
                      const file = event.target.files?.[0];
                      if (file)
                        addElement({
                          id: uid(),
                          type: "image",
                          x: draft.design_data.canvas.width * 0.2,
                          y: draft.design_data.canvas.height * 0.2,
                          width: draft.design_data.canvas.width * 0.6,
                          height: draft.design_data.canvas.height * 0.6,
                          src: await readImage(file),
                          fit: "contain",
                          radius: 0,
                          opacity: 1,
                        });
                      event.target.value = "";
                    }}
                  />
                </label>
              </div>
            </div>
            <div className="mt-6 space-y-2">
              <Label>Fundo</Label>
              <div className="flex gap-2">
                <Input
                  type="color"
                  className="w-14 p-1"
                  value={
                    draft.design_data.backgroundColor === "transparent"
                      ? "#ffffff"
                      : draft.design_data.backgroundColor
                  }
                  disabled={draft.design_data.backgroundColor === "transparent"}
                  onChange={(event) =>
                    updateDesign({ backgroundColor: event.target.value })
                  }
                />
                <label className="inline-flex h-9 flex-1 cursor-pointer items-center justify-center rounded-md border border-muted text-sm">
                  Imagem de fundo
                  <input
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={async (event) => {
                      const file = event.target.files?.[0];
                      if (file)
                        updateDesign({ background: await readImage(file) });
                    }}
                  />
                </label>
              </div>
              <label className="flex items-center gap-2 text-sm text-muted-foreground">
                <input
                  type="checkbox"
                  checked={draft.design_data.backgroundColor === "transparent"}
                  onChange={(event) =>
                    updateDesign({
                      backgroundColor: event.target.checked
                        ? "transparent"
                        : "#ffffff",
                    })
                  }
                />
                Fundo transparente
              </label>
              {draft.design_data.background && (
                <Button
                  variant="ghost"
                  className="w-full"
                  onClick={() => updateDesign({ background: null })}
                >
                  Remover imagem
                </Button>
              )}
            </div>
            <div className="mt-6 space-y-2">
              <Label>Tamanho (mm)</Label>
              <ToolbarCombobox
                value={
                  BADGE_SIZE_PRESETS.find(
                    (preset) =>
                      preset.width === draft.design_data.canvas.width &&
                      preset.height === draft.design_data.canvas.height,
                  )?.value
                }
                options={BADGE_SIZE_PRESETS}
                placeholder="Predefinições"
                className="w-full"
                triggerClassName="h-10"
                onChange={(value) => {
                  const preset = BADGE_SIZE_PRESETS.find(
                    (item) => item.value === value,
                  );
                  if (preset)
                    resizeCanvas({
                      width: preset.width,
                      height: preset.height,
                    });
                }}
              />
              <div className="grid grid-cols-2 gap-2">
                <Input
                  type="number"
                  value={pxToMm(draft.design_data.canvas.width)}
                  onChange={(e) =>
                    resizeCanvas({
                      ...draft.design_data.canvas,
                      width: mmToPx(Number(e.target.value)),
                    })
                  }
                />
                <Input
                  type="number"
                  value={pxToMm(draft.design_data.canvas.height)}
                  onChange={(e) =>
                    resizeCanvas({
                      ...draft.design_data.canvas,
                      height: mmToPx(Number(e.target.value)),
                    })
                  }
                />
              </div>
            </div>
            <div className="mt-6 space-y-2">
              <Label>Camadas</Label>
              <ul className="space-y-1">
                {[...draft.design_data.elements]
                  .reverse()
                  .map((element, reversedIndex) => {
                    const index =
                      draft.design_data.elements.length - 1 - reversedIndex;
                    const Icon =
                      element.type === "text"
                        ? Type
                        : element.type === "image"
                          ? ImageIcon
                          : BadgeCheck;
                    return (
                      <li
                        key={element.id}
                        className={cn(
                          "flex cursor-pointer items-center gap-2 rounded-md border px-2 py-1.5 text-sm",
                          selectedId === element.id
                            ? "border-ring bg-muted"
                            : "border-transparent hover:bg-muted/60",
                        )}
                        onClick={() => setSelectedId(element.id)}
                      >
                        <Icon className="size-3.5 shrink-0 text-muted-foreground" />
                        <span className="min-w-0 flex-1 truncate">
                          {layerLabel(element)}
                        </span>
                        <button
                          type="button"
                          aria-label="Trazer para frente"
                          disabled={
                            index === draft.design_data.elements.length - 1
                          }
                          className="rounded p-0.5 disabled:opacity-30"
                          onClick={(event) => {
                            event.stopPropagation();
                            moveElement(element.id, 1);
                          }}
                        >
                          <ArrowUp className="size-3.5" />
                        </button>
                        <button
                          type="button"
                          aria-label="Enviar para trás"
                          disabled={index === 0}
                          className="rounded p-0.5 disabled:opacity-30"
                          onClick={(event) => {
                            event.stopPropagation();
                            moveElement(element.id, -1);
                          }}
                        >
                          <ArrowDown className="size-3.5" />
                        </button>
                        {element.type !== "qr" ? (
                          <button
                            type="button"
                            aria-label="Excluir camada"
                            className="rounded p-0.5 hover:bg-destructive hover:text-destructive-foreground"
                            onClick={(event) => {
                              event.stopPropagation();
                              deleteElement(element.id);
                            }}
                          >
                            <Trash2 className="size-3.5" />
                          </button>
                        ) : null}
                      </li>
                    );
                  })}
              </ul>
            </div>
          </aside>
          <div className="flex min-w-90 flex-1 flex-col">
            <RichTextToolbar
              controller={textController}
              selectionStyles={textSelectionStyles}
              variableOptions={VARIABLE_OPTIONS}
            />
            <BadgeCanvas
              design={draft.design_data}
              selectedId={selectedId}
              onSelect={setSelectedId}
              onChangeElement={updateElement}
              textAdapter={textAdapter}
              onDeleteElement={(id) => {
                deleteElement(id);
              }}
            />
          </div>
          <aside className="w-80 shrink-0 overflow-y-auto border-l border-muted bg-card p-4">
            <h2 className="font-semibold">Propriedades</h2>
            {selected ? (
              <div className="mt-4 space-y-4">
                {selected.type === "text" && (
                  <div className="hidden">
                    <div className="space-y-2">
                      <Label>Conteúdo</Label>
                      <textarea
                        className="min-h-24 w-full rounded-md border border-muted bg-background p-2 text-sm"
                        value={selected.content}
                        onChange={(e) =>
                          updateElement(selected.id, {
                            content: e.target.value,
                          })
                        }
                      />
                    </div>
                    <div className="flex flex-wrap gap-1">
                      {VARIABLES.map(([value, label]) => (
                        <button
                          type="button"
                          key={value}
                          title={label}
                          className="rounded border border-muted px-2 py-1 text-xs"
                          onClick={() =>
                            updateElement(selected.id, {
                              content: `${selected.content}${value}`,
                            })
                          }
                        >
                          {label}
                        </button>
                      ))}
                    </div>
                    <div className="grid grid-cols-2 gap-2">
                      <div>
                        <Label>Tamanho</Label>
                        <Input
                          type="number"
                          value={selected.fontSize}
                          onChange={(e) =>
                            updateElement(selected.id, {
                              fontSize: Number(e.target.value),
                            })
                          }
                        />
                      </div>
                      <div>
                        <Label>Cor</Label>
                        <Input
                          type="color"
                          value={selected.color}
                          onChange={(e) =>
                            updateElement(selected.id, {
                              color: e.target.value,
                            })
                          }
                        />
                      </div>
                    </div>
                    <div className="grid grid-cols-2 gap-2">
                      <ToolbarCombobox
                        value={selected.fontWeight}
                        options={[
                          { value: "normal", label: "Normal" },
                          { value: "bold", label: "Negrito" },
                        ]}
                        placeholder="Peso"
                        onChange={(value) =>
                          updateElement(selected.id, {
                            fontWeight: value as "normal" | "bold",
                          })
                        }
                      />
                      <ToolbarCombobox
                        value={selected.align}
                        options={[
                          { value: "left", label: "Esquerda" },
                          { value: "center", label: "Centro" },
                          { value: "right", label: "Direita" },
                        ]}
                        placeholder="Alinhamento"
                        onChange={(value) =>
                          updateElement(selected.id, {
                            align: value as "left" | "center" | "right",
                          })
                        }
                      />
                    </div>
                  </div>
                )}
                {selected.type === "image" && (
                  <>
                    <div>
                      <Label>Ajuste</Label>
                      <ToolbarCombobox
                        value={selected.fit}
                        options={[
                          { value: "contain", label: "Conter" },
                          { value: "cover", label: "Cobrir" },
                          { value: "fill", label: "Esticar" },
                        ]}
                        placeholder="Ajuste"
                        className="w-full"
                        onChange={(value) =>
                          updateElement(selected.id, {
                            fit: value as "contain" | "cover" | "fill",
                          })
                        }
                      />
                    </div>
                    <div>
                      <Label>Arredondamento</Label>
                      <Input
                        type="number"
                        value={selected.radius}
                        onChange={(e) =>
                          updateElement(selected.id, {
                            radius: Number(e.target.value),
                          })
                        }
                      />
                    </div>
                  </>
                )}
                {selected.type === "qr" && (
                  <div className="space-y-3">
                    <p className="rounded-md bg-muted p-2.5 text-xs text-muted-foreground">
                      O QR Code de check-in é obrigatório. Ele pode ser movido,
                      redimensionado e estilizado, mas não excluído.
                    </p>
                    <ToolbarCombobox
                      options={[
                        { value: "square", label: "Quadrado" },
                        { value: "rounded", label: "Arredondado" },
                        { value: "dots", label: "Pontos" },
                      ]}
                      value={selected.style}
                      placeholder="Estilo do QR Code"
                      className="w-full"
                      onChange={(style) =>
                        updateElement(selected.id, {
                          style: style as "square" | "rounded" | "dots",
                        })
                      }
                    />
                    <div className="grid grid-cols-2 gap-2">
                      <div>
                        <Label>Cor</Label>
                        <Input
                          type="color"
                          value={selected.foreground}
                          onChange={(e) =>
                            updateElement(selected.id, {
                              foreground: e.target.value,
                            })
                          }
                        />
                      </div>
                      <div>
                        <Label>Fundo</Label>
                        <Input
                          type="color"
                          value={selected.background}
                          onChange={(e) =>
                            updateElement(selected.id, {
                              background: e.target.value,
                            })
                          }
                        />
                      </div>
                    </div>
                  </div>
                )}
                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <Label>Largura (mm)</Label>
                    <Input
                      type="number"
                      value={pxToMm(selected.width)}
                      onChange={(e) =>
                        updateElement(selected.id, {
                          width: mmToPx(Number(e.target.value)),
                        })
                      }
                    />
                  </div>
                  <div>
                    <Label>Altura (mm)</Label>
                    <Input
                      type="number"
                      value={pxToMm(selected.height)}
                      onChange={(e) =>
                        updateElement(selected.id, {
                          height: mmToPx(Number(e.target.value)),
                        })
                      }
                    />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-2">
                  <Button
                    variant="outline"
                    disabled={
                      draft.design_data.elements.indexOf(selected) ===
                      draft.design_data.elements.length - 1
                    }
                    onClick={() => {
                      const items = [...draft.design_data.elements];
                      const i = items.indexOf(selected);
                      [items[i], items[i + 1]] = [items[i + 1], items[i]];
                      updateDesign({ elements: items });
                    }}
                  >
                    <ArrowUp className="size-4" />
                    Frente
                  </Button>
                  <Button
                    variant="outline"
                    disabled={
                      draft.design_data.elements.indexOf(selected) === 0
                    }
                    onClick={() => {
                      const items = [...draft.design_data.elements];
                      const i = items.indexOf(selected);
                      [items[i], items[i - 1]] = [items[i - 1], items[i]];
                      updateDesign({ elements: items });
                    }}
                  >
                    <ArrowDown className="size-4" />
                    Atrás
                  </Button>
                </div>
                {selected.type !== "qr" ? (
                  <Button
                    variant="destructive"
                    className="w-full"
                    onClick={() => {
                      deleteElement(selected.id);
                    }}
                  >
                    <Trash2 className="size-4" />
                    Excluir elemento
                  </Button>
                ) : null}
              </div>
            ) : (
              <p className="mt-4 text-sm text-muted-foreground">
                Selecione um elemento no crachá para personalizá-lo.
              </p>
            )}
          </aside>
        </div>
      </div>
    </div>
  );
}
