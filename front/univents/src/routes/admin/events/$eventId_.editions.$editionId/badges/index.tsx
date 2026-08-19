import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { BadgeCheck, FileText, Plus, Printer } from "lucide-react";
import { useMemo, useState } from "react";
import {
  badgePrintQueryOptions,
  badgeTemplatesQueryOptions,
} from "@/features/badges/api";
import { useDeleteBadgeTemplateMutation } from "@/features/badges/api/mutations";
import type { BadgePrintItem, BadgeTemplate } from "@/features/badges/model";
import { badgeDesignSchema } from "@/features/badges/model";
import AdminBadgeCard from "@/features/badges/ui/AdminBadgeCard";
import { BadgePreview } from "@/features/badges/ui/badge-preview";
import { allTicketsQueryOptions } from "@/features/tickets/api";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
// Using AdminBadgeCard component for card rendering.

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/badges/",
)({ component: RouteComponent });

/* BadgeMiniature extracted to `AdminBadgeCard` */

function PrintableBadge({ badge }: { badge: BadgePrintItem }) {
  const design = badgeDesignSchema.safeParse(badge.design_data).success
    ? badgeDesignSchema.parse(badge.design_data)
    : null;
  const widthMm = design ? design.canvas.width / (96 / 25.4) : 85;
  const heightMm = design ? design.canvas.height / (96 / 25.4) : 54;

  return (
    <article
      className="overflow-hidden bg-transparent text-black shadow print:break-inside-avoid print:shadow-none"
      style={{ width: `${widthMm}mm`, height: `${heightMm}mm` }}
    >
      <BadgePreview
        badge={badge}
        className="relative h-full w-full rounded-none border-0 shadow-none"
      />
    </article>
  );
}

