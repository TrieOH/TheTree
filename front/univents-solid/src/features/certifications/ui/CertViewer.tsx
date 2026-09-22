import type { JSX } from "@solidjs/web";
import { For, Show, createMemo, createSignal } from "solid-js";

import DownloadIcon from "~icons/lucide/download";
import FileImageIcon from "~icons/lucide/file-image";
import FileTextIcon from "~icons/lucide/file-text";
import LoaderCircleIcon from "~icons/lucide/loader-circle";

import { Button, Dialog } from "@trieoh/ui-solid";
import { toast } from "@/shared/ui/toast";
import { DEFAULT_CERTIFICATE_CANVAS } from "../editor/constants";
import { useElementSize } from "../editor/hooks/use-element-size";
import { CertificateElementView } from "../editor/ui/elements/certificate-element-view";
import {
  type CertificateVariableValues,
  resolveCertificationTemplate,
} from "../editor/variables";
import {
  type CertificateExportFormat,
  downloadCertificate,
} from "../export/certificate-export";
import type { CertificationTemplateI } from "../model";

const Download = DownloadIcon as unknown as (props: { class?: string }) => JSX.Element;
const FileImage = FileImageIcon as unknown as (props: { class?: string }) => JSX.Element;
const FileText = FileTextIcon as unknown as (props: { class?: string }) => JSX.Element;
const LoaderCircle = LoaderCircleIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface CertViewerProps {
  template: CertificationTemplateI;
  variables?: CertificateVariableValues;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CertificateDownloadButtons(props: {
  canvasElement?: () => HTMLDivElement | undefined;
  templateName: string;
  prominent?: boolean;
}): JSX.Element {
  const [exporting, setExporting] = createSignal<CertificateExportFormat | null>(null);

  async function download(format: CertificateExportFormat) {
    const el = props.canvasElement?.();
    if (!el || exporting()) return;
    setExporting(format);
    try {
      await downloadCertificate(el, props.templateName, format);
    } catch {
      toast.error(
        `Não foi possível exportar o certificado em ${format.toUpperCase()}.`,
      );
    } finally {
      setExporting(null);
    }
  }

  return (
    <div
      class={
        props.prominent
          ? "flex flex-col gap-2 sm:flex-row sm:items-center"
          : "flex items-center gap-2"
      }
    >
      <Button
        type="button"
        size={props.prominent ? "md" : "sm"}
        variant="outline"
        disabled={exporting() !== null}
        onClick={() => void download("png")}
        class="cursor-pointer"
      >
        <Show
          when={exporting() === "png"}
          fallback={<FileImage class={props.prominent ? "size-5" : "size-4"} />}
        >
          <LoaderCircle
            class={props.prominent ? "size-5 animate-spin" : "size-4 animate-spin"}
          />
        </Show>
        <span>{props.prominent ? "Baixar PNG" : "PNG"}</span>
      </Button>
      <Button
        type="button"
        size={props.prominent ? "md" : "sm"}
        disabled={exporting() !== null}
        onClick={() => void download("pdf")}
        class="cursor-pointer"
      >
        <Show
          when={exporting() === "pdf"}
          fallback={<Download class={props.prominent ? "size-5" : "size-4"} />}
        >
          <LoaderCircle
            class={props.prominent ? "size-5 animate-spin" : "size-4 animate-spin"}
          />
        </Show>
        <span>{props.prominent ? "Baixar PDF" : "PDF"}</span>
      </Button>
    </div>
  );
}

export function CertificateTemplateStaticView(props: {
  template: CertificationTemplateI;
  variables?: CertificateVariableValues;
  setCanvasRef?: (el: HTMLDivElement) => void;
  overlay?: JSX.Element;
}): JSX.Element {
  const { ref, size } = useElementSize<HTMLDivElement>();
  const canvas = () =>
    props.template.design_data?.canvas ?? DEFAULT_CERTIFICATE_CANVAS;

  const scale = createMemo(() => {
    const s = size();
    const c = canvas();
    if (s.width <= 0 || s.height <= 0) return 0;
    return Math.max(0, Math.min(s.width / c.width, s.height / c.height));
  });

  const resolvedTemplate = createMemo(() =>
    resolveCertificationTemplate(props.template, props.variables ?? {}),
  );

  const backgroundUrl = () => resolvedTemplate().design_data?.background;

  return (
    <div
      ref={ref}
      class="flex h-full w-full items-center justify-center overflow-hidden"
    >
      <Show when={scale() > 0}>
        <div
          class="relative shrink-0 overflow-hidden bg-white shadow-[inset_0_0_0_1px_rgba(0,0,0,0.14),0_0_12px_rgba(0,0,0,0.10)]"
          style={{
            width: `${canvas().width * scale()}px`,
            height: `${canvas().height * scale()}px`,
          }}
        >
          <div
            ref={props.setCanvasRef}
            class="relative origin-top-left overflow-hidden bg-white"
            style={{
              width: `${canvas().width}px`,
              height: `${canvas().height}px`,
              transform: `scale(${scale()})`,
              ...(backgroundUrl()
                ? {
                    "background-image": `url(${backgroundUrl()})`,
                    "background-position": "center",
                    "background-repeat": "no-repeat",
                    "background-size": "cover",
                  }
                : {}),
            }}
          >
            <For each={resolvedTemplate().design_data?.elements ?? []}>
              {(element) => (
                <div
                  class="absolute overflow-hidden"
                  style={{
                    left: `${element.x}px`,
                    top: `${element.y}px`,
                    width: `${element.width}px`,
                    height: `${element.height}px`,
                  }}
                >
                  <CertificateElementView element={element} />
                </div>
              )}
            </For>
          </div>
          <Show when={props.overlay}>{props.overlay}</Show>
        </div>
      </Show>
    </div>
  );
}

export function CertViewer(props: CertViewerProps): JSX.Element {
  let canvasRef: HTMLDivElement | undefined;
  const [exporting, setExporting] = createSignal<CertificateExportFormat | null>(null);

  async function exportCertificate(format: CertificateExportFormat) {
    if (!canvasRef || exporting()) return;
    setExporting(format);
    try {
      await downloadCertificate(canvasRef, props.template.name, format);
    } catch {
      toast.error(
        `Não foi possível exportar o certificado em ${format.toUpperCase()}.`,
      );
    } finally {
      setExporting(null);
    }
  }

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={
        <div class="flex items-center gap-2">
          <FileText class="size-4 shrink-0 text-muted-foreground" />
          <span class="truncate">{props.template.name}</span>
        </div>
      }
      description="Pré-visualização no tamanho e proporção de emissão."
      class="sm:max-w-5xl w-[95vw] h-[85vh] p-0 flex flex-col [&>div:nth-child(2)]:flex-1 [&>div:nth-child(2)]:min-h-0 [&>div:nth-child(2)]:p-0"
      footer={
        <div class="flex w-full items-center justify-between">
          <Button
            type="button"
            variant="outline"
            onClick={() => props.onOpenChange(false)}
            class="cursor-pointer"
          >
            Fechar
          </Button>

          <div class="flex items-center gap-2">
            <Button
              type="button"
              size="sm"
              variant="outline"
              disabled={exporting() !== null}
              onClick={() => void exportCertificate("png")}
              class="cursor-pointer"
            >
              <Show
                when={exporting() === "png"}
                fallback={<FileImage class="size-4" />}
              >
                <LoaderCircle class="size-4 animate-spin" />
              </Show>
              <span>PNG</span>
            </Button>
            <Button
              type="button"
              size="sm"
              disabled={exporting() !== null}
              onClick={() => void exportCertificate("pdf")}
              class="cursor-pointer"
            >
              <Show
                when={exporting() === "pdf"}
                fallback={<Download class="size-4" />}
              >
                <LoaderCircle class="size-4 animate-spin" />
              </Show>
              <span>PDF</span>
            </Button>
          </div>
        </div>
      }
    >
      <div class="h-full w-full overflow-hidden bg-muted/60 p-3 sm:p-5">
        <CertificateTemplateStaticView
          template={props.template}
          variables={props.variables}
          setCanvasRef={(el) => (canvasRef = el)}
        />
      </div>
    </Dialog>
  );
}
