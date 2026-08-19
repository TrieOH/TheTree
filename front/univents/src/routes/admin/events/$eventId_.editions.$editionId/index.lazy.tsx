import { useQuery, useQueryClient } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import type { SortState } from "@trieoh/ui-base";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import { format } from "date-fns";
import { ptBR } from "date-fns/locale";
import {
  CalendarDays,
  CheckCircle2,
  ChevronRight,
  CircleAlert,
  Globe,
  LayoutGrid,
  ShoppingBag,
  Ticket,
  Users,
} from "lucide-react";
import { motion } from "motion/react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import { usePublishEditionMutation } from "@/features/editions/api/mutations";
import type { EditionI } from "@/features/editions/model";
import { EditionImageActions } from "@/features/editions/ui/EditionImageActions";
import {
  editionPurchasesQueryOptions,
  useRefundPurchaseMutation,
} from "@/features/purchases/api";
import { purchaseQueryKeys } from "@/features/purchases/api/query-keys";
import { formatDateRange } from "@/shared/lib/date";
import { cn } from "@/shared/lib/utils";
import { Badge } from "@/shared/ui/shadcn/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/shadcn/card";
import { AlertModal } from "@/widgets/ui/alert-modal";
import { QuickAction } from "@/widgets/ui/quick-action";

export const Route = createLazyFileRoute(
  "/admin/events/$eventId_/editions/$editionId/",
)({
  component: AdminEditionDetailRoute,
});

const statusConfig: Record<
  EditionI["status"],
  { label: string; className: string; dot: string }
> = {
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
};

