export function OrganizerView() {
  return (
    <div className="max-w-5xl mx-auto space-y-32">
      {/* Hero */}
      <section className="text-center max-w-2xl mx-auto space-y-6">
        <p className="text-lg md:text-xl text-neutral-600 leading-relaxed">
          A plataforma mais simples para criar e vender ingressos.
          Comece grátis, pague apenas quando vender.
        </p>
        <div className="flex gap-4 justify-center pt-4">
          <button className="px-6 py-3 bg-neutral-900 text-white rounded-full text-sm font-medium hover:bg-neutral-800 transition-colors">
            Criar Evento Grátis
          </button>
          <button className="px-6 py-3 border border-neutral-200 rounded-full text-sm font-medium hover:border-neutral-400 transition-colors">
            Ver Preços
          </button>
        </div>
      </section>

      {/* Métricas */}
      <section className="grid grid-cols-3 gap-8 border-y border-neutral-200 py-12">
        <div className="text-center">
          <div className="text-4xl font-semibold mb-1">10%</div>
          <div className="text-sm text-neutral-500">Taxa por venda</div>
        </div>
        <div className="text-center border-x border-neutral-200">
          <div className="text-4xl font-semibold mb-1">2 dias</div>
          <div className="text-sm text-neutral-500">Para receber</div>
        </div>
        <div className="text-center">
          <div className="text-4xl font-semibold mb-1">0</div>
          <div className="text-sm text-neutral-500">Mensalidade</div>
        </div>
      </section>

      {/* Recursos - Layout Editorial */}
      <section className="space-y-16">
        <div className="grid md:grid-cols-2 gap-12 items-center">
          <div className="space-y-4">
            <h3 className="text-2xl font-medium">Crie em minutos</h3>
            <p className="text-neutral-600 leading-relaxed text-lg">
              Interface minimalista para configurar datas, lotes e preços
              sem burocracia. Foque no que importa: seu evento.
            </p>
          </div>
          <div className="aspect-video bg-neutral-100 rounded-2xl" />
        </div>

        <div className="grid md:grid-cols-2 gap-12 items-center">
          <div className="aspect-video bg-neutral-100 rounded-2xl order-2 md:order-1" />
          <div className="space-y-4 order-1 md:order-2">
            <h3 className="text-2xl font-medium">Venda sem limites</h3>
            <p className="text-neutral-600 leading-relaxed text-lg">
              Checkout otimizado que converte. Aceite Pix, cartão e parcelamento.
              Relatórios em tempo real de vendas.
            </p>
          </div>
        </div>
      </section>

      {/* Depoimento */}
      <section className="bg-neutral-50 rounded-3xl p-12 md:p-16 text-center">
        <blockquote className="text-2xl md:text-3xl font-medium leading-relaxed mb-8">
          "Reduzimos o tempo de organização em 70%.
          Finalmente uma plataforma que entende produtores."
        </blockquote>
        <div className="text-sm text-neutral-500">
          — Ana Silva, Produtora de Eventos
        </div>
      </section>
    </div>
  )
}