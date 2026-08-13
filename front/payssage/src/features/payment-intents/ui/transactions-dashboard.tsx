import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { SortState } from "@trieoh/ui-base";
import { PaginatedContainer } from "@trieoh/ui-base";
import {
  Ban,
  CheckCircle2,
  CircleGauge,
  Clock3,
  FlaskConical,
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
import type { Intent, IntentStatus } from "../model";
import { intentStatuses } from "../model";

// import { mockIntents } from "../model/mock"

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
}: {
  title?: string;
  description?: string;
  intents: Intent[];
}) {
  const queryClient = useQueryClient();
  const [statusFilter, setStatusFilter] = useState<IntentStatus | "all">("all");
  const [environmentFilter, setEnvironmentFilter] = useState<
    "all" | "production" | "test"
  >("all");
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
  const succeeded = intents.filter((intent) => intent.status === "succeeded");
  const volume = succeeded.reduce(
    (total, intent) => total + intent.amount_cents,
    0,
  );
  const productionVolume = succeeded
    .filter((intent) => !intent.sandbox)
    .reduce((total, intent) => total + intent.amount_cents, 0);
  const testVolume = succeeded
    .filter((intent) => intent.sandbox)
    .reduce((total, intent) => total + intent.amount_cents, 0);
  const successRate = intents.length
    ? Math.round((succeeded.length / intents.length) * 100)
    : 0;
  const environmentIntents = intents.filter((intent) => {
    if (environmentFilter === "production") return !intent.sandbox;
    if (environmentFilter === "test") return intent.sandbox;
    return true;
  });
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
      label: "Captured volume",
      value: formatBRL(volume),
      detail: "All succeeded payments",
      icon: TrendingUp,
    },
    {
      label: "Production volume",
      value: formatBRL(productionVolume),
      detail: "Live succeeded payments",
      icon: CheckCircle2,
    },
    {
      label: "Test volume",
      value: formatBRL(testVolume),
      detail: "Sandbox succeeded payments",
      icon: FlaskConical,
    },
    {
      label: "Success rate",
      value: `${successRate}%`,
      detail: `${succeeded.length} of ${intents.length} intents`,
      icon: CircleGauge,
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
          {(["all", "production", "test"] as const).map((environment) => (
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
              {environment === "all"
                ? "All"
                : environment === "production"
                  ? "Production"
                  : "Test"}
            </Button>
          ))}
        </div>
      </div>

      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        {stats.map((stat) => (
          <div
            key={stat.label}
            className="rounded-lg bg-card p-4 ring-1 ring-foreground/10"
          >
            <div className="flex items-center justify-between text-muted-foreground">
              <span className="text-xs font-medium">{stat.label}</span>
              <stat.icon className="size-4" />
            </div>
            <p className="mt-3 text-2xl font-semibold tracking-tight">
              {stat.value}
            </p>
            <p className="mt-1 text-xs text-muted-foreground">{stat.detail}</p>
          </div>
        ))}
      </section>

      <section className="rounded-lg bg-card p-4 ring-1 ring-foreground/10">
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
        <div className="mt-5 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {intentStatuses.map((status) => {
            const count = environmentIntents.filter(
              (intent) => intent.status === status,
            ).length;
            const percentage = environmentIntents.length
              ? (count / environmentIntents.length) * 100
              : 0;
            const details = statusDetails[status];
            return (
              <div key={status}>
                <div className="flex justify-between text-xs">
                  <span className="text-muted-foreground">{details.label}</span>
                  <span className="font-medium">{count}</span>
                </div>
                <div className="mt-2 h-1.5 overflow-hidden rounded-full bg-muted">
                  <div
                    className={cn("h-full rounded-full", details.dot)}
                    style={{ width: `${percentage}%` }}
                  />
                </div>
              </div>
            );
          })}
        </div>
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
