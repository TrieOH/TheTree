import { useSubdomain } from '#/shared/hooks/use-subdomain'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/test/admin')({
  component: TestAdmin,
})

function TestAdmin() {
  const subdomain = useSubdomain()

  return (
    <section className="island-shell rise-in rounded-2xl bg-[rgba(23,58,64,0.05)] p-8">
      <div className="mb-4 inline-block rounded-full bg-(--lagoon-deep) px-3 py-1 text-xs font-bold text-white uppercase tracking-wider">
        Área Administrativa
      </div>
      <h2 className="mb-4 text-3xl font-bold">Painel Admin de {subdomain || 'Global'}</h2>
      <p className="mb-8 text-lg opacity-80">
        Esta rota é <code>/test/admin</code> e mostra dados específicos para a gestão.
      </p>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="rounded-xl bg-white p-6 shadow-sm">
          <p className="mb-1 text-sm opacity-60">Status do Subdomínio</p>
          <p className="text-xl font-bold">{subdomain ? 'Ativo' : 'Não Vinculado'}</p>
        </div>
        <div className="rounded-xl bg-white p-6 shadow-sm">
          <p className="mb-1 text-sm opacity-60">Permissões</p>
          <p className="text-xl font-bold">Total (Root)</p>
        </div>
      </div>
    </section>
  )
}
