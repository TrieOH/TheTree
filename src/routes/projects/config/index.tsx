import { requireAuth } from '@/features/auth/lib/route-guard';
import { navigationStore } from '@/features/navigation';
import ScopeTable from '@/features/scope/ui/ScopeTable';
import CustomTabs from '@/widgets/tabs/ui/CustomTabs';
import { createFileRoute, redirect } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store';
import { Database, Globe, LayoutDashboard, Shield, ShieldCheck, UserCog } from 'lucide-react';

export const Route = createFileRoute('/projects/config/')({
  beforeLoad: async (ctx) => {
    requireAuth(ctx)
    const currentProjectId = navigationStore.state.currentProjectId;
    if (typeof window !== 'undefined' && !currentProjectId) throw redirect({ to: '/projects' });
  },
  component: RouteComponent,
  staticData: {
    components: {
      header: "schemas"
    }
  },
})



function RouteComponent() {
  const currentProjectId = useStore(navigationStore, (state) => state.currentProjectId || "");
  const items = [
    { value: 'dashboard', label: 'Dashboard', icon: LayoutDashboard, content: <p></p> },
    { value: 'schema', label: 'Schema', icon: Database, content: <p>Editor de tabelas e campos...</p> },
    { value: 'permissions', label: 'Permissions', icon: Shield, content: <p>Gerenciamento de permissões...</p> },
    { 
      value: 'scope', 
      label: 'Scope', 
      icon: Globe, 
      content: <ScopeTable data={[
        {name: "dwd", id: "d", created_at: "2026-02-11T02:26:04+03:00", type: "dwdw", external_id: "dw", updated_at: "dw", project_id: "dwdw"}
      ]}
      project_id={currentProjectId}
      />,
    },
    { value: 'roles', label: 'Roles', icon: ShieldCheck, content: <p>Gerenciamento de roles...</p> },
    { value: 'users', label: 'Users', icon: UserCog, content: <p>Gerenciamento de usuários...</p> },
  ];
  return (
    <main className='flex justify-center items-center h-(--screen--minus-header)'>
      <CustomTabs items={items} />
    </main>
  );
}