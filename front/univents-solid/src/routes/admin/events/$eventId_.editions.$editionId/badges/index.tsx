import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import {
  Button,
  EmptyState,
  PaginatedContainer,
  type SortState,
  cn,
} from "@trieoh/ui-solid";
import { For, Show, createEffect, createMemo, createSignal } from "solid-js";

import BadgeCheckIcon from "~icons/lucide/badge-check";
import CalendarClockIcon from "~icons/lucide/calendar-clock";
import PlusIcon from "~icons/lucide/plus";
import PrinterIcon from "~icons/lucide/printer";
import QrCodeIcon from "~icons/lucide/qr-code";

import {
  badgeEmissionsQueryOptions,
  badgePrintQueryOptions,
  badgeTemplatesQueryOptions,
} from "@/features/badges/api";
import { useDeleteBadgeTemplateMutation } from "@/features/badges/api/mutations";
import type {
  BadgeEditionEmission,
  BadgePrintItem,
  BadgeTemplate,
} from "@/features/badges/model";
import { selectBadgePrintItems } from "@/features/badges/model/print-selection";
import {
  AdminBadgeCard,
  AdminCreateBadgeCard,
  BadgeSectionTabs,
  type BadgeSection,
  DateFilterDialog,
  PrintableBadge,
  PrintableQr,
  QrPrintDialog,
} from "@/features/badges/ui";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import type { EditionI } from "@/features/editions/model";
import { asUniventsProfile, profileDisplayName } from "@/features/profile/model/profile-data";
import { allTicketsQueryOptions } from "@/features/tickets/api";
import type { TicketType } from "@trieoh/univents-api/schemas";
import { printElement } from "@/shared/lib/print-element";
import { toast } from "@/shared/ui/toast";

const BadgeCheck = BadgeCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const CalendarClock = CalendarClockIcon as unknown as (props: { class?: string }) => JSX.Element;
const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;
const Printer = PrinterIcon as unknown as (props: { class?: string }) => JSX.Element;
const QrCode = QrCodeIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/badges/",
)({
  head: () => ({
    meta: [{ title: "Crachás - Admin Univents" }],
  }),
  component: AdminEditionBadgesRoute,
});

function AdminTemplateBadgeItem(props: {
  template: BadgeTemplate;
  index: number;
  animate?: boolean;
  ticketNames: () => Map<string, string>;
  location: () => string;
  onEdit: () => void;
  onDuplicate: () => void;
  onDelete: () => void;
}): JSX.Element {
  const ticketName = createMemo(() => {
    if (!props.template.ticket_type_id) return "Padrão da edição";
    return props.ticketNames().get(props.template.ticket_type_id) ?? "Ingresso associado";
  });

  return (
    <AdminBadgeCard
      item={props.template}
      kind="template"
      index={props.index}
      animate={props.animate}
      ticketName={ticketName()}
      location={props.location()}
      onEdit={props.onEdit}
      onDuplicate={props.onDuplicate}
      onDelete={props.onDelete}
    />
  );
}

function AdminEmissionBadgeItem(props: {
  badge: BadgePrintItem;
  index: number;
  animate?: boolean;
  participantNames: () => Record<string, string>;
  location: () => string;
}): JSX.Element {
  const participantName = createMemo(() => props.participantNames()[props.badge.user_id]);

  return (
    <AdminBadgeCard
      item={props.badge}
      kind="emission"
      index={props.index}
      animate={props.animate}
      participantName={participantName()}
      location={props.location()}
    />
  );
}

function PrintableQrItem(props: {
  badge: BadgePrintItem;
  printQrSizeMm: () => number;
  participantNames: () => Record<string, string>;
}): JSX.Element {
  const size = createMemo(() => (props.printQrSizeMm() / 25.4) * 96);
  const participant = createMemo(() => props.participantNames()[props.badge.user_id] ?? props.badge.user_id.slice(0, 8));

  return (
    <PrintableQr
      badge={props.badge}
      size={size()}
      participant={participant()}
    />
  );
}

function PrintableBadgeItem(props: {
  badge: BadgePrintItem;
  participantNames: () => Record<string, string>;
  location: () => string;
}): JSX.Element {
  const participantName = createMemo(() => props.participantNames()[props.badge.user_id] ?? "");

  return (
    <PrintableBadge
      badge={props.badge}
      participantName={participantName()}
      location={props.location()}
    />
  );
}

