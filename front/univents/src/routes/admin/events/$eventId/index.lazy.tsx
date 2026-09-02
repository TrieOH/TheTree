import { useHotkeys } from "@tanstack/react-hotkeys";
import { useQueries, useQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
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
import { buildEventOverviewMetrics } from "@/features/events/model/event-overview";
import { EventEditionsList } from "@/features/events/ui/EventEditionsList";
import { EventOverviewChecklist } from "@/features/events/ui/EventOverviewChecklist";
import { EventOverviewDashboard } from "@/features/events/ui/EventOverviewDashboard";
import { EventOverviewDialogs } from "@/features/events/ui/EventOverviewDialogs";
import { EventOverviewHeader } from "@/features/events/ui/EventOverviewHeader";
import {
  type EventQuickAction,
  EventQuickActions,
} from "@/features/events/ui/EventQuickActions";
import {
  useConnectEventSellerMutation,
  useDisconnectEventSellerMutation,
} from "@/features/payments/api/mutations";
import { EventPaymentPanel } from "@/features/payments/ui/EventPaymentPanel";
import { productsByEditionQueryOptions } from "@/features/products/api";
import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import { editionPurchasesQueryOptions } from "@/features/purchases/api";
import {
  allTicketsQueryOptions,
  attendeeCountQueryOptions,
} from "@/features/tickets/api";
import { useUploadQueue } from "@/features/upload-queue";

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
  const attendeeCountQueries = useQueries({
    queries: editions.map((edition) => attendeeCountQueryOptions(edition.id)),
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
  const overviewMetrics = buildEventOverviewMetrics({
    editions,
    purchasesByEdition: purchaseQueries.map((query) => query.data ?? []),
    attendeeCounts: attendeeCountQueries.map((query) => query.data?.count ?? 0),
    ticketCounts: ticketQueries.map((query) => query.data?.length ?? 0),
    productCounts: productQueries.map((query) => query.data?.length ?? 0),
    programCounts: programQueries.map((query) => query.data?.length ?? 0),
    occurrenceCounts: occurrenceQueries.map((query) => query.data?.length ?? 0),
  });
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

  const actions: EventQuickAction[] = [
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

  return (
    <>
      <div className="relative mx-auto flex w-full max-w-7xl flex-col gap-6 p-6 pb-28!">
        <EventOverviewHeader event={event} />

        <EventQuickActions actions={actions} />

        <EventOverviewDashboard metrics={overviewMetrics} />

        <EventEditionsList eventId={eventId} editions={editions} />

        <EventPaymentPanel
          connected={Boolean(event?.payssage_seller_id)}
          disabled={
            !event ||
            connectSellerMutation.isPending ||
            disconnectSellerMutation.isPending
          }
          onConnect={(provider) =>
            connectSellerMutation.mutate({ eventId, provider })
          }
          onDisconnect={() => setDisconnectSellerConfirmOpen(true)}
        />

        <EventOverviewChecklist
          editionCount={editions.length}
          hasLogo={Boolean(event?.logo_url)}
          hasBanner={Boolean(event?.banner_url)}
          hasDescription={Boolean(event?.description)}
          paymentConnected={Boolean(event?.payssage_seller_id)}
          logoUploading={isImageUploading("logo_url")}
          bannerUploading={isImageUploading("banner_url")}
          onAddLogo={() =>
            document.getElementById("event-logo-upload")?.click()
          }
          onAddBanner={() =>
            document.getElementById("event-banner-upload")?.click()
          }
          onEdit={() => setEditEventOpen(true)}
        />
      </div>

      <EventOverviewDialogs
        event={event}
        editOpen={editEventOpen}
        publishOpen={publishConfirmOpen}
        discontinueOpen={discontinueConfirmOpen}
        disconnectOpen={disconnectSellerConfirmOpen}
        publishing={publishEventMutation.isPending}
        discontinuing={discontinueEventMutation.isPending}
        disconnecting={disconnectSellerMutation.isPending}
        onEditOpenChange={setEditEventOpen}
        onPublishOpenChange={setPublishConfirmOpen}
        onDiscontinueOpenChange={setDiscontinueConfirmOpen}
        onDisconnectOpenChange={setDisconnectSellerConfirmOpen}
        onEdit={(values) =>
          event
            ? patchEventMutation.mutateAsync({ eventId, data: values })
            : false
        }
        onDisconnect={() =>
          disconnectSellerMutation.mutate(eventId, {
            onSuccess: () => setDisconnectSellerConfirmOpen(false),
          })
        }
        onPublish={() => {
          handlePublishEvent();
          setPublishConfirmOpen(false);
        }}
        onDiscontinue={() => {
          handleDiscontinueEvent();
          setDiscontinueConfirmOpen(false);
        }}
      />
    </>
  );
}
