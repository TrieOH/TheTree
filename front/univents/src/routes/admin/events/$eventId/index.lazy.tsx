import { useHotkeys } from "@tanstack/react-hotkeys";
import { useQueries, useQuery } from "@tanstack/react-query";
import { createLazyFileRoute, Link } from "@tanstack/react-router";
import {
  DashboardBarList,
  DashboardLineChart,
  type DashboardLineChartPoint,
  DashboardStatCard,
} from "@trieoh/ui-base";
import { format } from "date-fns";
import { ptBR } from "date-fns/locale";
import {
  Activity,
  Calendar,
  CalendarRange,
  CheckCircle2,
  ChevronRight,
  CircleAlert,
  Command,
  Copy,
  CreditCard,
  ExternalLink,
  Eye,
  Layers3,
  Package,
  Pencil,
  ShoppingBag,
  Ticket,
  Wallet,
  XCircle,
} from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import {
  allJoinedEventsQueryOptions,
  allOwnEventsQueryOptions,
} from "@/features/events/api";
import {
  useDiscontinueEventMutation,
  usePatchEventMutation,
  usePublishEventMutation,
} from "@/features/events/api/mutations";
import { EventVisualCard } from "@/features/events/ui/EventVisualCard";
import { ManageEventModal } from "@/features/events/ui/ManageEventModal";
import {
  useConnectEventSellerMutation,
  useDisconnectEventSellerMutation,
} from "@/features/payments/api/mutations";
import type { PaymentProviderI } from "@/features/payments/model";
import { productsByEditionQueryOptions } from "@/features/products/api";
import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import { editionPurchasesQueryOptions } from "@/features/purchases/api";
import { allTicketsQueryOptions } from "@/features/tickets/api";
import { useUploadQueue } from "@/features/upload-queue";
import { Badge } from "@/shared/ui/shadcn/badge";
import { Button } from "@/shared/ui/shadcn/button";
import { AlertModal } from "@/widgets/ui/alert-modal";
import { DashboardPanel } from "@/widgets/ui/dashboard-panel";
import { StepChecklist } from "@/widgets/ui/step-checklist";

const editionStatusConfig = {
  draft: {
    label: "Rascunho",
    className: "border-amber-500/20 bg-amber-500/10 text-amber-700",
    dot: "bg-amber-500",
  },
  future: {
    label: "Futura",
    className: "border-sky-500/20 bg-sky-500/10 text-sky-700",
    dot: "bg-sky-500",
  },
  active: {
    label: "Ativa",
    className: "border-emerald-500/20 bg-emerald-500/10 text-emerald-700",
    dot: "bg-emerald-500",
  },
  past: {
    label: "Encerrada",
    className: "border-slate-500/20 bg-slate-500/10 text-slate-600",
    dot: "bg-slate-500",
  },
} as const;

const paymentProviders: Array<{
  id: PaymentProviderI;
  name: string;
  description: string;
}> = [
  {
    id: "mercadopago",
    name: "Mercado Pago",
    description: "Receba por Pix e cartão de crédito.",
  },
];

export const Route = createLazyFileRoute("/admin/events/$eventId/")({
  component: EventOverviewRoute,
});

