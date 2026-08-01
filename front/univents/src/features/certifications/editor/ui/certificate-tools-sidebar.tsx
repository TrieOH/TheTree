import {
  ArrowDown,
  ArrowUp,
  Fingerprint,
  Image as ImageIcon,
  ImagePlus,
  PenTool,
  Trash2,
  Type,
} from "lucide-react";
import type { ChangeEvent } from "react";
import { useEffect, useRef, useState } from "react";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { Separator } from "@/shared/ui/shadcn/separator";
import type { CertificationTemplateElement } from "../../model";
import {
  CERTIFICATE_CANVAS_PRESETS,
  CERTIFICATE_IMAGE_ACCEPT,
} from "../constants";
import {
  createImageElement,
  createSignatureElement,
  createTextElement,
} from "../factories";
import { certificateEditorActions, useCertificateEditorState } from "../store";
import {
  isSupportedCertificateImage,
  loadCertificateImageDimensions,
  readCertificateFile,
} from "../utils";
import { ToolbarCombobox } from "./toolbar-combobox";

const LAYER_ICON: Record<CertificationTemplateElement["type"], typeof Type> = {
  hash: Fingerprint,
  text: Type,
  image: ImageIcon,
  signature: PenTool,
};

function getLayerLabel(element: CertificationTemplateElement): string {
  if (element.type === "hash") return "Hash de verificação";
  if (element.type === "image") return "Imagem";
  if (element.type === "signature") return `Assinatura: ${element.name}`;

  const text = element.paragraphs
    .flatMap((paragraph) => paragraph.runs.map((run) => run.text))
    .join("")
    .trim();
  return text.length > 0 ? text.slice(0, 24) : "Texto vazio";
}

interface CanvasDimensionInputProps {
  label: string;
  value: number;
  onCommit: (value: number) => void;
}

function CanvasDimensionInput({
  label,
  value,
  onCommit,
}: CanvasDimensionInputProps) {
  const [draft, setDraft] = useState(String(Math.round(value)));

  useEffect(() => setDraft(String(Math.round(value))), [value]);

  function commit() {
    const parsed = Number(draft);
    if (Number.isFinite(parsed) && parsed >= 320 && parsed <= 6000) {
      onCommit(parsed);
      return;
    }
    setDraft(String(Math.round(value)));
  }

  return (
    <div className="min-w-0 flex-1 space-y-1">
      <span className="text-xs text-muted-foreground">{label}</span>
      <Input
        type="number"
        min={320}
        max={6000}
        value={draft}
        onChange={(event) => setDraft(event.target.value)}
        onBlur={commit}
        onKeyDown={(event) => {
          if (event.key === "Enter") event.currentTarget.blur();
        }}
      />
    </div>
  );
}

