import { Download, Eye, FileImage, FileText, LoaderCircle } from "lucide-react";
import { forwardRef, useMemo, useRef, useState } from "react";
import { toast } from "sonner";
import { Button } from "@/shared/ui/shadcn/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/shared/ui/shadcn/dialog";
import { DEFAULT_CERTIFICATE_CANVAS } from "../editor/constants";
import { useElementSize } from "../editor/hooks/use-element-size";
import { CertificateElementView } from "../editor/ui/elements/certificate-element-view";
import {
  type CertificateVariableValues,
  resolveCertificationTemplate,
} from "../editor/variables";
import type { CertificateExportFormat } from "../export/certificate-export";
import type { CertificationTemplateI } from "../model";

interface CertViewerProps {
  template: CertificationTemplateI;
  triggerLabel?: string;
  variables?: CertificateVariableValues;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export function CertViewer({
  template,
  triggerLabel = "Visualizar",
  variables = {},
  open,
  onOpenChange,
}: CertViewerProps) {
  const canvasRef = useRef<HTMLDivElement>(null);
  const [exporting, setExporting] = useState<CertificateExportFormat | null>(
    null,
  );

  async function exportCertificate(format: CertificateExportFormat) {
    if (!canvasRef.current || exporting) return;
    setExporting(format);
    try {
      const { downloadCertificate } = await import(
        "../export/certificate-export"
      );
      await downloadCertificate(canvasRef.current, template.name, format);
    } catch {
      toast.error(
        `Não foi possível exportar o certificado em ${format.toUpperCase()}.`,
      );
    } finally {
      setExporting(null);
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      {open === undefined ? (
        <DialogTrigger
          render={<Button type="button" variant="outline" size="sm" />}
        >
          <Eye className="size-4" />
          {triggerLabel}
        </DialogTrigger>
      ) : null}
      <DialogContent
        className="z-100! grid-rows-[auto_minmax(0,1fr)] gap-0 overflow-hidden p-0 shadow-2xl"
        overlayClassName="!z-[99] bg-black/40 backdrop-blur-md"
        style={{
          width: "min(72rem, calc(100vw - 2rem))",
          maxWidth: "calc(100vw - 2rem)",
          height: "min(48rem, calc(100dvh - 2rem))",
        }}
      >
        <DialogHeader className="items-stretch justify-between gap-3 border-b bg-background px-5 py-3 pr-14 sm:flex-row sm:items-center">
          <div className="min-w-0 space-y-1">
            <DialogTitle className="flex items-center gap-2 truncate">
              <FileText className="size-4 shrink-0 text-muted-foreground" />
              {template.name}
            </DialogTitle>
            <DialogDescription>
              Pré-visualização no tamanho e proporção de emissão.
            </DialogDescription>
          </div>
          <div className="flex shrink-0 items-center gap-2">
            <Button
              type="button"
              size="sm"
              variant="outline"
              disabled={exporting !== null}
              onClick={() => void exportCertificate("png")}
            >
              {exporting === "png" ? (
                <LoaderCircle className="size-4 animate-spin" />
              ) : (
                <FileImage className="size-4" />
              )}
              PNG
            </Button>
            <Button
              type="button"
              size="sm"
              disabled={exporting !== null}
              onClick={() => void exportCertificate("pdf")}
            >
              {exporting === "pdf" ? (
                <LoaderCircle className="size-4 animate-spin" />
              ) : (
                <Download className="size-4" />
              )}
              PDF
            </Button>
          </div>
        </DialogHeader>
        <div className="min-h-0 overflow-hidden bg-muted/60 p-3 sm:p-5">
          <CertificateTemplateStaticView
            ref={canvasRef}
            template={template}
            variables={variables}
          />
        </div>
      </DialogContent>
    </Dialog>
  );
}

interface CertificateTemplateStaticViewProps {
  template: CertificationTemplateI;
  variables?: CertificateVariableValues;
}

export const CertificateTemplateStaticView = forwardRef<
  HTMLDivElement,
  CertificateTemplateStaticViewProps
>(function CertificateTemplateStaticView(
  { template, variables = {} },
  canvasRef,
) {
  const { ref, size } = useElementSize<HTMLDivElement>();
  const canvas = template.design_data.canvas ?? DEFAULT_CERTIFICATE_CANVAS;
  const scale = Math.max(
    0,
    Math.min(size.width / canvas.width, size.height / canvas.height),
  );
  const resolvedTemplate = useMemo(
    () => resolveCertificationTemplate(template, variables),
    [template, variables],
  );
  const backgroundUrl = resolvedTemplate.design_data.background;

  return (
    <div
      ref={ref}
      className="flex h-full w-full items-center justify-center overflow-hidden"
    >
      {scale > 0 ? (
        <div
          className="relative shrink-0 overflow-hidden bg-white shadow-2xl ring-1 ring-black/10"
          style={{
            width: canvas.width * scale,
            height: canvas.height * scale,
          }}
        >
          <div
            ref={canvasRef}
            className="relative origin-top-left overflow-hidden bg-white"
            style={{
              width: canvas.width,
              height: canvas.height,
              transform: `scale(${scale})`,
              backgroundImage: backgroundUrl
                ? `url(${JSON.stringify(backgroundUrl)})`
                : undefined,
              backgroundPosition: "center",
              backgroundRepeat: "no-repeat",
              backgroundSize: "cover",
            }}
          >
            {resolvedTemplate.design_data.elements.map((element) => (
              <div
                key={element.id}
                className="absolute overflow-hidden"
                style={{
                  left: element.x,
                  top: element.y,
                  width: element.width,
                  height: element.height,
                }}
              >
                <CertificateElementView element={element} />
              </div>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  );
});
