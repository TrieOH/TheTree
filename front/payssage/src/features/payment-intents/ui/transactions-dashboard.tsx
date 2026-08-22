import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { SortState } from "@trieoh/ui-base";
import {
  DashboardBarList,
  DashboardLineChart,
  DashboardStatCard,
  PaginatedContainer,
} from "@trieoh/ui-base";
import {
  Ban,
  CheckCircle2,
  CircleGauge,
  Clock3,
  FlaskConical,
  Percent,
  ReceiptText,
  TrendingUp,
  XCircle,
} from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { cn } from "#/shared/lib/utils";
import { Badge } from "#/shared/ui/shadcn/badge";
import { Button } from "#/shared/ui/shadcn/button";
import { cancelIntentFn } from "../api";
import type { Intent, IntentFeeDetail, IntentStatus } from "../model";
import { intentStatuses } from "../model";

const statusDetails: Record<
  IntentStatus,
  { label: string; icon: typeof Clock3; className: string; dot: string }
> = {
  processing: {
    label: "Processing",
    icon: Clock3,
    className: "bg-amber-500/10 text-amber-700 dark:text-amber-400",
    dot: "bg-amber-500",
  },
  pending: {
    label: "Pending",
    icon: Clock3,
    className: "bg-amber-500/10 text-amber-700 dark:text-amber-400",
    dot: "bg-amber-500",
  },
  succeeded: {
    label: "Succeeded",
    icon: CheckCircle2,
    className: "bg-emerald-500/10 text-emerald-700 dark:text-emerald-400",
    dot: "bg-emerald-500",
  },
  cancelled: {
    label: "Cancelled",
    icon: Ban,
    className: "bg-muted text-muted-foreground",
    dot: "bg-muted-foreground",
  },
  failed: {
    label: "Failed",
    icon: XCircle,
    className: "bg-destructive/10 text-destructive",
    dot: "bg-destructive",
  },
  rejected: {
    label: "Rejected",
    icon: XCircle,
    className: "bg-destructive/10 text-destructive",
    dot: "bg-destructive",
  },
  refunded: {
    label: "Refunded",
    icon: ReceiptText,
    className: "bg-blue-500/10 text-blue-700 dark:text-blue-400",
    dot: "bg-blue-500",
  },
};

const formatAmount = (intent: Intent) =>
  new Intl.NumberFormat(undefined, {
    style: "currency",
    currency: intent.currency,
  }).format(intent.amount_cents / 100);

const applicationFeeCents = (intent: Intent) => {
  const feeDetails = intent.provider_data.fee_details;
  if (!Array.isArray(feeDetails)) return undefined;

  const applicationFee = (feeDetails as IntentFeeDetail[]).find(
    (fee) => fee.type === "application_fee",
  )?.amount;
  return typeof applicationFee === "number"
    ? Math.round(applicationFee * 100)
    : undefined;
};

const formatDate = (value: string) =>
  new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));

const shortId = (value: string) =>
  value.length > 22 ? `${value.slice(0, 12)}…${value.slice(-6)}` : value;

const formatBRL = (amountCents: number) =>
  new Intl.NumberFormat(undefined, {
    style: "currency",
    currency: "BRL",
  }).format(amountCents / 100);

function IntentCard({
  intent,
  onCancel,
  isCancelling,
}: {
  intent: Intent;
  onCancel: (intent: Intent) => void;
  isCancelling: boolean;
}) {
  const status = statusDetails[intent.status];
  const canCancel =
    intent.status === "pending" || intent.status === "processing";
  return (
    <article className="w-full rounded-lg bg-card p-4 ring-1 ring-foreground/10">
      <div className="flex min-w-0 items-start justify-between gap-3">
        <div className="min-w-0">
          <p className="text-xs capitalize text-muted-foreground">
            {intent.provider.replaceAll("_", " ")}
          </p>
          <p className="mt-1 text-xl font-semibold tracking-tight">
            {formatAmount(intent)}
          </p>
        </div>
        <div className="flex shrink-0 items-center gap-2">
          {intent.sandbox && (
            <Badge variant="secondary" className="gap-1">
              <FlaskConical /> Test
            </Badge>
          )}
          <Badge className={cn("gap-1.5 border-0", status.className)}>
            <span className={cn("size-1.5 rounded-full", status.dot)} />
            {status.label}
          </Badge>
        </div>
      </div>
      <div className="mt-5 space-y-2 text-xs">
        <div className="flex min-w-0 items-center justify-between gap-4">
          <span className="text-muted-foreground">Intent</span>
          <span className="min-w-0 truncate font-mono" title={intent.id}>
            {shortId(intent.id)}
          </span>
        </div>
        <div className="flex min-w-0 items-center justify-between gap-4">
          <span className="text-muted-foreground">Wallet</span>
          <span className="min-w-0 truncate font-mono" title={intent.wallet_id}>
            {shortId(intent.wallet_id)}
          </span>
        </div>
        <div className="flex items-center justify-between gap-4">
          <span className="text-muted-foreground">Created</span>
          <span className="truncate text-right">
            {formatDate(intent.created_at)}
          </span>
        </div>
      </div>
      {canCancel && (
        <div className="mt-4 flex justify-end border-t pt-3">
          <Button
            size="sm"
            variant="destructive"
            disabled={isCancelling}
            onClick={() => onCancel(intent)}
          >
            <Ban /> {isCancelling ? "Cancelling..." : "Cancel intent"}
          </Button>
        </div>
      )}
    </article>
  );
}

