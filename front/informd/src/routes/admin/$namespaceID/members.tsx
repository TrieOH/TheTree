import {
  addMemberToNamespaceFn,
  allNamespacesMembersQueryOptions,
  removeMemberFromNamespaceFn
} from '#/features/namespaces/api'
import { memberAddToNamespaceSchema } from '#/features/namespaces/model'
import type { MemberAddToNamespaceI, NamespaceMemberI } from '#/features/namespaces/model';
import { MemberCard } from '#/features/namespaces/ui/member-card'
import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import { Button } from '#/shared/ui/shadcn/button'
import FormModal from '#/widgets/modal/form-modal'
import { ConfirmModal } from '#/widgets/modal/modal';
import { PaginatedContainer } from '#/widgets/pagination/paginated-container-grid'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { useMemo, useState } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/$namespaceID/members')({
  component: RouteComponent,
})

function RouteComponent() {
  const queryClient = useQueryClient()
  const auth = Route.useRouteContext().auth?.auth
  const user_id = auth?.profile()?.id || null
  const { namespaceID } = Route.useParams()
  const { data: members = [] } = useQuery(allNamespacesMembersQueryOptions(namespaceID))

  const [filter, setFilter] = useState('')
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [memberToRemove, setMemberToRemove] = useState<NamespaceMemberI | null>(null)

  const count = members.length

  const header = useMemo(() => (
    <div className="flex items-start justify-between">
      <div>
        <h1 className="text-lg font-semibold tracking-tight">Members</h1>
        <p className="text-sm text-muted-foreground">
          {count === 0
            ? 'No members yet in this namespace'
            : `${count} member${count !== 1 ? 's' : ''} in this namespace`}
        </p>
      </div>
    </div>
  ), [count])

  useLayoutHeader(header)

  const filteredMembers = members.filter((member) => {
    const search = filter.toLowerCase().trim()

    if (!search) return true

    return (
      member.role.toLowerCase().includes(search) ||
      member.user_id.toLowerCase().includes(search) ||
      member.added_by.toLowerCase().includes(search)
    )
  })

  const { mutate: addMemberToNamespace, isPending: isCreating } = useMutation({
    mutationFn: (data: MemberAddToNamespaceI) => addMemberToNamespaceFn(namespaceID, data),
    onSuccess: (response, variable) => {
      if (response.success) {
        if (user_id) {
          queryClient.setQueryData(
            allNamespacesMembersQueryOptions(namespaceID).queryKey,
            (oldData: NamespaceMemberI[] = []) => [...oldData, {
              user_id: variable.user_id,
              role: variable.role,
              added_by: user_id,
              added_at: new Date().toISOString(),
              namespace_id: namespaceID,
            }]
          )
        } else {
          queryClient.invalidateQueries(
            { queryKey: allNamespacesMembersQueryOptions(namespaceID).queryKey }
          )
        }

        setIsCreateOpen(false)
        toast.success(response.message || "Member added successfully")
      } else toast.error(response.message || "Failed to add member")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  const { mutate: removeMemberFromNamespace, isPending: isRemoving } = useMutation({
    mutationFn: (user_id: string) => removeMemberFromNamespaceFn(namespaceID, user_id),
    onSuccess: (response, user_id) => {
      if (response.success) {
        queryClient.setQueryData(
          allNamespacesMembersQueryOptions(namespaceID).queryKey,
          (oldData: NamespaceMemberI[] = []) => oldData.filter(member => member.user_id !== user_id)
        )
        toast.success("Member removed successfully")
      } else toast.error(response.message || "Failed to remove member")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <main>
      <PaginatedContainer<NamespaceMemberI>
        items={filteredMembers}
        className='w-full'
        layout='flex'
        pageSize={10}
        sortFields={[
          { key: "role", label: "Role" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by role, user or added by…"
        itemLabel="members"
        headerActions={
          <Button
            onClick={() => setIsCreateOpen(true)}
            size="icon"
            variant="outline"
            className="sm:w-auto px-3 rounded-sm"
          >
            <Plus size={16} />
            <span className="hidden sm:inline ml-2">Add Member</span>
          </Button>
        }
        renderItems={(slice) => slice.map(item => {
          return (
            <MemberCard key={item.user_id} data={item} onRemove={setMemberToRemove} />
          )
        })}
      />
      <FormModal<MemberAddToNamespaceI>
        title="Add Member"
        description="Invite a user to join this namespace."
        buttonTitle="Add Member"
        schema={memberAddToNamespaceSchema}
        formId="add-member-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={addMemberToNamespace}
        fields={[
          {
            name: 'user_id',
            label: 'User ID',
            type: 'text',
            placeholder: 'e.g. user123',
          },
          {
            name: 'role',
            label: 'Role',
            type: 'select',
            options: [
              { value: 'viewer', label: 'Viewer' },
              { value: 'editor', label: 'Editor' },
              { value: 'admin', label: 'Admin' },
              { value: 'owner', label: 'Owner' },
            ],
          }
        ]}
        disabled={isCreating}
      />
      <ConfirmModal
        title="Remove Member"
        description="Are you sure you want to remove this member from the namespace? This action cannot be undone."
        confirmText="Remove Member"
        variant='destructive'
        isOpen={memberToRemove !== null}
        onClose={() => setMemberToRemove(null)}
        onConfirm={() => {
          if (memberToRemove) {
            removeMemberFromNamespace(memberToRemove.user_id)
            setMemberToRemove(null)
          }
        }}
        isLoading={isRemoving}
      />
    </main>
  )
}
