import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import {
  BadgeCheck,
  CalendarClock,
  FileText,
  Plus,
  Printer,
  QrCode,
} from "lucide-react";
import QRCode from "qrcode";
import { useEffect, useMemo, useRef, useState } from "react";
import {
  badgeEmissionsQueryOptions,
  badgePrintQueryOptions,
  badgeTemplatesQueryOptions,
} from "@/features/badges/api";
import { useDeleteBadgeTemplateMutation } from "@/features/badges/api/mutations";
import type { BadgePrintItem, BadgeTemplate } from "@/features/badges/model";
import { badgeDesignSchema } from "@/features/badges/model";
import AdminBadgeCard from "@/features/badges/ui/AdminBadgeCard";
import { BadgePreview } from "@/features/badges/ui/badge-preview";
import { ToolbarCombobox } from "@/features/certifications/editor/ui/toolbar-combobox";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import { useActorDisplayNames } from "@/features/profile/api/actor-display-names";
import { allTicketsQueryOptions } from "@/features/tickets/api";
import { printElement } from "@/shared/lib/print-element";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/shared/ui/shadcn/dialog";
import { Input } from "@/shared/ui/shadcn/input";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/badges/",
)({ component: RouteComponent });

function PrintableBadge({
  badge,
  participantName,
  location,
}: {
  badge: BadgePrintItem;
  participantName: string;
  location: string;
}) {
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
        participantName={participantName}
        location={location}
        className="relative h-full w-full rounded-none border-0 shadow-none"
        style={{ width: "100%", height: "100%", aspectRatio: "auto" }}
      />
    </article>
  );
}