export function TransactionsDashboard({
  title = "Transactions",
  description = "Payment activity across your wallets.",
  intents,
  feeBps = 0,
  walletFees,
}: {
  title?: string;
  description?: string;
  intents: Intent[];
  feeBps?: number;
  walletFees?: Record<string, number>;
}) {
  const queryClient = useQueryClient();
  const [statusFilter, setStatusFilter] = useState<IntentStatus | "all">("all");
  const [environmentFilter, setEnvironmentFilter] = useState<
    "production" | "test"
  >("production");
  const [search, setSearch] = useState("");
  const [sort, setSort] = useState<SortState<Intent>>({
    field: "created_at",
    direction: "desc",
  });
  const {
    mutate: cancelIntent,
    isPending: isCancelling,
    variables: cancellingIntent,
  } = useMutation({
    mutationFn: (intent: Intent) => cancelIntentFn(intent.id),
    onSuccess: (response) => {
      if (!response.success) {
        toast.error(response.message || "Failed to cancel intent");
        return;
      }
      void queryClient.invalidateQueries({ queryKey: ["intents"] });
      toast.success("Intent cancelled");
    },
    onError: () => toast.error("Failed to cancel intent"),
  });
  const environmentIntents = intents.filter((intent) =>
    environmentFilter === "production" ? !intent.sandbox : intent.sandbox,
  );
  const captured = environmentIntents.filter(
    (intent) => intent.status === "succeeded" || intent.status === "refunded",
  );
  const refunded = environmentIntents.filter(
    (intent) => intent.status === "refunded",
  );
  const volume = environmentIntents.reduce(
    (total, intent) => total + intent.amount_cents,
    0,
  );
  const capturedVolume = environmentIntents
    .filter((intent) => intent.status === "succeeded")
    .reduce((total, intent) => total + intent.amount_cents, 0);
  const revenue = captured.reduce(
    (total, intent) =>
      total +
      (applicationFeeCents(intent) ??
        Math.round(
          (intent.amount_cents * (walletFees?.[intent.wallet_id] ?? feeBps)) /
            10_000,
        )),
    0,
  );
  const revenueFeeDetail = "Based on each transaction's fee";
  const successRate = environmentIntents.length
    ? Math.round((captured.length / environmentIntents.length) * 100)
    : 0;
  const refundedRate = captured.length
    ? Math.round((refunded.length / captured.length) * 100)
    : 0;
  const filtered = environmentIntents.filter((intent) => {
    const matchesStatus =
      statusFilter === "all" || intent.status === statusFilter;
    const term = search.trim().toLowerCase();
    const matchesSearch =
      !term ||
      [intent.id, intent.wallet_id, intent.provider, intent.status].some(
        (value) => value.toLowerCase().includes(term),
      );
    return matchesStatus && matchesSearch;
  });

  const stats = [
    {
      label: "Total volume",
      value: formatBRL(volume),
      detail: "All payment intents",
      icon: TrendingUp,
    },
    {
      label: "Captured volume",
      value: formatBRL(capturedVolume),
      detail: "Succeeded payments",
      icon: CheckCircle2,
    },
    {
      label: "Revenue",
      value: formatBRL(revenue),
      detail: revenueFeeDetail,
      icon: TrendingUp,
    },
    {
      label: "Success rate",
      value: `${successRate}%`,
      detail: `${captured.length} of ${environmentIntents.length} intents`,
      icon: CircleGauge,
    },
    {
      label: "Refunded",
      value: `${refundedRate}%`,
      detail: `${refunded.length} of ${captured.length} captured intents`,
      icon: Percent,
    },
  ];

  return (
    <div className="min-w-0 space-y-6">
      <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-start">
        <div>
          <p className="text-xs font-medium uppercase tracking-widest text-muted-foreground">
            Overview
          </p>
          <h1 className="mt-1 text-2xl font-semibold tracking-tight">
            {title}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">{description}</p>
        </div>
        <div
          className="flex w-fit rounded-md bg-muted p-1"
          role="radiogroup"
          aria-label="Transaction environment"
        >
          {(["production", "test"] as const).map((environment) => (
            <Button
              key={environment}
              size="sm"
              variant="ghost"
              className={cn(
                "h-8 px-3 text-xs",
                environmentFilter === environment &&
                  "bg-background text-foreground shadow-xs hover:bg-background",
              )}
              onClick={() => setEnvironmentFilter(environment)}
            >
              {environment === "production" ? "Production" : "Test"}
            </Button>
          ))}
        </div>
      </div>

      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-5">
        {stats.map((stat) => (
          <DashboardStatCard
            key={stat.label}
            label={stat.label}
            value={stat.value}
            hint={stat.detail}
            icon={stat.icon}
          />
        ))}
      </section>

      <section className="grid gap-4 xl:grid-cols-[1fr_1.8fr]">
        <section className="order-2 rounded-lg bg-card p-4 ring-1 ring-foreground/10 xl:order-1">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 className="font-semibold">Status overview</h2>
              <p className="text-xs text-muted-foreground">
                Current intent distribution
              </p>
            </div>
            <span className="text-xs text-muted-foreground">
              {environmentIntents.length} total
            </span>
          </div>
          <div className="mt-5">
            <DashboardBarList
              items={intentStatuses.map((status) => ({
                id: status,
                label: statusDetails[status].label,
                value: environmentIntents.filter(
                  (intent) => intent.status === status,
                ).length,
                color: statusDetails[status].dot,
              }))}
              total={environmentIntents.length}
              showPercentage
              emptyMessage="No intents found."
            />
          </div>
        </section>

        <section className="order-1 rounded-lg bg-card p-4 ring-1 ring-foreground/10 xl:order-2">
          <div className="mb-4">
            <h2 className="font-semibold">Payment volume</h2>
            <p className="text-xs text-muted-foreground">
              Accumulated payment volume over time.
            </p>
          </div>
          <DashboardLineChart
            points={environmentIntents.map((intent) => ({
              timestamp: intent.created_at,
              status: intent.status,
              totalCents: intent.amount_cents,
            }))}
          />
        </section>
      </section>

      <section className="space-y-3">
        <div>
          <h2 className="font-semibold">Recent transactions</h2>
          <p className="text-xs text-muted-foreground">
            Latest payment activity
          </p>
        </div>
        <PaginatedContainer<Intent>
          items={filtered}
          pageSize={6}
          layout="grid"
          gap="3"
          minItemWidth="16rem"
          sortFields={[
            { key: "created_at", label: "Created" },
            { key: "amount_cents", label: "Amount" },
            { key: "status", label: "Status" },
            { key: "provider", label: "Provider" },
          ]}
          sort={sort}
          onSortChange={setSort}
          filterValue={search}
          onFilterChange={setSearch}
          filterPlaceholder="Search intents..."
          itemLabel="transactions"
          headerActions={
            <div className="flex flex-wrap gap-1.5">
              <Button
                size="sm"
                variant={statusFilter === "all" ? "default" : "ghost"}
                onClick={() => setStatusFilter("all")}
              >
                All statuses
              </Button>
              {intentStatuses.map((status) => (
                <Button
                  key={status}
                  size="sm"
                  variant={statusFilter === status ? "default" : "ghost"}
                  onClick={() => setStatusFilter(status)}
                >
                  {statusDetails[status].label}
                </Button>
              ))}
            </div>
          }
          emptyState={
            <div className="py-12 text-center">
              <ReceiptText className="mx-auto size-6 text-muted-foreground" />
              <p className="mt-3 text-sm font-medium">No transactions found</p>
              <p className="mt-1 text-xs text-muted-foreground">
                Try another status or search term.
              </p>
            </div>
          }
          renderItems={(items) =>
            items.map((intent) => (
              <IntentCard
                key={intent.id}
                intent={intent}
                onCancel={cancelIntent}
                isCancelling={
                  isCancelling && cancellingIntent?.id === intent.id
                }
              />
            ))
          }
        />
      </section>
    </div>
  );
}
