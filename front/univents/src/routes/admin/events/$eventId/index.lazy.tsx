import { useHotkeys } from "@tanstack/react-hotkeys";
import { useQueries, useQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import { Calendar } from "lucide-react";
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
import { EventOverviewDashboard } from "@/features/events/ui/EventOverviewDashboard";
import {
  type EventQuickAction,
  EventQuickActions,
} from "@/features/events/ui/EventQuickActions";
import { EventVisualCard } from "@/features/events/ui/EventVisualCard";
import { ManageEventModal } from "@/features/events/ui/ManageEventModal";
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
import { Badge } from "@/shared/ui/shadcn/badge";
import { AlertModal } from "@/widgets/ui/alert-modal";
import { StepChecklist } from "@/widgets/ui/step-checklist";

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
