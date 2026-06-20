import { allOrganizationsQueryOptions, createOrganizationFn } from '@/features/organizations/api'
import type { OrganizationCreateI, OrganizationI } from '@/features/organizations/model'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'
import { PaginatedContainer } from '@/widgets/pagination/PaginatedContainer'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { useState } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/')({
  component: RouteComponent,
})

function RouteComponent() {
  const [_isCreateOpen, setIsCreateOpen] = useState(false)
  const queryClient = useQueryClient()
  const { data: _orgs = [] } = useQuery(allOrganizationsQueryOptions())

  const { mutate: _createOrganization, isPending: _isCreating } = useMutation({
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
        items={[]}
        layout='grid'
        minItemWidth='18rem'
        pageSize={10}
        sortFields={[
          { key: "name", label: "Name" },
        ]}
        gap='6'
        filterPlaceholder="Filter by name, owner, member..."
        itemLabel="namespaces"
        headerActions={
          <ShadowButton
            onClick={() => setIsCreateOpen(true)}
            variant="outline"
            className="sm:w-auto px-3 rounded-sm"
            value='Create Organization'
            leftIcon={<Plus size={16} />}
          />
        }
        renderItems={(slice) => slice.map(item => <div key={item.id} />)}
      />
    </div>
  )
}
