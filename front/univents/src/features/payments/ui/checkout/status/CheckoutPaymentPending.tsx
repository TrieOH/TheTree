import StatusScreen from "@/widgets/feedback/ui/StatusScreen"

interface CheckoutPaymentPendingProps {
  message: string | null
}

export default function CheckoutPaymentPending({ message }: CheckoutPaymentPendingProps) {
  return (
    <StatusScreen
      icon="…"
      iconClass="text-primary bg-primary/10"
      title="Pagamento em análise"
      description={message ?? "Seu pagamento está sendo processado. Você será notificado em breve."}
      actions={[]}
    />
  )
}