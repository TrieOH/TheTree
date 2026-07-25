import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { PaginatedContainer, useLayoutHeader } from "@trieoh/ui-base";
import { Plus } from "lucide-react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import {
  allCapabilitiesQueryOptions,
  createCapabilityFn,
} from "@/features/capabilities/api";
import {
  type CapabilityCreateI,
  type CapabilityI,
  capabilityCreateSchema,
} from "@/features/capabilities/model";
import { CapabilityCard } from "@/features/capabilities/ui/capability-card";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { FormModal } from "@/widgets/modal/FormModal";

export const Route = createFileRoute("/admin/projects/$projectID/capabilities")(
  {
    component: RouteComponent,
  },
);

function RouteComponent() {
  const queryClient = useQueryClient();
  const { projectID } = Route.useParams();

  const { data: capabilities = [] } = useQuery(
    allCapabilitiesQueryOptions(projectID),
  );

  const [filter, setFilter] = useState("");
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  const projectCapabilities = capabilities.filter(
    (capability) => capability.project_id === projectID,
  );

  const header = useMemo(
    () => (
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-lg font-semibold tracking-tight">Capabilities</h1>
          <p className="text-sm text-muted-foreground">
            {projectCapabilities.length === 0
              ? "No capabilities yet in this project"
              : `${projectCapabilities.length} capability${projectCapabilities.length !== 1 ? "s" : ""} in this project`}
          </p>
        </div>
      </div>
    ),
    [projectCapabilities.length],
  );

  useLayoutHeader(header);

  const filteredCapabilities = projectCapabilities.filter((capability) => {
    const search = filter.toLowerCase().trim();

    if (!search) return true;

    return (
      capability.resource.toLowerCase().includes(search) ||
      capability.action.toLowerCase().includes(search) ||
      `${capability.resource}:${capability.action}`
        .toLowerCase()
        .includes(search)
    );
  });

  const { mutate: createCapability, isPending: isCreating } = useMutation({
    mutationFn: (data: CapabilityCreateI) =>
      createCapabilityFn(projectID, data),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({
          queryKey: allCapabilitiesQueryOptions(projectID).queryKey,
        });
        setIsCreateOpen(false);
        toast.success(response.message || "Capability created successfully");
      } else toast.error(response.message || "Failed to create capability");
    },
    onError: (error: Error) => toast.error(error.message),
  });

  return (
    <div>
      <PaginatedContainer<CapabilityI>
        items={filteredCapabilities}
        layout="grid"
        pageSize={10}
        minItemWidth="15rem"
        sortFields={[
          { key: "resource", label: "Resource" },
          { key: "action", label: "Action" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by resource or action..."
        itemLabel="capabilities"
        headerActions={
          <ShadowButton
            onClick={() => setIsCreateOpen(true)}
            variant="outline"
            className="h-9 sm:w-auto px-3 rounded-sm"
            leftIcon={<Plus size={16} />}
            value="Add Capability"
          />
        }
        renderItems={(slice) =>
          slice.map((item) => <CapabilityCard key={item.id} data={item} />)
        }
      />

      <FormModal<CapabilityCreateI>
        title="Add Capability"
        description="Create a new capability for this project."
        submitLabel="Add Capability"
        schema={capabilityCreateSchema}
        formId="add-capability-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={createCapability}
        defaultValues={{ resource: "", action: "" }}
        isLoading={isCreating}
        fields={[
          {
            name: "resource",
            label: "Resource",
            type: "text",
            placeholder: "e.g. invoices",
            required: true,
          },
          {
            name: "action",
            label: "Action",
            type: "text",
            placeholder: "e.g. read",
            required: true,
          },
        ]}
      />
    </div>
  );
}
