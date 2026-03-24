import { Button } from "@/shared/ui/shadcn/button"

interface CheckoutPixPendingProps {
  qrCode: string
  qrCodeBase64: string
  totalCents: number
  onCancel: () => void;
}

function formatCurrency(cents: number) {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL"
  }).format(cents / 100)
}

export default function CheckoutPixPending({
  qrCode,
  qrCodeBase64,
  totalCents,
  onCancel
}: CheckoutPixPendingProps) {

  const handleCopy = () => navigator.clipboard.writeText(qrCode)

  return (
    <main className="w-full min-w-75 max-w-sm mx-auto px-3 py-8 flex flex-col items-center gap-6 text-center">
      <div>
        <h1 className="text-lg font-bold text-foreground">Pague com Pix</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Escaneie o QR Code ou copie o código abaixo
        </p>
      </div>

      {/* QR code image */}
      <div className="rounded-xl border border-border p-4 bg-white">
        <img
          src={`data:image/png;base64,${qrCodeBase64}`}
          alt="QR Code Pix"
          className="w-48 h-48"
        />
      </div>

      <p className="text-xl font-semibold text-foreground">{formatCurrency(totalCents)}</p>

      {/* Copy code */}
      <div className="w-full">
        <p className="text-xs text-muted-foreground mb-1">Código Pix copia e cola</p>
        <div className="flex items-center gap-2 rounded-md border border-border bg-muted/40 px-3 py-2">
          <code className="flex-1 text-xs truncate text-foreground">{qrCode}</code>
          <button
            onClick={handleCopy}
            className="text-xs text-primary hover:underline shrink-0"
          >
            Copiar
          </button>
        </div>
      </div>

      <p className="text-xs text-muted-foreground">
        Após o pagamento, a confirmação pode levar alguns instantes.
      </p>

      <Button
        variant="ghost"
        size="xl"
        onClick={onCancel}
        className="w-full text-muted-foreground hover:text-foreground"
      >
        Cancelar e voltar
      </Button>
    </main>
  )
}