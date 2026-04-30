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
import { UserPlus } from 'lucide-react';
import { promoteToClientSchema } from '#/features/permissions/model';
import type { PromoteToClientI } from '#/features/permissions/model';
import { checkSuperAdminPrivilegesFn, promoteUserToClientFn } from '#/features/permissions/api';

export const Route = createFileRoute('/admin/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { auth } = Route.useRouteContext()
  const userId = auth?.auth.profile()?.id || ''

  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [isPromoteOpen, setIsPromoteOpen] = useState(false)
  const queryClient = useQueryClient();

  const { data: namespaces = [], isLoading } = useQuery(allNamespacesQueryOptions(userId))

  const { data: isAdmin = false } = useQuery({
    queryKey: ['user', userId, 'super_admin'],
    queryFn: () => checkSuperAdminPrivilegesFn({ data: userId }),
    enabled: !!userId,
  })

  const { mutate: createNamespace, isPending: isPendingCreate } = useMutation({
    mutationFn: createNamespaceFn,
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allNamespacesQueryOptions(userId).queryKey,
          (old: NamespaceI[] = []) => [...old, response.data],
        )
        setIsCreateOpen(false)
        toast.success('Namespace created successfully')
      }
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

  if (isLoading) {
    return (
      <div className="space-y-8 animate-in fade-in duration-500">
        <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
          <div className="space-y-1">
            <div className="h-9 w-48 bg-muted animate-pulse rounded-sm" />
            <div className="h-5 w-64 bg-muted animate-pulse rounded-sm" />
          </div>
          <div className="h-10 w-full sm:w-36 bg-muted animate-pulse rounded-sm" />
        </div>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="h-35 border border-border bg-card p-5 space-y-4">
              <div className="space-y-2">
                <div className="h-6 w-3/4 bg-muted animate-pulse rounded-sm" />
                <div className="h-4 w-1/4 bg-muted animate-pulse rounded-sm" />
              </div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      {isAdmin && (
        <div className="flex justify-end">
          <Button
            variant="outline"
            className="rounded-sm gap-2"
            onClick={() => setIsPromoteOpen(true)}
          >
            <UserPlus className="w-4 h-4" />
            Authorize New Client
          </Button>
        </div>
      )}

      <NamespaceList
        openModal={() => setIsCreateOpen(true)}
        namespaces={namespaces}
      />

      <FormModal<NamespaceCreateI>
        title="Create Namespace"
        description="Give your namespace a name to identify it."
        buttonTitle="Create Project"
        schema={namespaceCreateSchema}
        formId="create-project-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={createNamespace}
        fields={[
          {
            name: "name",
            label: "e.g. My Team Namespace",
            type: "text",
          }
        ]}
        disabled={isPendingCreate}
      />

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
