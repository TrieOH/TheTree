import { useHotkeys } from "@tanstack/react-hotkeys";
import { useQueries, useQuery, useQueryClient } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import {
  DashboardBarList,
  DashboardLineChart,
  type DashboardLineChartPoint,
  DashboardStatCard,
} from "@trieoh/ui-base";
import {
  Activity,
  Calendar,
  CalendarRange,
  CheckCircle2,
  CircleAlert,
  Command,
  Copy,
  ExternalLink,
  Eye,
  Layers3,
  Package,
  ShoppingBag,
  Users,
} from "lucide-react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import {
  usePatchEditionMutation,
  usePublishEditionMutation,
} from "@/features/editions/api/mutations";
import { EditEditionModal } from "@/features/editions/ui/EditEditionModal";
import { EditionVisualCard } from "@/features/editions/ui/EditionVisualCard";
import {
  allJoinedEventsQueryOptions,
  allOwnEventsQueryOptions,
} from "@/features/events/api";
import { productsByEditionQueryOptions } from "@/features/products/api";
import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import {
  editionPurchasesQueryOptions,
  useRefundPurchaseMutation,
} from "@/features/purchases/api";
import { purchaseQueryKeys } from "@/features/purchases/api/query-keys";
import { allTicketsQueryOptions } from "@/features/tickets/api";
import { Badge } from "@/shared/ui/shadcn/badge";
import { Button } from "@/shared/ui/shadcn/button";
import { AlertModal } from "@/widgets/ui/alert-modal";
import { DashboardPanel } from "@/widgets/ui/dashboard-panel";
import { StepChecklist } from "@/widgets/ui/step-checklist";

export const Route = createLazyFileRoute(
  "/admin/events/$eventId_/editions/$editionId/",
)({
  component: AdminEditionDetailRoute,
});

