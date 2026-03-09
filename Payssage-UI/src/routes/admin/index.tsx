import { createFileRoute } from '@tanstack/react-router'
import { WorkspaceList } from '#/features/workspaces/ui/workspace-list'
import { allWorkspacesQueryOptions, createWorkspaceFn, disableWorkspaceSandboxModeFn, enableWorkspaceSandboxModeFn } from '#/features/workspaces/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import FormModal from '#/widgets/modal/form-modal';
import { workspaceCreateSchema } from '#/features/workspaces/model';
import { useState } from 'react';
import { toast } from 'sonner';
import type { WorkspaceCreateI, WorkspaceI } from '#/features/workspaces/model';

export const Route = createFileRoute('/admin/')({
  component: RouteComponent,
})


function RouteComponent() {
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const queryClient = useQueryClient();
  const { data: workspaces = [], isLoading } = useQuery(allWorkspacesQueryOptions())

  const { mutate: createWorkspace, isPending: isPendingCreate } = useMutation({
    mutationFn: createWorkspaceFn,
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allWorkspacesQueryOptions().queryKey,
          (old: WorkspaceI[] = []) => [...old, response.data],
        )
        setIsCreateOpen(false)
        toast.success('Workspace created successfully')
      }
    },
  })

  const { mutate: enableWorkspace } = useMutation({
    mutationFn: enableWorkspaceSandboxModeFn,
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allWorkspacesQueryOptions().queryKey,
          (old: WorkspaceI[] = []) =>
            old.map((ws) => (ws.id === response.data.id ? response.data : ws))
        )
        setIsCreateOpen(false)
        toast.success('Workspace is ready to production')
      }
    },
  })

  const { mutate: disableWorkspace } = useMutation({
    mutationFn: disableWorkspaceSandboxModeFn,
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allWorkspacesQueryOptions().queryKey,
          (old: WorkspaceI[] = []) =>
            old.map((ws) => (ws.id === response.data.id ? response.data : ws))
        )
        setIsCreateOpen(false)
        toast.success('Workspace is ready to test')
      }
    },
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
              <div className="flex justify-between items-center mt-auto">
                <div className="h-3 w-20 bg-muted animate-pulse rounded-sm" />
                <div className="h-3 w-24 bg-muted animate-pulse rounded-sm" />
              </div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-700">
      <WorkspaceList
        workspaces={workspaces}
        openModal={() => setIsCreateOpen(true)}
        handleEnableSandbox={enableWorkspace}
        handleDisableSandbox={disableWorkspace}
      />
      <FormModal<WorkspaceCreateI>
        title="Create Workspace"
        description="Give your workspace a name to identify it."
        buttonTitle="Create Workspace"
        schema={workspaceCreateSchema}
        formId="create-workspace-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={createWorkspace}
        fields={[
          {
            name: "name",
            label: "e.g. My Team Workspace",
            type: "text",
          }
        ]}
        disabled={isPendingCreate}
      />
    </div>
  )
}
