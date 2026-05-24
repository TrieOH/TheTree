import { allNamespacesQueryOptions, createNamespaceFn } from '#/features/namespaces/api';
import { namespaceCreateSchema, type NamespaceCreateI, type NamespaceI } from '#/features/namespaces/model';
import { NamespaceCard } from '#/features/namespaces/ui/namespace-card';
import { PaginatedContainer } from '#/widgets/pagination/paginated-container-grid';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react';
import { Plus } from 'lucide-react';
import { Button } from '#/shared/ui/shadcn/button';
import FormModal from '#/widgets/modal/form-modal';
import { toast } from 'sonner';
import type { FieldDefinition } from '#/shared/model/form-types';

export const Route = createFileRoute('/admin/')({
  component: RouteComponent,
})

const MOCK_NAMESPACES: NamespaceI[] = [
  {
    id: '1',
    name: 'Mock Namespace 1',
    owner_id: 'owner-1',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: '2',
    name: 'Mock Namespace 2',
    owner_id: 'owner-2',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
];

const NAMESPACE_FIELDS: FieldDefinition<NamespaceCreateI>[] = [
  {
    name: 'name',
    label: 'Namespace Name',
    placeholder: 'Enter namespace name...',
    type: 'text',
  },
];

function RouteComponent() {
  const [filter, setFilter] = useState('')
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const queryClient = useQueryClient()
  const { data: namespaces = [] } = useQuery(allNamespacesQueryOptions())

  const allNamespaces = [...MOCK_NAMESPACES, ...namespaces]

  const filteredNamespaces = allNamespaces.filter((namespace) =>
    namespace.name.toLowerCase().includes(filter.toLowerCase())
  )

  const { mutate: createNamespace, isPending: isCreating } = useMutation({
    mutationFn: (data: NamespaceCreateI) => createNamespaceFn(data),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(allNamespacesQueryOptions().queryKey, (oldData: NamespaceI[] = []) => {
          return [response.data, ...oldData];
        })
        setIsCreateOpen(false)
        toast.success(response.message || "Namespace created successfully")
      } else {
        toast.error(response.message || "Failed to create namespace")
      }
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <main className='flex flex-wrap gap-4 p-4'>
      <PaginatedContainer<NamespaceI>
        items={filteredNamespaces}
        className='w-full'
        layout='flex'
        pageSize={10}
        sortFields={[
          { key: "name", label: "Name" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by name…"
        itemLabel="namespaces"
        headerActions={
          <Button
            onClick={() => setIsCreateOpen(true)}
            size="icon"
            variant="outline"
            className="sm:w-auto px-3 rounded-sm"
          >
            <Plus size={16} />
            <span className="hidden sm:inline ml-2">Create Namespace</span>
          </Button>
        }
        renderItems={(slice) => slice.map(item => <NamespaceCard key={item.id} data={item} />)}
      />

      <FormModal<NamespaceCreateI>
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        title="Create Namespace"
        description="Namespaces allow you to group your resources and manage permissions."
        formId="create-namespace-form"
        buttonTitle="Create Namespace"
        fields={NAMESPACE_FIELDS}
        schema={namespaceCreateSchema}
        onSubmit={createNamespace}
        disabled={isCreating}
      />
    </main>
  )
}
