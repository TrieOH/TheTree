import { addMemberToFormFn, allFormsMembersQueryOptions, removeMemberFromFormFn } from '#/features/forms/api/member'
import { memberAddToFormSchema } from '#/features/forms/model/member'
import type { FormMemberI, MemberAddToFormI } from '#/features/forms/model/member';
import { MemberCard } from '#/features/forms/ui/member-card'
import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import { Button } from '#/shared/ui/shadcn/button'
import FormModal from '#/widgets/modal/form-modal'
import { ConfirmModal } from '#/widgets/modal/modal'
import { PaginatedContainer } from '#/widgets/pagination/paginated-container-grid'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { useMemo, useState } from 'react'
import { toast } from 'sonner'
import FormAdminHeader from '#/features/forms/ui/form-admin-header'
import {
  allNamespacesFormsQueryOptions,
  formResponseCountOnNamespaceQueryOptions,
} from '#/features/namespaces/api'
import {
  allUserFormsQueryOptions,
  formResponseCountQueryOptions,
} from '#/features/forms/api'
import type { FormI } from '#/features/forms/model'

export const Route = createFileRoute('/admin/form/$formID/members')({
  component: RouteComponent,
})

function RouteComponent() {
  const queryClient = useQueryClient()
  const auth = Route.useRouteContext().auth?.auth
  const { formID } = Route.useParams()
  const { namespaceID } = Route.useSearch()
  const user_id = auth?.profile()?.id || null
  const { data: members = [] } = useQuery(allFormsMembersQueryOptions(formID, namespaceID))

  const formsQueryKey = useMemo(() =>
    namespaceID
      ? allNamespacesFormsQueryOptions(namespaceID).queryKey
      : allUserFormsQueryOptions().queryKey,
    [namespaceID]
  )

  const formsQuery = useQuery({
    queryKey: formsQueryKey,
    queryFn: namespaceID
      ? allNamespacesFormsQueryOptions(namespaceID).queryFn
      : allUserFormsQueryOptions().queryFn,
  })

  const form = useMemo(() =>
    formsQuery.data?.find((f) => f.id === formID),
    [formsQuery.data, formID]
  )

  const countQuery = useQuery(
    namespaceID
      ? formResponseCountOnNamespaceQueryOptions(namespaceID, formID)
      : formResponseCountQueryOptions(formID)
  )

  const responseCount = countQuery.data?.count ?? 0

  const updateFormData = (updatedForm: FormI) => {
    queryClient.setQueryData(
      formsQueryKey,
      (oldData: FormI[] = []) => oldData.map(f => f.id === formID ? updatedForm : f)
    )
  }

  const [filter, setFilter] = useState('')
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [memberToRemove, setMemberToRemove] = useState<FormMemberI | null>(null)

  const count = members.length
  const header = useMemo(() => {
    if (!form) return null

    return (
      <FormAdminHeader
        title="Members"
        description={count === 0 ? 'No members yet in this form' : `${count} member${count !== 1 ? 's' : ''} in this form`}
        form={form}
        namespaceID={namespaceID}
        responseCount={responseCount}
        onUpdate={updateFormData}
      />
    )
  }, [form, count, responseCount, namespaceID])

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

  const { mutate: addMemberToForm, isPending: isCreating } = useMutation({
    mutationFn: (data: MemberAddToFormI) => addMemberToFormFn(data, formID, namespaceID),
    onSuccess: (response, variable) => {
      if (response.success) {
        if (user_id) {
          queryClient.setQueryData(
            allFormsMembersQueryOptions(formID, namespaceID).queryKey,
            (oldData: FormMemberI[] = []) => [...oldData, {
              user_id: variable.user_id,
              role: variable.role,
              added_by: user_id,
              added_at: new Date().toISOString(),
              form_id: formID
            }]
          )
        } else {
          queryClient.invalidateQueries(
            { queryKey: allFormsMembersQueryOptions(formID, namespaceID).queryKey }
          )
        }

        setIsCreateOpen(false)
        toast.success(response.message || "Member added successfully")
      } else toast.error(response.message || "Failed to add member")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  const { mutate: removeMemberFromForm, isPending: isRemoving } = useMutation({
    mutationFn: (rm_user_id: string) => removeMemberFromFormFn(rm_user_id, formID, namespaceID),
    onSuccess: (response, rm_user_id) => {
      if (response.success) {
        queryClient.setQueryData(
          allFormsMembersQueryOptions(formID, namespaceID).queryKey,
          (oldData: FormMemberI[] = []) => oldData.filter(member => member.user_id !== rm_user_id)
        )
        toast.success("Member removed successfully")
      } else toast.error(response.message || "Failed to remove member")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <div>
      <PaginatedContainer<FormMemberI>
        items={filteredMembers}
        layout='list'
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
      <FormModal<MemberAddToFormI>
        title="Add Member"
        description="Invite a user to join this form."
        buttonTitle="Add Member"
        schema={memberAddToFormSchema}
        formId="add-member-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={addMemberToForm}
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
        description="Are you sure you want to remove this member from the form? This action cannot be undone."
        confirmText="Remove Member"
        variant='destructive'
        isOpen={memberToRemove !== null}
        onClose={() => setMemberToRemove(null)}
        onConfirm={() => {
          if (memberToRemove) {
            removeMemberFromForm(memberToRemove.user_id)
            setMemberToRemove(null)
          }
        }}
        isLoading={isRemoving}
      />
    </div>
  )
}
