import type { JSX } from "@solidjs/web";
import { For, Show, createEffect, createSignal } from "solid-js";
import ArrowDownIcon from "~icons/lucide/arrow-down";
import ArrowUpIcon from "~icons/lucide/arrow-up";
import FingerprintIcon from "~icons/lucide/fingerprint";
import ImageIcon from "~icons/lucide/image";
import ImagePlusIcon from "~icons/lucide/image-plus";
import PenToolIcon from "~icons/lucide/pen-tool";
import Trash2Icon from "~icons/lucide/trash-2";
import TypeIcon from "~icons/lucide/type";

import { Button, cn } from "@trieoh/ui-solid";
import { ToolbarCombobox } from "@/features/editor/toolbar-combobox";
import type { CertificationTemplateElement } from "../../model";
import {
  CERTIFICATE_CANVAS_PRESETS,
  CERTIFICATE_IMAGE_ACCEPT,
  MAX_CERTIFICATE_CANVAS_SIZE,
  MIN_CERTIFICATE_CANVAS_SIZE,
} from "../constants";
import {
  createImageElement,
  createSignatureElement,
  createTextElement,
} from "../factories";
import { certificateEditorActions, certificateEditorStore } from "../store";
import type {
  CertificateCanvasSize,
  SignatureCertificateElement,
  TextCertificateElement,
} from "../types";
import {
  isSupportedCertificateImage,
  loadCertificateImageDimensions,
  readCertificateFile,
} from "../utils";

const ArrowDown = ArrowDownIcon as unknown as (props: { class?: string }) => JSX.Element;
const ArrowUp = ArrowUpIcon as unknown as (props: { class?: string }) => JSX.Element;
const Fingerprint = FingerprintIcon as unknown as (props: { class?: string }) => JSX.Element;
const Image = ImageIcon as unknown as (props: { class?: string }) => JSX.Element;
const ImagePlus = ImagePlusIcon as unknown as (props: { class?: string }) => JSX.Element;
const PenTool = PenToolIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash2 = Trash2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Type = TypeIcon as unknown as (props: { class?: string }) => JSX.Element;

function getLayerLabel(element: CertificationTemplateElement): string {
  if (element.type === "hash") return "Hash de verificação";
  if (element.type === "image") return "Imagem";
  if (element.type === "signature") {
    return `Assinatura: ${(element as SignatureCertificateElement).name}`;
  }

  const text = (element as TextCertificateElement).paragraphs
    .flatMap((paragraph) => paragraph.runs.map((run) => run.text))
    .join("")
    .trim();
  return text.length > 0 ? text.slice(0, 24) : "Texto vazio";
}

function CanvasDimensionInput(props: {
  label: string;
  value: number;
  onCommit: (value: number) => void;
}): JSX.Element {
  return (
    <div class="min-w-0 flex-1 space-y-1">
      <span class="text-xs text-muted-foreground">{props.label}</span>
      <input
        type="number"
        min={MIN_CERTIFICATE_CANVAS_SIZE.width}
        max={MAX_CERTIFICATE_CANVAS_SIZE.width}
        value={Number.isFinite(props.value) ? Math.round(props.value) : ""}
        class="h-9 w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        onChange={(event) => {
          const parsed = Number(event.currentTarget.value);
          if (
            Number.isFinite(parsed) &&
            parsed >= MIN_CERTIFICATE_CANVAS_SIZE.width &&
            parsed <= MAX_CERTIFICATE_CANVAS_SIZE.width
          ) {
            props.onCommit(parsed);
          } else {
            event.currentTarget.value = String(
              Number.isFinite(props.value) ? Math.round(props.value) : "",
            );
          }
        }}
        onKeyDown={(event) => {
          if (event.key === "Enter") event.currentTarget.blur();
        }}
      />
    </div>
  );
}

