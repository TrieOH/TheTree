import { z } from 'zod';
import { requireAuth } from '@/features/auth/lib/route-guard';
import { usersQueryOptions } from '@/features/user/api';
import UserTable from '@/features/user/ui/UserTable';
import CustomTabs from '@/widgets/tabs/ui/CustomTabs';
import { useSuspenseQuery, useQueryClient } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router'
import { KeySquare, UserCog } from 'lucide-react';
import APIKeyManager from '@/features/api-keys/ui/APIKeyManager';
import { useMemo } from 'react';

export const Route = createFileRoute('/projects/$projectId/config')({
  beforeLoad: requireAuth,
  loader: async ({ context: { queryClient }, params}) => {
    if (typeof window === 'undefined') return { }
    queryClient.prefetchQuery(usersQueryOptions(params.projectId || ""))
    return { }
  },
  validateSearch: z.object({ tab: z.string().optional().default('users') }),
  component: RouteComponent,
  staticData: {
    components: {
      header: "projects/config"
    }
  },
})

function RouteComponent() {
  const queryClient = useQueryClient();
  const { projectId: currentProjectId } = Route.useParams();
  const { data: users } = useSuspenseQuery(usersQueryOptions(currentProjectId))
  const { tab } = Route.useSearch();

  const items = useMemo(() => [
    {
      value: 'users',
      label: 'Users',
      icon: UserCog,
      content: <UserTable data={users} project_id={currentProjectId} />,
      onRefresh: () => queryClient.invalidateQueries(usersQueryOptions(currentProjectId))
    },
    {
      value: 'api-keys',
      label: 'API Keys',
      icon: KeySquare,
      content: <APIKeyManager publicKey={currentProjectId}/>
    },
  ], [currentProjectId, queryClient, users]);

  return (
    <main className='flex justify-center items-center h-(--screen--minus-header)'>
      <CustomTabs items={items} initialValue={tab}/>
    </main>
  );
}
