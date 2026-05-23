import { allNamespacesQueryOptions } from '#/features/namespaces/api';
import type { NamespaceI } from '#/features/namespaces/model';
import { NamespaceCard } from '#/features/namespaces/ui/namespace-card';
import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'

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

  const { data: namespaces = [] } = useQuery(allNamespacesQueryOptions())

  const allNamespaces = [...MOCK_NAMESPACES, ...namespaces]

  return (
    <main className='flex flex-wrap gap-4 p-4'>
      {allNamespaces.map(item => {
        return <NamespaceCard key={item.id} data={item} />
      })}
    </main>
  )
}
