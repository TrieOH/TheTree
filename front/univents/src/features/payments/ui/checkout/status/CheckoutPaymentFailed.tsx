import StatusScreen from "@/widgets/feedback/ui/StatusScreen"

interface CheckoutPaymentFailedProps {
  message: string | null
  onBack: () => void
}

export default function CheckoutPaymentFailed({
  message,
  onBack
}: CheckoutPaymentFailedProps) {
  return (
    <StatusScreen
      icon="✕"
      iconClass="text-destructive bg-destructive/10"
      title="Erro no Pagamento"
      description={message ?? "Pagamento recusado. Verifique suas credenciais e tente novamente."}
      actions={[{ label: "Voltar", onClick: onBack, variant: "outline" }]}
    />
  )
}