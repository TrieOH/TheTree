export function ParticipantView() {
  return (
    <div className="max-w-5xl mx-auto space-y-32">
      {/* Hero Text */}
      <section className="text-center max-w-2xl mx-auto space-y-6">
        <p className="text-lg md:text-xl text-neutral-600 leading-relaxed">
          Descubra eventos selecionados e compre ingressos sem complicação.
          De shows a conferências, tudo em um só lugar.
        </p>
        <div className="flex gap-4 justify-center pt-4">
          <button className="px-6 py-3 bg-neutral-900 text-white rounded-full text-sm font-medium hover:bg-neutral-800 transition-colors">
            Explorar Agenda
          </button>
          <button className="px-6 py-3 border border-neutral-200 rounded-full text-sm font-medium hover:border-neutral-400 transition-colors">
            Como Funciona
          </button>
        </div>
      </section>

      {/* Grid de Eventos Destaque */}
      <section className="space-y-8">
        <div className="flex justify-between items-end border-b border-neutral-200 pb-4">
          <h2 className="text-2xl font-medium">Em Alta</h2>
          <a href="#" className="text-sm text-neutral-500 hover:text-neutral-900">Ver Todos</a>
        </div>

        <div className="grid md:grid-cols-2 gap-6">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="group cursor-pointer">
              <div className="aspect-4/3 bg-neutral-100 rounded-2xl mb-4 overflow-hidden">
                <div className="w-full h-full bg-linear-to-br from-neutral-200 to-neutral-300 group-hover:scale-105 transition-transform duration-500" />
              </div>
              <div className="space-y-2">
                <div className="flex justify-between items-start">
                  <h3 className="text-lg font-medium group-hover:text-neutral-600 transition-colors">
                    Nome do Evento {i}
                  </h3>
                  <span className="text-sm font-medium text-neutral-500">R$ 120</span>
                </div>
                <p className="text-sm text-neutral-500">São Paulo • 24 Mar 2026</p>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Diferenciais - Minimalistas */}
      <section className="grid md:grid-cols-3 gap-12 border-t border-neutral-200 pt-16">
        <div className="space-y-3">
          <span className="text-sm font-medium text-neutral-400">01</span>
          <h3 className="text-lg font-medium">Sem Taxas Escondidas</h3>
          <p className="text-neutral-600 leading-relaxed">
            O preço que você vê é o preço final. Sem surpresas no checkout.
          </p>
        </div>
        <div className="space-y-3">
          <span className="text-sm font-medium text-neutral-400">02</span>
          <h3 className="text-lg font-medium">Ingresso Digital</h3>
          <p className="text-neutral-600 leading-relaxed">
            Receba instantaneamente no email ou app. Sem filas de impressão.
          </p>
        </div>
        <div className="space-y-3">
          <span className="text-sm font-medium text-neutral-400">03</span>
          <h3 className="text-lg font-medium">Garantia de Entrada</h3>
          <p className="text-neutral-600 leading-relaxed">
            Seu dinheiro de volta se o evento for cancelado.
          </p>
        </div>
      </section>
    </div>
  )
}