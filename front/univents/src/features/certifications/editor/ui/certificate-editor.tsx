import { useQuery } from "@tanstack/react-query";
import { Link, useNavigate } from "@tanstack/react-router";
import { ArrowLeft, Loader2, Monitor, Save } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { allSignaturesQueryOptions } from "../../../signatures/api";
import { useCreateCertificationTemplateMutation } from "../../api/mutations";
import { certificationTemplateCreateSchema } from "../../model";
import { certificateEditorActions, useCertificateEditorState } from "../store";
import { uploadCertificateAssets } from "../upload-assets";
import { CertificateCanvas } from "./certificate-canvas";
import { CertificatePropertiesPanel } from "./certificate-properties-panel";
import { CertificateTextToolbar } from "./certificate-text-toolbar";
import { CertificateToolsSidebar } from "./certificate-tools-sidebar";

interface CertificateEditorProps {
  eventId: string;
  editionId: string;
}

export function CertificateEditor({
  eventId,
  editionId,
}: CertificateEditorProps) {
  const navigate = useNavigate();
  const [uploadingAssets, setUploadingAssets] = useState(false);
  const title = useCertificateEditorState((state) => state.draft.title);
  const { data: signatures = [] } = useQuery(
    allSignaturesQueryOptions(eventId, editionId),
  );

  useEffect(() => {
    certificateEditorActions.reset();
    return () => certificateEditorActions.reset();
  }, []);

  useEffect(() => {
    certificateEditorActions.setAvailableSignatures(
      signatures.map((signature) => ({
        id: signature.id,
        url: signature.url,
        name: signature.title,
      })),
    );
  }, [signatures]);

  const createTemplate = useCreateCertificationTemplateMutation();

  async function saveTemplate() {
    const result = certificationTemplateCreateSchema.safeParse(
      certificateEditorActions.getDraft(),
    );
    if (!result.success) {
      toast.error(result.error.issues[0]?.message ?? "Template inválido");
      return;
    }
    setUploadingAssets(true);
    try {
      const data = await uploadCertificateAssets(
        result.data,
        eventId,
        editionId,
      );
      createTemplate.mutate(
        { eventId, editionId, data },
        {
          onSuccess: (response) => {
            if (!response.success) return;
            void navigate({
              to: "/admin/events/$eventId/editions/$editionId/certifications",
              params: { eventId, editionId },
            });
          },
        },
      );
    } catch {
      toast.error("Não foi possível enviar as imagens do certificado");
    } finally {
      setUploadingAssets(false);
    }
  }

  return (
    <div className="h-dvh min-h-0 bg-background text-foreground">
      <div className="flex h-full flex-col items-center justify-center gap-3 px-6 text-center lg:hidden!">
        <Monitor className="size-10 text-muted-foreground" />
        <p className="max-w-xs text-sm text-muted-foreground">
          O editor de certificados foi desenvolvido para telas maiores. Abra
          esta página em um computador para editar o template.
        </p>
      </div>

      <div className="hidden h-full min-h-0 flex-col lg:flex">
        <header className="flex h-16 shrink-0 items-center justify-between gap-4 border-b border-border/70 bg-card/95 px-4 text-card-foreground shadow-sm backdrop-blur">
          <div className="flex min-w-0 flex-1 items-center gap-3">
            <Link
              to="/admin/events/$eventId/editions/$editionId/certifications"
              params={{ eventId, editionId }}
              aria-label="Voltar para certificados"
              className="inline-flex size-8 shrink-0 items-center justify-center rounded-lg border border-border/60 bg-background/80 hover:bg-accent/10"
            >
              <ArrowLeft className="size-4" />
            </Link>
            <span className="hidden shrink-0 text-sm font-semibold tracking-tight xl:inline">
              Editor de certificados
            </span>
            <Input
              value={title}
              maxLength={160}
              onChange={(event) =>
                certificateEditorActions.setTitle(event.target.value)
              }
              placeholder="Título do certificado"
              aria-label="Título do certificado"
              className="max-w-sm border-border/60 bg-background/80 text-card-foreground placeholder:text-muted-foreground/70"
            />
          </div>

          <Button
            type="button"
            variant="default"
            disabled={uploadingAssets || createTemplate.isPending}
            onClick={() => void saveTemplate()}
          >
            {uploadingAssets || createTemplate.isPending ? (
              <Loader2 className="size-4 animate-spin" />
            ) : (
              <Save className="size-4" />
            )}
            Salvar
          </Button>
        </header>

        <div className="flex min-h-0 flex-1">
          <CertificateToolsSidebar />
          <div className="flex min-w-0 flex-1 flex-col">
            <CertificateTextToolbar />
            <CertificateCanvas />
          </div>
          <CertificatePropertiesPanel />
        </div>
      </div>
    </div>
  );
}
