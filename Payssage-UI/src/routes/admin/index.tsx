import { requireAuth } from '#/features/auths/lib/route-guard'
import { createFileRoute } from '@tanstack/react-router'
import { AdminLayout } from '#/features/admin/ui/admin-layout'
import { WorkspaceList } from '#/features/workspaces/ui/workspace-list'
import type { WorkspaceI } from '#/features/workspaces/model'

export const Route = createFileRoute('/admin/')({
  beforeLoad: (ctx) => {
    requireAuth(ctx);
  },
  component: RouteComponent,
})

const MOCK_WORKSPACES: WorkspaceI[] = [
  {
    id: 'ws_prod_01jhc83',
    name: 'Loja Principal',
    created_at: '2025-01-15T10:00:00Z',
    sandbox: false,
  },
  {
    id: 'ws_sbx_01jhc84',
    name: 'Ambiente de Testes',
    created_at: '2025-02-20T14:30:00Z',
    sandbox: true,
  },
  {
    id: 'ws_sbx_01jhc85',
    name: 'Desenvolvimento API',
    created_at: '2025-03-01T09:15:00Z',
    sandbox: true,
  },
]

function RouteComponent() {
  return (
    <AdminLayout>
      <div className="animate-in fade-in slide-in-from-bottom-4 duration-700">
        <WorkspaceList workspaces={MOCK_WORKSPACES} />
      </div>
    </AdminLayout>
  )
}