function AdminEditionDetailRoute() {
  const { eventId, editionId } = Route.useParams();
  const { data: editions = [], isPending } = useQuery(
    allAdminEditionsQueryOptions(eventId),
  );
  const queryClient = useQueryClient();
  const { data: purchases = [] } = useQuery(
    editionPurchasesQueryOptions(editionId),
  );
  const [publishConfirmOpen, setPublishConfirmOpen] = useState(false);
  const [refundPurchaseId, setRefundPurchaseId] = useState<string | null>(null);
  const [purchaseFilter, setPurchaseFilter] = useState("");
  const [purchaseSort, setPurchaseSort] = useState<SortState<EditionPurchase>>({
    field: "created_at",
    direction: "desc",
  });
  const edition = useMemo(
    () => editions.find((item) => item.id === editionId) ?? null,
    [editionId, editions],
  );

  const publishEditionMutation = usePublishEditionMutation();
  const refundMutation = useRefundPurchaseMutation();
  const filteredPurchases = useMemo(() => {
    const search = purchaseFilter.trim().toLowerCase();
    return purchases
      .filter((purchase) => {
        if (!search) return true;
        return [
          purchase.payer_email ?? "",
          purchase.status,
          purchase.purchase_id,
          ...purchase.attendees.map((attendee) => attendee.email),
        ].some((value) => value.toLowerCase().includes(search));
      })
      .sort((a, b) => {
        const direction = purchaseSort.direction === "asc" ? 1 : -1;
        if (purchaseSort.field === "total_cents") {
          return (a.total_cents - b.total_cents) * direction;
        }
        return (
          (new Date(a.created_at ?? 0).getTime() -
            new Date(b.created_at ?? 0).getTime()) *
          direction
        );
      });
  }, [purchaseFilter, purchaseSort, purchases]);

  const copyLink = () => {
    if (!edition) return;
    void navigator.clipboard.writeText(
      `${window.location.origin}/events/${edition.event_id}/editions/${edition.slug}`,
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

  const status = statusConfig[edition.status];
  const heroTagline =
    edition.tagline ??
    edition.description ??
    "Sem descrição cadastrada para esta edição.";
  const isDraft = edition.status === "draft";

  const metrics = [
    {
      label: "Período",
      value: formatDateRange(edition.starts_at, edition.ends_at),
      hint: "Datas de início e fim da edição",
    },
    {
      label: "Local",
      value: edition.location_name ?? "—",
      hint: edition.location_description ?? "Local físico da edição",
    },
    {
      label: "Inscrições abertas em",
      value: edition.registration_opens_at
        ? format(
            new Date(edition.registration_opens_at),
            "dd 'de' MMM 'de' yyyy",
            { locale: ptBR },
          )
        : "Não definido",
      hint: edition.registration_opens_at
        ? `Abertura em ${format(new Date(edition.registration_opens_at), "HH:mm", { locale: ptBR })}`
        : "Data de abertura das inscrições",
    },
  ];

  const checklist = [
    {
      label: "Banner cadastrado",
      done: Boolean(edition.banner_url),
    },
    {
      label: "Logo cadastrado",
      done: Boolean(edition.logo_url),
    },
    {
      label: "Descrição preenchida",
      done: Boolean(edition.description),
    },
    {
      label: "Tagline definida",
      done: Boolean(edition.tagline),
    },
    {
      label: "Local definido",
      done: Boolean(edition.location_name),
    },
  ];

  const actions = [
    ...(isDraft
      ? [
          {
            label: "Publicar edição",
            onClick: () => setPublishConfirmOpen(true),
            disabled: publishEditionMutation.isPending,
            variant: "default" as const,
          },
        ]
      : []),
    {
      label: "Copiar link público",
      onClick: copyLink,
      disabled: isDraft,
      variant: "default" as const,
    },
    {
      label: "E-mail de contato",
      onClick: () => {
        if (edition.contact_email) {
          void navigator.clipboard.writeText(edition.contact_email);
          toast.success("E-mail copiado");
        }
      },
      disabled: !edition.contact_email,
      variant: "default" as const,
    },
  ];
  const approvedPurchases = purchases.filter(
    (purchase) => purchase.status === "approved",
  ).length;
  const refundedPurchases = purchases.filter(
    (purchase) => purchase.status === "refunded",
  ).length;
  const revenue = purchases
    .filter((purchase) => purchase.status === "approved")
    .reduce((total, purchase) => total + purchase.total_cents, 0);

  return (
    <div className="relative space-y-6 p-6 pb-28!">
      <motion.section
        initial={{ opacity: 0, y: 12 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.35 }}
        className="overflow-hidden rounded-md border border-border/60 bg-card shadow-[0_1px_0_0_rgba(255,255,255,0.03),0_20px_40px_-24px_rgba(15,23,42,0.24)]"
      >
        {edition.banner_url ? (
          <div className="h-48 overflow-hidden bg-muted">
            <img
              src={edition.banner_url}
              alt={edition.name}
              className="h-full w-full object-cover"
            />
          </div>
        ) : null}

        <div className="relative flex flex-col gap-6 p-6">
          <div className="pointer-events-none absolute inset-x-0 top-0 h-px bg-linear-to-r from-transparent via-primary/20 to-transparent" />
          <div className="pointer-events-none absolute inset-x-0 bottom-0 h-24 bg-linear-to-t from-primary/5 via-transparent to-transparent" />

          <div className="space-y-5">
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="rounded-full px-3">
                <LayoutGrid className="size-3.5" />
                Overview
              </Badge>
            </div>

            <div className="flex flex-wrap items-center gap-2">
              <span className="text-xs text-muted-foreground">
                Imagens da edição
              </span>
              <EditionImageActions
                edition={edition}
                eventId={eventId}
                field="logo_url"
              />
              <EditionImageActions
                edition={edition}
                eventId={eventId}
                field="banner_url"
              />
            </div>

            <div className="space-y-3">
              <h1 className="max-w-3xl text-3xl font-semibold tracking-tight text-foreground md:text-4xl">
                {edition.name}
              </h1>
              <p className="max-w-2xl text-sm leading-6 text-muted-foreground md:text-base">
                {heroTagline}
              </p>
            </div>

            <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
              <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2.5 py-1">
                <CalendarDays className="size-3.5" />
                {edition.slug}
              </span>
              <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2.5 py-1">
                <Globe className="size-3.5" />
                Link público
              </span>
              <span
                className={cn(
                  "inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1",
                  status.className,
                )}
              >
                <span className={cn("size-1.5 rounded-full", status.dot)} />
                {status.label}
              </span>
            </div>
          </div>
        </div>
      </motion.section>

      <section className="grid gap-4 md:grid-cols-3">
        {metrics.map((metric) => (
          <Card
            key={metric.label}
            className="border-border/60 bg-card/95 shadow-sm transition-shadow hover:shadow-md"
          >
            <CardHeader className="pb-2">
              <CardDescription className="flex items-center gap-2 text-xs uppercase tracking-[0.22em]">
                <span className="size-1.5 rounded-full bg-primary/60" />
                {metric.label}
              </CardDescription>
              <CardTitle className="text-lg font-semibold tracking-tight">
                {metric.value}
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-0">
              <p className="text-xs text-muted-foreground">{metric.hint}</p>
            </CardContent>
          </Card>
        ))}
      </section>

      <section className="grid gap-3 sm:grid-cols-3">
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
            hint: `${refundedPurchases} reembolsada(s)`,
            Icon: Ticket,
          },
        ].map(({ label, value, hint, Icon }) => (
          <Card
            key={String(label)}
            className="border-border/60 bg-card/95 shadow-sm"
          >
            <CardContent className="flex items-center gap-3 p-4">
              <div className="flex size-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
                <Icon className="size-5" />
              </div>
              <div className="min-w-0">
                <p className="text-xs text-muted-foreground">{label}</p>
                <p className="truncate text-xl font-semibold tracking-tight">
                  {value}
                </p>
                <p className="truncate text-xs text-muted-foreground">{hint}</p>
              </div>
            </CardContent>
          </Card>
        ))}
      </section>

      <PaginatedContainer<EditionPurchase>
        items={filteredPurchases}
        layout="list"
        pageSize={6}
        gap="3"
        sort={purchaseSort}
        onSortChange={setPurchaseSort}
        sortFields={[
          { key: "created_at", label: "Data" },
          {
            key: "total_cents",
            label: "Valor",
            comparator: (a, b) => a.total_cents - b.total_cents,
          },
        ]}
        filterValue={purchaseFilter}
        onFilterChange={setPurchaseFilter}
        filterPlaceholder="Buscar por e-mail, status ou participante..."
        itemLabel="compras"
        emptyState={
          <EmptyState
            icon={ShoppingBag}
            eyebrow="Compras"
            title="Nenhuma compra encontrada"
            description="Os pedidos desta edição aparecerão aqui."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((purchase) => {
            const itemQuantity = purchase.items.reduce(
              (total, item) => total + item.quantity,
              0,
            );
            const contactEmail =
              purchase.payer_email ?? purchase.attendees[0]?.email ?? null;
            const contactLabel = purchase.payer_email
              ? "Pagador"
              : purchase.attendees[0]?.email
                ? "Contato do participante"
                : "Pagador";
            const displayEmail = contactEmail ?? "E-mail não informado";

            return (
              <Card
                key={purchase.purchase_id}
                className="border-border/60 bg-card/95 shadow-sm"
              >
                <CardContent className="flex flex-wrap items-center justify-between gap-4 p-4 sm:p-5">
                  <div className="flex min-w-0 items-start gap-3">
                    <div className="flex size-10 shrink-0 items-center justify-center rounded-full bg-muted text-sm font-semibold text-muted-foreground">
                      {displayEmail[0].toUpperCase()}
                    </div>
                    <div className="min-w-0 space-y-1">
                      <p className="text-[11px] uppercase tracking-wide text-muted-foreground">
                        {contactLabel}
                      </p>
                      <p className="truncate text-sm font-semibold">
                        {displayEmail}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {itemQuantity} {itemQuantity === 1 ? "item" : "itens"}
                        {" · "}
                        {purchase.attendees.length}{" "}
                        {purchase.attendees.length === 1
                          ? "participante atribuído"
                          : "participantes atribuídos"}
                      </p>
                      {purchase.attendees.length > 0 ? (
                        <p className="max-w-xl truncate text-xs text-foreground/70">
                          {purchase.attendees
                            .map((attendee) => attendee.name)
                            .join(", ")}
                        </p>
                      ) : (
                        <p className="text-xs italic text-muted-foreground">
                          Participantes ainda não atribuídos
                        </p>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <div className="text-right">
                      <p className="text-sm font-semibold tabular-nums">
                        {new Intl.NumberFormat("pt-BR", {
                          style: "currency",
                          currency: purchase.currency,
                        }).format(purchase.total_cents / 100)}
                      </p>
                      <p className="text-xs capitalize text-muted-foreground">
                        {purchase.status}
                      </p>
                    </div>
                    {purchase.status === "approved" ? (
                      <button
                        type="button"
                        className="text-sm font-medium text-destructive hover:underline"
                        onClick={() =>
                          setRefundPurchaseId(purchase.purchase_id)
                        }
                      >
                        Reembolsar
                      </button>
                    ) : null}
                  </div>
                </CardContent>
              </Card>
            );
          })
        }
      />

      <Card className="border-border/60 bg-card/95 shadow-sm">
        <CardHeader className="border-b border-border/60">
          <CardTitle>Checklist da edição</CardTitle>
          <CardDescription>
            Itens derivados dos dados já cadastrados na edição.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-3 py-5">
          {checklist.map((item) => (
            <div
              key={item.label}
              className="flex items-center justify-between rounded-2xl border border-border/60 bg-muted/15 px-4 py-3.5"
            >
              <div className="flex items-center gap-3">
                <div
                  className={cn(
                    "size-2 rounded-full",
                    item.done ? "bg-emerald-500" : "bg-amber-500",
                  )}
                />
                <span className="text-sm text-foreground">{item.label}</span>
              </div>
              {item.done ? (
                <CheckCircle2 className="size-4 text-emerald-500/70" />
              ) : (
                <CircleAlert className="size-4 text-amber-500/70" />
              )}
            </div>
          ))}
        </CardContent>
      </Card>

      <Card className="border-border/60 bg-card/95 shadow-sm">
        <CardHeader className="border-b border-border/60">
          <CardTitle>Ações rápidas</CardTitle>
          <CardDescription>
            Atalhos para as operações mais comuns da edição.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-3 py-5 sm:grid-cols-2 xl:grid-cols-3">
          {actions.map((action) => (
            <QuickAction
              key={action.label}
              onClick={action.onClick}
              disabled={action.disabled}
              variant={action.variant}
            >
              <span className="text-sm font-medium text-foreground">
                {action.label}
              </span>
              <ChevronRight className="size-4 text-muted-foreground" />
            </QuickAction>
          ))}
        </CardContent>
      </Card>

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
