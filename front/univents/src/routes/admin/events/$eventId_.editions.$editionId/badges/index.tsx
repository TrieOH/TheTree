import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { BadgeCheck, MoreVertical, Plus, Trash2 } from "lucide-react";
import { useMemo, useState } from "react";
import { badgeTemplatesQueryOptions } from "@/features/badges/api";
import { useDeleteBadgeTemplateMutation } from "@/features/badges/api/mutations";
import type { BadgeTemplate } from "@/features/badges/model";
import { allTicketsQueryOptions } from "@/features/tickets/api";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/shared/ui/shadcn/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/shared/ui/shadcn/dropdown-menu";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/badges/",
)({ component: RouteComponent });

function BadgeMiniature({ template }: { template: BadgeTemplate }) {
  const design = template.design_data;
  return (
    <div className="flex h-52 items-center justify-center overflow-hidden bg-muted/60 p-4">
      <div
        className="relative h-44 overflow-hidden shadow-md"
        style={{
          aspectRatio: `${design.canvas.width}/${design.canvas.height}`,
          backgroundColor: design.backgroundColor,
          backgroundImage: design.background
            ? `url(${design.background})`
            : undefined,
          backgroundSize: "cover",
        }}
      >
        {design.elements
          .filter((element) => element.type === "text")
          .map((element) => (
            <div
              key={element.id}
              className="absolute overflow-hidden"
              style={{
                left: `${(element.x / design.canvas.width) * 100}%`,
                top: `${(element.y / design.canvas.height) * 100}%`,
                width: `${(element.width / design.canvas.width) * 100}%`,
                height: `${(element.height / design.canvas.height) * 100}%`,
                color: element.paragraphs[0]?.runs[0]?.color,
                fontWeight: element.paragraphs[0]?.runs[0]?.bold
                  ? "bold"
                  : "normal",
                textAlign: element.paragraphs[0]?.align,
                fontSize: Math.max(
                  5,
                  (element.paragraphs[0]?.runs[0]?.fontSize ?? 24) * 0.17,
                ),
              }}
            >
              {element.paragraphs
                .flatMap((paragraph) => paragraph.runs.map((run) => run.text))
                .join("\n")}
            </div>
          ))}
      </div>
    </div>
  );
}

function RouteComponent() {
  const { eventId, editionId } = Route.useParams();
  const [filter, setFilter] = useState("");
  const { data: templates = [] } = useQuery(
    badgeTemplatesQueryOptions(editionId),
  );
  const { data: tickets = [] } = useQuery(allTicketsQueryOptions(editionId));
  const remove = useDeleteBadgeTemplateMutation();
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
  return (
    <div className="p-6 pb-28">
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
          slice.map((template) => (
            <Card key={template.id} className="overflow-hidden p-0">
              <BadgeMiniature template={template} />
              <CardHeader className="px-4">
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle className="text-base">{template.name}</CardTitle>
                    <p className="mt-1 text-xs text-muted-foreground">
                      {template.ticket_type_id
                        ? (ticketNames.get(template.ticket_type_id) ??
                          "Ingresso associado")
                        : "Padrão da edição"}
                    </p>
                  </div>
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      render={
                        <Button
                          variant="ghost"
                          size="icon"
                          aria-label="Abrir ações do template"
                        />
                      }
                    >
                      <MoreVertical className="size-4" />
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem
                        variant="destructive"
                        onClick={() =>
                          remove.mutate({ templateId: template.id })
                        }
                      >
                        <Trash2 className="size-4" />
                        Excluir
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </CardHeader>
              <CardContent className="px-4 text-xs text-muted-foreground">
                {template.design_data.canvas.width} ×{" "}
                {template.design_data.canvas.height}px ·{" "}
                {template.design_data.elements.length} elementos
              </CardContent>
              <CardFooter className="border-t px-4 py-3">
                <span className="text-xs text-muted-foreground">
                  Criado em{" "}
                  {new Date(template.created_at).toLocaleDateString("pt-BR")}
                </span>
              </CardFooter>
            </Card>
          ))
        }
      />
    </div>
  );
}
