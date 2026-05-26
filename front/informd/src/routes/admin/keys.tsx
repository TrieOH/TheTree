import { allApiKeysQueryOptions, createApiKeyFn, revokeApiKeyFn } from '#/features/keys/api'
import { apiKeyCreateSchema } from '#/features/keys/model'
import type { ApiKeyCreateI, ApiKeyCreateResponseI, ApiKeyI } from '#/features/keys/model';
import { ApiKeyCreatedModal } from '#/features/keys/ui/api-key-created-modal'
import { APIKeyCard } from '#/features/keys/ui/key-card'
import type { FieldDefinition } from '#/shared/model/form-types'
import { Button } from '#/shared/ui/shadcn/button'
import FormModal from '#/widgets/modal/form-modal'
import { ConfirmModal } from '#/widgets/modal/modal'
import { PaginatedContainer } from '#/widgets/pagination/paginated-container-grid'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { useState } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/keys')({
  component: RouteComponent,
})

const KEYS_FIELDS: FieldDefinition<ApiKeyCreateI>[] = [
  {
    name: 'name',
    label: 'API Key Name',
    placeholder: 'Enter API key name...',
    type: 'text',
  },
];

const MOCK_API_KEYS: ApiKeyI[] = [
  {
    id: '1',
    name: 'My First API Key',
    prefix: 'abc123',
    revoked_at: undefined,
    created_at: new Date().toISOString(),
  },
  {
    id: '2',
    name: 'Old API Key',
    prefix: 'def456',
    revoked_at: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(), // revoked 1 day ago
    created_at: new Date(Date.now() - 1000 * 60 * 60 * 24 * 30).toISOString(), // created 30 days ago
  },
];

function RouteComponent() {
  const queryClient = useQueryClient()
  const [filter, setFilter] = useState('')
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [keyToRemove, setKeyToRemove] = useState<ApiKeyI | null>(null)
  const [newlyCreatedKey, setNewlyCreatedKey] = useState<ApiKeyCreateResponseI | null>(null)

  const filteredApiKeys = MOCK_API_KEYS.filter((key) => key.name.toLowerCase().includes(filter.toLowerCase()))

  const { mutate: createApiKey, isPending: isCreating } = useMutation({
    mutationFn: (data: ApiKeyCreateI) => createApiKeyFn(data),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allApiKeysQueryOptions().queryKey,
          (oldData: ApiKeyI[] = []) => [response.data.APIKeyResponse, ...oldData]
        );
        setIsCreateOpen(false)
        setNewlyCreatedKey(response.data)
        toast.success(response.message || "API key created successfully")
      } else toast.error(response.message || "Failed to create API key")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  const { mutate: revokeApiKey, isPending: isRemoving } = useMutation({
    mutationFn: (id: string) => revokeApiKeyFn(id),
    onSuccess: (response, id) => {
      if (response.success) {
        queryClient.setQueryData(
          allApiKeysQueryOptions().queryKey,
          (old: ApiKeyI[] = []) =>
            old.filter((key) => key.id !== id)
        );
        setKeyToRemove(null)
        toast.success(response.message)
      } else toast.error(response.message)
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <div className='flex flex-wrap p-4'>
      <PaginatedContainer<ApiKeyI>
        items={filteredApiKeys}
        layout='list'
        pageSize={10}
        sortFields={[
          { key: "name", label: "Name" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by name…"
        itemLabel="api keys"
        headerActions={
          <Button
            onClick={() => setIsCreateOpen(true)}
            size="icon"
            variant="outline"
            className="sm:w-auto px-3 rounded-sm"
          >
            <Plus size={16} />
            <span className="hidden sm:inline ml-2">Add API Key</span>
          </Button>
        }
        renderItems={(slice) => slice.map(item => {
          return (
            <APIKeyCard key={item.id} data={item} onRevoke={setKeyToRemove} />
          )
        })}
      />
      <FormModal<ApiKeyCreateI>
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        title="Create API Key"
        description="API keys allow you to authenticate and authorize requests to the API."
        formId="create-api-key-form"
        buttonTitle="Create API Key"
        fields={KEYS_FIELDS}
        schema={apiKeyCreateSchema}
        onSubmit={createApiKey}
        disabled={isCreating}
      />

      {/* New Key Result Modal */}
      <ApiKeyCreatedModal
        apiKey={newlyCreatedKey}
        isOpen={!!newlyCreatedKey}
        onClose={() => setNewlyCreatedKey(null)}
      />

      {/* Revoke Confirmation Modal */}
      <ConfirmModal
        title="Revoke API Key"
        description="Are you sure you want to revoke this API key? This action will immediately invalidate the key and cannot be undone."
        confirmText="Revoke Key"
        variant='destructive'
        isOpen={keyToRemove !== null}
        onClose={() => setKeyToRemove(null)}
        onConfirm={() => {
          if (keyToRemove) {
            revokeApiKey(keyToRemove.id)
            setKeyToRemove(null)
          }
        }}
        isLoading={isRemoving}
      />
    </div>
  )
}
