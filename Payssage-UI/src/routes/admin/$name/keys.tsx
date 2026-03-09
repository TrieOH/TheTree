import { createFileRoute, useParams } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { ConfirmModal } from '#/widgets/modal/modal'
import { apiKeyCreateSchema } from '#/features/keys/model'
import { toast } from 'sonner'
import FormModal from '#/widgets/modal/form-modal'
import type { ApiKeyCreateI, ApiKeyCreateResponseI, ApiKeyI } from '#/features/keys/model'
import { useState } from 'react'
import { allWorkspaceApiKeysQueryOptions, createApiKeyOnWorkspaceFn, revokeApiKeyOnWorkspaceFn } from '#/features/keys/api'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { KeyList } from '#/features/keys/ui/key-list'
import { ApiKeyCreatedModal } from '#/features/keys/ui/api-key-created-modal'

export const Route = createFileRoute('/admin/$name/keys')({
  component: RouteComponent,
})


function RouteComponent() {
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [revokeKeyId, setRevokeKeyId] = useState<string | null>(null)
  const [newlyCreatedKey, setNewlyCreatedKey] = useState<ApiKeyCreateResponseI | null>(null)

  const { name } = useParams({ from: '/admin/$name' })
  const queryClient = useQueryClient();
  const { data: keys = [], isLoading } = useQuery(allWorkspaceApiKeysQueryOptions(name))

  const { mutate: createApiKey, isPending: isPendingCreate } = useMutation({
    mutationFn: (data: ApiKeyCreateI) => createApiKeyOnWorkspaceFn(data, name),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allWorkspaceApiKeysQueryOptions(name).queryKey,
          (old: ApiKeyI[] = []) => [response.data, ...old],
        )
        setIsCreateOpen(false)
        setNewlyCreatedKey(response.data)
        toast.success("API Key created successfully")
      }
    },
  })

  const { mutate: revokeApiKey } = useMutation({
    mutationFn: (id: string) => revokeApiKeyOnWorkspaceFn(id, name),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allWorkspaceApiKeysQueryOptions(name).queryKey,
          (old: ApiKeyI[] = []) =>
            old.map((ws) =>
              ws.id === revokeKeyId
                ? { ...ws, revoked_at: new Date().toISOString() }
                : ws
            )
        );
        setRevokeKeyId(null)
        toast.success(response.message)
      }
    },
  })

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div className="space-y-1">
          <h2 className="text-2xl md:text-3xl font-black uppercase tracking-tighter">API Keys</h2>
          <p className="text-muted-foreground text-sm uppercase tracking-wider font-bold opacity-70">Programmatic access for your workspace.</p>
        </div>

        <Button
          onClick={() => setIsCreateOpen(true)}
          className="rounded-none gap-2 h-10 font-black uppercase tracking-widest transition-all"
        >
          <Plus className="w-4 h-4" />
          New Key
        </Button>
      </div>

      <KeyList
        keys={keys}
        isLoading={isLoading}
        onRevoke={setRevokeKeyId}
      />

      {/* Create Modal */}
      <FormModal<ApiKeyCreateI>
        title="Create API Key"
        description="Give your key a name to identify it later."
        buttonTitle="Generate Key"
        schema={apiKeyCreateSchema}
        formId="create-key-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={createApiKey}
        fields={[
          {
            name: "name",
            label: "e.g. Production Mobile App",
            type: "text",
          }
        ]}
        disabled={isPendingCreate}
      />

      {/* New Key Result Modal */}
      <ApiKeyCreatedModal
        apiKey={newlyCreatedKey}
        isOpen={!!newlyCreatedKey}
        onClose={() => setNewlyCreatedKey(null)}
      />

      {/* Revoke Confirmation Modal */}
      {revokeKeyId && <ConfirmModal
        isOpen={!!revokeKeyId}
        onClose={() => setRevokeKeyId(null)}
        onConfirm={() => revokeApiKey(revokeKeyId)}
        title="Revoke API Key"
        description="Are you sure you want to revoke this API key? This action will immediately invalidate the key and cannot be undone."
        confirmText="Revoke Key"
        variant="destructive"
      />}
    </div>
  )
}