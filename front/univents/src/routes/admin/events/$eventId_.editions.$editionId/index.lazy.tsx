import { useHotkeys } from "@tanstack/react-hotkeys";
import { useQueries, useQuery, useQueryClient } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import {
  Activity,
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
import { usePublishEditionMutation } from "@/features/editions/api/mutations";
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
import { Button } from "@/shared/ui/shadcn/button";
import { AlertModal } from "@/widgets/ui/alert-modal";
import { DashboardBarList } from "@/widgets/ui/dashboard-bar-list";
import {
  DashboardLineChart,
  type DashboardPurchasePoint,
} from "@/widgets/ui/dashboard-line-chart";
import { DashboardPanel } from "@/widgets/ui/dashboard-panel";
import { DashboardStatCard } from "@/widgets/ui/dashboard-stat-card";
import { StepChecklist } from "@/widgets/ui/step-checklist";

export const Route = createLazyFileRoute(
  "/admin/events/$eventId_/editions/$editionId/",
)({
  component: AdminEditionDetailRoute,
});

function AdminEditionDetailRoute() {
  const { eventId, editionId } = Route.useParams();
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
  const [refundPurchaseId, setRefundPurchaseId] = useState<string | null>(null);

  const edition = useMemo(
    () => editions.find((item) => item.id === editionId) ?? null,
    [editionId, editions],
  );
  const eventSlug = [...ownEvents, ...joinedEvents].find(
    (event) => event.id === eventId,
  )?.slug;

  const publishEditionMutation = usePublishEditionMutation();
  const refundMutation = useRefundPurchaseMutation();

  const copyLink = () => {
    if (!edition || !eventSlug) return;
    void navigator.clipboard.writeText(
      `${window.location.origin}/events/${eventSlug}/editions/${edition.slug}`,
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
            window.location.href = `/events/${eventSlug}/editions/${edition.slug}`;
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

  const checklist = [
    {
      label: "Banner cadastrado",
      description: "Imagem principal exibida no topo da edição.",
      done: Boolean(edition.banner_url),
    },
    {
      label: "Logo cadastrado",
      description: "Identifica a edição nos cards e páginas públicas.",
      done: Boolean(edition.logo_url),
    },
    {
      label: "Descrição preenchida",
      description: "Apresente a edição para quem ainda não a conhece.",
      done: Boolean(edition.description),
    },
    {
      label: "Tagline definida",
      description: "Resumo curto exibido junto ao nome da edição.",
      done: Boolean(edition.tagline),
    },
    {
      label: "Local definido",
      description: "Local onde a edição será realizada.",
      done: Boolean(edition.location_name),
    },
  ];

  const actions = [
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
              window.location.href = `/events/${eventSlug}/editions/${edition.slug}`;
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
  const purchaseTimeline: DashboardPurchasePoint[] = purchases
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

      <div
        className="order-1 mt-12 space-y-3 rounded-xl border border-border/60 bg-muted/20 p-3 sm:mt-14"
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
                className="h-9 gap-1.5 px-2 text-[11px] shadow-xs"
                disabled={action.disabled}
                onClick={action.onClick}
                title={`${action.label} · ${action.shortcut}`}
                aria-label={action.label}
              >
                <Icon className="size-4" />
                <span>{action.label.replace(" edição", "")}</span>
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
            label: "Receita aprovada",
            value: new Intl.NumberFormat("pt-BR", {
              style: "currency",
              currency: "BRL",
            }).format(revenue / 100),
            hint: "Compras aprovadas",
            Icon: ShoppingBag,
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

      <DashboardPanel
        title="Status das compras"
        description={`${purchases.length} compra${purchases.length === 1 ? "" : "s"} registrada${purchases.length === 1 ? "" : "s"}.`}
        icon={ShoppingBag}
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
      >
        <DashboardLineChart purchases={purchaseTimeline} />
      </DashboardPanel>

      <StepChecklist
        title="Checklist da edição"
        items={checklist.map((item) => ({
          id: item.label,
          title: item.label,
          description: item.description,
          completed: item.done,
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
