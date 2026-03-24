import StatusScreen from "@/widgets/feedback/ui/StatusScreen"

interface CheckoutReservationFailedProps {
  message: string | null
  onBack: () => void
}

export default function CheckoutReservationFailed({
  message,
  onBack
}: CheckoutReservationFailedProps) {
  return (
    <StatusScreen
      icon="✕"
      iconClass="text-destructive bg-destructive/10"
      title="Itens indisponíveis"
      description={message ?? "Não foi possível reservar os itens selecionados."}
      actions={[{ label: "Voltar", onClick: onBack, variant: "outline" }]}
    />
  )
}