import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import {
  Show,
  createEffect,
  createMemo,
  createSignal,
} from "solid-js";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import {
  usePatchEditionMutation,
  usePublishEditionMutation,
} from "@/features/editions/api/mutations";
import {
  type EditionOverviewMetrics,
  buildEditionOverviewMetrics,
} from "@/features/editions/model/edition-overview";
import { EditionOverviewChecklist } from "@/features/editions/ui/EditionOverviewChecklist";
import { EditionOverviewDashboard } from "@/features/editions/ui/EditionOverviewDashboard";
import { EditionOverviewDialogs } from "@/features/editions/ui/EditionOverviewDialogs";
import { EditionOverviewHeader } from "@/features/editions/ui/EditionOverviewHeader";
import {
  type EditionQuickAction,
  EditionQuickActions,
} from "@/features/editions/ui/EditionQuickActions";
import type { ManageEditionValues } from "@/features/editions/ui/ManageEditionDialog";
import {
  allJoinedEventsQueryOptions,
  allOwnEventsQueryOptions,
} from "@/features/events/api";
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
import { toast } from "@/shared/ui/toast";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/",
)({
  head: ({ params }) => ({
    meta: [{ title: `Edição ${params.editionId} - Admin - Univents` }],
  }),
  component: AdminEditionDetailRoute,
});

function AdminEditionDetailRoute(): JSX.Element {
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
    () => editions().find((item) => item.id === editionId()) ?? null,
  );

  const eventSlug = createMemo(() => {
    const own = ownEventsQuery().data ?? [];
    const joined = joinedEventsQuery().data ?? [];
    return [...own, ...joined].find((event) => event.id === eventId())?.slug;
  });

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
    if (!slug) return;
    void navigator.clipboard.writeText(
      `${window.location.origin}/events/${slug}`,
    );
    toast.success("Link copiado");
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
    try {
      await publishMutation.mutateAsync({
        eventId: eventId(),
        editionId: editionId(),
      });
      setPublishConfirmOpen(false);
    } catch {
      // Handled
    }
  };

  const metrics = createMemo<EditionOverviewMetrics | null>(() => {
    const ed = edition();
    if (!ed) return null;

    const purchases = purchasesQuery().data ?? [];
    const attendeeCount = attendeeCountQuery().data?.count ?? 0;
    const tickets = ticketsQuery().data ?? [];
    const products = productsQuery().data ?? [];
    const programs = programsQuery().data ?? [];
    const occurrences = occurrencesQuery().data ?? [];

    return buildEditionOverviewMetrics({
      edition: ed,
      purchases,
      attendeeCount,
      ticketCount: tickets.length,
      productCount: products.length,
      programCount: programs.length,
      occurrenceCount: occurrences.length,
    });
  });

  const actions = (): EditionQuickAction[] => {
    const ed = edition();
    const slug = eventSlug();
    const isDraft = ed?.status === "draft";

    return [
      {
        label: "Editar edição",
        shortcut: "Mod+E",
        onClick: () => setEditEditionOpen(true),
        disabled: !ed,
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
        label: "Copiar link público",
        shortcut: "Mod+Shift+C",
        onClick: copyLink,
        disabled: !slug || isDraft,
        variant: "default",
      },
      ...(!isDraft && slug
        ? [
          {
            label: "Abrir página pública",
            shortcut: "Mod+Shift+O",
            to: "/events/$slug" as const,
            params: { slug },
            variant: "default" as const,
          },
        ]
        : []),
    ];
  };

  // Keyboard shortcuts
  createEffect(
    () => ({ ed: edition(), slug: eventSlug() }),
    ({ ed, slug }) => {
      const handleKeyDown = (e: KeyboardEvent) => {
        const isMod = e.metaKey || e.ctrlKey;
        if (!isMod) return;

        if (e.key.toLowerCase() === "e" && !e.shiftKey) {
          e.preventDefault();
          setEditEditionOpen(true);
        } else if (e.key.toLowerCase() === "p" && !e.shiftKey) {
          if (ed?.status === "draft") {
            e.preventDefault();
            setPublishConfirmOpen(true);
          }
        } else if (e.key.toLowerCase() === "c" && e.shiftKey) {
          if (slug && ed?.status !== "draft") {
            e.preventDefault();
            copyLink();
          }
        } else if (e.key.toLowerCase() === "o" && e.shiftKey) {
          if (slug && ed?.status !== "draft") {
            e.preventDefault();
            void navigate({
              to: "/events/$slug",
              params: { slug },
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
