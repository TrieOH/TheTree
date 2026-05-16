import { useSubdomain } from '#/shared/hooks/use-subdomain'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/test/')({
  component: TestIndex,
})

function TestIndex() {
  const subdomain = useSubdomain()

  return (
    <section className="island-shell rise-in rounded-2xl p-8">
      <h2 className="mb-4 text-3xl font-bold">Página de Início do Teste</h2>
      <p className="text-lg opacity-80">
        Esta é a rota principal de teste <code>/test</code>.
      </p>
      <div className="mt-6 rounded-lg bg-white/40 p-6">
        <p className="mb-2 font-medium">Informações do Cliente:</p>
        <p className="text-2xl font-semibold text-(--lagoon-deep)">
          {subdomain ? `Bem-vindo, ${subdomain}!` : 'Nenhum cliente identificado via subdomínio.'}
        </p>
      </div>
    </section>
  )
}
