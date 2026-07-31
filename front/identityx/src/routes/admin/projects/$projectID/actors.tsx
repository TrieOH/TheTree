import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import {
  EmptyState,
  PaginatedContainer,
  useLayoutHeader,
} from "@trieoh/ui-base";
import {
  Bot,
  Cpu,
  Fingerprint,
  KeyRound,
  Plus,
  Shield,
  User2,
} from "lucide-react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { allActorsQueryOptions, createActorFn } from "@/features/actor/api";
import {
  type ActorCreateI,
  type ActorI,
  actorCreateSchema,
} from "@/features/actor/model";
import { ActorCard } from "@/features/actor/ui/actor-card";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { FormModal } from "@/widgets/modal/FormModal";

export const Route = createFileRoute("/admin/projects/$projectID/actors")({
  component: RouteComponent,
});

function RouteComponent() {
  const queryClient = useQueryClient();
  const { projectID } = Route.useParams();
  const { organizationID } = Route.useSearch();
  const [filter, setFilter] = useState("");
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  const { data: actors = [] } = useQuery(
    allActorsQueryOptions(projectID, organizationID),
  );

  const header = useMemo(
    () => (
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-lg font-semibold tracking-tight">Actors</h1>
          <p className="text-sm text-muted-foreground">
            {actors.length === 0
              ? "No actors yet in this project"
              : `${actors.length} actor${actors.length !== 1 ? "s" : ""} in this project`}
          </p>
        </div>
      </div>
    ),
    [actors.length],
  );

  useLayoutHeader(header);

  const { mutate: createActor, isPending: isCreating } = useMutation({
    mutationFn: (data: ActorCreateI) =>
      createActorFn(projectID, data, organizationID),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({
          queryKey: allActorsQueryOptions(projectID, organizationID).queryKey,
        });
        setIsCreateOpen(false);
        toast.success(response.message || "Actor created successfully");
      } else toast.error(response.message || "Failed to create actor");
    },
    onError: (error: Error) => toast.error(error.message),
  });

  const filteredActors = actors.filter((actor) => {
    const search = filter.toLowerCase().trim();
    if (!search) return true;

    return (
      actor.id.toLowerCase().includes(search) ||
      actor.type.toLowerCase().includes(search) ||
      actor.auth_method.toLowerCase().includes(search) ||
      (actor.email?.toLowerCase().includes(search) ?? false)
    );
  });

  return (
    <div>
      <PaginatedContainer<ActorI>
        items={filteredActors}
        layout="grid"
        pageSize={10}
        minItemWidth="15rem"
        sortFields={[
          { key: "created_at", label: "Created At" },
          { key: "type", label: "Type" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by id, email, type or auth method..."
        itemLabel="actors"
        headerActions={
          <ShadowButton
            onClick={() => setIsCreateOpen(true)}
            variant="outline"
            className="h-9 sm:w-auto px-3 rounded-sm"
            leftIcon={<Plus size={16} />}
            value="Create Actor"
          />
        }
        renderItems={(slice) =>
          slice.map((item) => <ActorCard key={item.id} data={item} />)
        }
        emptyState={
          <EmptyState
            icon={Bot}
            title="No actors"
            description="No actors found for this project."
          />
        }
      />

      <FormModal<ActorCreateI>
        title="Create Actor"
        description="Create a new actor for this project."
        submitLabel="Create Actor"
        schema={actorCreateSchema}
        formId="create-actor-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={createActor}
        defaultValues={{
          auth_method: "password",
          type: "human",
          email: undefined,
        }}
        isLoading={isCreating}
        fields={[
          {
            name: "auth_method",
            label: "Auth Method",
            type: "option-picker",
            options: [
              { value: "password", label: "Password", icon: KeyRound },
              { value: "api_key", label: "API Key", icon: Fingerprint },
            ],
            required: true,
          },
          {
            name: "type",
            label: "Type",
            type: "option-picker",
            options: [
              { value: "human", label: "Human", icon: User2 },
              { value: "service", label: "Service", icon: Shield },
              { value: "machine", label: "Machine", icon: Cpu },
            ],
            required: true,
          },
          {
            name: "email",
            label: "Email",
            type: "text",
            placeholder: "e.g. user@example.com",
          },
        ]}
      />
    </div>
  );
}
