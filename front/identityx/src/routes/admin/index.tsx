import { allOrganizationsQueryOptions, createOrganizationFn } from '@/features/organizations/api'
import { organizationCreateSchema, type OrganizationCreateI, type OrganizationI } from '@/features/organizations/model'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'
import { PaginatedContainer } from '@/widgets/pagination/PaginatedContainer'
import { FormModal } from '@/widgets/modal/FormModal'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { useState } from 'react'
import { toast } from 'sonner'
import OrganizationCard from '@/features/organizations/ui/organization-card'

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

  const { mutate: createOrganization, isPending: isCreating } = useMutation({
    mutationFn: (data: OrganizationCreateI) => createOrganizationFn(data),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(allOrganizationsQueryOptions().queryKey, (oldData: OrganizationI[] = []) => {
          return [response.data, ...oldData];
        })
        setIsCreateOpen(false)
        toast.success(response.message || "Namespace created successfully")
      } else toast.error(response.message || "Failed to create namespace")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <div className='flex flex-wrap p-4'>
      <PaginatedContainer<OrganizationI>
        items={filteredOrgs}
        layout='grid'
        minItemWidth='16rem'
        pageSize={10}
        sortFields={[
          { key: "name", label: "Name" },
          { key: "slug", label: "Slug" },
        ]}
        gap='6'
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by name, slug..."
        itemLabel="organizations"
        headerActions={
          <ShadowButton
            onClick={() => setIsCreateOpen(true)}
            variant="outline"
            className="sm:w-auto px-3 rounded-sm"
            value='Create Organization'
            leftIcon={<Plus size={16} />}
          />
        }
        renderItems={(slice) => slice.map(item => <OrganizationCard data={item} key={item.id} />)}
      />

      <FormModal<OrganizationCreateI>
        formId='org-form'
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        title="Create Organization"
        description="Enter the details to create a new organization."
        schema={organizationCreateSchema}
        defaultValues={{ name: '', slug: '' }}
        fields={[
          { name: 'name', label: 'Name', placeholder: 'John Doe Goods' },
          { name: 'slug', label: 'Slug', placeholder: 'jd-goods' },
        ]}
        onSubmit={(data) => createOrganization(data)}
        submitLabel="Create Organization"
        isLoading={isCreating}
      />
    </div>
  )
}
