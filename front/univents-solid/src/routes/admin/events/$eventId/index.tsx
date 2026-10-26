import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useQuery, type QueryClient } from "@trieoh/front-core-solid";
import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import { createMemo, createSignal } from "solid-js";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import { eventKeys } from "@/features/events/api/query-keys";
import {
  allJoinedEventsQueryOptions,
  allOwnEventsQueryOptions,
} from "@/features/events/api";
import {
  useDiscontinueEventMutation,
  usePatchEventMutation,
  usePublishEventMutation,
} from "@/features/events/api/mutations";
import {
  type EventOverviewMetrics,
  buildEventOverviewMetrics,
} from "@/features/events/model/event-overview";
import { EventEditionsList } from "@/features/events/ui/EventEditionsList";
import { EventOverviewChecklist } from "@/features/events/ui/EventOverviewChecklist";
import { EventOverviewDashboard } from "@/features/events/ui/EventOverviewDashboard";
import { EventOverviewDialogs } from "@/features/events/ui/EventOverviewDialogs";
import { EventOverviewHeader } from "@/features/events/ui/EventOverviewHeader";
import {
  type EventQuickAction,
  EventQuickActions,
} from "@/features/events/ui/EventQuickActions";
import type { ManageEventValues } from "@/features/events/ui/ManageEventDialog";
import type { PaymentProviderI } from "@/features/payments/api";
import {
  useConnectEventSellerMutation,
  useDisconnectEventSellerMutation,
} from "@/features/payments/api/mutations";
import { EventPaymentPanel } from "@/features/payments/ui/EventPaymentPanel";
import { productsQueryOptions } from "@/features/products/api";
import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import { editionPurchasesQueryOptions } from "@/features/purchases/api";
import {
  allTicketsQueryOptions,
  attendeeCountQueryOptions,
} from "@/features/tickets/api";
import { createHotkeys } from "@/shared/lib/hotkeys";
import { useUploadQueue } from "@/features/upload-queue";
import { handleShare } from "@/shared/lib/share";
import { appQueryClient } from "@/shared/lib/query-client";
import { GridSkeleton } from "@/shared/ui/GridSkeleton";
import { QueryBoundary, combineQueries } from "@/shared/ui/QueryBoundary";
import { warmQuery } from "@/shared/lib/route-data";

export const Route = createFileRoute("/admin/events/$eventId/")({
  loader: ({ params }) => warmQuery(appQueryClient, allAdminEditionsQueryOptions(params.eventId)),
  head: ({ params }) => ({
    meta: [{ title: `Evento ${params.eventId} - Admin - Univents` }],
  }),
  component: AdminEventOverviewRoute,
});

