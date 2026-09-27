import type { JSX } from "@solidjs/web";
import { For, Show, createMemo, createSignal } from "solid-js";
import { Button } from "@trieoh/ui-solid";
import { ToolbarCombobox, normalizeHexColor } from "@/features/editor";
import type { CertificationTemplateElement } from "../../model";
import {
  CERTIFICATE_IMAGE_ACCEPT,
  CERTIFICATE_VARIABLES,
  MIN_CERTIFICATE_ELEMENT_SIZE,
} from "../constants";
import { certificateEditorActions, certificateEditorStore } from "../store";
import type {
  HashCertificateElement,
  ImageCertificateElement,
  SignatureCertificateElement,
  TextCertificateElement,
} from "../types";
import { isSupportedCertificateImage, readCertificateFile } from "../utils";

interface FieldProps {
  label: string;
  children: JSX.Element;
}

function Field(props: FieldProps): JSX.Element {
  return (
    <div class="block space-y-1">
      <span class="text-xs text-muted-foreground">{props.label}</span>
      {props.children}
    </div>
  );
}

interface NumberPropertyProps {
  label: string;
  value: number;
  min?: number;
  max?: number;
  onCommit: (value: number) => void;
}

function NumberProperty(props: NumberPropertyProps): JSX.Element {
  return (
    <Field label={props.label}>
      <input
        type="number"
        value={Number.isFinite(props.value) ? Math.round(props.value) : ""}
        min={props.min}
        max={props.max}
        class="h-9 w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        onChange={(event) => {
          const parsed = Number(event.currentTarget.value);
          if (
            Number.isFinite(parsed) &&
            (props.min === undefined || parsed >= props.min) &&
            (props.max === undefined || parsed <= props.max)
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
    </Field>
  );
}

function PositionSizeProperties(props: {
  element: CertificationTemplateElement;
}): JSX.Element {
  return (
    <div class="grid grid-cols-2 gap-2">
      <NumberProperty
        label="Posição X"
        value={props.element.x}
        onCommit={(x) =>
          certificateEditorActions.updateElementBounds(props.element.id, { x })
        }
      />
      <NumberProperty
        label="Posição Y"
        value={props.element.y}
        onCommit={(y) =>
          certificateEditorActions.updateElementBounds(props.element.id, { y })
        }
      />
      <NumberProperty
        label="Largura"
        value={props.element.width}
        min={MIN_CERTIFICATE_ELEMENT_SIZE.width}
        onCommit={(width) =>
          certificateEditorActions.updateElementBounds(props.element.id, { width })
        }
      />
      <NumberProperty
        label="Altura"
        value={props.element.height}
        min={MIN_CERTIFICATE_ELEMENT_SIZE.height}
        onCommit={(height) =>
          certificateEditorActions.updateElementBounds(props.element.id, { height })
        }
      />
    </div>
  );
}

function updateHash(
  id: string,
  patch: Partial<Omit<HashCertificateElement, "id" | "type">>,
) {
  certificateEditorActions.updateElement(id, (element: CertificationTemplateElement) =>
    element.type === "hash" ? { ...element, ...patch } : element,
  );
}

function HashProperties(props: { element: HashCertificateElement }): JSX.Element {
  return (
    <div class="space-y-4">
      <p class="rounded-md bg-muted p-2.5 text-xs leading-relaxed text-muted-foreground">
        O código e o endereço de validação são preenchidos automaticamente. O
        bloco pode ser estilizado e movido, mas não excluído.
      </p>
      <Field label="Rótulo do hash">
        <input
          class="h-9 w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
          value={props.element.hashLabel}
          onInput={(event) =>
            updateHash(props.element.id, { hashLabel: event.currentTarget.value })
          }
        />
      </Field>
      <Field label="Texto do link">
        <input
          class="h-9 w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
          value={props.element.linkLabel}
          onInput={(event) =>
            updateHash(props.element.id, { linkLabel: event.currentTarget.value })
          }
        />
      </Field>
      <div class="grid grid-cols-2 gap-2">
        <NumberProperty
          label="Tamanho"
          value={props.element.fontSize}
          min={6}
          max={200}
          onCommit={(fontSize) => updateHash(props.element.id, { fontSize })}
        />
        <Field label="Alinhamento">
          <select
            class="h-9 w-full rounded-md border border-border bg-background px-3 text-sm focus:outline-none"
            value={props.element.align}
            onChange={(event) =>
              updateHash(props.element.id, {
                align: event.currentTarget.value as HashCertificateElement["align"],
              })
            }
          >
            <option value="left">Esquerda</option>
            <option value="center">Centro</option>
            <option value="right">Direita</option>
          </select>
        </Field>
      </div>
      <Field label="Cor">
        <input
          type="color"
          value={normalizeHexColor(props.element.color)}
          class="h-9 w-full cursor-pointer rounded-md border border-border bg-background p-1"
          onInput={(event) =>
            updateHash(props.element.id, { color: event.currentTarget.value })
          }
        />
      </Field>
      <div class="border-t border-border" />
      <PositionSizeProperties element={props.element} />
    </div>
  );
}

function updateImage(
  id: string,
  patch: Partial<Omit<ImageCertificateElement, "id" | "type">>,
) {
  certificateEditorActions.updateElement(id, (element: CertificationTemplateElement) =>
    element.type === "image" ? { ...element, ...patch } : element,
  );
}

function ImageProperties(props: { element: ImageCertificateElement }): JSX.Element {
  let inputRef: HTMLInputElement | undefined;
  const [error, setError] = createSignal<string | null>(null);

  async function replaceImage(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    input.value = "";
    if (!file) return;
    if (!isSupportedCertificateImage(file)) {
      setError("Use uma imagem PNG, JPEG ou WebP.");
      return;
    }

    try {
      updateImage(props.element.id, { src: await readCertificateFile(file) });
      setError(null);
    } catch {
      setError("Não foi possível carregar a imagem.");
    }
  }

  return (
    <div class="space-y-4">
      <div class="overflow-hidden rounded-md border bg-white p-2">
        <img
          src={props.element.src}
          alt=""
          class="max-h-32 w-full object-contain"
        />
      </div>
      <Button
        type="button"
        variant="outline"
        class="w-full cursor-pointer"
        onClick={() => inputRef?.click()}
      >
        Substituir imagem
      </Button>
      <input
        ref={(el) => (inputRef = el)}
        type="file"
        accept={CERTIFICATE_IMAGE_ACCEPT}
        class="hidden"
        onChange={(event) => void replaceImage(event)}
      />
      <Show when={error()}>
        <p class="text-xs text-destructive">{error()}</p>
      </Show>
      <Field label="Ajuste">
        <select
          class="h-9 w-full rounded-md border border-border bg-background px-3 text-sm focus:outline-none"
          value={props.element.fit}
          onChange={(event) =>
            updateImage(props.element.id, {
              fit: event.currentTarget.value as ImageCertificateElement["fit"],
            })
          }
        >
          <option value="cover">Preencher</option>
          <option value="contain">Ajustar</option>
          <option value="fill">Esticar</option>
        </select>
      </Field>
      <div class="grid grid-cols-2 gap-2">
        <NumberProperty
          label="Arredondamento"
          value={props.element.radius}
          min={0}
          max={400}
          onCommit={(radius) => updateImage(props.element.id, { radius })}
        />
        <NumberProperty
          label="Opacidade (%)"
          value={props.element.opacity * 100}
          min={0}
          max={100}
          onCommit={(opacity) =>
            updateImage(props.element.id, { opacity: opacity / 100 })
          }
        />
      </div>
      <div class="border-t border-border" />
      <PositionSizeProperties element={props.element} />
    </div>
  );
}

function updateSignature(
  id: string,
  patch: Partial<Omit<SignatureCertificateElement, "id" | "type">>,
) {
  certificateEditorActions.updateElement(id, (element: CertificationTemplateElement) =>
    element.type === "signature" ? { ...element, ...patch } : element,
  );
}

function SignatureProperties(props: {
  element: SignatureCertificateElement;
}): JSX.Element {
  const signatures = () => certificateEditorStore.availableSignatures();

  return (
    <div class="space-y-4">
      <Field label="Trocar assinatura">
        <div class="grid grid-cols-2 gap-2">
          <For each={signatures()}>
            {(signature) => (
              <button
                type="button"
                title={signature.name}
                onClick={() =>
                  updateSignature(props.element.id, {
                    signatureId: signature.id,
                    src: signature.url,
                    name: signature.name,
                  })
                }
                class={
                  "overflow-hidden rounded-md border text-left cursor-pointer transition-colors " +
                  (props.element.signatureId === signature.id
                    ? "border-ring ring-2 ring-ring/30"
                    : "hover:border-ring")
                }
              >
                <span class="flex h-14 items-center justify-center bg-white p-1.5">
                  <img
                    src={signature.url}
                    alt=""
                    class="max-h-full max-w-full object-contain"
                  />
                </span>
                <span class="block truncate border-t px-1.5 py-1 text-[11px]">
                  {signature.name}
                </span>
              </button>
            )}
          </For>
        </div>
      </Field>
      <Field label="Ajuste">
        <ToolbarCombobox
          value={props.element.fit}
          options={[
            { value: "contain", label: "Ajustar" },
            { value: "cover", label: "Preencher" },
            { value: "fill", label: "Esticar" },
          ]}
          placeholder="Ajuste"
          class="w-full"
          onChange={(val) =>
            updateSignature(props.element.id, {
              fit: val as SignatureCertificateElement["fit"],
            })
          }
        />
      </Field>
      <div class="grid grid-cols-2 gap-2">
        <NumberProperty
          label="Arredondamento"
          value={props.element.radius}
          min={0}
          max={400}
          onCommit={(radius) => updateSignature(props.element.id, { radius })}
        />
        <NumberProperty
          label="Opacidade (%)"
          value={props.element.opacity * 100}
          min={0}
          max={100}
          onCommit={(opacity) =>
            updateSignature(props.element.id, { opacity: opacity / 100 })
          }
        />
      </div>
      <div class="border-t border-border" />
      <PositionSizeProperties element={props.element} />
    </div>
  );
}

function TextProperties(props: { element: TextCertificateElement }): JSX.Element {
  const controller = () => certificateEditorStore.richTextController();

  return (
    <div class="space-y-4">
      <p class="rounded-md bg-muted p-2.5 text-xs leading-relaxed text-muted-foreground">
        Dê um duplo clique no texto para editar. A formatação da seleção aparece
        na barra acima do canvas.
      </p>
      <div class="border-t border-border" />
      <PositionSizeProperties element={props.element} />
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
          <For each={CERTIFICATE_VARIABLES}>
            {(variable) => (
              <Button
                type="button"
                variant="outline"
                class="h-auto w-full justify-start whitespace-normal px-3 py-2 text-left text-xs cursor-pointer"
                disabled={!controller()}
                title={variable.description}
                onClick={() => controller()?.insertText(variable.token)}
              >
                <span class="min-w-0">
                  <span class="block truncate font-medium">{variable.label}</span>
                  <span class="mt-0.5 block text-[11px] font-normal leading-tight text-muted-foreground">
                    {variable.description}
                  </span>
                </span>
              </Button>
            )}
          </For>
        </div>
      </div>
    </div>
  );
}

export function CertificatePropertiesPanel(): JSX.Element {
  const selectedElementId = () =>
    certificateEditorStore.selectedElementId();
  const elements = () =>
    certificateEditorStore.draft().design_data.elements ?? [];
  const selectedElement = createMemo(() =>
    elements().find(
      (item: CertificationTemplateElement) => item.id === selectedElementId(),
    ) ?? null,
  );

  return (
    <aside class="w-72 shrink-0 overflow-y-auto border-l border-border bg-card p-4 text-card-foreground">
      <h2 class="mb-4 text-xs font-semibold tracking-wide text-muted-foreground uppercase">
        Propriedades
      </h2>
      <Show
        when={selectedElement()}
        fallback={
          <p class="text-sm text-muted-foreground">
            Selecione um elemento no certificado para editar suas propriedades.
          </p>
        }
      >
        {(el) => (
          <>
            <Show when={el().type === "hash"}>
              <HashProperties element={el() as HashCertificateElement} />
            </Show>
            <Show when={el().type === "text"}>
              <TextProperties element={el() as TextCertificateElement} />
            </Show>
            <Show when={el().type === "image"}>
              <ImageProperties element={el() as ImageCertificateElement} />
            </Show>
            <Show when={el().type === "signature"}>
              <SignatureProperties element={el() as SignatureCertificateElement} />
            </Show>
          </>
        )}
      </Show>
    </aside>
  );
}
