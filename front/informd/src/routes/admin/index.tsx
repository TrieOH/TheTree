import { allNamespacesQueryOptions, createNamespaceFn } from '#/features/namespaces/api';
import { namespaceCreateSchema } from '#/features/namespaces/model'
import type { NamespaceI, NamespaceCreateI } from '#/features/namespaces/model';
import NamespaceList from '#/features/namespaces/ui/namespace-list'
import FormModal from '#/widgets/modal/form-modal'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { toast } from 'sonner'
import { Button } from '#/shared/ui/shadcn/button';
import { UserPlus, Plus } from 'lucide-react';
import { promoteToClientSchema } from '#/features/permissions/model';
import type { PromoteToClientI } from '#/features/permissions/model';
import { checkSuperAdminPrivilegesFn, promoteUserToClientFn } from '#/features/permissions/api';
import { allUserFormsQueryOptions, createFormOnUserFn } from '#/features/forms/api'
import { formCreateSchema } from '#/features/forms/model'
import type { FormCreateI, FormI } from '#/features/forms/model'
import { FormList } from '#/features/forms/ui/form-list'
import { allUserApiKeysQueryOptions, createApiKeyFn, revokeApiKeyFn } from '#/features/keys/api'
import type { ApiKeyCreateI, ApiKeyCreateResponseI, ApiKeyI } from '#/features/keys/model'
import { apiKeyCreateSchema } from '#/features/keys/model'
import { KeyList } from '#/features/keys/ui/key-list'
import { ApiKeyCreatedModal } from '#/features/keys/ui/api-key-created-modal'
import { ConfirmModal } from '#/widgets/modal/modal'
import { cn } from '#/shared/lib/utils'

export const Route = createFileRoute('/admin/')({
  component: RouteComponent,
})

type TabType = 'namespaces' | 'forms' | 'keys'

