import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useQuery, useQueryClient } from "@trieoh/front-core-solid";
import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import {
  Show,
  createEffect,
  createMemo,
  createSignal,
} from "solid-js";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import type { EditionI } from "@/features/editions/model";
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
import { useUploadQueue } from "@/features/upload-queue";
import { toast } from "@/shared/ui/toast";

export const Route = createFileRoute("/admin/events/$eventId/")({
  head: ({ params }) => ({
    meta: [{ title: `Evento ${params.eventId} - Admin - Univents` }],
  }),
  component: AdminEventOverviewRoute,
});

function AdminEventOverviewRoute(): JSX.Element {
  const params = Route.useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const uploadQueue = useUploadQueue();

  const eventId = () => params().eventId;

  const ownedQuery = useQuery(() => allOwnEventsQueryOptions());
  const joinedQuery = useQuery(() => allJoinedEventsQueryOptions());
  const editionsQuery = useQuery(() => allAdminEditionsQueryOptions(eventId()));

  const [publishConfirmOpen, setPublishConfirmOpen] = createSignal(false);
  const [discontinueConfirmOpen, setDiscontinueConfirmOpen] =
    createSignal(false);
  const [disconnectSellerConfirmOpen, setDisconnectSellerConfirmOpen] =
    createSignal(false);
  const [editEventOpen, setEditEventOpen] = createSignal(false);

  const event = createMemo(() => {
    const owned = ownedQuery().data ?? [];
    const joined = joinedQuery().data ?? [];
    return [...owned, ...joined].find((item) => item.id === eventId()) ?? null;
  });

  const editions = createMemo(() => editionsQuery().data ?? []);

  const [metrics, setMetrics] = createSignal<EventOverviewMetrics | null>(null);

  createEffect(
    () => editions(),
    (currentEditions) => {
      let cancelled = false;

      if (currentEditions.length === 0) {
        setMetrics(
          buildEventOverviewMetrics({
            editions: [],
            purchasesByEdition: [],
            attendeeCounts: [],
            ticketCounts: [],
            productCounts: [],
            programCounts: [],
            occurrenceCounts: [],
          }),
        );
        return;
      }

      const loadMetrics = async (eds: EditionI[]) => {
        try {
          const results = await Promise.all(
            eds.map((edition) =>
              Promise.all([
                queryClient
                  .fetchQuery(editionPurchasesQueryOptions(edition.id))
                  .catch(() => []),
                queryClient
                  .fetchQuery(attendeeCountQueryOptions(edition.id))
                  .catch(() => ({ count: 0 })),
                queryClient
                  .fetchQuery(allTicketsQueryOptions(edition.id))
                  .catch(() => []),
                queryClient
                  .fetchQuery(productsQueryOptions(edition.id))
                  .catch(() => []),
                queryClient
                  .fetchQuery(programsQueryOptions(edition.id))
                  .catch(() => []),
                queryClient
                  .fetchQuery(occurrencesQueryOptions(edition.id))
                  .catch(() => []),
              ]),
            ),
          );
          if (cancelled) return;

          const purchasesByEdition = results.map(
            (r) => (r[0] as EditionPurchase[]) ?? [],
          );
          const attendeeCounts = results.map(
            (r) => ((r[1] as { count?: number })?.count ?? 0),
          );
          const ticketCounts = results.map(
            (r) => ((r[2] as unknown[]) ?? []).length,
          );
          const productCounts = results.map(
            (r) => ((r[3] as unknown[]) ?? []).length,
          );
          const programCounts = results.map(
            (r) => ((r[4] as unknown[]) ?? []).length,
          );
          const occurrenceCounts = results.map(
            (r) => ((r[5] as unknown[]) ?? []).length,
          );

          const calculated = buildEventOverviewMetrics({
            editions: eds,
            purchasesByEdition,
            attendeeCounts,
            ticketCounts,
            productCounts,
            programCounts,
            occurrenceCounts,
          });

          setMetrics(calculated);
        } catch {
          // Handled silently
        }
      };

      void loadMetrics(currentEditions);

      return () => {
        cancelled = true;
      };
    },
  );

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
    void navigator.clipboard.writeText(
      `${window.location.origin}/events/${ev.slug}`,
    );
    toast.success("Link copiado");
  };

  const handleEditEvent = async (values: ManageEventValues) => {
    const ev = event();
    if (!ev) return false;
    try {
      await patchMutation.mutateAsync({
        eventId: ev.id,
        data: values,
      });
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

  // Keyboard shortcuts
  createEffect(
    () => event(),
    (ev) => {
      const handleKeyDown = (e: KeyboardEvent) => {
        const isMod = e.metaKey || e.ctrlKey;
        if (!isMod) return;

        if (e.key.toLowerCase() === "e" && !e.shiftKey) {
          e.preventDefault();
          setEditEventOpen(true);
        } else if (e.key.toLowerCase() === "p" && !e.shiftKey) {
          if (ev?.status === "draft") {
            e.preventDefault();
            setPublishConfirmOpen(true);
          }
        } else if (e.key.toLowerCase() === "c" && e.shiftKey) {
          e.preventDefault();
          copyLink();
        } else if (e.key.toLowerCase() === "d" && e.shiftKey) {
          if (ev?.status === "active") {
            e.preventDefault();
            setDiscontinueConfirmOpen(true);
          }
        } else if (e.key.toLowerCase() === "o" && e.shiftKey) {
          if (ev?.status === "active" && ev) {
            e.preventDefault();
            void navigate({
              to: "/events/$slug",
              params: { slug: ev.slug },
            });
          }
        }
      };

      window.addEventListener("keydown", handleKeyDown);
      return () => {
        window.removeEventListener("keydown", handleKeyDown);
      };
    },
  );

  return (
    <div class="relative mx-auto flex w-full max-w-7xl flex-col gap-6">
      <EventOverviewHeader event={event()} />

      <EventQuickActions actions={actions()} />

      <Show when={metrics()}>
        {(m) => <EventOverviewDashboard metrics={m()} />}
      </Show>

      <EventEditionsList eventId={eventId()} editions={editions()} />

      <EventPaymentPanel
        connected={Boolean(event()?.payssage_seller_id)}
        disabled={!event() || isConnectingSeller() || isDisconnectingSeller()}
        onConnect={handleConnectSeller}
        onDisconnect={() => setDisconnectSellerConfirmOpen(true)}
      />

      <EventOverviewChecklist
        editionCount={editions().length}
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
