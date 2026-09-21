import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { For, Show, createMemo, createSignal } from "solid-js";

import AwardIcon from "~icons/lucide/award";
import PlusIcon from "~icons/lucide/plus";

import { useQuery } from "@trieoh/front-core-solid";
import type { SortState } from "@trieoh/ui-solid";
import { Button, EmptyState, PaginatedContainer } from "@trieoh/ui-solid";
import { allCertificationTemplatesQueryOptions } from "@/features/certifications/api";
import { useDeleteCertificationTemplateMutation } from "@/features/certifications/api/mutations";
import type { CertificationTemplateI } from "@/features/certifications/model";
import {
  AdminCertificationTemplateCard,
  AdminCreateCertificationCard,
  CertViewer,
  CertificationEmissionErrorsList,
  CertificationList,
} from "@/features/certifications/ui";
import {
  type CertificationSection,
  CertificationSectionTabs,
} from "@/features/certifications/ui/CertificationSectionTabs";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import { toast } from "@/shared/ui/toast";

const Award = AwardIcon as unknown as (props: { class?: string }) => JSX.Element;
const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/certifications/",
)({
  head: () => ({
    meta: [{ title: "Certificados - Admin Univents" }],
  }),
  component: AdminEditionCertificationsRoute,
});

function AdminEditionCertificationsRoute(): JSX.Element {
  const params = Route.useParams();
  const navigate = useNavigate();

  const eventId = () => params().eventId;
  const editionId = () => params().editionId;

  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<CertificationTemplateI>>({
    field: "name",
    direction: "asc",
  });
  const [viewingTemplate, setViewingTemplate] =
    createSignal<CertificationTemplateI | null>(null);
  const [activeSection, setActiveSection] =
    createSignal<CertificationSection>("templates");

  const editionsQuery = useQuery(() => allAdminEditionsQueryOptions(eventId()));
  const templatesQuery = useQuery(() =>
    allCertificationTemplatesQueryOptions(editionId()),
  );

  const edition = createMemo(() => {
    const list = editionsQuery().data;
    if (!Array.isArray(list)) return undefined;
    return list.find((item) => item.id === editionId());
  });

  const templates = () =>
    (templatesQuery().data ?? []) as CertificationTemplateI[];

  const filteredTemplates = createMemo(() => {
    const staticSearch = filter().trim().toLowerCase();
    if (!staticSearch) return templates();

    return templates().filter((template) =>
      [template.name, template.description ?? ""].some((value) =>
        value.toLowerCase().includes(staticSearch),
      ),
    );
  });

  const deleteTemplateMutation = useDeleteCertificationTemplateMutation();

  const handleOpenEditor = (templateId?: string, duplicate?: boolean) => {
    navigate({
      to: "/admin/events/$eventId/editions/$editionId/certifications/editor",
      params: { eventId: eventId(), editionId: editionId() },
      search: {
        templateId: templateId ?? undefined,
        duplicate: duplicate ? true : undefined,
      },
    });
  };

  const handleDeleteTemplate = async (template: CertificationTemplateI) => {
    try {
      await deleteTemplateMutation.mutateAsync({
        editionId: editionId(),
        templateId: template.id,
      });
      toast.success("Template excluído com sucesso.");
    } catch {
      toast.error("Erro ao excluir template.");
    }
  };

  return (
    <div class="mx-auto max-w-7xl space-y-6">
      {/* Header with Title & Description */}
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-foreground">
            Certificados da Edição
          </h1>
          <p class="text-xs sm:text-sm text-muted-foreground">
            Crie templates visuais e gerencie a emissão de certificados para participantes desta edição.
          </p>
        </div>
      </div>

      {/* Navigation Tabs */}
      <CertificationSectionTabs
        active={activeSection()}
        onChange={setActiveSection}
      />

      {/* Content: Templates */}
      <Show when={activeSection() === "templates"}>
        <PaginatedContainer<CertificationTemplateI>
          items={filteredTemplates()}
          layout="grid"
          minItemWidth="16rem"
          maxRows={(columns) => (columns === 1 ? 8 : 4)}
          gap="2"
          sort={sort()}
          onSortChange={setSort}
          sortFields={[
            {
              key: "name",
              label: "Nome",
              ascLabel: "A → Z",
              descLabel: "Z → A",
            },
            {
              key: "created_at",
              label: "Data de criação",
              ascLabel: "Mais antigos primeiro",
              descLabel: "Mais recentes primeiro",
              comparator: (a, b) =>
                new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
            },
          ]}
          filterValue={filter()}
          onFilterChange={setFilter}
          filterPlaceholder="Buscar por nome do template..."
          itemLabel="templates"
          emptyState={
            <EmptyState
              class="border-0 bg-transparent px-0 py-4 shadow-none"
              icon={<Award class="size-6 text-foreground/70" />}
              eyebrow="Certificações"
              title="Nenhum template cadastrado"
              description={
                filter()
                  ? "Nenhum template corresponde à busca informada."
                  : "Crie o primeiro template para começar a emitir certificados nesta edição."
              }
              action={
                <Button
                  variant="outline"
                  class="gap-2 rounded-sm py-4 cursor-pointer"
                  onClick={() => handleOpenEditor()}
                >
                  <Plus class="size-4" />
                  Novo template
                </Button>
              }
            />
          }
          renderItems={(slice, options) => (
            <>
              <AdminCreateCertificationCard
                index={0}
                animate={options.animate}
                onCreate={() => handleOpenEditor()}
              />
              <For each={slice}>
                {(template, index) => (
                  <AdminCertificationTemplateCard
                    template={template}
                    index={index() + 1}
                    animate={options.animate}
                    onEdit={() => handleOpenEditor(template.id)}
                    onView={() => setViewingTemplate(template)}
                    onDuplicate={() => handleOpenEditor(template.id, true)}
                    onDelete={() => void handleDeleteTemplate(template)}
                  />
                )}
              </For>
            </>
          )}
        />
      </Show>

      {/* Content: Certificates */}
      <Show when={activeSection() === "certificates"}>
        <section class="w-full">
          <CertificationList
            eventId={eventId()}
            editionId={editionId()}
            editionName={edition()?.name}
          />
        </section>
      </Show>

      {/* Content: Emission Errors */}
      <Show when={activeSection() === "errors"}>
        <section class="w-full">
          <CertificationEmissionErrorsList editionId={editionId()} />
        </section>
      </Show>

      {/* Certificate Viewer Modal */}
      <Show when={viewingTemplate()}>
        {(tmpl) => (
          <CertViewer
            template={tmpl()}
            open={viewingTemplate() !== null}
            onOpenChange={(open) => {
              if (!open) setViewingTemplate(null);
            }}
            variables={{
              participant_name: "Nome civil do participante",
              event_name: "Nome do evento",
              edition_name: edition()?.name ?? "Nome da edição",
              activity_name: edition()?.name ?? "Nome da edição",
              participation_type: "edição",
              location: edition()?.location_name ?? "Local da edição",
              workload_hours: "Carga horária conforme presença",
              participation_date: "Data da participação",
              certified_at: new Date().toLocaleDateString("pt-BR"),
              cert_hash: "UNIV-2026-CERT-SAMPLE",
              verify_url:
                typeof window !== "undefined" ? window.location.origin : "",
            }}
          />
        )}
      </Show>
    </div>
  );
}