function RouteComponent() {
  const { eventId, editionId } = Route.useParams();
  const navigate = useNavigate();
  const remove = useDeleteBadgeTemplateMutation();
  const [filter, setFilter] = useState("");
  const [emissionFilter, setEmissionFilter] = useState("");
  const [activeSection, setActiveSection] = useState<"templates" | "emissions">(
    "templates",
  );
  const { data: templates = [] } = useQuery(
    badgeTemplatesQueryOptions(editionId),
  );
  const printQuery = useQuery(badgePrintQueryOptions(editionId));
  const { data: tickets = [] } = useQuery(allTicketsQueryOptions(editionId));
  const filtered = useMemo(
    () =>
      templates.filter((item) =>
        item.name.toLowerCase().includes(filter.trim().toLowerCase()),
      ),
    [filter, templates],
  );
  const ticketNames = new Map(
    tickets.map((ticket) => [ticket.id, ticket.name]),
  );
  const printItems = printQuery.data ?? [];
  const filteredPrintItems = useMemo(() => {
    const search = emissionFilter.trim().toLowerCase();
    if (!search) return printItems;
    return printItems.filter((item) =>
      [item.event_name, item.edition_name, item.ticket_name, item.template_name]
        .filter(Boolean)
        .some((value) => value?.toLowerCase().includes(search)),
    );
  }, [emissionFilter, printItems]);
  return (
    <>
      <div className="hidden print:block">
        <div className="flex flex-wrap content-start gap-[4mm] p-[10mm]">
          {printItems?.map((badge) => (
            <PrintableBadge key={badge.emission_id} badge={badge} />
          ))}
        </div>
      </div>
      <div className="p-6 pb-28 print:hidden">
        <nav className="mb-6 flex gap-1 border-b border-border">
          {(["templates", "emissions"] as const).map((section) => (
            <button
              key={section}
              type="button"
              onClick={() => setActiveSection(section)}
              className={cn(
                "inline-flex items-center gap-2 whitespace-nowrap border-b-2 px-3 py-2.5 text-sm font-medium",
                activeSection === section
                  ? "border-primary text-primary"
                  : "border-transparent text-muted-foreground",
              )}
            >
              {section === "templates" ? (
                <>
                  <FileText className="size-4" /> Templates
                </>
              ) : (
                <>
                  <BadgeCheck className="size-4" /> Crachás emitidos
                </>
              )}
            </button>
          ))}
        </nav>
        {activeSection === "templates" ? (
          <PaginatedContainer<BadgeTemplate>
            items={filtered}
            layout="grid"
            minItemWidth="16rem"
            pageSize={8}
            gap="6"
            filterValue={filter}
            onFilterChange={setFilter}
            filterPlaceholder="Buscar templates..."
            itemLabel="templates"
            headerActions={
              <Link
                to="/admin/events/$eventId/editions/$editionId/badges/editor"
                params={{ eventId, editionId }}
                search={{ templateId: "", duplicate: false }}
                className={cn(
                  "inline-flex h-9 items-center justify-center gap-2 rounded-lg bg-primary px-4 text-sm font-medium text-primary-foreground",
                )}
              >
                <Plus className="size-4" />
                Novo template
              </Link>
            }
            emptyState={
              <EmptyState
                icon={BadgeCheck}
                eyebrow="Crachás"
                title="Nenhum template encontrado"
                description="Crie um template totalmente personalizado ou parta do padrão do sistema."
                className="border-0 bg-transparent"
              />
            }
            renderItems={(slice) =>
              slice.map((template, i) => (
                <AdminBadgeCard
                  key={template.id}
                  item={template}
                  kind="template"
                  index={i}
                  ticketName={
                    template.ticket_type_id
                      ? (ticketNames.get(template.ticket_type_id) ??
                        "Ingresso associado")
                      : "Padrão da edição"
                  }
                  onView={() =>
                    void navigate({
                      to: "/admin/events/$eventId/editions/$editionId/badges/editor",
                      params: { eventId, editionId },
                      search: { templateId: template.id, duplicate: false },
                    })
                  }
                  onEdit={() =>
                    void navigate({
                      to: "/admin/events/$eventId/editions/$editionId/badges/editor",
                      params: { eventId, editionId },
                      search: { templateId: template.id, duplicate: false },
                    })
                  }
                  onDuplicate={() =>
                    void navigate({
                      to: "/admin/events/$eventId/editions/$editionId/badges/editor",
                      params: { eventId, editionId },
                      search: { templateId: template.id, duplicate: true },
                    })
                  }
                  onDelete={() => remove.mutate({ templateId: template.id })}
                />
              ))
            }
          />
        ) : (
          <PaginatedContainer<BadgePrintItem>
            items={filteredPrintItems}
            layout="grid"
            minItemWidth="16rem"
            pageSize={8}
            gap="6"
            filterValue={emissionFilter}
            onFilterChange={setEmissionFilter}
            filterPlaceholder="Buscar crachá..."
            itemLabel="crachás"
            headerActions={
              <Button
                variant="default"
                disabled={printQuery.isFetching}
                onClick={async () => {
                  const result = await printQuery.refetch();
                  if (result.data) {
                    requestAnimationFrame(() => window.print());
                  }
                }}
                className={cn(
                  "inline-flex h-9 items-center justify-center gap-2 rounded-lg",
                )}
              >
                <Printer className="size-4" />
                {printQuery.isFetching ? "Preparando…" : "Imprimir crachás"}
              </Button>
            }
            emptyState={
              <EmptyState
                title="Nenhum crachá emitido"
                description="As emissões aparecerão aqui."
              />
            }
            renderItems={(slice) =>
              slice.map((badge, i) => (
                <AdminBadgeCard
                  key={badge.emission_id}
                  item={badge}
                  kind="emission"
                  index={i}
                  onView={() => undefined}
                />
              ))
            }
          />
        )}
        {false && (
          <PaginatedContainer<BadgeTemplate>
            items={filtered}
            layout="grid"
            minItemWidth="16rem"
            pageSize={8}
            gap="6"
            filterValue={filter}
            onFilterChange={setFilter}
            filterPlaceholder="Buscar templates..."
            itemLabel="templates"
            headerActions={
              <Link
                to="/admin/events/$eventId/editions/$editionId/badges/editor"
                params={{ eventId, editionId }}
                search={{ templateId: "", duplicate: false }}
                className={cn(
                  "inline-flex h-9 items-center justify-center gap-2 rounded-lg bg-primary px-4 text-sm font-medium text-primary-foreground",
                )}
              >
                <Plus className="size-4" />
                Novo template
              </Link>
            }
            emptyState={
              <EmptyState
                icon={BadgeCheck}
                eyebrow="Crachás"
                title="Nenhum template encontrado"
                description="Crie um template totalmente personalizado ou parta do padrão do sistema."
                className="border-0 bg-transparent"
              />
            }
            renderItems={(slice) =>
              slice.map((template, i) => (
                <AdminBadgeCard
                  key={template.id}
                  item={template}
                  kind="template"
                  index={i}
                  ticketName={
                    template.ticket_type_id
                      ? (ticketNames.get(template.ticket_type_id) ??
                        "Ingresso associado")
                      : "Padrão da edição"
                  }
                  onView={() =>
                    void navigate({
                      to: "/admin/events/$eventId/editions/$editionId/badges/editor",
                      params: { eventId, editionId },
                      search: { templateId: template.id, duplicate: false },
                    })
                  }
                  onEdit={() =>
                    void navigate({
                      to: "/admin/events/$eventId/editions/$editionId/badges/editor",
                      params: { eventId, editionId },
                      search: { templateId: template.id, duplicate: false },
                    })
                  }
                  onDelete={() => remove.mutate({ templateId: template.id })}
                />
              ))
            }
          />
        )}
      </div>
    </>
  );
}
