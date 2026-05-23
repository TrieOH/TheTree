import { allNamespacesQueryOptions } from '#/features/namespaces/api';
import type { NamespaceI } from '#/features/namespaces/model';
import { NamespaceCard } from '#/features/namespaces/ui/namespace-card';
import { PaginatedContainer } from '#/widgets/pagination/paginated-container-grid';
import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react';

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

function RouteComponent() {
  const [filter, setFilter] = useState('')
  const { data: namespaces = [] } = useQuery(allNamespacesQueryOptions())

  const allNamespaces = [...MOCK_NAMESPACES, ...namespaces]

  const filteredNamespaces = allNamespaces.filter((namespace) =>
    namespace.name.toLowerCase().includes(filter.toLowerCase())
  )

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
        renderItems={(slice) => slice.map(item => <NamespaceCard key={item.id} data={item} />)}
      />
    </main>
  )
}