function AdminEventOverviewRoute(): JSX.Element {
  const params = Route.useParams();
  const navigate = useNavigate();
  const uploadQueue = useUploadQueue();

  const eventId = () => params().eventId;

  const eventsQuery = combineQueries(
    useQuery(() => allOwnEventsQueryOptions()),
    useQuery(() => allJoinedEventsQueryOptions()),
  );
  const editionsQuery = useQuery(() => allAdminEditionsQueryOptions(eventId()));

  const [publishConfirmOpen, setPublishConfirmOpen] = createSignal(false);
  const [discontinueConfirmOpen, setDiscontinueConfirmOpen] =
    createSignal(false);
  const [disconnectSellerConfirmOpen, setDisconnectSellerConfirmOpen] =
    createSignal(false);
  const [editEventOpen, setEditEventOpen] = createSignal(false);

  const event = createMemo(() => {
    const [owned = [], joined = []] = eventsQuery().data ?? [];
    return [...owned, ...joined].find((item) => item.id === eventId()) ?? null;
  });

  const editions = createMemo(() => editionsQuery().data ?? []);

  // Totals across every edition. Each edition's data goes through the query
  // cache, so this mostly reuses what other pages already loaded. A failed
  // request fails the whole thing (instead of silently counting as zero), and
  // it always recomputes on mount since the pieces it reads may have changed.
  const metricsQuery = useQuery(() => ({
    queryKey: [...eventKeys.all, eventId(), "overview-metrics", editions().map((edition) => edition.id)] as const,
    enabled: editionsQuery().data !== undefined,
    staleTime: 0,
    queryFn: async ({ client }: { client: QueryClient }): Promise<EventOverviewMetrics> => {
      const current = editions();
      const results = await Promise.all(
        current.map((edition) =>
          Promise.all([
            client.fetchQuery(editionPurchasesQueryOptions(edition.id)),
            client.fetchQuery(attendeeCountQueryOptions(edition.id)),
            client.fetchQuery(allTicketsQueryOptions(edition.id)),
            client.fetchQuery(productsQueryOptions(edition.id)),
            client.fetchQuery(programsQueryOptions(edition.id)),
            client.fetchQuery(occurrencesQueryOptions(edition.id)),
          ]),
        ),
      );

      return buildEventOverviewMetrics({
        editions: current,
        purchasesByEdition: results.map((r) => r[0] as EditionPurchase[]),
        attendeeCounts: results.map((r) => (r[1] as { count?: number })?.count ?? 0),
        ticketCounts: results.map((r) => r[2].length),
        productCounts: results.map((r) => r[3].length),
        programCounts: results.map((r) => r[4].length),
        occurrenceCounts: results.map((r) => r[5].length),
      });
    },
  }));

  const isImageUploading = (field: "logo_url" | "banner_url") =>
    uploadQueue.tasks.some(
      (task) =>
        task.owner.type === "event" &&
        task.owner.id === eventId() &&
        task.association?.handlerKey === "event-image" &&
        task.association.input?.field === field &&
        !["completed", "failed", "rejected"].includes(task.status),
    );

  const publishMutation = usePublishEventMutation();
  const discontinueMutation = useDiscontinueEventMutation();
  const patchMutation = usePatchEventMutation();
  const connectSellerMutation = useConnectEventSellerMutation();
  const disconnectSellerMutation = useDisconnectEventSellerMutation();

  const isPublishing = () => publishMutation.result().status === "pending";
  const isDiscontinuing = () =>
    discontinueMutation.result().status === "pending";
  const isConnectingSeller = () =>
    connectSellerMutation.result().status === "pending";
  const isDisconnectingSeller = () =>
    disconnectSellerMutation.result().status === "pending";

  const isPublished = () => event()?.status === "active";

  const copyLink = () => {
    const ev = event();
    if (!ev) return;
    void handleShare(
      ev.full_name,
      `${window.location.origin}/events/${ev.slug}`,
    );
  };

  const handleEditEvent = async (values: ManageEventValues) => {
    const ev = event();
    if (!ev) return false;
    try {
      await patchMutation.mutateAsync({
        eventId: ev.id,
        data: values,
      });
      setEditEventOpen(false);
      return true;
    } catch {
      return false;
    }
  };

  const handlePublish = async () => {
    try {
      await publishMutation.mutateAsync(eventId());
      setPublishConfirmOpen(false);
    } catch {
      // Handled
    }
  };

  const handleDiscontinue = async () => {
    try {
      await discontinueMutation.mutateAsync(eventId());
      setDiscontinueConfirmOpen(false);
    } catch {
      // Handled
    }
  };

  const handleDisconnectSeller = async () => {
    try {
      await disconnectSellerMutation.mutateAsync(eventId());
      setDisconnectSellerConfirmOpen(false);
    } catch {
      // Handled
    }
  };

  const handleConnectSeller = (provider: PaymentProviderI) => {
    void connectSellerMutation.mutateAsync({ eventId: eventId(), provider });
  };

  const actions = (): EventQuickAction[] => {
    const ev = event();
    return [
      {
        label: "Editar evento",
        shortcut: "Mod+E",
        onClick: () => setEditEventOpen(true),
        disabled: !ev,
        variant: "default",
      },
      ...(ev?.status === "draft"
        ? [
          {
            label: "Publicar evento",
            shortcut: "Mod+P",
            onClick: () => setPublishConfirmOpen(true),
            disabled: isPublishing(),
            variant: "default" as const,
          },
        ]
        : []),
      {
        label: "Copiar link público",
        shortcut: "Mod+Shift+C",
        onClick: copyLink,
        disabled: !ev,
        variant: "default",
      },
      ...(isPublished()
        ? [
          {
            label: "Descontinuar evento",
            shortcut: "Mod+Shift+D",
            onClick: () => setDiscontinueConfirmOpen(true),
            disabled: isDiscontinuing(),
            variant: "destructive" as const,
          },
          {
            label: "Abrir painel público",
            shortcut: "Mod+Shift+O",
            to: "/events/$slug" as const,
            params: { slug: ev?.slug ?? "" },
            variant: "default" as const,
          },
        ]
        : []),
    ];
  };

  // Keyboard shortcuts via TanStack Hotkeys
  createHotkeys(
    () => [
      {
        hotkey: "Mod+E",
        callback: () => setEditEventOpen(true),
        options: { enabled: Boolean(event()) },
      },
      {
        hotkey: "Mod+P",
        callback: () => setPublishConfirmOpen(true),
        options: { enabled: event()?.status === "draft" },
      },
      {
        hotkey: "Mod+Shift+C",
        callback: () => copyLink(),
        options: { enabled: Boolean(event()) },
      },
      {
        hotkey: "Mod+Shift+D",
        callback: () => setDiscontinueConfirmOpen(true),
        options: { enabled: isPublished() },
      },
      {
        hotkey: "Mod+Shift+O",
        callback: () => {
          const ev = event();
          if (ev) {
            void navigate({
              to: "/events/$slug",
              params: { slug: ev.slug },
            });
          }
        },
        options: { enabled: Boolean(isPublished() && event()) },
      },
    ],
    { ignoreInputs: true, preventDefault: true },
  );

  return (
    <div class="relative mx-auto flex w-full max-w-7xl flex-col gap-6">
      <EventOverviewHeader event={event()} />

      <EventQuickActions actions={actions()} />

      <QueryBoundary
        query={metricsQuery}
        what="as métricas do evento"
        skeleton={<div class="h-80 animate-pulse rounded-xl bg-muted" />}
      >
        {(metrics) => <EventOverviewDashboard metrics={metrics()} />}
      </QueryBoundary>

      <QueryBoundary
        query={editionsQuery}
        what="as edições"
        skeleton={<GridSkeleton minItemWidth="16rem" itemClass="h-28" />}
      >
        {(loadedEditions) => (
          <EventEditionsList eventId={eventId()} editions={loadedEditions()} />
        )}
      </QueryBoundary>

      <EventPaymentPanel
        connected={Boolean(event()?.payssage_seller_id)}
        disabled={!event() || isConnectingSeller() || isDisconnectingSeller()}
        onConnect={handleConnectSeller}
        onDisconnect={() => setDisconnectSellerConfirmOpen(true)}
      />

      <EventOverviewChecklist
        editionCount={editionsQuery().data?.length ?? 0}
        hasLogo={Boolean(event()?.logo_url)}
        hasBanner={Boolean(event()?.banner_url)}
        hasDescription={Boolean(event()?.description)}
        paymentConnected={Boolean(event()?.payssage_seller_id)}
        logoUploading={isImageUploading("logo_url")}
        bannerUploading={isImageUploading("banner_url")}
        onAddLogo={() => document.getElementById("event-logo-upload")?.click()}
        onAddBanner={() =>
          document.getElementById("event-banner-upload")?.click()
        }
        onEdit={() => setEditEventOpen(true)}
      />

      <EventOverviewDialogs
        event={event()}
        editOpen={editEventOpen()}
        publishOpen={publishConfirmOpen()}
        discontinueOpen={discontinueConfirmOpen()}
        disconnectOpen={disconnectSellerConfirmOpen()}
        publishing={isPublishing()}
        discontinuing={isDiscontinuing()}
        disconnecting={isDisconnectingSeller()}
        onEditOpenChange={setEditEventOpen}
        onPublishOpenChange={setPublishConfirmOpen}
        onDiscontinueOpenChange={setDiscontinueConfirmOpen}
        onDisconnectOpenChange={setDisconnectSellerConfirmOpen}
        onEdit={handleEditEvent}
        onPublish={handlePublish}
        onDiscontinue={handleDiscontinue}
        onDisconnect={handleDisconnectSeller}
      />
    </div>
  );
}
