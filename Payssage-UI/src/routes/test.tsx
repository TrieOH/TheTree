import { createFileRoute, Outlet, Link } from '@tanstack/react-router'
import { useSubdomain } from '../lib/use-subdomain'

export const Route = createFileRoute('/test')({
  component: TestLayout,
})

function TestLayout() {
  const subdomain = useSubdomain()

  return (
    <main className="page-wrap px-4 pb-8 pt-14">
      <div className="mb-6 flex items-center justify-between rounded-xl bg-[rgba(79,184,178,0.1)] p-4">
        <div>
          <span className="text-sm font-medium opacity-70">Subdomínio detectado: </span>
          <code className="rounded bg-white/50 px-2 py-1 font-bold text-[var(--lagoon-deep)]">
            {subdomain || 'nenhum'}
          </code>
        </div>
        <nav className="flex gap-4">
          <Link to="/test" className="text-sm font-semibold hover:underline">Início Teste</Link>
          <Link to="/test/admin" className="text-sm font-semibold hover:underline">Admin Teste</Link>
        </nav>
      </div>

      <Outlet />
    </main>
  )
}
