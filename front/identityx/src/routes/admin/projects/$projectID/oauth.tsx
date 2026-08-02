import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import {
  EmptyState,
  PaginatedContainer,
  useLayoutHeader,
} from "@trieoh/ui-base";
import { Code2, Plus } from "lucide-react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import {
  allOAuthProvidersQueryOptions,
  createOAuthProviderFn,
  deleteOAuthProviderFn,
  setOAuthProviderEnabledFn,
  updateOAuthProviderFn,
} from "@/features/oauth/api";
import {
  type OAuthProviderCreateI,
  type OAuthProviderI,
  type OAuthProviderUpdateI,
  oauthProviderCreateSchema,
  oauthProviderUpdateSchema,
} from "@/features/oauth/model";
import { OAuthProviderCard } from "@/features/oauth/ui/oauth-provider-card";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { FormModal } from "@/widgets/modal/FormModal";
import { ConfirmModal } from "@/widgets/modal/modal";

export const Route = createFileRoute("/admin/projects/$projectID/oauth")({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectID } = Route.useParams();
  const queryClient = useQueryClient();
  const query = allOAuthProvidersQueryOptions(projectID);
  const { data: providers = [] } = useQuery(query);
  const [editing, setEditing] = useState<OAuthProviderI | null>(null);
  const [creating, setCreating] = useState(false);
  const [deleting, setDeleting] = useState<OAuthProviderI | null>(null);
  const [filter, setFilter] = useState("");
  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: query.queryKey });
  const mutationOptions = {
    onSuccess: invalidate,
    onError: (e: Error) => toast.error(e.message),
  };
  const create = useMutation({
    mutationFn: (data: OAuthProviderCreateI) =>
      createOAuthProviderFn(projectID, data),
    ...mutationOptions,
    onSuccess: () => {
      invalidate();
      setCreating(false);
      toast.success("OAuth provider configured");
    },
  });
  const availableProviders = (["google", "github"] as const).filter(
    (provider) => !providers.some((item) => item.provider === provider),
  );
  const defaultProvider = availableProviders[0] ?? "google";
  const update = useMutation({
    mutationFn: (data: OAuthProviderUpdateI) => {
      if (!editing) throw new Error("No OAuth provider selected");
      return updateOAuthProviderFn(editing.id, data);
    },
    ...mutationOptions,
    onSuccess: () => {
      invalidate();
      setEditing(null);
      toast.success("OAuth provider updated");
    },
  });
  const toggle = useMutation({
    mutationFn: (p: OAuthProviderI) =>
      setOAuthProviderEnabledFn(p.id, !p.enabled),
    ...mutationOptions,
    onSuccess: () => {
      invalidate();
      toast.success("OAuth provider status updated");
    },
  });
  const remove = useMutation({
    mutationFn: (p: OAuthProviderI) => deleteOAuthProviderFn(p.id),
    ...mutationOptions,
    onSuccess: () => {
      invalidate();
      toast.success("OAuth provider removed");
    },
  });

  const filteredProviders = providers.filter((provider) => {
    const search = filter.trim().toLowerCase();
    return (
      !search ||
      provider.provider.toLowerCase().includes(search) ||
      provider.client_id.toLowerCase().includes(search)
    );
  });

  const header = useMemo(
    () => (
      <div>
        <div>
          <h1 className="text-lg font-semibold">OAuth Providers</h1>
          <p className="text-sm text-muted-foreground">
            Configure social login for this project.
          </p>
        </div>
      </div>
    ),
    [],
  );

  useLayoutHeader(header);
  return (
    <div className="space-y-3">
      <PaginatedContainer<OAuthProviderI>
        items={filteredProviders}
        layout="list"
        pageSize={10}
        itemLabel="OAuth providers"
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Search providers or client IDs…"
        sortFields={[
          { key: "provider", label: "Provider" },
          { key: "enabled", label: "Status" },
          { key: "updated_at", label: "Last updated" },
        ]}
        headerActions={
          <ShadowButton
            onClick={() => setCreating(true)}
            variant="outline"
            disabled={availableProviders.length === 0}
            className="h-9 sm:w-auto px-3 rounded-sm"
            leftIcon={<Plus size={16} />}
            value="Add provider"
          />
        }
        renderItems={(slice) =>
          slice.map((p) => (
            <OAuthProviderCard
              key={p.id}
              data={p}
              onEdit={setEditing}
              onToggle={(p) => toggle.mutate(p)}
              onDelete={setDeleting}
            />
          ))
        }
        emptyState={
          <EmptyState
            icon={Code2}
            title="No OAuth providers"
            description="Add Google or GitHub to enable social login for this project."
          />
        }
      />
      <FormModal<OAuthProviderCreateI>
        isOpen={creating}
        onClose={() => setCreating(false)}
        title="Add OAuth provider"
        description="The client secret is encrypted and never returned."
        schema={oauthProviderCreateSchema}
        fields={[
          {
            name: "provider",
            label: "Provider",
            type: "option-picker",
            required: true,
            options: availableProviders.map((provider) => ({
              label: provider === "google" ? "Google" : "GitHub",
              value: provider,
            })),
          },
          { name: "client_id", label: "Client ID", required: true },
          { name: "client_secret", label: "Client secret", required: true },
        ]}
        defaultValues={{
          provider: defaultProvider,
          client_id: "",
          client_secret: "",
        }}
        formId={`create-oauth-provider-${availableProviders.join("-")}`}
        submitLabel="Save provider"
        isLoading={create.isPending}
        onSubmit={async (data) => {
          await create.mutateAsync(data);
        }}
      />
      <ConfirmModal
        isOpen={!!deleting}
        onClose={() => setDeleting(null)}
        onConfirm={() => {
          if (!deleting) return;
          remove.mutate(deleting, { onSuccess: () => setDeleting(null) });
        }}
        title="Remove OAuth provider?"
        description={`This will remove the ${deleting?.provider ?? "selected"} configuration and disable new social logins for this project.`}
        confirmText="Remove provider"
        variant="destructive"
        isLoading={remove.isPending}
      />
      <FormModal<OAuthProviderUpdateI>
        isOpen={!!editing}
        onClose={() => setEditing(null)}
        title="Edit OAuth provider"
        description="Leave the secret empty to keep the current one."
        schema={oauthProviderUpdateSchema}
        fields={[
          { name: "client_id", label: "Client ID", required: true },
          { name: "client_secret", label: "New client secret" },
        ]}
        defaultValues={{
          client_id: editing?.client_id ?? "",
          client_secret: "",
        }}
        formId="edit-oauth-provider"
        submitLabel="Save changes"
        isLoading={update.isPending}
        onSubmit={async (data) => {
          await update.mutateAsync(data);
        }}
      />
    </div>
  );
}
