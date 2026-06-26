import { allApiKeysQueryOptions, revokeApiKeyFn, rotateApiKeyFn } from '@/features/api-keys/api'
import { apiKeyCreateSchema, type ApiKeyCreateI, type ApiKeyI, type CreateApiKeyResponseI } from '@/features/api-keys/model'
import { ApiKeyCard } from '@/features/api-keys/ui/api-key-card'
import { ApiKeyCreatedDisplay } from '@/features/api-keys/ui/api-key-created-display'
import { useLayoutHeader } from '@/shared/lib/hooks/layout-context'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'
import { EmptyState } from '@/shared/ui/placeholders/EmptyState'
import { FormModal } from '@/widgets/modal/FormModal'
import { Modal } from '@/widgets/modal/modal'
import { PaginatedContainer } from '@/widgets/pagination/PaginatedContainer'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { KeySquare, Plus } from 'lucide-react'
import { useMemo, useState } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/projects/$projectID/')({
  component: RouteComponent,
})

// Mock data for development
const MOCK_API_KEYS: ApiKeyI[] = [
  {
    id: 'ak_01j2x',
    actor_id: 'user_001',
    project_id: 'proj_001',
    name: 'Production API Key',
    key_prefix: 'trieoh_pk_prod',
    key_hash: 'hash_abc123',
    metadata: null,
    expires_at: '2027-06-26T00:00:00Z',
    last_used_at: '2026-06-25T14:30:00Z',
    created_at: '2026-01-15T10:00:00Z',
  },
  {
    id: 'ak_02j3y',
    actor_id: 'user_001',
    project_id: 'proj_001',
    name: 'Staging API Key',
    key_prefix: 'trieoh_pk_stag',
    key_hash: 'hash_def456',
    metadata: null,
    revoked_at: '2026-05-20T08:00:00Z',
    created_at: '2026-02-10T12:00:00Z',
  },
  {
    id: 'ak_03j4z',
    actor_id: 'user_002',
    project_id: 'proj_001',
    name: 'Development Key',
    key_prefix: 'trieoh_pk_dev',
    key_hash: 'hash_ghi789',
    metadata: null,
    expires_at: '2026-12-31T23:59:59Z',
    last_used_at: '2026-06-26T09:15:00Z',
    created_at: '2026-03-20T16:00:00Z',
  },
  {
    id: 'ak_04j5a',
    actor_id: 'user_002',
    project_id: 'proj_001',
    name: 'CI/CD Pipeline Key',
    key_prefix: 'trieoh_pk_cicd',
    key_hash: 'hash_jkl012',
    metadata: null,
    created_at: '2026-04-05T09:00:00Z',
  },
]

