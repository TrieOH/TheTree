import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AlertCircle, ArrowLeft, Plus, Webhook } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { cn } from "#/shared/lib/utils";
import { Badge } from "#/shared/ui/shadcn/badge";
import { Button } from "#/shared/ui/shadcn/button";
import { Input } from "#/shared/ui/shadcn/input";
import { Label } from "#/shared/ui/shadcn/label";
import {
  createWebhookEndpointFn,
  deleteWebhookEndpointFn,
  webhookDeliveriesQueryOptions,
  webhookEndpointsQueryOptions,
  webhookEventsQueryOptions,
} from "../api";
import type {
  WebhookDelivery,
  WebhookEndpoint,
  WebhookEndpointCreateRequest,
} from "../model";
import { webhookEndpointCreateSchema } from "../model";
import { WebhookCreatedModal } from "./webhook-created-modal";
import { WebhookList } from "./webhook-list";

const formatDate = (value: string | null) =>
  value
    ? new Intl.DateTimeFormat(undefined, {
        dateStyle: "medium",
        timeStyle: "short",
      }).format(new Date(value))
    : "Never";

export function WebhooksDashboard({ walletId }: { walletId: string }) {
  const queryClient = useQueryClient();
  const [view, setView] = useState<"endpoints" | "events">("endpoints");
  const [selectedEndpoint, setSelectedEndpoint] =
    useState<WebhookEndpoint | null>(null);
  const [createdEndpoint, setCreatedEndpoint] =
    useState<WebhookEndpoint | null>(null);
  const endpointsOptions = webhookEndpointsQueryOptions(walletId);
  const eventsOptions = webhookEventsQueryOptions(walletId);
  const { data: endpoints = [], isLoading } = useQuery(endpointsOptions);
  const { data: events = [] } = useQuery(eventsOptions);
  const deliveriesOptions = webhookDeliveriesQueryOptions(
    selectedEndpoint?.id ?? "",
  );
  const { data: deliveries = [], isLoading: isLoadingDeliveries } = useQuery({
    ...deliveriesOptions,
    enabled: Boolean(selectedEndpoint),
  });
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<WebhookEndpointCreateRequest>({
    resolver: zodResolver(webhookEndpointCreateSchema),
    defaultValues: { name: "", url: "" },
  });

  const { mutate: createEndpoint, isPending: isCreating } = useMutation({
    mutationFn: (payload: WebhookEndpointCreateRequest) =>
      createWebhookEndpointFn(walletId, payload),
    onSuccess: (response) => {
      if (!response.success)
        return toast.error(response.message || "Failed to create endpoint");
      void queryClient.invalidateQueries({
        queryKey: endpointsOptions.queryKey,
      });
      setCreatedEndpoint(response.data);
      reset();
    },
    onError: () => toast.error("Failed to create endpoint"),
  });
  const { mutate: deleteEndpoint, isPending: isDeleting } = useMutation({
    mutationFn: (endpointId: string) => deleteWebhookEndpointFn(endpointId),
    onSuccess: (response) => {
      if (!response.success)
        return toast.error(response.message || "Failed to delete endpoint");
      void queryClient.invalidateQueries({
        queryKey: endpointsOptions.queryKey,
      });
      toast.success("Webhook endpoint deleted");
    },
    onError: () => toast.error("Failed to delete endpoint"),
  });

  if (selectedEndpoint) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => setSelectedEndpoint(null)}>
          <ArrowLeft /> Back to endpoints
        </Button>
        <div>
          <h2 className="font-semibold">{selectedEndpoint.name}</h2>
          <p className="text-xs text-muted-foreground">
            Deliveries to {selectedEndpoint.url}
          </p>
        </div>
        <DeliveryList deliveries={deliveries} isLoading={isLoadingDeliveries} />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex gap-1 rounded-md bg-muted p-1 w-fit">
        {(["endpoints", "events"] as const).map((item) => (
          <Button
            key={item}
            size="sm"
            variant="ghost"
            className={cn(
              "capitalize",
              view === item && "bg-background shadow-xs",
            )}
            onClick={() => setView(item)}
          >
            {item}
          </Button>
        ))}
      </div>

      {view === "endpoints" ? (
        <>
          <section className="rounded-sm border bg-card p-4">
            <div className="mb-4">
              <h2 className="text-sm font-semibold">Create endpoint</h2>
              <p className="text-xs text-muted-foreground">
                Register a URL for this wallet's webhook deliveries.
              </p>
            </div>
            <form
              className="grid gap-4 md:grid-cols-[1fr_2fr_auto] md:items-start"
              onSubmit={handleSubmit((data) => createEndpoint(data))}
            >
              <FormField label="Name" error={errors.name?.message}>
                <Input placeholder="Payments" {...register("name")} />
              </FormField>
              <FormField label="URL" error={errors.url?.message}>
                <Input
                  type="url"
                  placeholder="https://example.com/webhooks"
                  {...register("url")}
                />
              </FormField>
              <Button type="submit" disabled={isCreating} className="mt-6">
                <Plus /> {isCreating ? "Creating..." : "Create"}
              </Button>
            </form>
          </section>
          <WebhookList
            webhooks={endpoints}
            isLoading={isLoading}
            onDelete={deleteEndpoint}
            onViewDeliveries={setSelectedEndpoint}
          />
          {isDeleting && (
            <p className="text-xs text-muted-foreground">
              Deleting endpoint...
            </p>
          )}
        </>
      ) : (
        <div className="space-y-3">
          {events.map((event) => (
            <article key={event.id} className="rounded-sm border bg-card p-4">
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div>
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary">{event.event_type}</Badge>
                    <span className="text-xs capitalize text-muted-foreground">
                      {event.provider}
                    </span>
                  </div>
                  <p className="mt-2 font-mono text-xs">
                    Intent: {event.intent_id}
                  </p>
                  <p className="mt-1 font-mono text-xs text-muted-foreground">
                    External: {event.external_id}
                  </p>
                </div>
                <span className="text-xs text-muted-foreground">
                  {formatDate(event.received_at)}
                </span>
              </div>
              <pre className="mt-4 max-h-72 overflow-auto rounded-sm bg-muted p-3 text-xs">
                {JSON.stringify(event.payload, null, 2)}
              </pre>
            </article>
          ))}
          {!events.length && <EmptyState label="No webhook events found" />}
        </div>
      )}

      <WebhookCreatedModal
        webhook={createdEndpoint}
        isOpen={Boolean(createdEndpoint)}
        onClose={() => setCreatedEndpoint(null)}
      />
    </div>
  );
}