function AdminEditionBadgesRoute(): JSX.Element {
  const navigate = useNavigate();
  const params = Route.useParams();
  const eventId = () => params().eventId;
  const editionId = () => params().editionId;
  const { auth } = useAuth();

  // Queries
  const templatesQuery = useQuery(() => badgeTemplatesQueryOptions(editionId()));
  const emissionsQuery = useQuery(() => badgeEmissionsQueryOptions(editionId()));
  const printQuery = useQuery(() => badgePrintQueryOptions(editionId()));
  const editionsQuery = useQuery(() => allAdminEditionsQueryOptions(eventId()));
  const ticketsQuery = useQuery(() => allTicketsQueryOptions(editionId()));

  // Mutations
  const deleteMutation = useDeleteBadgeTemplateMutation();

  // State
  const [activeSection, setActiveSection] = createSignal<BadgeSection>("templates");
  const [templateFilter, setTemplateFilter] = createSignal("");
  const [emissionFilter, setEmissionFilter] = createSignal("");
  const [templateSort, setTemplateSort] = createSignal<SortState<BadgeTemplate>>({
    field: "name",
    direction: "asc",
  });
  const [emissionSort, setEmissionSort] = createSignal<SortState<BadgePrintItem>>({
    field: "template_name",
    direction: "asc",
  });
  const [printedAfter, setPrintedAfter] = createSignal("");
  const [printMode, setPrintMode] = createSignal<"badges" | "qrs">("badges");
  const [printQrSizeMm, setPrintQrSizeMm] = createSignal(48);
  const [isPrinting, setIsPrinting] = createSignal(false);

  // Dialogs
  const [dateDialogOpen, setDateDialogOpen] = createSignal(false);
  const [qrDialogOpen, setQrDialogOpen] = createSignal(false);

  let printRootRef: HTMLDivElement | undefined;

  // Data helpers
  const templates = createMemo(() => (templatesQuery().data ?? []) as BadgeTemplate[]);
  const emissions = createMemo(() => (emissionsQuery().data ?? []) as BadgeEditionEmission[]);
  const printItems = createMemo(() => (printQuery().data ?? []) as BadgePrintItem[]);
  const editions = createMemo(() => (editionsQuery().data ?? []) as EditionI[]);
  const tickets = createMemo(() => (ticketsQuery().data ?? []) as TicketType[]);

  const ticketNames = createMemo(() => {
    const map = new Map<string, string>();
    for (const t of tickets()) {
      map.set(t.id, t.name);
    }
    return map;
  });

  const location = createMemo(() => {
    const ed = editions().find((e) => e.id === editionId());
    return ed?.location_name ?? "";
  });

  const filteredTemplates = createMemo(() => {
    const list = templates();
    const query = templateFilter().trim().toLowerCase();
    if (!query) return list;
    const result: BadgeTemplate[] = [];
    for (const item of list) {
      if (item.name.toLowerCase().includes(query)) {
        result.push(item);
      }
    }
    return result;
  });

  const printableItems = createMemo(() =>
    selectBadgePrintItems(printItems(), emissions(), printedAfter()),
  );

  const printActorIds = createMemo(() => [
    ...new Set(printableItems().map((item) => item.user_id)),
  ]);

  const [participantNames, setParticipantNames] = createSignal<Record<string, string>>({});

  createEffect(
    () => printActorIds(),
    (ids) => {
      if (!ids || ids.length === 0) return;
      void (async () => {
        const result: Record<string, string> = {};
        await Promise.all(
          ids.map(async (id) => {
            try {
              const res = await auth.getActorProfile(id);
              if (res.success && res.data) {
                const p = asUniventsProfile(res.data.profile ?? {});
                result[id] = profileDisplayName(p);
              }
            } catch {
              // ignore
            }
          }),
        );
        setParticipantNames((prev) => ({ ...prev, ...result }));
      })();
    },
  );

  const filteredPrintItems = createMemo(() => {
    const list = printableItems();
    const query = emissionFilter().trim().toLowerCase();
    if (!query) return list;
    const result: BadgePrintItem[] = [];
    for (const item of list) {
      const match =
        (item.event_name && item.event_name.toLowerCase().includes(query)) ||
        (item.edition_name && item.edition_name.toLowerCase().includes(query)) ||
        (item.ticket_name && item.ticket_name.toLowerCase().includes(query)) ||
        (item.template_name && item.template_name.toLowerCase().includes(query));
      if (match) {
        result.push(item);
      }
    }
    return result;
  });

  // Handlers
  const handleOpenEditor = (templateId?: string, duplicate?: boolean) => {
    navigate({
      to: "/admin/events/$eventId/editions/$editionId/badges/editor",
      params: { eventId: eventId(), editionId: editionId() },
      search: {
        templateId: templateId ?? undefined,
        duplicate: duplicate ? true : undefined,
      },
    });
  };

  const handleDeleteTemplate = async (template: BadgeTemplate) => {
    try {
      await deleteMutation.mutateAsync({
        editionId: editionId(),
        templateId: template.id,
      });
      toast.success("Template excluído com sucesso.");
    } catch {
      toast.error("Erro ao excluir template. Tente novamente.");
    }
  };

  const handlePrintBadges = async () => {
    if (printableItems().length === 0 || !printRootRef) return;
    setPrintMode("badges");
    setIsPrinting(true);
    try {
      // Allow DOM to update printable elements
      await new Promise((resolve) => setTimeout(resolve, 150));
      await printElement(printRootRef, "Crachás");
    } catch {
      toast.error("Erro ao preparar impressão.");
    } finally {
      setIsPrinting(false);
    }
  };

  const handlePrintQrs = async (sizeMm: number) => {
    if (printableItems().length === 0 || !printRootRef) return;
    setPrintQrSizeMm(sizeMm);
    setPrintMode("qrs");
    setIsPrinting(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 150));
      await printElement(printRootRef, "QR Codes dos Crachás");
    } catch {
      toast.error("Erro ao preparar impressão dos QR codes.");
    } finally {
      setIsPrinting(false);
    }
  };

  return (
    <>
      {/* Hidden container for printing */}
      <div class="hidden print:block">
        <div ref={(el) => (printRootRef = el)}>
          <Show
            when={printMode() === "badges"}
            fallback={
              <For each={printableItems()}>
                {(badge) => (
                  <PrintableQrItem
                    badge={badge}
                    printQrSizeMm={printQrSizeMm}
                    participantNames={participantNames}
                  />
                )}
              </For>
            }
          >
            <For each={printableItems()}>
              {(badge) => (
                <PrintableBadgeItem
                  badge={badge}
                  participantNames={participantNames}
                  location={location}
                />
              )}
            </For>
          </Show>
        </div>
      </div>

      {/* Main Admin UI */}
      <div class="mx-auto max-w-7xl space-y-6 print:hidden">
        {/* Header with Title & Description */}
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-foreground">
              Crachás da Edição
            </h1>
            <p class="text-xs sm:text-sm text-muted-foreground">
              Crie templates visuais e gerencie as credenciais e impressões dos participantes.
            </p>
          </div>
        </div>

        {/* Navigation Tabs */}
        <BadgeSectionTabs
          active={activeSection()}
          onChange={setActiveSection}
          emissionsCount={emissions().length > 0 ? emissions().length : undefined}
        />

        {/* Content: Templates */}
        <Show when={activeSection() === "templates"}>
          <PaginatedContainer<BadgeTemplate>
            items={filteredTemplates()}
            layout="grid"
            minItemWidth="16rem"
            maxRows={(columns) => (columns === 1 ? 8 : 4)}
            gap="2"
            sort={templateSort()}
            onSortChange={setTemplateSort}
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
            filterValue={templateFilter()}
            onFilterChange={setTemplateFilter}
            filterPlaceholder="Buscar por nome do template..."
            itemLabel="templates"
            emptyState={
              <EmptyState
                class="border-0 bg-transparent px-0 py-4 shadow-none"
                icon={<BadgeCheck class="size-6 text-foreground/70" />}
                eyebrow="Crachás da edição"
                title="Nenhum template cadastrado"
                description={
                  templateFilter()
                    ? "Nenhum template corresponde à busca informada."
                    : "Crie o primeiro modelo para começar a emitir crachás nesta edição."
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
                <AdminCreateBadgeCard
                  index={0}
                  animate={options.animate}
                  onCreate={() => handleOpenEditor()}
                />
                <For each={slice}>
                  {(template, index) => (
                    <AdminTemplateBadgeItem
                      template={template}
                      index={index() + 1}
                      animate={options.animate}
                      ticketNames={ticketNames}
                      location={location}
                      onEdit={() => handleOpenEditor(template.id)}
                      onDuplicate={() => handleOpenEditor(template.id, true)}
                      onDelete={() => handleDeleteTemplate(template)}
                    />
                  )}
                </For>
              </>
            )}
          />
        </Show>

        {/* Content: Emissions */}
        <Show when={activeSection() === "emissions"}>
          <PaginatedContainer<BadgePrintItem>
            items={filteredPrintItems()}
            layout="grid"
            minItemWidth="16rem"
            maxRows={(columns) => (columns === 1 ? 8 : 4)}
            gap="2"
            sort={emissionSort()}
            onSortChange={setEmissionSort}
            sortFields={[
              {
                key: "template_name",
                label: "Template",
                ascLabel: "A → Z",
                descLabel: "Z → A",
              },
              {
                key: "user_id",
                label: "Participante",
                ascLabel: "A → Z",
                descLabel: "Z → A",
              },
            ]}
            filterValue={emissionFilter()}
            onFilterChange={setEmissionFilter}
            filterPlaceholder="Buscar por evento, edição ou template..."
            itemLabel="crachás"
            headerActions={
              <div class="flex flex-wrap items-center gap-1.5 sm:gap-2">
                <Button
                  variant={printedAfter() ? "default" : "outline"}
                  size="sm"
                  onClick={() => setDateDialogOpen(true)}
                  class="relative h-9 gap-1.5 px-2.5 sm:px-3 text-xs cursor-pointer"
                >
                  <CalendarClock class="size-3.5 sm:size-4" />
                  <span class="hidden sm:inline">Filtrar data</span>
                  <span class="sm:hidden">Data</span>
                  <Show when={printedAfter()}>
                    <span class="size-1.5 rounded-full bg-primary-foreground" />
                  </Show>
                </Button>

                <Button
                  variant="default"
                  size="sm"
                  disabled={isPrinting() || printableItems().length === 0}
                  onClick={handlePrintBadges}
                  class="h-9 gap-1.5 px-2.5 sm:px-3 text-xs cursor-pointer"
                >
                  <Printer class="size-3.5 sm:size-4" />
                  <span>{isPrinting() && printMode() === "badges" ? "Preparando..." : "Crachás"}</span>
                </Button>

                <Button
                  variant="outline"
                  size="sm"
                  disabled={isPrinting() || printableItems().length === 0}
                  onClick={() => setQrDialogOpen(true)}
                  class="h-9 gap-1.5 px-2.5 sm:px-3 text-xs cursor-pointer"
                >
                  <QrCode class="size-3.5 sm:size-4" />
                  <span class="hidden sm:inline">Imprimir </span>
                  <span>QR codes</span>
                </Button>
              </div>
            }
            emptyState={
              <EmptyState
                class="border-0 bg-transparent px-0 py-4 shadow-none"
                icon={<BadgeCheck class="size-6 text-foreground/70" />}
                eyebrow="Crachás emitidos"
                title="Nenhum crachá emitido"
                description="Os crachás serão gerados automaticamente quando os participantes confirmarem suas inscrições."
              />
            }
            renderItems={(slice, options) => (
              <For each={slice}>
                {(badge, index) => (
                  <AdminEmissionBadgeItem
                    badge={badge}
                    index={index()}
                    animate={options.animate}
                    participantNames={participantNames}
                    location={location}
                  />
                )}
              </For>
            )}
          />
        </Show>
      </div>

      {/* Dialogs */}
      <DateFilterDialog
        open={dateDialogOpen()}
        onOpenChange={setDateDialogOpen}
        printedAfter={printedAfter()}
        onApply={setPrintedAfter}
      />

      <QrPrintDialog
        open={qrDialogOpen()}
        onOpenChange={setQrDialogOpen}
        loading={isPrinting()}
        onConfirmPrint={handlePrintQrs}
      />
    </>
  );
}
