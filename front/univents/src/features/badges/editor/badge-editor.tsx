import { useQuery } from "@tanstack/react-query";
import { Link, useNavigate } from "@tanstack/react-router";
import {
  ArrowDown,
  ArrowLeft,
  ArrowUp,
  ImageIcon,
  Loader2,
  Monitor,
  Save,
  Trash2,
  Type,
} from "lucide-react";
import { useCallback, useMemo, useState } from "react";
import { toast } from "sonner";
import { allTicketsQueryOptions } from "@/features/tickets/api";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { Label } from "@/shared/ui/shadcn/label";
import type {
  CertificateRichTextController,
  CertificateTextSelectionStyles,
} from "../../certifications/editor/store";
import { RichTextToolbar } from "../../certifications/editor/ui/certificate-text-toolbar";
import { ToolbarCombobox } from "../../certifications/editor/ui/toolbar-combobox";
import { useCreateBadgeTemplateMutation } from "../api/mutations";
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

const BADGE_SIZE_PRESETS = [
  { value: "portrait", label: "Vertical (638×1011)", width: 638, height: 1011 },
  {
    value: "landscape",
    label: "Horizontal (1011×638)",
    width: 1011,
    height: 638,
  },
  { value: "square", label: "Quadrado (800×800)", width: 800, height: 800 },
] as const;

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
}: {
  eventId: string;
  editionId: string;
}) {
  const navigate = useNavigate();
  const createMutation = useCreateBadgeTemplateMutation();
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

  async function save() {
    const parsed = badgeTemplateCreateSchema.safeParse(draft);
    if (!parsed.success)
      return toast.error(
        parsed.error.issues[0]?.message ?? "Template inválido",
      );
    setUploading(true);
    try {
      const data = await uploadBadgeAssets(parsed.data, eventId, editionId);
      createMutation.mutate(
        { editionId, data },
        {
          onSuccess: () =>
            void navigate({
              to: "/admin/events/$eventId/editions/$editionId/badges",
              params: { eventId, editionId },
            }),
        },
      );
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
            disabled={createMutation.isPending || uploading}
          >
            {createMutation.isPending || uploading ? (
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
                onChange={(value) =>
                  setDraft({
                    ...draft,
                    ticket_type_id: value || null,
                  })
                }
              />
              <p className="text-xs text-muted-foreground">
                Só pode existir um template por tipo de ingresso.
              </p>
            </div>
            <div className="mt-6">
              <Label>Adicionar</Label>
              <div className="mt-2 grid grid-cols-2 gap-2">
                <Button
                  variant="outline"
                  onClick={() =>
                    addElement({
                      id: uid(),
                      type: "text",
                      x: 100,
                      y: 300,
                      width: 438,
                      height: 80,
                      paragraphs: [
                        {
                          align: "center",
                          lineHeight: 1.25,
                          runs: [
                            {
                              text: "Novo texto",
                              fontSize: 32,
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
                <label className="col-span-2 inline-flex h-9 cursor-pointer items-center justify-center gap-2 rounded-md border border-muted text-sm font-medium hover:bg-accent">
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
                          x: 159,
                          y: 180,
                          width: 320,
                          height: 220,
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
                  value={draft.design_data.backgroundColor}
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
              <Label>Tamanho (px)</Label>
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
                    updateDesign({
                      canvas: { width: preset.width, height: preset.height },
                    });
                }}
              />
              <div className="grid grid-cols-2 gap-2">
                <Input
                  type="number"
                  value={draft.design_data.canvas.width}
                  onChange={(e) =>
                    updateDesign({
                      canvas: {
                        ...draft.design_data.canvas,
                        width: Number(e.target.value),
                      },
                    })
                  }
                />
                <Input
                  type="number"
                  value={draft.design_data.canvas.height}
                  onChange={(e) =>
                    updateDesign({
                      canvas: {
                        ...draft.design_data.canvas,
                        height: Number(e.target.value),
                      },
                    })
                  }
                />
              </div>
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
                updateDesign({
                  elements: draft.design_data.elements.filter(
                    (item) => item.id !== id,
                  ),
                });
                setSelectedId(null);
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
                    <Label>Largura</Label>
                    <Input
                      type="number"
                      value={Math.round(selected.width)}
                      onChange={(e) =>
                        updateElement(selected.id, {
                          width: Number(e.target.value),
                        })
                      }
                    />
                  </div>
                  <div>
                    <Label>Altura</Label>
                    <Input
                      type="number"
                      value={Math.round(selected.height)}
                      onChange={(e) =>
                        updateElement(selected.id, {
                          height: Number(e.target.value),
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
                      updateDesign({
                        elements: draft.design_data.elements.filter(
                          (item) => item.id !== selected.id,
                        ),
                      });
                      setSelectedId(null);
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