export function CertificateToolsSidebar(): JSX.Element {
  let imageInputRef: HTMLInputElement | undefined;
  let backgroundInputRef: HTMLInputElement | undefined;

  const [imageError, setImageError] = createSignal<string | null>(null);
  const [backgroundError, setBackgroundError] = createSignal<string | null>(null);
  const [readingImage, setReadingImage] = createSignal(false);
  const [readingBackground, setReadingBackground] = createSignal(false);
  const [backgroundSize, setBackgroundSize] =
    createSignal<CertificateCanvasSize | null>(null);

  const canvas = () => certificateEditorStore.state.canvas;
  const name = () => certificateEditorStore.state.draft.name;
  const kind = () => certificateEditorStore.state.draft.kind;
  const description = () =>
    certificateEditorStore.state.draft.description ?? "";
  const signatures = () => certificateEditorStore.state.availableSignatures;
  const backgroundUrl = () =>
    certificateEditorStore.state.draft.design_data.background;
  const elements = () =>
    certificateEditorStore.state.draft.design_data.elements ?? [];
  const selectedElementId = () =>
    certificateEditorStore.state.selectedElementId;

  const selectedPreset = () => {
    const c = canvas();
    const match = CERTIFICATE_CANVAS_PRESETS.find(
      (p) => p.size.width === c.width && p.size.height === c.height,
    );
    return match?.id;
  };

  createEffect(
    () => backgroundUrl(),
    (url) => {
      let active = true;
      if (!url) {
        setBackgroundSize(null);
        return;
      }
      void loadCertificateImageDimensions(url)
        .then((dim) => {
          if (active) setBackgroundSize(dim);
        })
        .catch(() => {
          if (active) setBackgroundSize(null);
        });
      return () => {
        active = false;
      };
    },
  );

  async function addImage(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    input.value = "";
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
        createImageElement(src, canvas(), naturalSize),
      );
    } catch {
      setImageError("Não foi possível carregar a imagem.");
    } finally {
      setReadingImage(false);
    }
  }

  async function setBackground(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    input.value = "";
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
    <aside class="flex w-64 shrink-0 flex-col gap-6 overflow-x-hidden overflow-y-auto border-r border-border bg-card p-4 text-card-foreground">
      {/* Informações */}
      <section class="space-y-3">
        <h2 class="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          Informações
        </h2>
        <div class="space-y-1.5">
          <span class="text-xs text-muted-foreground">Nome</span>
          <input
            value={name()}
            maxlength={160}
            class="h-9 w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
            onInput={(event) =>
              certificateEditorActions.setName(event.currentTarget.value)
            }
            placeholder="Nome do certificado"
          />
        </div>
        <div class="space-y-1.5">
          <span class="text-xs text-muted-foreground">Tipo</span>
          <ToolbarCombobox
            value={kind()}
            options={[
              { value: "edition_attendance", label: "Presença na edição" },
              { value: "program_attendance", label: "Presença na atividade" },
            ]}
            placeholder="Selecione o tipo"
            class="w-full"
            onChange={(value) =>
              certificateEditorActions.setKind(
                value as "edition_attendance" | "program_attendance",
              )
            }
          />
        </div>
        <div class="space-y-1.5">
          <span class="text-xs text-muted-foreground">Descrição</span>
          <input
            value={description()}
            maxlength={500}
            class="h-9 w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
            onInput={(event) =>
              certificateEditorActions.setDescription(
                event.currentTarget.value,
              )
            }
            placeholder="Opcional"
          />
        </div>
      </section>

      {/* Adicionar */}
      <section class="space-y-2.5">
        <h2 class="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          Adicionar
        </h2>
        <div class="grid grid-cols-2 gap-2">
          <Button
            type="button"
            variant="outline"
            class="h-16 flex-col gap-1 cursor-pointer"
            onClick={() =>
              certificateEditorActions.addElement(createTextElement(canvas()))
            }
          >
            <Type class="size-4" />
            <span class="text-xs">Texto</span>
          </Button>
          <Button
            type="button"
            variant="outline"
            class="h-16 flex-col gap-1 cursor-pointer"
            disabled={readingImage()}
            onClick={() => imageInputRef?.click()}
          >
            <ImagePlus class="size-4" />
            <span class="text-xs">
              {readingImage() ? "Carregando…" : "Imagem"}
            </span>
          </Button>
          <input
            ref={(el) => (imageInputRef = el)}
            type="file"
            accept={CERTIFICATE_IMAGE_ACCEPT}
            class="hidden"
            onChange={(event) => void addImage(event)}
          />
        </div>
        <Show when={imageError()}>
          <p class="text-xs text-destructive" role="alert">
            {imageError()}
          </p>
        </Show>
      </section>

      {/* Assinaturas */}
      <section class="space-y-2.5">
        <div class="flex items-center gap-1.5 text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          <PenTool class="size-3.5" />
          <span>Assinaturas</span>
        </div>
        <Show
          when={signatures().length > 0}
          fallback={
            <p class="rounded-md border border-dashed border-border p-2.5 text-center text-xs text-muted-foreground">
              Nenhuma assinatura disponível
            </p>
          }
        >
          <div class="grid grid-cols-2 gap-2">
            <For each={signatures()}>
              {(signature) => (
                <button
                  type="button"
                  onClick={() =>
                    certificateEditorActions.addElement(
                      createSignatureElement(signature, canvas()),
                    )
                  }
                  class="overflow-hidden rounded-md border border-border bg-popover text-left transition-colors hover:border-ring hover:bg-muted cursor-pointer"
                  title={`Adicionar assinatura de ${signature.name}`}
                >
                  <span class="flex h-16 items-center justify-center bg-white p-2">
                    <img
                      src={signature.url}
                      alt=""
                      class="max-h-full max-w-full object-contain"
                    />
                  </span>
                  <span class="block truncate border-t border-border px-2 py-1.5 text-xs">
                    {signature.name}
                  </span>
                </button>
              )}
            </For>
          </div>
        </Show>
      </section>

      <div class="border-t border-border" />

      {/* Fundo */}
      <section class="space-y-2.5">
        <h2 class="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          Fundo
        </h2>
        <div class="flex gap-2">
          <Button
            type="button"
            variant="outline"
            class="min-w-0 flex-1 cursor-pointer"
            disabled={readingBackground()}
            onClick={() => backgroundInputRef?.click()}
          >
            {readingBackground()
              ? "Carregando…"
              : backgroundUrl()
                ? "Trocar imagem"
                : "Adicionar imagem"}
          </Button>
          <Show when={backgroundUrl()}>
            <Button
              type="button"
              variant="ghost"
              class="cursor-pointer"
              onClick={() => certificateEditorActions.setBackgroundUrl(null)}
            >
              Remover
            </Button>
          </Show>
        </div>
        <input
          ref={(el) => (backgroundInputRef = el)}
          type="file"
          accept={CERTIFICATE_IMAGE_ACCEPT}
          class="hidden"
          onChange={(event) => void setBackground(event)}
        />
        <Show when={backgroundSize()}>
          {(size) => (
            <Button
              type="button"
              variant="outline"
              class="h-auto w-full whitespace-normal px-3 py-2 text-center text-xs leading-tight cursor-pointer"
              disabled={
                size().width < 320 ||
                size().height < 320 ||
                size().width > 6000 ||
                size().height > 6000
              }
              onClick={() => certificateEditorActions.setCanvasSize(size())}
            >
              <span class="block">
                Usar tamanho da imagem
                <span class="block text-[11px] text-muted-foreground">
                  {size().width}×{size().height}
                </span>
              </span>
            </Button>
          )}
        </Show>
        <Show when={backgroundError()}>
          <p class="text-xs text-destructive" role="alert">
            {backgroundError()}
          </p>
        </Show>
      </section>

      <div class="border-t border-border" />

      {/* Tamanho do certificado */}
      <section class="space-y-2.5">
        <h2 class="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          Tamanho do certificado
        </h2>
        <ToolbarCombobox
          value={selectedPreset()}
          options={CERTIFICATE_CANVAS_PRESETS.map((preset) => ({
            value: preset.id,
            label: `${preset.label} (${preset.size.width}×${preset.size.height})`,
          }))}
          placeholder="Predefinições"
          class="w-full"
          onChange={(value) => {
            const preset = CERTIFICATE_CANVAS_PRESETS.find(
              (item) => item.id === value,
            );
            if (preset) certificateEditorActions.setCanvasSize(preset.size);
          }}
        />
        <div class="flex items-end gap-2">
          <CanvasDimensionInput
            label="Largura"
            value={canvas().width}
            onCommit={(width) =>
              certificateEditorActions.setCanvasSize({
                width,
                height: canvas().height,
              })
            }
          />
          <span class="pb-2 text-xs text-muted-foreground">×</span>
          <CanvasDimensionInput
            label="Altura"
            value={canvas().height}
            onCommit={(height) =>
              certificateEditorActions.setCanvasSize({
                width: canvas().width,
                height,
              })
            }
          />
        </div>
        <p class="text-xs leading-relaxed text-muted-foreground">
          Os elementos e as fontes são redimensionados proporcionalmente com o
          canvas.
        </p>
      </section>

      <div class="border-t border-border" />

      {/* Camadas */}
      <section class="space-y-2.5">
        <h2 class="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          Camadas
        </h2>
        <ul class="space-y-1">
          <For each={[...elements()].reverse()}>
            {(element, reversedIndex) => {
              const index = () => elements().length - 1 - reversedIndex();
              const selected = () => selectedElementId() === element.id;

              return (
                <li
                  class={cn(
                    "flex cursor-pointer items-center gap-2 rounded-md border px-2 py-1.5 text-sm transition-colors",
                    selected()
                      ? "border-ring bg-muted"
                      : "border-transparent hover:bg-muted/60",
                  )}
                  onClick={() =>
                    certificateEditorActions.selectElement(element.id)
                  }
                >
                  <Show when={element.type === "hash"}>
                    <Fingerprint class="size-3.5 shrink-0 text-muted-foreground" />
                  </Show>
                  <Show when={element.type === "text"}>
                    <Type class="size-3.5 shrink-0 text-muted-foreground" />
                  </Show>
                  <Show when={element.type === "image"}>
                    <Image class="size-3.5 shrink-0 text-muted-foreground" />
                  </Show>
                  <Show when={element.type === "signature"}>
                    <PenTool class="size-3.5 shrink-0 text-muted-foreground" />
                  </Show>

                  <span class="min-w-0 flex-1 truncate text-xs">
                    {getLayerLabel(element)}
                  </span>

                  <button
                    type="button"
                    title="Trazer para frente"
                    aria-label="Trazer para frente"
                    disabled={index() === elements().length - 1}
                    class="rounded p-0.5 text-muted-foreground hover:bg-background disabled:opacity-30 cursor-pointer"
                    onClick={(event) => {
                      event.stopPropagation();
                      certificateEditorActions.bringForward(element.id);
                    }}
                  >
                    <ArrowUp class="size-3.5" />
                  </button>

                  <button
                    type="button"
                    title="Enviar para trás"
                    aria-label="Enviar para trás"
                    disabled={index() === 0}
                    class="rounded p-0.5 text-muted-foreground hover:bg-background disabled:opacity-30 cursor-pointer"
                    onClick={(event) => {
                      event.stopPropagation();
                      certificateEditorActions.sendBackward(element.id);
                    }}
                  >
                    <ArrowDown class="size-3.5" />
                  </button>

                  <Show when={element.type !== "hash"}>
                    <button
                      type="button"
                      title="Excluir"
                      aria-label="Excluir camada"
                      class="rounded p-0.5 text-muted-foreground hover:bg-destructive hover:text-destructive-foreground cursor-pointer"
                      onClick={(event) => {
                        event.stopPropagation();
                        certificateEditorActions.removeElement(element.id);
                      }}
                    >
                      <Trash2 class="size-3.5" />
                    </button>
                  </Show>
                </li>
              );
            }}
          </For>
        </ul>
      </section>
    </aside>
  );
}