function FormField({
  label,
  error,
  children,
}: {
  label: string;
  error?: string;
  children: React.ReactNode;
}) {
  return (
    <div className="space-y-2">
      <Label>{label}</Label>
      {children}
      {error && (
        <p className="flex items-center gap-1 text-xs text-destructive">
          <AlertCircle className="size-3" />
          {error}
        </p>
      )}
    </div>
  );
}

function DeliveryList({
  deliveries,
  isLoading,
}: {
  deliveries: WebhookDelivery[];
  isLoading: boolean;
}) {
  if (isLoading)
    return (
      <p className="text-sm text-muted-foreground">Loading deliveries...</p>
    );
  if (!deliveries.length) return <EmptyState label="No deliveries found" />;
  return (
    <div className="space-y-3">
      {deliveries.map((delivery) => (
        <article key={delivery.id} className="rounded-sm border bg-card p-4">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div className="flex items-center gap-2">
              <DeliveryStatus status={delivery.status} />
              <span className="font-mono text-xs">{delivery.id}</span>
            </div>
            <span className="text-xs text-muted-foreground">
              {formatDate(delivery.created_at)}
            </span>
          </div>
          <dl className="mt-4 grid gap-3 text-xs sm:grid-cols-3">
            <div>
              <dt className="text-muted-foreground">Attempts</dt>
              <dd>{delivery.attempts}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Last attempt</dt>
              <dd>{formatDate(delivery.last_attempted_at)}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">HTTP status</dt>
              <dd>{delivery.response_status ?? "No response"}</dd>
            </div>
          </dl>
          {delivery.response_body && (
            <pre className="mt-4 max-h-52 overflow-auto rounded-sm bg-muted p-3 text-xs">
              {delivery.response_body}
            </pre>
          )}
        </article>
      ))}
    </div>
  );
}

function DeliveryStatus({ status }: { status: WebhookDelivery["status"] }) {
  return (
    <Badge
      className={cn(
        status === "delivered" && "bg-emerald-600",
        status === "failed" && "bg-destructive",
      )}
    >
      {status}
    </Badge>
  );
}

function EmptyState({ label }: { label: string }) {
  return (
    <div className="flex flex-col items-center justify-center border border-dashed py-16 text-muted-foreground">
      <Webhook className="size-7" />
      <p className="mt-2 text-sm">{label}</p>
    </div>
  );
}
