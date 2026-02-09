import { requireAuth } from '@/features/auth/lib/route-guard';
import { navigationStore } from '@/features/navigation';
import { Tabs, TabsList, TabsTrigger } from '@/shared/ui/shadcn/tabs';
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
  console.log(currentProjectId)
  return (
    <main>
      <Tabs defaultValue="dashboard" className="w-full">
        <TabsList className="flex h-auto w-full justify-start gap-1 overflow-x-auto rounded-lg bg-muted p-1 md:grid md:w-fit md:grid-cols-6">
          <TabsTrigger 
            value="dashboard" 
            className="flex shrink-0 items-center gap-2 px-3 py-2 text-xs md:text-sm"
          >
            <LayoutDashboard className="h-4 w-4" />
            <span className="hidden sm:inline">Dashboard</span>
            <span className="sm:hidden">Dash</span>
          </TabsTrigger>
          <TabsTrigger 
            value="schema" 
            className="flex shrink-0 items-center gap-2 px-3 py-2 text-xs md:text-sm"
          >
            <Database className="h-4 w-4" />
            <span>Schema</span>
          </TabsTrigger>
          <TabsTrigger 
            value="permissions" 
            className="flex shrink-0 items-center gap-2 px-3 py-2 text-xs md:text-sm"
          >
            <Shield className="h-4 w-4" />
            <span className="hidden sm:inline">Permissions</span>
            <span className="sm:hidden">Perms</span>
          </TabsTrigger>

          <TabsTrigger 
            value="scope" 
            className="flex shrink-0 items-center gap-2 px-3 py-2 text-xs md:text-sm"
          >
            <Globe className="h-4 w-4" />
            <span>Scope</span>
          </TabsTrigger>

          <TabsTrigger 
            value="roles" 
            className="flex shrink-0 items-center gap-2 px-3 py-2 text-xs md:text-sm"
          >
            <ShieldCheck className="h-4 w-4" />
            <span>Roles</span>
          </TabsTrigger>

          <TabsTrigger 
            value="users" 
            className="flex shrink-0 items-center gap-2 px-3 py-2 text-xs md:text-sm"
          >
            <UserCog className="h-4 w-4" />
            <span>Users</span>
          </TabsTrigger>
        </TabsList>
      </Tabs>
    </main>
  )
}