function AdminEditionDetailRoute() {
  const { eventId, editionId } = Route.useParams();
  const navigate = Route.useNavigate();
  const { data: editions = [], isPending } = useQuery(
    allAdminEditionsQueryOptions(eventId),
  );
  const { data: ownEvents = [] } = useQuery(allOwnEventsQueryOptions());
  const { data: joinedEvents = [] } = useQuery(allJoinedEventsQueryOptions());
  const queryClient = useQueryClient();
  const { data: purchases = [] } = useQuery(
    editionPurchasesQueryOptions(editionId),
  );
  const [
    { data: tickets = [] },
    { data: products = [] },
    { data: programs = [] },
    { data: occurrences = [] },
  ] = useQueries({
    queries: [
      allTicketsQueryOptions(editionId),
      productsByEditionQueryOptions(editionId),
      programsQueryOptions(editionId),
      occurrencesQueryOptions(editionId),
    ],
  });
  const [publishConfirmOpen, setPublishConfirmOpen] = useState(false);
  const [editEditionOpen, setEditEditionOpen] = useState(false);
  const [refundPurchaseId, setRefundPurchaseId] = useState<string | null>(null);

  const edition = useMemo(
    () => editions.find((item) => item.id === editionId) ?? null,
    [editionId, editions],
  );
  const eventSlug = [...ownEvents, ...joinedEvents].find(
    (event) => event.id === eventId,
  )?.slug;

  const publishEditionMutation = usePublishEditionMutation();
  const patchEditionMutation = usePatchEditionMutation();
  const refundMutation = useRefundPurchaseMutation();

  const copyLink = () => {
    if (!edition || !eventSlug) return;
    void navigator.clipboard.writeText(
      `${window.location.origin}/events/${eventSlug}`,
    );
    toast.success("Link copiado");
  };

  const handlePublishEdition = () => {
    if (!edition) return;
    publishEditionMutation.mutate({ eventId, editionId });
  };

  const handleRefund = async () => {
    if (!refundPurchaseId) return;
    try {
      await refundMutation.mutateAsync(refundPurchaseId);
      toast.success("Reembolso solicitado");
      await queryClient.invalidateQueries({
        queryKey: purchaseQueryKeys.edition(editionId),
      });
      setRefundPurchaseId(null);
    } catch {
      toast.error("Não foi possível solicitar o reembolso.");
    }
  };

  useHotkeys(
    [
      {
        hotkey: "Mod+P",
        callback: () => setPublishConfirmOpen(true),
        options: { enabled: edition?.status === "draft" },
      },
      {
        hotkey: "Mod+E",
        callback: () => setEditEditionOpen(true),
        options: { enabled: Boolean(edition) },
      },
      {
        hotkey: "Mod+Shift+C",
        callback: copyLink,
        options: {
          enabled: Boolean(eventSlug && edition && edition.status !== "draft"),
        },
      },
      {
        hotkey: "Mod+Shift+O",
        callback: () => {
          if (edition && eventSlug && edition.status !== "draft") {
            void navigate({ to: "/events/$slug", params: { slug: eventSlug } });
          }
        },
        options: {
          enabled: Boolean(eventSlug && edition && edition.status !== "draft"),
        },
      },
    ],
    { ignoreInputs: true, preventDefault: true },
  );

  if (isPending) {
    return (
      <div className="flex min-h-72 items-center justify-center text-sm text-muted-foreground">
        Carregando edição...
      </div>
    );
  }

  if (!edition) {
    return (
      <div className="flex min-h-72 items-center justify-center text-sm text-muted-foreground">
        Edição não encontrada.
      </div>
    );
  }

  const isDraft = edition.status === "draft";
  const editionStatus = {
    draft: { label: "Rascunho", className: "bg-amber-500/10 text-amber-700" },
    future: { label: "Futura", className: "bg-sky-500/10 text-sky-700" },
    active: { label: "Ativa", className: "bg-emerald-500/10 text-emerald-700" },
    past: { label: "Encerrada", className: "bg-slate-500/10 text-slate-700" },
    archived: {
      label: "Arquivada",
      className: "bg-slate-500/10 text-slate-700",
    },
  }[edition.status] ?? {
    label: edition.status,
    className: "bg-muted text-muted-foreground",
  };
  const createdDate = new Date(edition.created_at)
    .toLocaleDateString("pt-BR", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    })
    .replace(".", "");

  const checklist = [
    {
      label: "Banner cadastrado",
      description: "Imagem principal exibida no topo da edição.",
      done: Boolean(edition.banner_url),
      action: edition.banner_url
        ? undefined
        : {
            label: "Adicionar",
            onClick: () =>
              document
                .getElementById(`edition-${edition.id}-banner-upload`)
                ?.click(),
          },
    },
    {
      label: "Logo cadastrado",
      description: "Identifica a edição nos cards e páginas públicas.",
      done: Boolean(edition.logo_url),
      action: edition.logo_url
        ? undefined
        : {
            label: "Adicionar",
            onClick: () =>
              document
                .getElementById(`edition-${edition.id}-logo-upload`)
                ?.click(),
          },
    },
    {
      label: "Descrição preenchida",
      description: "Apresente a edição para quem ainda não a conhece.",
      done: Boolean(edition.description),
      action: { label: "Editar", onClick: () => setEditEditionOpen(true) },
    },
    {
      label: "Tagline definida",
      description: "Resumo curto exibido junto ao nome da edição.",
      done: Boolean(edition.tagline),
      action: { label: "Editar", onClick: () => setEditEditionOpen(true) },
    },
    {
      label: "Local definido",
      description: "Local onde a edição será realizada.",
      done: Boolean(edition.location_name),
      action: { label: "Editar", onClick: () => setEditEditionOpen(true) },
    },
  ];

  const actions = [
    {
      label: "Editar edição",
      shortcut: "Mod+E",
      onClick: () => setEditEditionOpen(true),
      disabled: false,
      variant: "default" as const,
    },
    ...(isDraft
      ? [
          {
            label: "Publicar edição",
            shortcut: "Mod+P",
            onClick: () => setPublishConfirmOpen(true),
            disabled: publishEditionMutation.isPending,
            variant: "default" as const,
          },
        ]
      : []),
    {
      label: "Copiar link público",
      shortcut: "Mod+Shift+C",
      onClick: copyLink,
      disabled: isDraft,
      variant: "default" as const,
    },
    ...(!isDraft
      ? [
          {
            label: "Abrir página pública",
            shortcut: "Mod+Shift+O",
            onClick: () => {
              if (!eventSlug) return;
              void navigate({
                to: "/events/$slug",
                params: { slug: eventSlug },
              });
            },
            disabled: false,
            variant: "default" as const,
          },
        ]
      : []),
  ];

  const actionIcon = (label: string) => {
    if (label.includes("Publicar")) return Eye;
    if (label.includes("link")) return Copy;
    return ExternalLink;
  };
  const approvedPurchases = purchases.filter(
    (purchase) => purchase.status === "approved",
  ).length;
  const refundedPurchases = purchases.filter(
    (purchase) => purchase.status === "refunded",
  ).length;
  const revenue = purchases
    .filter((purchase) => purchase.status === "approved")
    .reduce((total, purchase) => total + purchase.total_cents, 0);
  const purchaseTimeline: DashboardLineChartPoint[] = purchases
    .filter((purchase) => purchase.created_at)
    .map((purchase) => ({
      timestamp: purchase.created_at ?? "",
      status: purchase.status,
      totalCents: purchase.total_cents,
    }));
  const purchaseStatusBars = [
    ["Aprovadas", "approved", "bg-emerald-500"],
    ["Pendentes", "pending", "bg-amber-500"],
    ["Reembolsadas", "refunded", "bg-sky-500"],
    ["Expiradas", "expired", "bg-slate-400"],
    ["Canceladas", "cancelled", "bg-rose-500"],
  ].map(([label, status, color]) => ({
    id: status,
    label,
    value: purchases.filter((purchase) => purchase.status === status).length,
    color,
  }));

  return (
    <div className="relative mx-auto max-w-7xl space-y-6 p-6 pb-28!">
      <EditionVisualCard edition={edition} eventId={eventId} />

      <div className="space-y-1 px-1 text-center md:text-left">
        <h1 className="text-xl font-medium tracking-tight text-foreground/90">
          {edition.name}
        </h1>
        <p className="flex items-center justify-center gap-1.5 text-xs text-muted-foreground md:justify-start">
          <Calendar className="size-3.5" />
          Criada em {createdDate}
        </p>
        <div>
          <Badge
            className={`w-fit border-0 px-2 py-0.5 text-xs font-normal ${editionStatus.className}`}
          >
            {editionStatus.label}
          </Badge>
        </div>
      </div>

      <div
        className="order-1 space-y-3 rounded-xl border border-border/60 bg-muted/20 p-3"
        role="toolbar"
        aria-label="Atalhos da edição"
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
            return (
              <Button
                key={action.label}
                size="default"
                variant="outline"
                className="h-9 shrink-0 flex-row gap-1.5 border border-border px-2 text-xs sm:h-14! sm:min-w-28! sm:flex-col! sm:gap-1 sm:px-2 sm:py-1.5 sm:text-[11px] sm:leading-tight"
                disabled={action.disabled}
                onClick={action.onClick}
                title={`${action.label} · ${action.shortcut}`}
                aria-label={action.label}
              >
                <span className="flex items-center gap-1.5">
                  <Icon className="size-4" />
                  {action.label.replace(" edição", "")}
                </span>
                <kbd className="hidden rounded border border-border/70 bg-muted/70 px-1 py-0.5 font-mono text-[9px] text-muted-foreground sm:inline-block">
                  {action.shortcut.replace("Mod", "⌘/Ctrl")}
                </kbd>
              </Button>
            );
          })}
        </div>
      </div>

      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        {[
          {
            label: "Receita aprovada",
            value: new Intl.NumberFormat("pt-BR", {
              style: "currency",
              currency: "BRL",
            }).format(revenue / 100),
            hint: "Compras aprovadas",
            Icon: ShoppingBag,
          },
          {
            label: "Compras",
            value: purchases.length,
            hint: "Pedidos registrados",
            Icon: Users,
          },
          {
            label: "Aprovadas",
            value: approvedPurchases,
            hint: "Prontas para uso",
            Icon: CheckCircle2,
          },
          {
            label: "Reembolsadas",
            value: refundedPurchases,
            hint: "Compras reembolsadas",
            Icon: CircleAlert,
          },
        ].map(({ label, value, hint, Icon }) => (
          <DashboardStatCard
            key={String(label)}
            label={label}
            value={value}
            hint={hint}
            icon={Icon}
          />
        ))}
      </section>

      <DashboardPanel
        title="Catálogo"
        description="Ingressos, produtos e programação desta edição."
        icon={Package}
      >
        <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          {[
            {
              label: "Ingressos",
              value: tickets.length,
              hint: "Tipos cadastrados",
              icon: Users,
            },
            {
              label: "Produtos",
              value: products.length,
              hint: "Produtos cadastrados",
              icon: Package,
            },
            {
              label: "Programas",
              value: programs.length,
              hint: "Atividades cadastradas",
              icon: CalendarRange,
            },
            {
              label: "Ocorrências",
              value: occurrences.length,
              hint: "Horários cadastrados",
              icon: Layers3,
            },
          ].map((metric) => (
            <DashboardStatCard key={metric.label} {...metric} />
          ))}
        </div>
      </DashboardPanel>

      <section className="grid gap-4 xl:grid-cols-[1fr_1.8fr]">
        <DashboardPanel
          title="Status das compras"
          description={`${purchases.length} compra${purchases.length === 1 ? "" : "s"} registrada${purchases.length === 1 ? "" : "s"}.`}
          icon={ShoppingBag}
          className="h-full rounded-lg bg-card p-5 ring-1 ring-foreground/10"
        >
          <DashboardBarList
            items={purchaseStatusBars}
            maxValue={Math.max(
              ...purchaseStatusBars.map((item) => item.value),
              1,
            )}
            emptyMessage="Nenhuma compra registrada."
          />
        </DashboardPanel>

        <DashboardPanel
          title="Lucro"
          description="Crescimento acumulado das compras da edição."
          icon={Activity}
          className="h-full rounded-lg bg-card p-5 ring-1 ring-foreground/10"
        >
          <DashboardLineChart points={purchaseTimeline} />
        </DashboardPanel>
      </section>

      <StepChecklist
        title="Checklist da edição"
        items={checklist.map((item) => ({
          id: item.label,
          title: item.label,
          description: item.description,
          completed: item.done,
          action: item.action,
        }))}
        className="w-full"
        mobileInline
      />

      <AlertModal
        open={publishConfirmOpen}
        onOpenChange={setPublishConfirmOpen}
        title="Publicar edição?"
        description="Depois de publicar, a edição ficará visível ao público."
        confirmLabel="Publicar edição"
        variant="default"
        loading={publishEditionMutation.isPending}
        onConfirm={async () => {
          handlePublishEdition();
          setPublishConfirmOpen(false);
        }}
      />
      <EditEditionModal
        open={editEditionOpen}
        edition={edition}
        onOpenChange={setEditEditionOpen}
        onUpdate={(values) =>
          patchEditionMutation.mutateAsync({
            eventId,
            editionId,
            data: values,
          })
        }
      />
      <AlertModal
        open={refundPurchaseId !== null}
        onOpenChange={(open) => {
          if (!open && !refundMutation.isPending) setRefundPurchaseId(null);
        }}
        title="Solicitar reembolso?"
        description="O pagamento será encaminhado para reembolso. A compra ficará como reembolsada após a confirmação do webhook."
        confirmLabel="Solicitar reembolso"
        variant="destructive"
        loading={refundMutation.isPending}
        onConfirm={handleRefund}
      />
    </div>
  );
}
