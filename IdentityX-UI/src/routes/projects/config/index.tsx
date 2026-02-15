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
import { Database, Globe, LayoutDashboard, Shield, ShieldCheck, UserCog } from 'lucide-react';

export const Route = createFileRoute('/projects/config/')({
  beforeLoad: async (ctx) => {
    requireAuth(ctx)
    const currentProjectId = navigationStore.state.currentProjectId;
    if (typeof window !== 'undefined' && !currentProjectId) throw redirect({ to: '/projects' });
  },
  loader: async ({ context: { queryClient }}) => {
    const id = navigationStore.state.currentProjectId;
    queryClient.prefetchQuery(usersQueryOptions(id || ""))
    return { }
  },
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
      value: 'users', 
      label: 'Users', 
      icon: UserCog, 
      content: <UserTable data={users}/>
    },
    { value: 'permissions', label: 'Permissions', icon: Shield, content: <p>Gerenciamento de permissões...</p> },
    { value: 'roles', label: 'Roles', icon: ShieldCheck, content: <p>Gerenciamento de roles...</p> },
  ];
  return (
    <main className='flex justify-center items-center h-(--screen--minus-header)'>
      <CustomTabs items={items} />
    </main>
  );
}