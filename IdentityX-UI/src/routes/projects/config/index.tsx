import { z } from 'zod';
import { requireAuth } from '@/features/auth/lib/route-guard';
import { navigationStore } from '@/features/navigation';
import SchemaTable from '@/features/schema/ui/SchemaTable';
import ScopeTable from '@/features/scope/ui/ScopeTable';
import { usersQueryOptions } from '@/features/user/api';
import UserTable from '@/features/user/ui/UserTable';
import CustomTabs from '@/widgets/tabs/ui/CustomTabs';
import { useSuspenseQuery, useQueryClient } from '@tanstack/react-query';
import { createFileRoute, redirect } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store';
import { Database, Globe, KeySquare, Shield, ShieldCheck, UserCog } from 'lucide-react';
import PermissionTable from '@/features/permission/ui/PermissionTable';
import RoleTable from '@/features/role/ui/RoleTable';
import APIKeyManager from '@/features/api-keys/ui/APIKeyManager';
import { schemasQueryOptions } from '@/features/schema/api';
import { scopesQueryOptions } from '@/features/scope/api';
import { roleQueryOptions } from '@/features/role/api';
import { permissionsQueryOptions } from '@/features/permission/api';
import { useMemo } from 'react';

export const Route = createFileRoute('/projects/config/')({
  beforeLoad: async (ctx) => {
    requireAuth(ctx)
    const currentProjectId = navigationStore.state.currentProjectId;
    if (typeof window !== 'undefined' && !currentProjectId) throw redirect({ to: '/projects' });
  },
  loader: async ({ context: { queryClient }}) => {
    if (typeof window === 'undefined') return { }
    
    const id = navigationStore.state.currentProjectId;
    queryClient.prefetchQuery(usersQueryOptions(id || ""))
    return { }
  },
  validateSearch: z.object({ tab: z.string().optional().default('schema') }),
  component: RouteComponent,
  staticData: {
    components: {
      header: "projects/config"
    }
  },
})


function RouteComponent() {
  const queryClient = useQueryClient();
  const currentProjectId = useStore(navigationStore, (state) => state.currentProjectId || "");
  const { data: users } = useSuspenseQuery(usersQueryOptions(currentProjectId))
  const { tab } = Route.useSearch();

  const items = useMemo(() => [
    {
      value: 'schema',
      label: 'Schema',
      icon: Database,
      content: <SchemaTable project_id={currentProjectId}/>,
      onRefresh: () => {
        queryClient.invalidateQueries(schemasQueryOptions(currentProjectId));
        queryClient.invalidateQueries({ queryKey: ['latestSchemaVersion', currentProjectId] });
        queryClient.invalidateQueries({ queryKey: ['currentSchemaVersion', currentProjectId] });
        queryClient.invalidateQueries({ queryKey: ['schemaVersionById', currentProjectId] });
      }
    },
    {
      value: 'scope',
      label: 'Scope',
      icon: Globe,
      content: <ScopeTable project_id={currentProjectId} />,
      onRefresh: () => queryClient.invalidateQueries(scopesQueryOptions(currentProjectId))
    },
    { 
      value: 'roles', 
      label: 'Roles', 
      icon: ShieldCheck, 
      content: <RoleTable project_id={currentProjectId}/>,
      onRefresh: () => {
        queryClient.invalidateQueries(roleQueryOptions(currentProjectId));
        queryClient.invalidateQueries({ queryKey: ['rolePermissions', currentProjectId] });
      }
    },
    { 
      value: 'permissions', 
      label: 'Permissions', 
      icon: Shield, 
      content: <PermissionTable project_id={currentProjectId} />,
      onRefresh: () => queryClient.invalidateQueries(permissionsQueryOptions(currentProjectId))
    },
    {
      value: 'users',
      label: 'Users',
      icon: UserCog,
      content: <UserTable data={users} project_id={currentProjectId} />,
      onRefresh: () => {
        // Invalidate the main user list
        queryClient.invalidateQueries(usersQueryOptions(currentProjectId));
        
        // Invalidate all potential sub-data for all users in this project
        queryClient.invalidateQueries({ queryKey: ['userRoles', currentProjectId] });
        queryClient.invalidateQueries({ queryKey: ['userPermissions', currentProjectId] });
        queryClient.invalidateQueries({ queryKey: ['rolePermissions', currentProjectId] });
        
        // Also refresh core configuration that users depend on
        queryClient.invalidateQueries(scopesQueryOptions(currentProjectId));
        queryClient.invalidateQueries(roleQueryOptions(currentProjectId));
        queryClient.invalidateQueries(permissionsQueryOptions(currentProjectId));
      }
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