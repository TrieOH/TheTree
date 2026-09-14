import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { createMemo, createSignal, Show } from "solid-js";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import type { EditionI } from "@/features/editions/model";
import {
  allTicketsQueryOptions,
  attendeeCountQueryOptions,
} from "@/features/tickets/api";
import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import { productsByEditionQueryOptions } from "@/features/products/api";
import { editionPurchasesQueryOptions } from "@/features/purchases/api";
import {
  allJoinedEventsQueryOptions,
  allOwnEventsQueryOptions,
} from "@/features/events/api";
import {
  usePatchEditionMutation,
  usePublishEditionMutation,
} from "@/features/editions/api/mutations";
import { useUploadQueue } from "@/features/upload-queue/hooks/use-upload-queue";
import { buildEditionOverviewMetrics } from "@/features/editions/model/edition-overview";
import { handleShare } from "@/shared/lib/share";
import { createHotkeys } from "@/shared/lib/hotkeys";
import type { ManageEditionValues } from "@/features/editions/ui/ManageEditionDialog";

import { EditionOverviewHeader } from "@/features/editions/ui/EditionOverviewHeader";
import {
  type EditionQuickAction,
  EditionQuickActions,
} from "@/features/editions/ui/EditionQuickActions";
import { EditionOverviewDashboard } from "@/features/editions/ui/EditionOverviewDashboard";
import { EditionOverviewChecklist } from "@/features/editions/ui/EditionOverviewChecklist";
import { EditionOverviewDialogs } from "@/features/editions/ui/EditionOverviewDialogs";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/",
)({
  head: () => ({
    meta: [{ title: "Edição - Admin - Univents" }],
  }),
  component: AdminEditionDetailView,
});

