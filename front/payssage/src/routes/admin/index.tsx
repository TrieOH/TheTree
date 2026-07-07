import { createFileRoute } from '@tanstack/react-router'
import { allOrganizationsQueryOptions, createOrganizationFn } from '#/features/organizations/api'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import FormModal from '#/widgets/modal/form-modal'
import { organizationCreateSchema } from '#/features/organizations/model'
import { useState } from 'react'
import { toast } from 'sonner'
import type { OrganizationCreateI, OrganizationI } from '#/features/organizations/model'
import { Button } from '#/shared/ui/shadcn/button'
import { Plus } from 'lucide-react'
import OrganizationCard from '#/features/organizations/ui/organization-card'
import { PaginatedContainer } from '@trieoh/ui-base'

export const Route = createFileRoute('/admin/')({
  component: RouteComponent,
})


function RouteComponent() {
  const queryClient = useQueryClient()
  const [filter, setFilter] = useState('')
  const [isCreateOpen, setIsCreateOpen] = useState(false)

  const { data: orgs = [] } = useQuery(allOrganizationsQueryOptions())

  const filteredOrgs = orgs.filter((org) => {
    const search = filter.toLowerCase()

    return (
      org.name.toLowerCase().includes(search) ||
      org.slug.includes(search)
    )
  })

  const { mutate: createOrganization, isPending: isPendingCreate } = useMutation({
    mutationFn: createOrganizationFn,
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allOrganizationsQueryOptions().queryKey,
          (old: OrganizationI[] = []) => [response.data, ...old],
        )
        setIsCreateOpen(false)
        toast.success(response.message || 'Organization created successfully')
      } else toast.error(response.message || "Failed to create namespace")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <div className="flex flex-wrap p-4">
      <PaginatedContainer<OrganizationI>
        items={filteredOrgs}
        layout="grid"
        minItemWidth="16rem"
        pageSize={10}
        sortFields={[
          { key: 'name', label: 'Name' },
          { key: 'slug', label: 'Slug' },
        ]}
        gap="6"
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by name, slug..."
        itemLabel="organizations"
        headerActions={
          <Button
            onClick={() => setIsCreateOpen(true)}
            className="rounded-sm gap-2"
          >
            <Plus className="w-4 h-4" />
            Create Organization
          </Button>
        }
        renderItems={(slice) => slice.map((item) => <OrganizationCard data={item} key={item.id} />)}
      />

      <FormModal<OrganizationCreateI>
        title="Create Organization"
        description="Give your organization a name to identify it."
        buttonTitle="Create Organization"
        schema={organizationCreateSchema}
        formId="create-organization-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={createOrganization}
        fields={[
          {
            name: "name",
            label: "e.g. My Team Organization",
            type: "text",
          },
          {
            name: "slug",
            label: "e.g. my-team",
            type: "text",
          }
        ]}
        disabled={isPendingCreate}
      />
    </div>
  )
}