function EventOverviewRoute() {
  const { eventId } = Route.useParams();
  const navigate = Route.useNavigate();
  const { data: ownedEvents = [] } = useQuery(allOwnEventsQueryOptions());
  const { data: joinedEvents = [] } = useQuery(allJoinedEventsQueryOptions());
  const { data: editions = [] } = useQuery(
    allAdminEditionsQueryOptions(eventId),
  );
  const purchaseQueries = useQueries({
    queries: editions.map((edition) =>
      editionPurchasesQueryOptions(edition.id),
    ),
  });
  const ticketQueries = useQueries({
    queries: editions.map((edition) => allTicketsQueryOptions(edition.id)),
  });
  const productQueries = useQueries({
    queries: editions.map((edition) =>
      productsByEditionQueryOptions(edition.id),
    ),
  });
  const programQueries = useQueries({
    queries: editions.map((edition) => programsQueryOptions(edition.id)),
  });
  const occurrenceQueries = useQueries({
    queries: editions.map((edition) => occurrencesQueryOptions(edition.id)),
  });
  const { tasks } = useUploadQueue();
  const [publishConfirmOpen, setPublishConfirmOpen] = useState(false);
  const [discontinueConfirmOpen, setDiscontinueConfirmOpen] = useState(false);
  const [disconnectSellerConfirmOpen, setDisconnectSellerConfirmOpen] =
    useState(false);
  const [editEventOpen, setEditEventOpen] = useState(false);
  const event =
    [...ownedEvents, ...joinedEvents].find((item) => item.id === eventId) ??
    null;
  const isPublished = event?.status === "active";
  const eventStatus = event
    ? {
        draft: {
          label: "Rascunho",
          className: "bg-amber-500/10 text-amber-700",
        },
        active: {
          label: "Ativo",
          className: "bg-emerald-500/10 text-emerald-700",
        },
        archived: {
          label: "Arquivado",
          className: "bg-slate-500/10 text-slate-700",
        },
        discontinued: {
          label: "Descontinuado",
          className: "bg-rose-500/10 text-rose-700",
        },
      }[event.status]
    : null;
  const createdDate = event
    ? new Date(event.created_at)
        .toLocaleDateString("pt-BR", {
          day: "2-digit",
          month: "short",
          year: "numeric",
        })
        .replace(".", "")
    : "";
  const purchases = purchaseQueries.flatMap((query) => query.data ?? []);
  const approvedPurchases = purchases.filter(
    (purchase) => purchase.status === "approved",
  );
  const refundedPurchases = purchases.filter(
    (purchase) => purchase.status === "refunded",
  );
  const revenue = approvedPurchases.reduce(
    (total, purchase) => total + purchase.total_cents,
    0,
  );
  const ticketCount = ticketQueries.reduce(
    (total, query) => total + (query.data?.length ?? 0),
    0,
  );
  const productCount = productQueries.reduce(
    (total, query) => total + (query.data?.length ?? 0),
    0,
  );
  const programCount = programQueries.reduce(
    (total, query) => total + (query.data?.length ?? 0),
    0,
  );
  const occurrenceCount = occurrenceQueries.reduce(
    (total, query) => total + (query.data?.length ?? 0),
    0,
  );
  const editionSales = editions.map((edition, index) => {
    const editionPurchases = purchaseQueries[index]?.data ?? [];
    const approvedRevenue = editionPurchases
      .filter((purchase) => purchase.status === "approved")
      .reduce((total, purchase) => total + purchase.total_cents, 0);
    return {
      name: edition.name,
      purchases: editionPurchases.length,
      revenue: approvedRevenue,
    };
  });
  const maxEditionRevenue = Math.max(
    ...editionSales.map((edition) => edition.revenue),
    1,
  );
  const purchaseStatuses = [
    { label: "Aprovadas", status: "approved", color: "bg-emerald-500" },
    { label: "Pendentes", status: "pending", color: "bg-amber-500" },
    { label: "Reembolsadas", status: "refunded", color: "bg-sky-500" },
    { label: "Expiradas", status: "expired", color: "bg-slate-400" },
    { label: "Canceladas", status: "cancelled", color: "bg-rose-500" },
  ] as const;
  const revenueBars = editionSales.slice(0, 6).map((edition) => ({
    id: edition.name,
    label: edition.name,
    value: edition.revenue,
    detail: new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: "BRL",
    }).format(edition.revenue / 100),
  }));
  const statusBars = purchaseStatuses.map((item) => ({
    id: item.status,
    label: item.label,
    value: purchases.filter((purchase) => purchase.status === item.status)
      .length,
    color: item.color,
  }));
  const purchaseTimeline: DashboardLineChartPoint[] = purchases
    .filter((purchase) => purchase.created_at)
    .map((purchase) => ({
      timestamp: purchase.created_at ?? "",
      status: purchase.status,
      totalCents: purchase.total_cents,
    }));
  const summaryMetrics = [
    {
      label: "Receita aprovada",
      value: new Intl.NumberFormat("pt-BR", {
        style: "currency",
        currency: "BRL",
      }).format(revenue / 100),
      hint: "Somente compras aprovadas",
      icon: Wallet,
    },
    {
      label: "Compras",
      value: purchases.length,
      hint: `${approvedPurchases.length} aprovada(s)`,
      icon: ShoppingBag,
    },
    {
      label: "Aprovadas",
      value: approvedPurchases.length,
      hint: "Prontas para uso",
      icon: CheckCircle2,
    },
    {
      label: "Reembolsadas",
      value: refundedPurchases.length,
      hint: "Compras reembolsadas",
      icon: CircleAlert,
    },
  ];
  const isImageUploading = (field: "logo_url" | "banner_url") =>
    tasks.some(
      (task) =>
        task.owner.type === "event" &&
        task.owner.id === eventId &&
        task.association?.handlerKey === "event-image" &&
        task.association.input?.field === field &&
        !["completed", "failed", "rejected"].includes(task.status),
    );

  const publishEventMutation = usePublishEventMutation();
  const discontinueEventMutation = useDiscontinueEventMutation();
  const patchEventMutation = usePatchEventMutation();
  const connectSellerMutation = useConnectEventSellerMutation();
  const disconnectSellerMutation = useDisconnectEventSellerMutation();

  const copyLink = () => {
    if (!event) return;
    void navigator.clipboard.writeText(
      `${window.location.origin}/events/${event.slug}`,
    );
    toast.success("Link copiado");
  };

  const handlePublishEvent = () => {
    if (!event || isPublished) return;
    publishEventMutation.mutate(eventId);
  };

  const handleDiscontinueEvent = () => {
    if (!event || event.status !== "active") return;
    discontinueEventMutation.mutate(eventId);
  };

  const checklist = [
    {
      label: "Edição criada",
      description:
        "Crie uma edição para publicar datas, catálogo e programação.",
      done: editions.length > 0,
    },
    {
      label: "Logo cadastrado",
      description: "Identifica o evento nos cards e páginas públicas.",
      done: Boolean(event?.logo_url),
      action: event?.logo_url
        ? undefined
        : {
            label: "Adicionar",
            disabled: isImageUploading("logo_url"),
            onClick: () =>
              document.getElementById("event-logo-upload")?.click(),
          },
    },
    {
      label: "Banner cadastrado",
      description: "Imagem principal exibida no topo do evento.",
      done: Boolean(event?.banner_url),
      action: event?.banner_url
        ? undefined
        : {
            label: "Adicionar",
            disabled: isImageUploading("banner_url"),
            onClick: () =>
              document.getElementById("event-banner-upload")?.click(),
          },
    },
    {
      label: "Descrição preenchida",
      description: "Apresente o evento para quem ainda não o conhece.",
      done: Boolean(event?.description),
      action: {
        label: "Editar",
        onClick: () => setEditEventOpen(true),
      },
    },
    ...(editions.length > 0
      ? [
          {
            label: "Pagamento conectado",
            description: "Necessário para vender ingressos ou produtos.",
            done: Boolean(event?.payssage_seller_id),
          },
        ]
      : []),
  ];

  const actions = [
    {
      label: "Editar evento",
      shortcut: "Mod+E",
      onClick: () => setEditEventOpen(true),
      disabled: !event,
      variant: "default" as const,
    },
    ...(event?.status === "draft"
      ? [
          {
            label: "Publicar evento",
            shortcut: "Mod+P",
            onClick: () => setPublishConfirmOpen(true),
            disabled: publishEventMutation.isPending,
            variant: "default" as const,
          },
        ]
      : []),
    {
      label: "Copiar link público",
      shortcut: "Mod+Shift+C",
      onClick: copyLink,
      disabled: !event,
      variant: "default" as const,
    },
    ...(isPublished
      ? [
          {
            label: "Descontinuar evento",
            shortcut: "Mod+Shift+D",
            onClick: () => setDiscontinueConfirmOpen(true),
            disabled: discontinueEventMutation.isPending,
            variant: "destructive" as const,
          },
          {
            label: "Abrir painel público",
            shortcut: "Mod+Shift+O",
            to: "/events/$slug" as const,
            params: { slug: event?.slug ?? "" },
            variant: "default" as const,
          },
        ]
      : []),
  ];

  useHotkeys(
    [
      {
        hotkey: "Mod+E",
        callback: () => setEditEventOpen(true),
        options: { enabled: Boolean(event) },
      },
      {
        hotkey: "Mod+P",
        callback: () => setPublishConfirmOpen(true),
        options: { enabled: event?.status === "draft" },
      },
      {
        hotkey: "Mod+Shift+C",
        callback: copyLink,
        options: { enabled: Boolean(event) },
      },
      {
        hotkey: "Mod+Shift+D",
        callback: () => setDiscontinueConfirmOpen(true),
        options: { enabled: isPublished },
      },
      {
        hotkey: "Mod+Shift+O",
        callback: () => {
          if (event) {
            void navigate({
              to: "/events/$slug",
              params: { slug: event.slug },
            });
          }
        },
        options: { enabled: isPublished && Boolean(event) },
      },
    ],
    { ignoreInputs: true, preventDefault: true },
  );

  const actionIcon = (label: string) => {
    if (label.includes("Editar")) return Pencil;
    if (label.includes("Copiar")) return Copy;
    if (label.includes("Publicar")) return Eye;
    if (label.includes("Descontinuar")) return XCircle;
    return ExternalLink;
  };

  return (
    <>
      <div className="relative mx-auto flex w-full max-w-7xl flex-col gap-6 p-6 pb-28!">
        {event ? <EventVisualCard event={event} /> : null}

        {event && eventStatus ? (
          <div className="space-y-1 px-1 text-center md:text-left">
            <h1 className="text-xl font-medium tracking-tight text-foreground/90">
              {event.full_name}
            </h1>
            <p className="flex items-center justify-center gap-1.5 text-xs text-muted-foreground md:justify-start">
              <Calendar className="size-3.5" />
              Criado em {createdDate}
            </p>
            <div>
              <Badge
                className={`w-fit border-0 px-2 py-0.5 text-xs font-normal ${eventStatus.className}`}
              >
                {eventStatus.label}
              </Badge>
            </div>
          </div>
        ) : null}

        <div
          className="order-1 space-y-3 rounded-xl border border-border/60 bg-muted/20 p-3"
          role="toolbar"
          aria-label="Atalhos do evento"
        >
          <div className="flex items-center justify-between gap-3 px-1">
            <div className="flex min-w-0 items-center gap-2">
              <div className="flex size-7 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary">
                <Command className="size-3.5" />
              </div>
              <div className="min-w-0">
                <p className="text-xs font-semibold text-foreground">
                  Ações rápidas
                </p>
                <p className="truncate text-[11px] text-muted-foreground">
                  Atalhos para as tarefas mais usadas
                </p>
              </div>
            </div>
            <span className="hidden shrink-0 text-[11px] text-muted-foreground sm:inline">
              Ctrl/⌘ + tecla
            </span>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            {actions.map((action) => {
              const Icon = actionIcon(action.label);
              if ("to" in action && action.to) {
                return (
                  <Link
                    key={action.label}
                    to={action.to}
                    params={action.params}
                    aria-disabled={action.disabled}
                    title={`${action.label} · ${action.shortcut}`}
                    aria-label={action.label}
                    className="inline-flex h-9 shrink-0 flex-row items-center justify-center gap-1.5 rounded-md border border-border bg-background px-2 text-xs font-medium text-foreground shadow-xs transition-colors hover:bg-muted aria-disabled:pointer-events-none aria-disabled:opacity-50 sm:h-14! sm:min-w-28! sm:flex-col! sm:gap-1 sm:px-2 sm:py-1.5 sm:text-[11px] sm:leading-tight"
                  >
                    <span className="flex items-center gap-1.5">
                      <Icon className="size-4" />
                      <span>
                        {action.label
                          .replace(" evento", "")
                          .replace(" público", "")
                          .replace(" conta", "")}
                      </span>
                    </span>
                    <kbd className="hidden rounded border border-border/70 bg-muted px-1 py-0.5 font-mono text-[9px] text-muted-foreground sm:inline-block">
                      {action.shortcut.replace("Mod", "⌘/Ctrl")}
                    </kbd>
                  </Link>
                );
              }

              return (
                <Button
                  key={action.label}
                  size="default"
                  variant={
                    action.variant === "destructive" ? "destructive" : "outline"
                  }
                  className="h-9 shrink-0 flex-row gap-1.5 px-2 text-xs sm:h-14! sm:min-w-28! sm:flex-col! sm:gap-1 sm:px-2 sm:py-1.5 sm:text-[11px] sm:leading-tight"
                  disabled={action.disabled}
                  onClick={action.onClick}
                  title={`${action.label} · ${action.shortcut}`}
                  aria-label={action.label}
                >
                  <span className="flex items-center gap-1.5">
                    <Icon className="size-4" />
                    <span>
                      {action.label
                        .replace(" evento", "")
                        .replace(" público", "")
                        .replace(" conta", "")}
                    </span>
                  </span>
                  <kbd className="hidden rounded border border-border/70 bg-muted/70 px-1 py-0.5 font-mono text-[9px] text-muted-foreground sm:inline-block">
                    {action.shortcut.replace("Mod", "⌘/Ctrl")}
                  </kbd>
                </Button>
              );
            })}
          </div>
        </div>

        <section className="order-2 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          {summaryMetrics.map((metric) => (
            <DashboardStatCard
              key={metric.label}
              label={metric.label}
              value={metric.value}
              hint={metric.hint}
              icon={metric.icon}
            />
          ))}
        </section>

        <DashboardPanel
          title="Catálogo"
          description="Conteúdo publicado nas edições deste evento."
          icon={Package}
          className="order-4"
        >
          <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
            {[
              {
                label: "Ingressos",
                value: ticketCount,
                hint: "Tipos cadastrados",
                icon: Ticket,
              },
              {
                label: "Produtos",
                value: productCount,
                hint: "Produtos cadastrados",
                icon: Package,
              },
              {
                label: "Programas",
                value: programCount,
                hint: "Atividades cadastradas",
                icon: CalendarRange,
              },
              {
                label: "Ocorrências",
                value: occurrenceCount,
                hint: "Horários da programação",
                icon: Layers3,
              },
            ].map((metric) => (
              <DashboardStatCard
                key={metric.label}
                label={metric.label}
                value={metric.value}
                hint={metric.hint}
                icon={metric.icon}
              />
            ))}
          </div>
        </DashboardPanel>

        <section className="order-5 grid gap-4 xl:grid-cols-[1fr_1.35fr]">
          <DashboardPanel
            title="Status das compras"
            description={`${purchases.length} compra${purchases.length === 1 ? "" : "s"} no total.`}
            icon={ShoppingBag}
            className="rounded-lg bg-card p-5 ring-1 ring-foreground/10"
          >
            <div className="mt-2">
              <DashboardBarList
                items={statusBars}
                emptyMessage="Nenhuma compra registrada."
              />
            </div>
          </DashboardPanel>
          <DashboardPanel
            title="Receita por edição"
            description="Receita aprovada comparada entre as edições."
            icon={Wallet}
            className="rounded-lg bg-card p-5 ring-1 ring-foreground/10"
          >
            <div className="mt-2">
              <DashboardBarList
                items={revenueBars}
                maxValue={maxEditionRevenue}
                emptyMessage="Ainda não há edições para comparar."
              />
            </div>
          </DashboardPanel>
        </section>

        <DashboardPanel
          title="Lucro"
          description="Crescimento acumulado no período selecionado."
          icon={Activity}
          className="order-6"
        >
          <DashboardLineChart points={purchaseTimeline} />
        </DashboardPanel>

        <section className="order-7 space-y-3">
          <div className="flex min-w-0 items-center justify-between gap-3 px-1">
            <div className="flex min-w-0 flex-1 items-center gap-3">
              <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <Layers3 className="size-5" />
              </div>
              <div className="min-w-0 flex-1">
                <h2 className="text-base font-semibold tracking-tight">
                  Edições
                </h2>
                <p className="truncate text-xs text-muted-foreground">
                  Acesse rapidamente cada edição do evento.
                </p>
              </div>
            </div>
            <span className="flex size-8 shrink-0 items-center justify-center rounded-full border border-border/60 bg-muted/50 p-0 text-xs font-medium text-muted-foreground">
              {editions.length}
            </span>
          </div>
          <div className="overflow-hidden rounded-lg border border-border/60 bg-card/95 shadow-sm">
            {editions.length === 0 ? (
              <p className="p-5 text-sm text-muted-foreground">
                Nenhuma edição cadastrada.
              </p>
            ) : (
              editions.slice(0, 6).map((edition, index) => {
                const editionStatus = editionStatusConfig[edition.status];
                return (
                  <Link
                    key={edition.id}
                    to="/admin/events/$eventId/editions/$editionId"
                    params={{ eventId, editionId: edition.id }}
                    className={`group flex items-center justify-between gap-4 px-5 py-4 transition-colors hover:bg-muted/30 ${index > 0 ? "border-t border-border/60" : ""}`}
                  >
                    <div className="flex min-w-0 items-center gap-3">
                      <span className="flex size-8 shrink-0 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">
                        {index + 1}
                      </span>
                      <div className="min-w-0">
                        <p className="truncate text-sm font-semibold">
                          {edition.name}
                        </p>
                        <p className="mt-1 truncate text-xs text-muted-foreground">
                          {format(new Date(edition.starts_at), "dd MMM yyyy", {
                            locale: ptBR,
                          })}
                          {edition.location_name
                            ? ` · ${edition.location_name}`
                            : ""}
                        </p>
                      </div>
                    </div>
                    <div className="flex shrink-0 items-center gap-3">
                      <span
                        className={`hidden items-center gap-1.5 rounded-full border px-2 py-0.5 text-[11px] sm:inline-flex ${editionStatus.className}`}
                      >
                        <span
                          className={`size-1.5 rounded-full ${editionStatus.dot}`}
                        />
                        {editionStatus.label}
                      </span>
                      <ChevronRight className="size-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
                    </div>
                  </Link>
                );
              })
            )}
          </div>
        </section>

        <section className="order-3 space-y-3">
          <div className="flex items-center gap-3 px-1">
            <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <CreditCard className="size-5" />
            </div>
            <div className="min-w-0">
              <h2 className="text-base font-semibold tracking-tight">
                Pagamentos
              </h2>
              <p className="truncate text-xs text-muted-foreground">
                Conta que receberá as vendas deste evento.
              </p>
            </div>
          </div>
          <div className="flex flex-wrap gap-3">
            {paymentProviders.map((provider) => {
              const connected = Boolean(event?.payssage_seller_id);

              return (
                <div
                  key={provider.id}
                  className="flex w-full max-w-md items-center gap-3 rounded-xl bg-card px-3 py-3 ring-1 ring-foreground/10 shadow-xs"
                >
                  <div className="flex size-12 shrink-0 items-center justify-center rounded-lg bg-muted/50 p-2">
                    <img
                      src="/mercado-pago.svg"
                      alt="Mercado Pago"
                      className="size-full object-contain"
                    />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="truncate text-sm font-semibold">
                      {provider.name}
                    </p>
                    <Badge
                      className="mt-1"
                      variant={connected ? "default" : "secondary"}
                    >
                      {connected ? "Conectado" : "Não conectado"}
                    </Badge>
                  </div>
                  <div className="shrink-0">
                    <Button
                      className="h-9 text-xs"
                      variant={connected ? "outline" : "default"}
                      disabled={
                        !event ||
                        connectSellerMutation.isPending ||
                        disconnectSellerMutation.isPending
                      }
                      onClick={() =>
                        connected
                          ? setDisconnectSellerConfirmOpen(true)
                          : connectSellerMutation.mutate({
                              eventId,
                              provider: provider.id,
                            })
                      }
                    >
                      {connected ? "Desconectar conta" : "Conectar"}
                    </Button>
                  </div>
                </div>
              );
            })}
          </div>
        </section>

        <StepChecklist
          title="Event checklist"
          items={checklist.map((item) => ({
            id: item.label,
            title: item.label,
            description: item.description,
            completed: item.done,
            action: item.action,
          }))}
          className="order-8 w-full sm:fixed sm:right-4 sm:top-24 sm:z-40 sm:w-auto!"
          mobileInline
        />
      </div>

      <ManageEventModal
        key={event?.id ?? "event"}
        open={editEventOpen}
        onOpenChange={setEditEventOpen}
        event={event}
        onCreate={(values) =>
          event
            ? patchEventMutation.mutateAsync({ eventId, data: values })
            : false
        }
      />

      <AlertModal
        open={disconnectSellerConfirmOpen}
        onOpenChange={setDisconnectSellerConfirmOpen}
        title="Desconectar Mercado Pago?"
        description="Este evento deixará de receber novos pagamentos até uma conta ser conectada novamente."
        confirmLabel="Desconectar"
        variant="destructive"
        loading={disconnectSellerMutation.isPending}
        onConfirm={() => {
          disconnectSellerMutation.mutate(eventId, {
            onSuccess: () => setDisconnectSellerConfirmOpen(false),
          });
        }}
      />

      <AlertModal
        open={publishConfirmOpen}
        onOpenChange={setPublishConfirmOpen}
        title="Publicar evento?"
        description="Depois de publicar, o painel público ficará disponível para o evento."
        confirmLabel="Publicar evento"
        variant="default"
        loading={publishEventMutation.isPending}
        onConfirm={async () => {
          handlePublishEvent();
          setPublishConfirmOpen(false);
        }}
      />

      <AlertModal
        open={discontinueConfirmOpen}
        onOpenChange={setDiscontinueConfirmOpen}
        title="Descontinuar evento?"
        description="O evento deixará de ser ativo e a data de atualização será atualizada."
        confirmLabel="Descontinuar evento"
        variant="destructive"
        loading={discontinueEventMutation.isPending}
        onConfirm={async () => {
          handleDiscontinueEvent();
          setDiscontinueConfirmOpen(false);
        }}
      />
    </>
  );
}
