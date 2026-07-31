import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { FileText, Plus } from "lucide-react";
import { useMemo, useState } from "react";
import { allCertificationTemplatesQueryOptions } from "@/features/certifications/api";
import { useDeleteCertificationTemplateMutation } from "@/features/certifications/api/mutations";
import type { CertificationTemplateI } from "@/features/certifications/model";
import { AdminCertificationTemplateCard } from "@/features/certifications/ui/AdminCertificationTemplateCard";
import {
  CertificationEmissionErrorsList,
  CertificationList,
} from "@/features/certifications/ui/CertificationLists";
import {
  type CertificationSection,
  CertificationSectionTabs,
} from "@/features/certifications/ui/CertificationSectionTabs";
import { CertViewer } from "@/features/certifications/ui/CertViewer";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import { cn } from "@/shared/lib/utils";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/shared/ui/shadcn/alert-dialog";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/certifications/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId, editionId } = Route.useParams();
  const navigate = useNavigate();
  const [filter, setFilter] = useState("");
  const [deletingTemplate, setDeletingTemplate] =
    useState<CertificationTemplateI | null>(null);
  const [viewingTemplate, setViewingTemplate] =
    useState<CertificationTemplateI | null>(null);
  const [activeSection, setActiveSection] =
    useState<CertificationSection>("templates");

  const { data: editions = [] } = useQuery(
    allAdminEditionsQueryOptions(eventId),
  );
  const { data: templates = [] } = useQuery(
    allCertificationTemplatesQueryOptions(editionId),
  );

  const edition = editions.find((item) => item.id === editionId) ?? null;

  const filteredTemplates = useMemo(() => {
    const search = filter.trim().toLowerCase();
    if (!search) return templates;

    return templates.filter((template) =>
      [template.name, template.description ?? ""].some((value) =>
        value.toLowerCase().includes(search),
      ),
    );
  }, [filter, templates]);
  const deleteTemplateMutation = useDeleteCertificationTemplateMutation();

  return (
    <div className="flex flex-wrap gap-6 p-6 pb-28!">
      <div className="w-full">
        <CertificationSectionTabs
          active={activeSection}
          onChange={setActiveSection}
        />
      </div>
      {activeSection === "templates" ? (
        <PaginatedContainer<CertificationTemplateI>
          items={filteredTemplates}
          layout="grid"
          minItemWidth="16rem"
          pageSize={6}
          gap="6"
          filterValue={filter}
          onFilterChange={setFilter}
          filterPlaceholder="Buscar por título ou URL..."
          itemLabel="templates"
          headerActions={
            <Link
              to="/admin/events/$eventId/editions/$editionId/certifications/editor"
              params={{ eventId, editionId }}
              search={{ templateId: "" }}
              className={cn(
                "inline-flex h-9 items-center justify-center gap-2 rounded-lg px-4 text-sm font-medium",
                "bg-primary text-primary-foreground shadow-sm transition-colors hover:bg-primary/90",
                "sm:min-w-40 sm:px-5",
              )}
            >
              <Plus className="size-4 shrink-0" />
              <span className="whitespace-nowrap">Novo template</span>
            </Link>
          }
          emptyState={
            <EmptyState
              icon={FileText}
              eyebrow="Certificações"
              title="Nenhum template encontrado"
              description="Crie o primeiro template para começar a emitir certificados nessa edição."
              className="border-0 bg-transparent px-0 py-4 shadow-none"
            />
          }
          renderItems={(slice) =>
            slice.map((template, index) => {
              return (
                <AdminCertificationTemplateCard
                  key={template.id}
                  template={template}
                  index={index}
                  onEdit={() => {
                    void navigate({
                      to: "/admin/events/$eventId/editions/$editionId/certifications/editor",
                      params: { eventId, editionId },
                      search: { templateId: template.id },
                    });
                  }}
                  onView={() => setViewingTemplate(template)}
                  onDelete={() => setDeletingTemplate(template)}
                />
              );
            })
          }
        />
      ) : activeSection === "certificates" ? (
        <section className="w-full">
          <CertificationList eventId={eventId} editionId={editionId} />
        </section>
      ) : (
        <section className="w-full">
          <CertificationEmissionErrorsList
            eventId={eventId}
            editionId={editionId}
          />
        </section>
      )}
      {viewingTemplate ? (
        <CertViewer
          template={viewingTemplate}
          open
          onOpenChange={(open) => {
            if (!open) setViewingTemplate(null);
          }}
          variables={{
            activity_name: edition?.name ?? "Nome da edição",
            certified_at: "DD/MM/AAAA",
            cert_hash: "HASH-DE-EXEMPLO",
            verify_url: window.location.href,
          }}
        />
      ) : null}
      <AlertDialog
        open={deletingTemplate !== null}
        onOpenChange={(open) => {
          if (!open) setDeletingTemplate(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Excluir template?</AlertDialogTitle>
            <AlertDialogDescription>
              O template “{deletingTemplate?.name}” será removido
              permanentemente.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancelar</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              onClick={() => {
                if (!deletingTemplate) return;
                deleteTemplateMutation.mutate(
                  { templateId: deletingTemplate.id },
                  { onSuccess: () => setDeletingTemplate(null) },
                );
              }}
            >
              Excluir template
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
