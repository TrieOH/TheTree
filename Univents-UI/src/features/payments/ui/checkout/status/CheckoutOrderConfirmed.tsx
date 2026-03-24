import StatusScreen from "@/widgets/feedback/ui/StatusScreen";

export default function CheckoutOrderConfirmed() {
  return (
    <StatusScreen
      icon="✓"
      iconClass="text-emerald-600 bg-emerald-500/10"
      title="Pedido confirmado!"
      description="Seu pagamento foi aprovado. Em breve você receberá um e-mail de confirmação."
      actions={[]}
    />
  )
}