import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import type { SortState } from "@trieoh/ui-base";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { Calendar, Plus } from "lucide-react";
import { useState } from "react";
import { allAdminActivitiesQueryOptions } from "@/features/activities/api";
import {
  useCompleteActivityMutation,
  useCreateActivityMutation,
  usePublishActivityMutation,
  useUpdateActivityMutation,
} from "@/features/activities/api/mutations";
import type { ActivityI } from "@/features/activities/model";
import AdminActivityCard from "@/features/activities/ui/AdminActivityCard";
import { ManageActivityModal } from "@/features/activities/ui/ManageActivityModal";
import { Button } from "@/shared/ui/shadcn/button";
import { AlertModal } from "@/widgets/ui/alert-modal";

const STATUS_SORT_ORDER: Record<ActivityI["status"], number> = {
  draft: 0,
  published: 1,
  ongoing: 2,
  completed: 3,
  canceled: 4,
};

export const Route = createLazyFileRoute(
  "/admin/events/$eventId_/editions/$editionId/activities/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId, editionId } = Route.useParams();
  const { data: activities = [] } = useQuery(
    allAdminActivitiesQueryOptions(eventId, editionId),
  );
  const createActivityMutation = useCreateActivityMutation();
  const updateActivityMutation = useUpdateActivityMutation();
  const publishActivityMutation = usePublishActivityMutation();
  const completeActivityMutation = useCompleteActivityMutation();
  const [filter, setFilter] = useState("");
  const [sort, setSort] = useState<SortState<ActivityI>>({
    field: "starts_at",
    direction: "asc",
  });
  const [modalState, setModalState] = useState<{
    open: boolean;
    activity?: ActivityI;
  }>({
    open: false,
  });
  const [publishingActivity, setPublishingActivity] =
    useState<ActivityI | null>(null);
  const [completingActivity, setCompletingActivity] =
    useState<ActivityI | null>(null);
  const handlePublishActivity = (currentActivity: ActivityI) => {
    setPublishingActivity(currentActivity);
  };
  const handleCompleteActivity = (currentActivity: ActivityI) => {
    setCompletingActivity(currentActivity);
  };

  const filteredActivities = [...activities]
    .filter((activity) => {
      const search = filter.trim().toLowerCase();
      if (!search) return true;

      return [
        activity.title,
        activity.location,
        activity.presenter_name ?? "",
        activity.status,
        activity.difficulty,
      ].some((value) => value.toLowerCase().includes(search));
    })
    .sort((a, b) => {
      const direction = sort.direction === "asc" ? 1 : -1;

      if (sort.field === "starts_at") {
        return (
          (new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime()) *
          direction
        );
      }

      if (sort.field === "status") {
        return (
          (STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status]) *
          direction
        );
      }

      if (sort.field === "difficulty") {
        return (
          String(a.difficulty).localeCompare(String(b.difficulty)) * direction
        );
      }

      return (
        String(a[sort.field]).localeCompare(String(b[sort.field])) * direction
      );
    });

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <PaginatedContainer<ActivityI>
        items={filteredActivities}
        layout="grid"
        minItemWidth="16rem"
        pageSize={8}
        gap="6"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          {
            key: "starts_at",
            label: "Início",
            comparator: (a, b) =>
              new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime(),
          },
          { key: "title", label: "Título" },
          { key: "location", label: "Local" },
          {
            key: "status",
            label: "Status",
            comparator: (a, b) =>
              STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status],
          },
          { key: "difficulty", label: "Dificuldade" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por título, local, palestrante ou status..."
        itemLabel="atividades"
        headerActions={
          <Button
            type="button"
            onClick={() => setModalState({ open: true, activity: undefined })}
            className="h-9 gap-2"
          >
            <Plus className="size-4" />
            Nova atividade
          </Button>
        }
        emptyState={
          <EmptyState
            icon={Calendar}
            eyebrow="Atividades"
            title="Nenhuma atividade encontrada"
            description="Crie a primeira atividade para começar a organizar essa edição."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((activity, idx) => (
            <AdminActivityCard
              key={activity.id}
              activity={activity}
              index={idx}
              onManage={(currentActivity) =>
                setModalState({ open: true, activity: currentActivity })
              }
              onPublish={
                activity.status === "draft" ? handlePublishActivity : undefined
              }
              onComplete={
                activity.status === "ongoing"
                  ? handleCompleteActivity
                  : undefined
              }
            />
          ))
        }
      />

      <ManageActivityModal
        key={modalState.activity?.id ?? "activity-create"}
        open={modalState.open}
        activity={modalState.activity}
        onOpenChange={(open) => {
          if (open) {
            setModalState((prev) => ({ ...prev, open }));
            return;
          }

          setModalState({ open: false, activity: undefined });
        }}
        onCreate={async (values) => {
          const res = await createActivityMutation.mutateAsync({
            eventId,
            editionId,
            data: values,
          });

          return res.success ? res.data : false;
        }}
        onUpdate={async (activityId, values) => {
          const res = await updateActivityMutation.mutateAsync({
            eventId,
            editionId,
            activityId,
            data: values,
          });

          return res.success ? res.data : false;
        }}
      />

      <AlertModal
        open={Boolean(publishingActivity)}
        onOpenChange={() => setPublishingActivity(null)}
        title="Publicar atividade?"
        description={
          publishingActivity
            ? `Ao publicar "${publishingActivity.title}", ela ficará disponível para os participantes.`
            : undefined
        }
        confirmLabel="Publicar atividade"
        variant="default"
        loading={publishActivityMutation.isPending}
        onConfirm={async () => {
          if (!publishingActivity) return;
          await publishActivityMutation.mutateAsync({
            eventId,
            editionId,
            activityId: publishingActivity.id,
          });
          setPublishingActivity(null);
        }}
      />

      <AlertModal
        open={Boolean(completingActivity)}
        onOpenChange={() => setCompletingActivity(null)}
        title="Concluir atividade?"
        description={
          completingActivity
            ? `Ao concluir "${completingActivity.title}", ela ficará marcada como finalizada.`
            : undefined
        }
        confirmLabel="Concluir atividade"
        variant="default"
        loading={completeActivityMutation.isPending}
        onConfirm={async () => {
          if (!completingActivity) return;
          await completeActivityMutation.mutateAsync({
            eventId,
            editionId,
            activityId: completingActivity.id,
          });
          setCompletingActivity(null);
        }}
      />
    </div>
  );
}
