import type { ChangeEvent, ReactNode } from "react";
import { useEffect, useRef, useState } from "react";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { Separator } from "@/shared/ui/shadcn/separator";
import type { CertificationTemplateElement } from "../../model";
import {
  CERTIFICATE_IMAGE_ACCEPT,
  MIN_CERTIFICATE_ELEMENT_SIZE,
} from "../constants";
import { certificateEditorActions, useCertificateEditorState } from "../store";
import type {
  HashCertificateElement,
  ImageCertificateElement,
  SignatureCertificateElement,
  TextCertificateElement,
} from "../types";
import { isSupportedCertificateImage, readCertificateFile } from "../utils";
import { ToolbarCombobox } from "./toolbar-combobox";

interface FieldProps {
  label: string;
  children: ReactNode;
}

function Field({ label, children }: FieldProps) {
  return (
    <div className="block space-y-1">
      <span className="text-xs text-muted-foreground">{label}</span>
      {children}
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

function NumberProperty({
  label,
  value,
  min,
  max,
  onCommit,
}: NumberPropertyProps) {
  const [draft, setDraft] = useState(String(Math.round(value)));
  useEffect(() => setDraft(String(Math.round(value))), [value]);

  function commit() {
    const parsed = Number(draft);
    if (
      Number.isFinite(parsed) &&
      (min === undefined || parsed >= min) &&
      (max === undefined || parsed <= max)
    ) {
      onCommit(parsed);
      return;
    }
    setDraft(String(Math.round(value)));
  }

  return (
    <Field label={label}>
      <Input
        type="number"
        value={draft}
        min={min}
        max={max}
        onChange={(event) => setDraft(event.target.value)}
        onBlur={commit}
        onKeyDown={(event) => {
          if (event.key === "Enter") event.currentTarget.blur();
        }}
      />
    </Field>
  );
}

function PositionSizeProperties({
  element,
}: {
  element: CertificationTemplateElement;
}) {
  return (
    <div className="grid grid-cols-2 gap-2">
      <NumberProperty
        label="Posição X"
        value={element.x}
        onCommit={(x) =>
          certificateEditorActions.updateElementBounds(element.id, { x })
        }
      />
      <NumberProperty
        label="Posição Y"
        value={element.y}
        onCommit={(y) =>
          certificateEditorActions.updateElementBounds(element.id, { y })
        }
      />
      <NumberProperty
        label="Largura"
        value={element.width}
        min={MIN_CERTIFICATE_ELEMENT_SIZE.width}
        onCommit={(width) =>
          certificateEditorActions.updateElementBounds(element.id, { width })
        }
      />
      <NumberProperty
        label="Altura"
        value={element.height}
        min={MIN_CERTIFICATE_ELEMENT_SIZE.height}
        onCommit={(height) =>
          certificateEditorActions.updateElementBounds(element.id, { height })
        }
      />
    </div>
  );
}

function updateHash(
  id: string,
  patch: Partial<Omit<HashCertificateElement, "id" | "type">>,
) {
  certificateEditorActions.updateElement(id, (element) =>
    element.type === "hash" ? { ...element, ...patch } : element,
  );
}

function HashProperties({ element }: { element: HashCertificateElement }) {
  return (
    <div className="space-y-4">
      <p className="rounded-md bg-muted p-2.5 text-xs leading-relaxed text-muted-foreground">
        O código e o endereço de validação são preenchidos automaticamente. O
        bloco pode ser estilizado e movido, mas não excluído.
      </p>
      <Field label="Rótulo do hash">
        <Input
          value={element.hashLabel}
          onChange={(event) =>
            updateHash(element.id, { hashLabel: event.target.value })
          }
        />
      </Field>
      <Field label="Texto do link">
        <Input
          value={element.linkLabel}
          onChange={(event) =>
            updateHash(element.id, { linkLabel: event.target.value })
          }
        />
      </Field>
      <div className="grid grid-cols-2 gap-2">
        <NumberProperty
          label="Tamanho"
          value={element.fontSize}
          min={6}
          max={200}
          onCommit={(fontSize) => updateHash(element.id, { fontSize })}
        />
        <Field label="Alinhamento">
          <select
            className="h-10 w-full rounded-lg border border-input bg-background px-3 text-sm"
            value={element.align}
            onChange={(event) =>
              updateHash(element.id, {
                align: event.target.value as HashCertificateElement["align"],
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
        <Input
          type="color"
          value={element.color.slice(0, 7)}
          className="p-1"
          onChange={(event) =>
            updateHash(element.id, { color: event.target.value })
          }
        />
      </Field>
      <Separator />
      <PositionSizeProperties element={element} />
    </div>
  );
}

function updateImage(
  id: string,
  patch: Partial<Omit<ImageCertificateElement, "id" | "type">>,
) {
  certificateEditorActions.updateElement(id, (element) =>
    element.type === "image" ? { ...element, ...patch } : element,
  );
}

function ImageProperties({ element }: { element: ImageCertificateElement }) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [error, setError] = useState<string | null>(null);

  async function replaceImage(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    event.target.value = "";
    if (!file) return;
    if (!isSupportedCertificateImage(file)) {
      setError("Use uma imagem PNG, JPEG ou WebP.");
      return;
    }

    try {
      updateImage(element.id, { src: await readCertificateFile(file) });
      setError(null);
    } catch {
      setError("Não foi possível carregar a imagem.");
    }
  }

  return (
    <div className="space-y-4">
      <div className="overflow-hidden rounded-md border bg-white p-2">
        <img
          src={element.src}
          alt=""
          className="max-h-32 w-full object-contain"
        />
      </div>
      <Button
        type="button"
        variant="outline"
        className="w-full"
        onClick={() => inputRef.current?.click()}
      >
        Substituir imagem
      </Button>
      <input
        ref={inputRef}
        type="file"
        accept={CERTIFICATE_IMAGE_ACCEPT}
        className="hidden"
        onChange={(event) => void replaceImage(event)}
      />
      {error ? <p className="text-xs text-destructive">{error}</p> : null}
      <Field label="Ajuste">
        <select
          className="h-10 w-full rounded-lg border border-input bg-background px-3 text-sm"
          value={element.fit}
          onChange={(event) =>
            updateImage(element.id, {
              fit: event.target.value as ImageCertificateElement["fit"],
            })
          }
        >
          <option value="cover">Preencher</option>
          <option value="contain">Ajustar</option>
          <option value="fill">Esticar</option>
        </select>
      </Field>
      <div className="grid grid-cols-2 gap-2">
        <NumberProperty
          label="Arredondamento"
          value={element.radius}
          min={0}
          max={400}
          onCommit={(radius) => updateImage(element.id, { radius })}
        />
        <NumberProperty
          label="Opacidade (%)"
          value={element.opacity * 100}
          min={0}
          max={100}
          onCommit={(opacity) =>
            updateImage(element.id, { opacity: opacity / 100 })
          }
        />
      </div>
      <Separator />
      <PositionSizeProperties element={element} />
    </div>
  );
}

function updateSignature(
  id: string,
  patch: Partial<Omit<SignatureCertificateElement, "id" | "type">>,
) {
  certificateEditorActions.updateElement(id, (element) =>
    element.type === "signature" ? { ...element, ...patch } : element,
  );
}

function SignatureProperties({
  element,
}: {
  element: SignatureCertificateElement;
}) {
  const signatures = useCertificateEditorState(
    (state) => state.availableSignatures,
  );

  return (
    <div className="space-y-4">
      <Field label="Trocar assinatura">
        <div className="grid grid-cols-2 gap-2">
          {signatures.map((signature) => (
            <button
              key={signature.id}
              type="button"
              title={signature.name}
              onClick={() =>
                updateSignature(element.id, {
                  signatureId: signature.id,
                  src: signature.url,
                  name: signature.name,
                })
              }
              className={
                "overflow-hidden rounded-md border text-left " +
                (element.signatureId === signature.id
                  ? "border-ring ring-2 ring-ring/30"
                  : "hover:border-ring")
              }
            >
              <span className="flex h-14 items-center justify-center bg-white p-1.5">
                <img
                  src={signature.url}
                  alt=""
                  className="max-h-full max-w-full object-contain"
                />
              </span>
              <span className="block truncate border-t px-1.5 py-1 text-[11px]">
                {signature.name}
              </span>
            </button>
          ))}
        </div>
      </Field>
      <Field label="Ajuste">
        <ToolbarCombobox
          value={element.fit}
          options={[
            { value: "contain", label: "Ajustar" },
            { value: "cover", label: "Preencher" },
            { value: "fill", label: "Esticar" },
          ]}
          placeholder="Ajuste"
          className="w-full"
          onChange={(event) =>
            updateSignature(element.id, {
              fit: event as SignatureCertificateElement["fit"],
            })
          }
        />
      </Field>
      <div className="grid grid-cols-2 gap-2">
        <NumberProperty
          label="Arredondamento"
          value={element.radius}
          min={0}
          max={400}
          onCommit={(radius) => updateSignature(element.id, { radius })}
        />
        <NumberProperty
          label="Opacidade (%)"
          value={element.opacity * 100}
          min={0}
          max={100}
          onCommit={(opacity) =>
            updateSignature(element.id, { opacity: opacity / 100 })
          }
        />
      </div>
      <Separator />
      <PositionSizeProperties element={element} />
    </div>
  );
}

function TextProperties({ element }: { element: TextCertificateElement }) {
  return (
    <div className="space-y-4">
      <p className="rounded-md bg-muted p-2.5 text-xs leading-relaxed text-muted-foreground">
        Dê um duplo clique no texto para editar. A formatação da seleção aparece
        na barra acima do canvas, onde também é possível inserir informações
        dinâmicas.
      </p>
      <Separator />
      <PositionSizeProperties element={element} />
    </div>
  );
}

export function CertificatePropertiesPanel() {
  const selectedElementId = useCertificateEditorState(
    (state) => state.selectedElementId,
  );
  const elements = useCertificateEditorState(
    (state) => state.draft.design_data.elements,
  );
  const element =
    elements.find((item) => item.id === selectedElementId) ?? null;

  return (
    <aside className="w-72 shrink-0 overflow-y-auto border-l border-border bg-card p-4 text-card-foreground">
      <h2 className="mb-4 text-xs font-semibold tracking-wide text-muted-foreground uppercase">
        Propriedades
      </h2>
      {!element ? (
        <p className="text-sm text-muted-foreground">
          Selecione um elemento no certificado para editar suas propriedades.
        </p>
      ) : element.type === "hash" ? (
        <HashProperties element={element} />
      ) : element.type === "text" ? (
        <TextProperties element={element} />
      ) : element.type === "image" ? (
        <ImageProperties element={element} />
      ) : (
        <SignatureProperties element={element} />
      )}
    </aside>
  );
}
