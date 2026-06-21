import { addMemberToProjectFn, allProjectMembersQueryOptions, removeMemberFromProjectFn } from '@/features/project/api'
import { memberAddToProjectSchema, type MemberAddToProjectI, type ProjectMemberI } from '@/features/project/model'
import { MemberCard } from '@/features/project/ui/member-card'
import { useLayoutHeader } from '@/shared/lib/hooks/layout-context'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'
import { FormModal } from '@/widgets/modal/FormModal'
import { Modal } from '@/widgets/modal/modal'
import { PaginatedContainer } from '@/widgets/pagination/PaginatedContainer'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Crown, Plus, Shield, User2 } from 'lucide-react'
import { useId, useMemo, useState } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/projects/$projectID/members')({
  component: RouteComponent,
})

function RouteComponent() {
  const queryClient = useQueryClient()
  const { projectID } = Route.useParams()
  const { organizationID } = Route.useSearch()

  const { data: members = [] } = useQuery(allProjectMembersQueryOptions(projectID, organizationID))

  const [filter, setFilter] = useState('')
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [memberToRemove, setMemberToRemove] = useState<ProjectMemberI | null>(null)
  const [confirmEmail, setConfirmEmail] = useState('')
  const confirmEmailInputId = useId()

  const count = members.length

  const header = useMemo(() => (
    <div className="flex items-start justify-between">
      <div>
        <h1 className="text-lg font-semibold tracking-tight">Members</h1>
        <p className="text-sm text-muted-foreground">
          {count === 0
            ? 'No members yet in this organization'
            : `${count} member${count !== 1 ? 's' : ''} in this organization`}
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
      member.actor_id.toLowerCase().includes(search)
    )
  })

  const { mutate: addMemberToProject, isPending: isCreating } = useMutation({
    mutationFn: (data: MemberAddToProjectI) => addMemberToProjectFn(projectID, data, organizationID),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({
          queryKey: allProjectMembersQueryOptions(projectID, organizationID).queryKey
        })
        setIsCreateOpen(false)
        toast.success(response.message || "Member added successfully")
      } else toast.error(response.message || "Failed to add member")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  const { mutate: removeMemberFromProject, isPending: isRemoving } = useMutation({
    mutationFn: (email: string) => removeMemberFromProjectFn(projectID, email, organizationID),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({
          queryKey: allProjectMembersQueryOptions(projectID, organizationID).queryKey
        })
        setMemberToRemove(null)
        setConfirmEmail('')
        toast.success("Member removed successfully")
      } else toast.error(response.message || "Failed to remove member")
    },
    onError: (error: Error) => toast.error(error.message)
  })


  return (
    <div>
      <PaginatedContainer<ProjectMemberI>
        items={filteredMembers}
        layout="list"
        pageSize={10}
        sortFields={[
          { key: "role", label: "Role" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by role or actor…"
        itemLabel="members"
        headerActions={
          <ShadowButton
            onClick={() => setIsCreateOpen(true)}
            variant="outline"
            className="h-9 sm:w-auto px-3 rounded-sm"
            leftIcon={<Plus size={16} />}
            value="Add Member"
          />
        }
        renderItems={(slice) => slice.map(item => {
          return (
            <MemberCard key={item.actor_id} data={item} onRemove={setMemberToRemove} />
          )
        })}
      />

      <FormModal<MemberAddToProjectI>
        title="Add Member"
        description="Invite a user to join this project."
        submitLabel="Add Member"
        schema={memberAddToProjectSchema}
        formId="add-member-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={addMemberToProject}
        defaultValues={{ actor_email: '', role: 'member' }}
        isLoading={isCreating}
        fields={[
          {
            name: 'actor_email',
            label: 'Actor Email',
            type: 'text',
            placeholder: 'e.g. member@example.com',
            required: true,
          },
          {
            name: 'role',
            label: 'Role',
            type: 'option-picker',
            options: [
              { value: 'member', label: 'Member', icon: User2 },
              { value: 'admin', label: 'Admin', icon: Shield },
              { value: 'owner', label: 'Owner', icon: Crown },
            ],
            required: true,
          }
        ]}
      />

      <Modal
        isOpen={memberToRemove !== null}
        onClose={() => {
          setMemberToRemove(null)
          setConfirmEmail('')
        }}
        title="Remove Member"
        description="To remove this member from the organization, please confirm by typing their email address. This action cannot be undone."
        footer={
          <div className="flex flex-col-reverse sm:flex-row justify-end gap-2 w-full">
            <ShadowButton
              variant="ghost"
              onClick={() => {
                setMemberToRemove(null)
                setConfirmEmail('')
              }}
              className="rounded-sm font-medium text-xs"
              disabled={isRemoving}
              value="Cancel"
            />
            <ShadowButton
              variant="destructive"
              onClick={() => {
                if (confirmEmail.trim()) {
                  removeMemberFromProject(confirmEmail.trim())
                }
              }}
              className="rounded-sm font-bold text-xs px-6"
              disabled={isRemoving || !confirmEmail.trim()}
              value={isRemoving ? 'Removing...' : 'Remove Member'}
            />
          </div>
        }
      >
        <div className="space-y-4">
          <div className="text-xs text-muted-foreground">
            Removing member: <span className="font-semibold text-foreground">{memberToRemove?.actor_id}</span>
          </div>
          <div className="space-y-2">
            <label htmlFor={confirmEmailInputId} className="text-xs font-medium text-foreground">
              Member Email Address
            </label>
            <input
              id={confirmEmailInputId}
              type="email"
              value={confirmEmail}
              onChange={(e) => setConfirmEmail(e.target.value)}
              placeholder="e.g. member@example.com"
              className="flex h-10 w-full rounded-sm border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-hidden focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            />
          </div>
        </div>
      </Modal>
    </div>
  )
}