export function CertificateToolsSidebar() {
  const canvas = useCertificateEditorState((state) => state.canvas);
  const name = useCertificateEditorState((state) => state.draft.name);
  const kind = useCertificateEditorState((state) => state.draft.kind);
  const description = useCertificateEditorState(
    (state) => state.draft.description ?? "",
  );
  const selectedCanvasPreset = CERTIFICATE_CANVAS_PRESETS.find(
    (preset) =>
      preset.size.width === canvas.width &&
      preset.size.height === canvas.height,
  );
  const signatures = useCertificateEditorState(
    (state) => state.availableSignatures,
  );
  const backgroundUrl = useCertificateEditorState(
    (state) => state.draft.design_data.background,
  );
  const elements = useCertificateEditorState(
    (state) => state.draft.design_data.elements,
  );
  const selectedElementId = useCertificateEditorState(
    (state) => state.selectedElementId,
  );
  const imageInputRef = useRef<HTMLInputElement>(null);
  const backgroundInputRef = useRef<HTMLInputElement>(null);
  const [imageError, setImageError] = useState<string | null>(null);
  const [backgroundError, setBackgroundError] = useState<string | null>(null);
  const [readingImage, setReadingImage] = useState(false);
  const [readingBackground, setReadingBackground] = useState(false);
  const [backgroundSize, setBackgroundSize] = useState<{
    width: number;
    height: number;
  } | null>(null);

  useEffect(() => {
    let active = true;
    if (!backgroundUrl) {
      setBackgroundSize(null);
      return;
    }

    void loadCertificateImageDimensions(backgroundUrl)
      .then((size) => {
        if (active) setBackgroundSize(size);
      })
      .catch(() => {
        if (active) setBackgroundSize(null);
      });
    return () => {
      active = false;
    };
  }, [backgroundUrl]);

  async function addImage(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    event.target.value = "";
    if (!file) return;

    if (!isSupportedCertificateImage(file)) {
      setImageError("Use uma imagem PNG, JPEG ou WebP.");
      return;
    }

    setReadingImage(true);
    setImageError(null);
    try {
      const src = await readCertificateFile(file);
      const naturalSize = await loadCertificateImageDimensions(src).catch(
        () => undefined,
      );
      certificateEditorActions.addElement(
        createImageElement(src, canvas, naturalSize),
      );
    } catch {
      setImageError("Não foi possível carregar a imagem.");
    } finally {
      setReadingImage(false);
    }
  }

  async function setBackground(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    event.target.value = "";
    if (!file) return;

    if (!isSupportedCertificateImage(file)) {
      setBackgroundError("Use uma imagem PNG, JPEG ou WebP.");
      return;
    }

    setReadingBackground(true);
    setBackgroundError(null);
    try {
      const src = await readCertificateFile(file);
      certificateEditorActions.setBackgroundUrl(src);
    } catch {
      setBackgroundError("Não foi possível carregar a imagem de fundo.");
    } finally {
      setReadingBackground(false);
    }
  }

  return (
    <aside className="flex w-64 shrink-0 flex-col gap-6 overflow-x-hidden overflow-y-auto border-r border-border bg-card p-4 text-card-foreground">
      <section className="space-y-3">
        <h2 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
          Informações
        </h2>
        <div className="space-y-1.5">
          <span className="text-xs text-muted-foreground">Nome</span>
          <Input
            value={name}
            maxLength={160}
            onChange={(event) =>
              certificateEditorActions.setName(event.target.value)
            }
            placeholder="Nome do certificado"
          />
        </div>
        <div className="space-y-1.5">
          <span className="text-xs text-muted-foreground">Tipo</span>
          <ToolbarCombobox
            value={kind}
            options={[
              { value: "edition_attendance", label: "Presença na edição" },
              { value: "program_attendance", label: "Presença na atividade" },
            ]}
            placeholder="Selecione o tipo"
            searchPlaceholder="Buscar tipo..."
            onChange={(value) =>
              certificateEditorActions.setKind(
                value as "edition_attendance" | "program_attendance",
              )
            }
            className="w-full"
            triggerClassName="h-9"
          />
        </div>
        <div className="space-y-1.5">
          <span className="text-xs text-muted-foreground">Descrição</span>
          <Input
            value={description}
            maxLength={500}
            onChange={(event) =>
              certificateEditorActions.setDescription(event.target.value)
            }
            placeholder="Opcional"
          />
        </div>
      </section>
      <section className="space-y-2.5">
        <h2 className="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          Adicionar
        </h2>
        <div className="grid grid-cols-2 gap-2">
          <Button
            type="button"
            variant="outline"
            className="h-16 flex-col gap-1"
            onClick={() =>
              certificateEditorActions.addElement(createTextElement(canvas))
            }
          >
            <Type className="size-4" />
            Texto
          </Button>
          <Button
            type="button"
            variant="outline"
            className="h-16 flex-col gap-1"
            disabled={readingImage}
            onClick={() => imageInputRef.current?.click()}
          >
            <ImagePlus className="size-4" />
            {readingImage ? "Carregando…" : "Imagem"}
          </Button>
          <input
            ref={imageInputRef}
            type="file"
            accept={CERTIFICATE_IMAGE_ACCEPT}
            className="hidden"
            onChange={(event) => void addImage(event)}
          />
        </div>
        {imageError ? (
          <p className="text-xs text-destructive" role="alert">
            {imageError}
          </p>
        ) : null}
      </section>

      <section className="space-y-2.5">
        <div className="flex items-center gap-1.5 text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          <PenTool className="size-3.5" />
          Assinaturas
        </div>
        {signatures.length === 0 ? (
          <p className="rounded-md border border-dashed p-2.5 text-center text-xs text-muted-foreground">
            Nenhuma assinatura disponível
          </p>
        ) : (
          <div className="grid grid-cols-2 gap-2">
            {signatures.map((signature) => (
              <button
                key={signature.id}
                type="button"
                onClick={() =>
                  certificateEditorActions.addElement(
                    createSignatureElement(signature, canvas),
                  )
                }
                className="overflow-hidden rounded-md border bg-popover text-left transition-colors hover:border-ring hover:bg-muted focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/50 focus-visible:outline-none"
                title={`Adicionar assinatura de ${signature.name}`}
              >
                <span className="flex h-16 items-center justify-center bg-white p-2">
                  <img
                    src={signature.url}
                    alt=""
                    className="max-h-full max-w-full object-contain"
                  />
                </span>
                <span className="block truncate border-t px-2 py-1.5 text-xs">
                  {signature.name}
                </span>
              </button>
            ))}
          </div>
        )}
      </section>

      <Separator />

      <section className="space-y-2.5">
        <h2 className="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          Fundo
        </h2>
        <div className="flex gap-2">
          <Button
            type="button"
            variant="outline"
            className="min-w-0 flex-1"
            disabled={readingBackground}
            onClick={() => backgroundInputRef.current?.click()}
          >
            {readingBackground
              ? "Carregando…"
              : backgroundUrl
                ? "Trocar imagem"
                : "Adicionar imagem"}
          </Button>
          {backgroundUrl ? (
            <Button
              type="button"
              variant="ghost"
              onClick={() => certificateEditorActions.setBackgroundUrl(null)}
            >
              Remover
            </Button>
          ) : null}
        </div>
        <input
          ref={backgroundInputRef}
          type="file"
          accept={CERTIFICATE_IMAGE_ACCEPT}
          className="hidden"
          onChange={(event) => void setBackground(event)}
        />
        {backgroundSize ? (
          <Button
            type="button"
            variant="outline"
            className="h-auto w-full whitespace-normal px-3 py-2 text-center text-xs leading-tight"
            disabled={
              backgroundSize.width < 320 ||
              backgroundSize.height < 320 ||
              backgroundSize.width > 6000 ||
              backgroundSize.height > 6000
            }
            onClick={() =>
              certificateEditorActions.setCanvasSize(backgroundSize)
            }
          >
            <span className="block">
              Usar tamanho da imagem
              <span className="block text-[11px] text-muted-foreground">
                {backgroundSize.width}×{backgroundSize.height}
              </span>
            </span>
          </Button>
        ) : null}
        {backgroundError ? (
          <p className="text-xs text-destructive" role="alert">
            {backgroundError}
          </p>
        ) : null}
      </section>

      <Separator />

      <section className="space-y-2.5">
        <h2 className="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          Tamanho do certificado
        </h2>
        <ToolbarCombobox
          aria-label="Predefinição de tamanho"
          value={selectedCanvasPreset?.id}
          options={CERTIFICATE_CANVAS_PRESETS.map((preset) => ({
            value: preset.id,
            label: `${preset.label} (${preset.size.width}×${preset.size.height})`,
          }))}
          placeholder="Predefinições"
          className="w-full"
          triggerClassName="h-10 text-sm"
          dropdownClassName="w-full"
          onChange={(value) => {
            const preset = CERTIFICATE_CANVAS_PRESETS.find(
              (item) => item.id === value,
            );
            if (preset) certificateEditorActions.setCanvasSize(preset.size);
          }}
        />
        <div className="flex items-end gap-2">
          <CanvasDimensionInput
            label="Largura"
            value={canvas.width}
            onCommit={(width) =>
              certificateEditorActions.setCanvasSize({
                width,
                height: canvas.height,
              })
            }
          />
          <span className="pb-3 text-xs text-muted-foreground">×</span>
          <CanvasDimensionInput
            label="Altura"
            value={canvas.height}
            onCommit={(height) =>
              certificateEditorActions.setCanvasSize({
                width: canvas.width,
                height,
              })
            }
          />
        </div>
        <p className="text-xs leading-relaxed text-muted-foreground">
          Os elementos e as fontes são redimensionados proporcionalmente com o
          canvas.
        </p>
      </section>

      <Separator />

      <section className="space-y-2.5">
        <h2 className="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          Camadas
        </h2>
        <ul className="space-y-1">
          {[...elements].reverse().map((element, reversedIndex) => {
            const index = elements.length - 1 - reversedIndex;
            const Icon = LAYER_ICON[element.type];
            const selected = selectedElementId === element.id;

            return (
              <li
                key={element.id}
                className={
                  "flex cursor-pointer items-center gap-2 rounded-md border px-2 py-1.5 text-sm " +
                  (selected
                    ? "border-ring bg-muted"
                    : "border-transparent hover:bg-muted/60")
                }
                onClick={() =>
                  certificateEditorActions.selectElement(element.id)
                }
              >
                <Icon className="size-3.5 shrink-0 text-muted-foreground" />
                <span className="min-w-0 flex-1 truncate">
                  {getLayerLabel(element)}
                </span>
                <button
                  type="button"
                  title="Trazer para frente"
                  aria-label="Trazer para frente"
                  disabled={index === elements.length - 1}
                  className="rounded p-0.5 text-muted-foreground hover:bg-background disabled:opacity-30"
                  onClick={(event) => {
                    event.stopPropagation();
                    certificateEditorActions.bringForward(element.id);
                  }}
                >
                  <ArrowUp className="size-3.5" />
                </button>
                <button
                  type="button"
                  title="Enviar para trás"
                  aria-label="Enviar para trás"
                  disabled={index === 0}
                  className="rounded p-0.5 text-muted-foreground hover:bg-background disabled:opacity-30"
                  onClick={(event) => {
                    event.stopPropagation();
                    certificateEditorActions.sendBackward(element.id);
                  }}
                >
                  <ArrowDown className="size-3.5" />
                </button>
                {element.type !== "hash" ? (
                  <button
                    type="button"
                    title="Excluir"
                    aria-label="Excluir camada"
                    className="rounded p-0.5 text-muted-foreground hover:bg-destructive hover:text-destructive-foreground"
                    onClick={(event) => {
                      event.stopPropagation();
                      certificateEditorActions.removeElement(element.id);
                    }}
                  >
                    <Trash2 className="size-3.5" />
                  </button>
                ) : null}
              </li>
            );
          })}
        </ul>
      </section>
    </aside>
  );
}
