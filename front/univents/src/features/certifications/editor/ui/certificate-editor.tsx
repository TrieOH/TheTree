import { useQuery } from "@tanstack/react-query";
import { Link, useNavigate } from "@tanstack/react-router";
import { ArrowLeft, Monitor, Save } from "lucide-react";
import { useEffect } from "react";
import { toast } from "sonner";
import {
  useCreateCertificationTemplateMutation,
  useUpdateCertificationTemplateMutation,
} from "@/features/certifications/api/mutations";
import { Button } from "@/shared/ui/shadcn/button";
import { allSignaturesQueryOptions } from "../../../signatures/api";
import { certificationTemplateQueryOptions } from "../../api";
import { certificationTemplateCreateSchema } from "../../model";
import { certificateEditorActions } from "../store";
import { uploadCertificateAssets } from "../upload-assets";
import { CertificateCanvas } from "./certificate-canvas";
import { CertificatePropertiesPanel } from "./certificate-properties-panel";
import { CertificateTextToolbar } from "./certificate-text-toolbar";
import { CertificateToolsSidebar } from "./certificate-tools-sidebar";

interface CertificateEditorProps {
  eventId: string;
  editionId: string;
  templateId?: string;
  duplicate?: boolean;
}

export function CertificateEditor({
  eventId,
  editionId,
  templateId,
  duplicate = false,
}: CertificateEditorProps) {
  const { data: signatures = [] } = useQuery(
    allSignaturesQueryOptions(editionId),
  );
  const templateQuery = useQuery({
    ...certificationTemplateQueryOptions(templateId ?? ""),
    enabled: Boolean(templateId),
  });

  useEffect(() => {
    certificateEditorActions.reset();
    if (templateQuery.data) {
      certificateEditorActions.loadDraft(templateQuery.data);
      if (duplicate)
        certificateEditorActions.setName(`${templateQuery.data.name} (cópia)`);
    }
    return () => certificateEditorActions.reset();
  }, [duplicate, templateQuery.data]);

  useEffect(() => {
    certificateEditorActions.setAvailableSignatures(
      signatures.map((signature) => ({
        id: signature.id,
        url: signature.image_url,
        name: signature.signatory_name,
      })),
    );
  }, [signatures]);

  const navigate = useNavigate();
  const createTemplate = useCreateCertificationTemplateMutation();
  const updateTemplate = useUpdateCertificationTemplateMutation();

  async function saveTemplate() {
    const result = certificationTemplateCreateSchema.safeParse(
      certificateEditorActions.getDraft(),
    );
    if (!result.success) {
      toast.error(result.error.issues[0]?.message ?? "Template inválido");
      return;
    }

    const data = await uploadCertificateAssets(result.data, eventId, editionId);
    const onSuccess = () =>
      void navigate({
        to: "/admin/events/$eventId/editions/$editionId/certifications",
        params: { eventId, editionId },
      });
    if (templateId && !duplicate)
      updateTemplate.mutate({ templateId, data }, { onSuccess });
    else createTemplate.mutate({ editionId, data }, { onSuccess });
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
          </div>

          <Button type="button" variant="default" onClick={saveTemplate}>
            <Save className="size-4" />
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
