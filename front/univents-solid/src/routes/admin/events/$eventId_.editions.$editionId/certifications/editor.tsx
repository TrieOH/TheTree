import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { CertificateEditor } from "@/features/certifications/editor/CertificateEditor";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/certifications/editor",
)({
  validateSearch: (
    search: Record<string, unknown>,
  ): { templateId?: string; duplicate?: boolean } => ({
    templateId:
      typeof search.templateId === "string" && search.templateId
        ? search.templateId
        : undefined,
    duplicate:
      search.duplicate === true || search.duplicate === "true"
        ? true
        : undefined,
  }),
  head: () => ({
    meta: [{ title: "Editor de Certificados - Admin Univents" }],
  }),
  component: AdminCertificateEditorRoute,
});

function AdminCertificateEditorRoute(): JSX.Element {
  const params = Route.useParams();
  const search = Route.useSearch();

  return (
    <CertificateEditor
      eventId={params().eventId}
      editionId={params().editionId}
      templateId={search().templateId}
      duplicate={Boolean(search().duplicate)}
    />
  );
}