function RouteComponent() {
  const { auth } = Route.useRouteContext()
  const userId = auth?.auth.profile()?.id || ''

  const [activeTab, setActiveTab] = useState<TabType>('namespaces')
  const [isCreateNamespaceOpen, setIsCreateNamespaceOpen] = useState(false)
  const [isCreateFormOpen, setIsCreateFormOpen] = useState(false)
  const [isCreateKeyOpen, setIsCreateKeyOpen] = useState(false)
  const [isPromoteOpen, setIsPromoteOpen] = useState(false)
  const [revokeKeyId, setRevokeKeyId] = useState<string | null>(null)
  const [newlyCreatedKey, setNewlyCreatedKey] = useState<ApiKeyCreateResponseI | null>(null)

  const queryClient = useQueryClient();

  const { data: namespaces = [], isLoading: isLoadingNamespaces } = useQuery(allNamespacesQueryOptions(userId))
  const { data: forms = [], isLoading: isLoadingForms } = useQuery(allUserFormsQueryOptions(userId))
  const { data: keys = [], isLoading: isLoadingKeys } = useQuery(allUserApiKeysQueryOptions(userId))

  const { data: isAdmin = false } = useQuery({
    queryKey: ['user', userId, 'super_admin'],
    queryFn: () => checkSuperAdminPrivilegesFn({ data: userId }),
    enabled: !!userId,
  })

  const { mutate: createNamespace, isPending: isPendingCreateNamespace } = useMutation({
    mutationFn: createNamespaceFn,
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allNamespacesQueryOptions(userId).queryKey,
          (old: NamespaceI[] = []) => [...old, response.data],
        )
        setIsCreateNamespaceOpen(false)
        toast.success('Namespace created successfully')
      }
    },
    onError: (error: Error) => toast.error(error.message)
  })

  const { mutate: createForm, isPending: isPendingCreateForm } = useMutation({
    mutationFn: (data: FormCreateI) => createFormOnUserFn(data),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allUserFormsQueryOptions(userId).queryKey,
          (old: FormI[] = []) => [...old, response.data],
        )
        setIsCreateFormOpen(false)
        toast.success('Form created successfully')
      }
    },
    onError: (error: Error) => toast.error(error.message)
  })

  const { mutate: createApiKey, isPending: isPendingCreateKey } = useMutation({
    mutationFn: (data: ApiKeyCreateI) => createApiKeyFn(data),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allUserApiKeysQueryOptions(userId).queryKey,
          (old: ApiKeyI[] = []) => [response.data, ...old],
        )
        setIsCreateKeyOpen(false)
        setNewlyCreatedKey(response.data)
        toast.success(response.message || "API Key created successfully")
      }
    },
    onError: (error: Error) => toast.error(error.message)
  })

  const { mutate: revokeApiKey } = useMutation({
    mutationFn: (id: string) => revokeApiKeyFn(id),
    onSuccess: (response, id) => {
      if (response.success) {
        queryClient.setQueryData(
          allUserApiKeysQueryOptions(userId).queryKey,
          (old: ApiKeyI[] = []) =>
            old.filter((key) => key.id !== id)
        );
        setRevokeKeyId(null)
        toast.success(response.message)
      } else toast.error(response.message)
    },
    onError: (error: Error) => toast.error(error.message)
  })

  const { mutate: promoteUser, isPending: isPendingPromote } = useMutation({
    mutationFn: (data: Omit<PromoteToClientI, 'requesterId'>) => promoteUserToClientFn({
      data: {
        userId: data.userId,
        requesterId: userId
      }
    }),
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message)
        setIsPromoteOpen(false)
      } else toast.warning(response.message)
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to promote user')
  })

  const tabs = [
    { id: 'namespaces', label: 'Namespaces', count: namespaces.length },
    { id: 'forms', label: 'Personal Forms', count: forms.length },
    { id: 'keys', label: 'API Keys', count: keys.length },
  ]

  const isLoading = isLoadingNamespaces || isLoadingForms || isLoadingKeys

  if (isLoading) {
    return (
      <div className="space-y-8 animate-in fade-in duration-500">
        <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
          <div className="h-9 w-48 bg-muted animate-pulse rounded-sm" />
          <div className="h-10 w-full sm:w-36 bg-muted animate-pulse rounded-sm" />
        </div>
        <div className="grid gap-6 grid-cols-[repeat(auto-fill,minmax(min(100%,320px),1fr))]">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="h-40 border-2 border-border/50 bg-card/50 animate-pulse rounded-none" />
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-10 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-8 pb-4 border-b-2 border-border/20">
        <div className="flex border-2 border-border bg-muted/10 p-1 self-start overflow-x-auto no-scrollbar max-w-full">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as TabType)}
              className={cn(
                "px-6 py-3 text-[10px] font-black uppercase tracking-[0.2em] transition-all relative shrink-0",
                activeTab === tab.id 
                  ? "bg-primary text-primary-foreground shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] -translate-x-1 -translate-y-1" 
                  : "text-muted-foreground hover:text-foreground hover:bg-muted"
              )}
            >
              {tab.label}
              {tab.count > 0 && (
                <span className={cn(
                  "ml-2 text-[9px] font-bold px-1.5 py-0.5 border",
                  activeTab === tab.id ? "bg-primary-foreground text-primary border-transparent" : "bg-muted border-border"
                )}>
                  {tab.count}
                </span>
              )}
            </button>
          ))}
        </div>

        {isAdmin && (
          <Button
            variant="outline"
            className="rounded-none h-12 gap-2 font-black uppercase tracking-[0.15em] text-[10px] border-2 border-primary/40 hover:bg-primary/5 hover:border-primary transition-all self-start sm:self-auto"
            onClick={() => setIsPromoteOpen(true)}
          >
            <UserPlus className="w-4 h-4" />
            Authorize New Client
          </Button>
        )}
      </div>

      <div className="min-h-100 animate-in fade-in slide-in-from-top-2 duration-500">
        {activeTab === 'namespaces' && (
          <NamespaceList
            openModal={() => setIsCreateNamespaceOpen(true)}
            namespaces={namespaces}
          />
        )}

        {activeTab === 'forms' && (
          <FormList
            forms={forms}
            openModal={() => setIsCreateFormOpen(true)}
          />
        )}

        {activeTab === 'keys' && (
          <div className="space-y-8">
            <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
              <div className="space-y-1">
                <h2 className="text-2xl md:text-3xl font-black uppercase tracking-tighter">
                  API Keys
                </h2>
                <p className="text-muted-foreground text-sm uppercase tracking-wider font-bold opacity-70">
                  Programmatic access for your account.
                </p>
              </div>

              <Button
                onClick={() => setIsCreateKeyOpen(true)}
                className="rounded-none gap-2 h-10 font-black uppercase tracking-widest transition-all"
              >
                <Plus className="w-4 h-4" />
                New Key
              </Button>
            </div>

            <KeyList
              keys={keys}
              isLoading={isLoadingKeys}
              onRevoke={setRevokeKeyId}
            />
          </div>
        )}
      </div>

      {/* Modals */}
      <FormModal<NamespaceCreateI>
        title="Create Namespace"
        description="Give your namespace a name to identify it."
        buttonTitle="Create Project"
        schema={namespaceCreateSchema}
        formId="create-project-form"
        isOpen={isCreateNamespaceOpen}
        onClose={() => setIsCreateNamespaceOpen(false)}
        onSubmit={createNamespace}
        fields={[
          {
            name: "name",
            label: "e.g. My Team Namespace",
            type: "text",
          }
        ]}
        disabled={isPendingCreateNamespace}
      />

      <FormModal<FormCreateI>
        title="Create Form"
        description="Give your personal form a title to identify it."
        buttonTitle="Create Form"
        schema={formCreateSchema}
        formId="create-personal-form"
        isOpen={isCreateFormOpen}
        onClose={() => setIsCreateFormOpen(false)}
        onSubmit={createForm}
        fields={[
          {
            name: 'title',
            label: 'Form Title',
            type: 'text',
            placeholder: 'e.g. Personal Contact Form',
          },
        ]}
        disabled={isPendingCreateForm}
      />

      <FormModal<ApiKeyCreateI>
        title="Create API Key"
        description="Give your key a name to identify it later."
        buttonTitle="Generate Key"
        schema={apiKeyCreateSchema}
        formId="create-key-form"
        isOpen={isCreateKeyOpen}
        onClose={() => setIsCreateKeyOpen(false)}
        onSubmit={createApiKey}
        fields={[
          {
            name: "name",
            label: "e.g. My API Key",
            type: "text",
          }
        ]}
        disabled={isPendingCreateKey}
      />

      <ApiKeyCreatedModal
        apiKey={newlyCreatedKey}
        isOpen={!!newlyCreatedKey}
        onClose={() => setNewlyCreatedKey(null)}
      />

      {revokeKeyId && <ConfirmModal
        isOpen={!!revokeKeyId}
        onClose={() => setRevokeKeyId(null)}
        onConfirm={() => revokeApiKey(revokeKeyId)}
        title="Revoke API Key"
        description="Are you sure you want to revoke this API key? This action will immediately invalidate the key and cannot be undone."
        confirmText="Revoke Key"
        variant="destructive"
      />}

      <FormModal<{ userId: string }>
        title="Authorize Client"
        description="Allow a user to create and manage their own namespaces."
        buttonTitle="Grant Client Access"
        schema={promoteToClientSchema.omit({ requesterId: true })}
        formId="promote-user-form"
        isOpen={isPromoteOpen}
        onClose={() => setIsPromoteOpen(false)}
        onSubmit={promoteUser}
        fields={[
          {
            name: "userId",
            label: "User Identity ID",
            type: "text",
            placeholder: "Enter the user's ID to authorize",
          }
        ]}
        disabled={isPendingPromote}
      />
    </div>
  )
}