function AdminEditionDetailView(): JSX.Element {
  const params = Route.useParams();
  const navigate = useNavigate();
  const uploadQueue = useUploadQueue();

  const eventId = () => params().eventId;
  const editionId = () => params().editionId;

  const editionsQuery = useQuery(() => allAdminEditionsQueryOptions(eventId()));
  const ownEventsQuery = useQuery(() => allOwnEventsQueryOptions());
  const joinedEventsQuery = useQuery(() => allJoinedEventsQueryOptions());

  const purchasesQuery = useQuery(() => editionPurchasesQueryOptions(editionId()));
  const attendeeCountQuery = useQuery(() => attendeeCountQueryOptions(editionId()));
  const ticketsQuery = useQuery(() => allTicketsQueryOptions(editionId()));
  const productsQuery = useQuery(() => productsByEditionQueryOptions(editionId()));
  const programsQuery = useQuery(() => programsQueryOptions(editionId()));
  const occurrencesQuery = useQuery(() => occurrencesQueryOptions(editionId()));

  const [publishConfirmOpen, setPublishConfirmOpen] = createSignal(false);
  const [editEditionOpen, setEditEditionOpen] = createSignal(false);

  const editions = createMemo(() => editionsQuery().data ?? []);
  const edition = createMemo(
    () => editions().find((item: EditionI) => item.id === editionId()) ?? null,
  );

  const parentEvent = createMemo(() => {
    const own = ownEventsQuery().data ?? [];
    const joined = joinedEventsQuery().data ?? [];
    return [...own, ...joined].find((event) => event.id === eventId());
  });

  const eventSlug = createMemo(() => parentEvent()?.slug);

  const publishMutation = usePublishEditionMutation();
  const patchMutation = usePatchEditionMutation();

  const isPublishing = () => publishMutation.result().status === "pending";

  const isImageUploading = (field: "logo_url" | "banner_url") =>
    uploadQueue.tasks.some(
      (task) =>
        task.owner.type === "edition" &&
        task.owner.id === editionId() &&
        task.association?.handlerKey === "edition-image" &&
        task.association.input?.field === field &&
        !["completed", "failed", "rejected"].includes(task.status),
    );

  const copyLink = () => {
    const slug = eventSlug();
    const ed = edition();
    if (!slug) return;
    void handleShare(
      ed?.name ?? "Edição",
      `${window.location.origin}/events/${slug}`,
    );
  };

  const handleEditEdition = async (values: ManageEditionValues) => {
    const ed = edition();
    if (!ed) return false;
    try {
      await patchMutation.mutateAsync({
        eventId: eventId(),
        editionId: ed.id,
        data: values,
      });
      return true;
    } catch {
      return false;
    }
  };

  const handlePublish = async () => {
    const ed = edition();
    if (!ed) return;
    try {
      await publishMutation.mutateAsync({
        eventId: eventId(),
        editionId: ed.id,
      });
      setPublishConfirmOpen(false);
    } catch {
      // Handled
    }
  };

  const metrics = createMemo(() => {
    const ed = edition();
    if (!ed) return null;

    const purchases = purchasesQuery().data ?? [];
    const attendeesData = attendeeCountQuery().data;
    const attendees = typeof attendeesData === "number"
      ? attendeesData
      : (attendeesData?.count ?? 0);
    const tickets = ticketsQuery().data ?? [];
    const products = productsQuery().data ?? [];
    const programs = programsQuery().data ?? [];
    const occurrences = occurrencesQuery().data ?? [];

    return buildEditionOverviewMetrics({
      edition: ed,
      purchases,
      attendeeCount: attendees,
      ticketCount: tickets.length,
      productCount: products.length,
      programCount: programs.length,
      occurrenceCount: occurrences.length,
    });
  });

  const actions = (): EditionQuickAction[] => {
    const ed = edition();
    const isDraft = ed?.status === "draft";
    return [
      {
        label: "Editar edição",
        shortcut: "Mod+E",
        onClick: () => setEditEditionOpen(true),
        disabled: false,
        variant: "default",
      },
      ...(isDraft
        ? [
          {
            label: "Publicar edição",
            shortcut: "Mod+P",
            onClick: () => setPublishConfirmOpen(true),
            disabled: isPublishing(),
            variant: "default" as const,
          },
        ]
        : []),
      {
        label: "Compartilhar",
        shortcut: "Mod+Shift+C",
        onClick: () => copyLink(),
        disabled: !eventSlug() || !ed || ed.status === "draft",
        variant: "default",
      },
      ...(!isDraft && eventSlug()
        ? [
          {
            label: "Abrir painel público",
            shortcut: "Mod+Shift+O",
            to: "/events/$slug" as const,
            params: { slug: eventSlug() ?? "" },
            variant: "default" as const,
          },
        ]
        : []),
    ];
  };

  // Keyboard shortcuts
  createHotkeys(
    () => [
      {
        hotkey: "Mod+P",
        callback: () => setPublishConfirmOpen(true),
        options: { enabled: edition()?.status === "draft" },
      },
      {
        hotkey: "Mod+E",
        callback: () => setEditEditionOpen(true),
        options: { enabled: Boolean(edition()) },
      },
      {
        hotkey: "Mod+Shift+C",
        callback: () => copyLink(),
        options: {
          enabled: Boolean(eventSlug() && edition() && edition()?.status !== "draft"),
        },
      },
      {
        hotkey: "Mod+Shift+O",
        callback: () => {
          const slug = eventSlug();
          if (edition() && slug && edition()?.status !== "draft") {
            void navigate({ to: "/events/$slug", params: { slug } });
          }
        },
        options: {
          enabled: Boolean(eventSlug() && edition() && edition()?.status !== "draft"),
        },
      },
    ],
    { ignoreInputs: true, preventDefault: true },
  );

  return (
    <Show
      when={!editionsQuery().isPending}
      fallback={
        <div class="flex min-h-72 items-center justify-center text-sm text-muted-foreground">
          Carregando edição...
        </div>
      }
    >
      <Show
        when={edition()}
        fallback={
          <div class="flex min-h-72 items-center justify-center text-sm text-muted-foreground">
            Edição não encontrada.
          </div>
        }
      >
        {(ed) => (
          <div class="relative mx-auto flex w-full max-w-7xl flex-col gap-6 pb-28">
            <EditionOverviewHeader edition={ed()} eventId={eventId()} />

            <EditionQuickActions actions={actions()} />

            <Show when={metrics()}>
              {(m) => <EditionOverviewDashboard metrics={m()} />}
            </Show>

            <EditionOverviewChecklist
              hasLogo={Boolean(ed().logo_url)}
              hasBanner={Boolean(ed().banner_url)}
              hasDescription={Boolean(ed().description)}
              hasTagline={Boolean(ed().tagline)}
              hasLocation={Boolean(ed().location_name)}
              logoUploading={isImageUploading("logo_url")}
              bannerUploading={isImageUploading("banner_url")}
              onAddLogo={() =>
                document.getElementById(`edition-${ed().id}-logo-upload`)?.click()
              }
              onAddBanner={() =>
                document.getElementById(`edition-${ed().id}-banner-upload`)?.click()
              }
              onEdit={() => setEditEditionOpen(true)}
            />

            <EditionOverviewDialogs
              edition={ed()}
              editOpen={editEditionOpen()}
              publishOpen={publishConfirmOpen()}
              publishing={isPublishing()}
              onEditOpenChange={setEditEditionOpen}
              onPublishOpenChange={setPublishConfirmOpen}
              onEdit={handleEditEdition}
              onPublish={handlePublish}
            />
          </div>
        )}
      </Show>
    </Show>
  );
}
