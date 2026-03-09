import { createFileRoute } from '@tanstack/react-router'
import { WorkspaceList } from '#/features/workspaces/ui/workspace-list'
import type { WorkspaceI } from '#/features/workspaces/model'

export const Route = createFileRoute('/admin/')({
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
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-700">
      <WorkspaceList workspaces={MOCK_WORKSPACES} />
    </div>
  )
}