function PrintableQr({
  badge,
  size,
  participant,
}: {
  badge: BadgePrintItem;
  size: number;
  participant: string;
}) {
  const matrix = QRCode.create(badge.action_url).modules;
  const margin = 2;
  const viewSize = matrix.size + margin * 2;
  return (
    <article
      className="flex break-inside-avoid flex-col items-center gap-2 text-center text-black"
      style={{ width: size }}
    >
      <svg
        role="img"
        aria-label={`QR Code de ${participant}`}
        viewBox={`0 0 ${viewSize} ${viewSize}`}
        style={{ width: size, height: size }}
        shapeRendering="crispEdges"
      >
        <rect width="100%" height="100%" fill="white" />
        {Array.from({ length: matrix.size * matrix.size }, (_, index) => {
          const row = Math.floor(index / matrix.size);
          const column = index % matrix.size;
          return matrix.get(row, column) ? (
            <rect
              key={`${row}-${column}`}
              x={column + margin}
              y={row + margin}
              width="1"
              height="1"
            />
          ) : null;
        })}
      </svg>
      <strong className="max-w-full truncate text-sm">{participant}</strong>
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
  const [printMode, setPrintMode] = useState<"badges" | "qrs">("badges");
  const [printedAfter, setPrintedAfter] = useState("");
  const [printPending, setPrintPending] = useState(false);
  const [qrSizeMm, setQrSizeMm] = useState(48);
  const [printQrSizeMm, setPrintQrSizeMm] = useState(48);
  const [customQrSize, setCustomQrSize] = useState("48");
  const [dateDialogOpen, setDateDialogOpen] = useState(false);
  const [qrDialogOpen, setQrDialogOpen] = useState(false);
  const printRootRef = useRef<HTMLDivElement>(null);
  const presetQrSizes = [30, 48, 60, 210];
  const selectedPreset = presetQrSizes.includes(Number(customQrSize))
    ? customQrSize
    : "";
  const { data: templates = [] } = useQuery(
    badgeTemplatesQueryOptions(editionId),
  );
  const printQuery = useQuery(badgePrintQueryOptions(editionId));
  const { data: emissions = [] } = useQuery(
    badgeEmissionsQueryOptions(editionId),
  );
  const { data: editions = [] } = useQuery(
    allAdminEditionsQueryOptions(eventId),
  );
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
  const printableItems = useMemo(() => {
    if (!printedAfter) return printItems;
    const timestamp = new Date(printedAfter).getTime();
    const emissionIds = new Set(
      emissions
        .filter(
          (emission) => new Date(emission.emitted_at).getTime() >= timestamp,
        )
        .map((emission) => emission.id),
    );
    return printItems.filter((item) => emissionIds.has(item.emission_id));
  }, [emissions, printItems, printedAfter]);
  const printActorIds = [
    ...new Set(printableItems.map((item) => item.user_id)),
  ];
  const { data: participantNames = {}, isPending: participantNamesPending } =
    useActorDisplayNames(printActorIds);
  const location =
    editions.find((edition) => edition.id === editionId)?.location_name ?? "";
  const filteredPrintItems = useMemo(() => {
    const search = emissionFilter.trim().toLowerCase();
    if (!search) return printableItems;
    return printableItems.filter((item) =>
      [item.event_name, item.edition_name, item.ticket_name, item.template_name]
        .filter(Boolean)
        .some((value) => value?.toLowerCase().includes(search)),
    );
  }, [emissionFilter, printableItems]);

  async function printQrs() {
    const sizeMm = customQrSize ? Number(customQrSize) : qrSizeMm;
    if (!Number.isFinite(sizeMm) || sizeMm < 20 || sizeMm > 210) return;
    setPrintQrSizeMm(sizeMm);
    const result = await printQuery.refetch();
    if (result.data) {
      setQrDialogOpen(false);
      setPrintMode("qrs");
      setPrintPending(true);
    }
  }

  useEffect(() => {
    if (
      !printPending ||
      qrDialogOpen ||
      !printQuery.data ||
      participantNamesPending
    )
      return;

    let cancelled = false;
    const prepare = async () => {
      await new Promise<void>((resolve) =>
        requestAnimationFrame(() => requestAnimationFrame(() => resolve())),
      );
      if (cancelled || !printRootRef.current) return;
      await printElement(printRootRef.current, "Crachás");
      setPrintPending(false);
    };
    void prepare();
    return () => {
      cancelled = true;
    };
  }, [participantNamesPending, printPending, printQuery.data, qrDialogOpen]);
  return (
    <>
      <div
        ref={printRootRef}
        data-badge-print-root
        className="hidden print:block"
      >
        <div className="flex flex-wrap content-start gap-[4mm] p-[10mm] print:p-0">
          {printMode === "badges"
            ? printableItems.map((badge) => (
                <PrintableBadge
                  key={badge.emission_id}
                  badge={badge}
                  participantName={participantNames[badge.user_id]}
                  location={location}
                />
              ))
            : printableItems.map((badge) => (
                <PrintableQr
                  key={badge.emission_id}
                  badge={badge}
                  size={(printQrSizeMm / 25.4) * 96}
                  participant={participantNames[badge.user_id]}
                />
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
                  location={location}
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
                  onDelete={() =>
                    remove.mutate({ editionId, templateId: template.id })
                  }
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
              <div className="flex flex-wrap items-center gap-2">
                <Dialog open={dateDialogOpen} onOpenChange={setDateDialogOpen}>
                  <DialogTrigger
                    render={
                      <Button
                        variant={printedAfter ? "default" : "outline"}
                        className="h-9"
                      >
                        <CalendarClock className="size-4" />
                        Filtrar data
                      </Button>
                    }
                  />
                  {dateDialogOpen && (
                    <DialogContent className="w-[calc(100%-2rem)] sm:max-w-sm">
                      <DialogHeader>
                        <DialogTitle>Filtrar crachás por data</DialogTitle>
                        <DialogDescription>
                          Exiba e imprima somente os crachás gerados a partir da
                          data e hora escolhidas.
                        </DialogDescription>
                      </DialogHeader>
                      <label
                        htmlFor="badges-printed-after"
                        className="space-y-2 text-sm font-medium"
                      >
                        Gerados a partir de
                        <Input
                          id="badges-printed-after"
                          type="datetime-local"
                          value={printedAfter}
                          onChange={(event) =>
                            setPrintedAfter(event.target.value)
                          }
                          className="mt-2 h-10 w-full min-w-0 text-base sm:text-sm"
                        />
                      </label>
                      <div className="flex gap-2">
                        {printedAfter ? (
                          <Button
                            variant="outline"
                            className="flex-1"
                            onClick={() => setPrintedAfter("")}
                          >
                            Limpar
                          </Button>
                        ) : null}
                        <Button
                          className="flex-1"
                          onClick={() => setDateDialogOpen(false)}
                        >
                          Aplicar
                        </Button>
                      </div>
                    </DialogContent>
                  )}
                </Dialog>
                <Button
                  variant="default"
                  disabled={
                    printQuery.isFetching || printableItems.length === 0
                  }
                  onClick={async () => {
                    const result = await printQuery.refetch();
                    if (result.data) {
                      setPrintMode("badges");
                      setPrintPending(true);
                    }
                  }}
                  className={cn(
                    "inline-flex h-9 items-center justify-center gap-2 rounded-lg",
                  )}
                >
                  <Printer className="size-4" />
                  {printQuery.isFetching ? "Preparando…" : "Imprimir crachás"}
                </Button>
                <Dialog open={qrDialogOpen} onOpenChange={setQrDialogOpen}>
                  <DialogTrigger
                    render={
                      <Button
                        variant="outline"
                        className="h-9"
                        disabled={
                          printQuery.isFetching || printableItems.length === 0
                        }
                      >
                        <QrCode className="size-4" />
                        Imprimir QR codes
                      </Button>
                    }
                  />
                  {qrDialogOpen && (
                    <DialogContent className="print:hidden sm:max-w-md">
                      <DialogHeader>
                        <DialogTitle>
                          Configurar impressão dos QR codes
                        </DialogTitle>
                        <DialogDescription>
                          Escolha o tamanho exato de cada QR code na impressão.
                        </DialogDescription>
                      </DialogHeader>
                      <div className="space-y-2">
                        <span className="text-sm font-medium">
                          Tamanho do QR code
                        </span>
                        <ToolbarCombobox
                          value={selectedPreset}
                          options={[
                            { value: "30", label: "Pequeno - 30 × 30 mm" },
                            { value: "48", label: "Médio - 48 × 48 mm" },
                            { value: "60", label: "Grande - 60 × 60 mm" },
                            {
                              value: "210",
                              label: "Largura A4 - 210 × 210 mm",
                            },
                          ]}
                          placeholder="Selecione o tamanho"
                          searchPlaceholder="Buscar tamanho..."
                          onChange={(value) => {
                            setQrSizeMm(Number(value));
                            setCustomQrSize(value);
                          }}
                          className="w-full"
                          triggerClassName="h-10"
                        />
                        <p className="pt-2 text-xs text-muted-foreground">
                          Ou informe um tamanho personalizado entre 20 e 210 mm.
                        </p>
                        <Input
                          type="number"
                          min={20}
                          max={210}
                          step={1}
                          value={customQrSize}
                          onChange={(event) => {
                            const value = event.target.value;
                            setCustomQrSize(value);
                            if (value !== "") setQrSizeMm(Number(value));
                          }}
                          placeholder="Ex.: 200"
                          className="h-10"
                        />
                      </div>
                      <Button
                        className="h-10 w-full"
                        disabled={
                          printQuery.isFetching ||
                          (customQrSize !== "" &&
                            (!Number.isFinite(Number(customQrSize)) ||
                              Number(customQrSize) < 20 ||
                              Number(customQrSize) > 210))
                        }
                        onClick={() => void printQrs()}
                      >
                        {printQuery.isFetching
                          ? "Preparando impressão…"
                          : "Confirmar e imprimir"}
                      </Button>
                      {customQrSize !== "" &&
                      Number(customQrSize) >= 20 &&
                      Number(customQrSize) < 30 ? (
                        <p className="text-xs text-amber-600">
                          QR codes menores que 30 mm podem não ser lidos por
                          alguns scanners.
                        </p>
                      ) : null}
                    </DialogContent>
                  )}
                </Dialog>
              </div>
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
                  participantName={participantNames[badge.user_id]}
                  location={location}
                  onView={() => undefined}
                />
              ))
            }
          />
        )}
      </div>
    </>
  );
}
