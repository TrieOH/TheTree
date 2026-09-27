import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { BadgeEditor } from "@/features/badges/editor/BadgeEditor";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/badges/editor",
)(({
  validateSearch: (search: Record<string, unknown>): { templateId?: string; duplicate?: boolean } => ({
    templateId: typeof search.templateId === "string" && search.templateId ? search.templateId : undefined,
    duplicate: search.duplicate === true || search.duplicate === "true" ? true : undefined,
  }),
  head: () => ({
    meta: [{ title: "Editor de Crachás - Admin Univents" }],
  }),
  component: AdminBadgeEditorRoute,
}));

function AdminBadgeEditorRoute(): JSX.Element {
  const params = Route.useParams();
  const search = Route.useSearch();

  return (
    <BadgeEditor
      eventId={params().eventId}
      editionId={params().editionId}
      templateId={search().templateId}
      duplicate={Boolean(search().duplicate)}
    />
  );
}