function RouteComponent() {
  const queryClient = useQueryClient()
  const { projectID } = Route.useParams()

  // const { data: apiKeys = [] } = useQuery(allApiKeysQueryOptions(projectID))
  const { data: apiKeys = MOCK_API_KEYS } = useQuery({
    ...allApiKeysQueryOptions(projectID),
    enabled: false, // Disabled while using mock data
    initialData: MOCK_API_KEYS,
  })

  const [filter, setFilter] = useState('')
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [createdKey, setCreatedKey] = useState<CreateApiKeyResponseI | null>(null)
  const [keyToRevoke, setKeyToRevoke] = useState<ApiKeyI | null>(null)

  const filteredApiKeys = apiKeys.filter((key) => {
    const search = filter.toLowerCase().trim()
    if (!search) return true
    return (
      key.name.toLowerCase().includes(search) ||
      key.key_prefix.toLowerCase().includes(search)
    )
  })

  const count = apiKeys.filter(k => !k.revoked_at).length

  const header = useMemo(() => (
    <div className="flex items-start justify-between">
      <div>
        <h1 className="text-lg font-semibold tracking-tight">Api Keys</h1>
        <p className="text-sm text-muted-foreground">
          {count === 0
            ? 'No API keys yet for this project'
            : `${count} active API key${count !== 1 ? 's' : ''} in this project`}
        </p>
      </div>
    </div>
  ), [count])

  useLayoutHeader(header)

  const { mutate: createApiKey, isPending: isCreating } = useMutation({
    mutationFn: (data: ApiKeyCreateI) => rotateApiKeyFn(projectID, data),
    onSuccess: (response) => {
      if (response.success) {
        setCreatedKey(response.data)
        queryClient.invalidateQueries({
          queryKey: allApiKeysQueryOptions(projectID).queryKey,
        })
        toast.success(response.message || "API key created successfully")
      } else {
        toast.error(response.message || "Failed to create API key")
      }
    },
    onError: (error: Error) => toast.error(error.message),
  })

  const { mutate: revokeApiKey, isPending: isRevoking } = useMutation({
    mutationFn: (key_id: string) => revokeApiKeyFn(projectID, key_id),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({
          queryKey: allApiKeysQueryOptions(projectID).queryKey,
        })
        setKeyToRevoke(null)
        toast.success(response.message || "API key revoked successfully")
      } else toast.error(response.message || "Failed to revoke API key")
    },
    onError: (error: Error) => toast.error(error.message),
  })

  return (
    <div>
      <PaginatedContainer<ApiKeyI>
        items={filteredApiKeys}
        layout="list"
        pageSize={10}
        sortFields={[
          { key: "name", label: "Name" },
          { key: "created_at", label: "Created At" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by name or key prefix…"
        itemLabel="API keys"
        headerActions={
          <ShadowButton
            onClick={() => setIsCreateOpen(true)}
            variant="outline"
            className="h-9 sm:w-auto px-3 rounded-sm"
            leftIcon={<Plus size={16} />}
            value="Create API Key"
          />
        }
        renderItems={(slice) => slice.map(item => {
          return (
            <ApiKeyCard key={item.id} data={item} onRevoke={setKeyToRevoke} />
          )
        })}
        emptyState={
          <EmptyState
            icon={KeySquare}
            title="No API keys"
            description="No API keys found for this project. Create one to get started."
            action={
              <ShadowButton
                onClick={() => setIsCreateOpen(true)}
                variant="default"
                className="px-4 rounded-sm"
                leftIcon={<Plus size={16} />}
                value="Create API Key"
              />
            }
          />
        }
      />

      {/* Create API Key Modal */}
      <FormModal<ApiKeyCreateI>
        title="Create API Key"
        description="Generate a new API key for this project."
        submitLabel="Create Key"
        schema={apiKeyCreateSchema}
        formId="create-api-key-form"
        isOpen={isCreateOpen && !createdKey}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={createApiKey}
        defaultValues={{ name: '', create_for_service_account: false, expires_at: undefined }}
        isLoading={isCreating}
        fields={[
          {
            name: 'name',
            label: 'Key Name',
            type: 'text',
            placeholder: 'e.g. Production Key',
            required: true,
          },
          {
            name: 'create_for_service_account',
            label: 'Create for service account',
            type: 'option-picker',
            options: [
              { value: 'false', label: 'No - Personal use' },
              { value: 'true', label: 'Yes - Service account' },
            ],
            required: true,
          },
          {
            name: 'expires_at',
            label: 'Expires At (optional)',
            type: 'text',
            placeholder: 'e.g. 2027-12-31T23:59:59Z',
          },
        ]}
      />

      {/* Created Key Display Modal */}
      <Modal
        isOpen={createdKey !== null}
        onClose={() => {
          setCreatedKey(null)
          setIsCreateOpen(false)
        }}
        title="API Key Created"
        description="Your new API key has been generated."
      >
        {createdKey && (
          <ApiKeyCreatedDisplay
            name={createdKey.key?.name ?? 'API Key'}
            rawKey={createdKey.raw_key}
            onClose={() => {
              setCreatedKey(null)
              setIsCreateOpen(false)
            }}
          />
        )}
      </Modal>

      {/* Revoke Confirmation Modal */}
      <Modal
        isOpen={keyToRevoke !== null}
        onClose={() => setKeyToRevoke(null)}
        title="Revoke API Key"
        description="Are you sure you want to revoke this API key? Any services using this key will immediately lose access. This action cannot be undone."
        footer={
          <div className="flex flex-col-reverse sm:flex-row justify-end gap-2 w-full">
            <ShadowButton
              variant="ghost"
              onClick={() => setKeyToRevoke(null)}
              className="rounded-sm font-medium text-xs"
              disabled={isRevoking}
              value="Cancel"
            />
            <ShadowButton
              variant="destructive"
              onClick={() => {
                if (keyToRevoke) {
                  revokeApiKey(keyToRevoke.id)
                }
              }}
              className="rounded-sm font-bold text-xs px-6"
              disabled={isRevoking}
              value={isRevoking ? 'Revoking...' : 'Revoke Key'}
            />
          </div>
        }
      >
        <div className="space-y-4">
          <div className="text-xs text-muted-foreground">
            Revoking key: <span className="font-semibold text-foreground">{keyToRevoke?.name}</span>
          </div>
          <div className="flex items-center gap-2 p-3 rounded-sm bg-muted border border-border">
            <span className="text-xs font-mono text-muted-foreground">
              {keyToRevoke?.key_prefix}...
            </span>
          </div>
        </div>
      </Modal>
    </div>
  )
}