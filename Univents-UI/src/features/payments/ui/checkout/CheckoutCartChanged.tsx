interface CheckoutCartChangedProps {
  onResume: () => void
  onUseNew: () => void
  onCancel: () => void
}

export default function CheckoutCartChanged({
  onResume,
  onUseNew,
  onCancel,
}: CheckoutCartChangedProps) {
  return (
    <main className="w-full min-w-75 max-w-sm mx-auto px-3 py-16 flex flex-col items-center gap-5 text-center">
      <div className="w-14 h-14 rounded-full flex items-center justify-center text-xl bg-amber-500/10 text-amber-500">
        ?
      </div>

      <div className="space-y-1">
        <h1 className="text-lg font-bold text-foreground">Carrinho alterado</h1>
        <p className="text-sm text-muted-foreground">
          Você tem uma reserva ativa, mas seu carrinho foi alterado. O que deseja fazer?
        </p>
      </div>

      <div className="flex flex-col gap-2 w-full">
        <button
          onClick={onUseNew}
          className="w-full rounded-md bg-primary text-primary-foreground px-4 py-2.5 text-sm font-medium hover:bg-primary/90 transition-colors"
        >
          Usar carrinho novo
        </button>
        <button
          onClick={onResume}
          className="w-full rounded-md border border-border px-4 py-2.5 text-sm text-foreground hover:bg-muted/50 transition-colors"
        >
          Continuar reserva anterior
        </button>
        <button
          onClick={onCancel}
          className="w-full rounded-md px-4 py-2.5 text-sm text-muted-foreground hover:bg-muted/50 transition-colors"
        >
          Cancelar
        </button>
      </div>
    </main>
  )
}