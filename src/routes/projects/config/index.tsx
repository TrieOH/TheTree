import { z } from 'zod';
import { requireAuth } from '@/features/auth/lib/route-guard';
import { navigationStore } from '@/features/navigation';
import SchemaTable from '@/features/schema/ui/SchemaTable';
import ScopeTable from '@/features/scope/ui/ScopeTable';
import { usersQueryOptions } from '@/features/user/api';
import UserTable from '@/features/user/ui/UserTable';
import CustomTabs from '@/widgets/tabs/ui/CustomTabs';
import { useSuspenseQuery } from '@tanstack/react-query';
import { createFileRoute, redirect } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store';
import { Database, Globe, KeySquare, LayoutDashboard, Shield, ShieldCheck, UserCog } from 'lucide-react';
import PermissionTable from '@/features/permission/ui/PermissionTable';
import RoleTable from '@/features/role/ui/RoleTable';
import APIKeyManager from '@/features/api-keys/ui/APIKeyManager';

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
  validateSearch: z.object({ tab: z.string().optional().default('dashboard') }),
  component: RouteComponent,
  staticData: {
    components: {
      header: "projects/config"
    }
  },
})


function RouteComponent() {
  const currentProjectId = useStore(navigationStore, (state) => state.currentProjectId || "");
  const { data: users } = useSuspenseQuery(usersQueryOptions(currentProjectId))
  const { tab } = Route.useSearch();

  const items = [
    { value: 'dashboard', label: 'Dashboard', icon: LayoutDashboard, content: <p></p> },
    {
      value: 'schema',
      label: 'Schema',
      icon: Database,
      content: <SchemaTable project_id={currentProjectId}/>
    },
    {
      value: 'scope',
      label: 'Scope',
      icon: Globe,
      content: <ScopeTable project_id={currentProjectId} />,
    },
    { 
      value: 'roles', 
      label: 'Roles', 
      icon: ShieldCheck, 
      content: <RoleTable project_id={currentProjectId}/>
    },
    { 
      value: 'permissions', 
      label: 'Permissions', 
      icon: Shield, 
      content: <PermissionTable project_id={currentProjectId} />
    },
    {
      value: 'users',
      label: 'Users',
      icon: UserCog,
      content: <UserTable data={users} project_id={currentProjectId} />
    },
    {
      value: 'api-keys',
      label: 'API Keys',
      icon: KeySquare,
      content: <APIKeyManager publicKey={currentProjectId}/>
    },
  ];

  return (
    <main className='flex justify-center items-center h-(--screen--minus-header)'>
      <CustomTabs items={items} initialValue={tab}/>
    </main>
  );
}