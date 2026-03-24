import StatusScreen from "@/widgets/feedback/ui/StatusScreen";

interface CheckoutErrorProps {
  message: string | null
  onRetry: () => void
  onBack: () => void
}


export default function CheckoutError({ message, onRetry, onBack }: CheckoutErrorProps) {
  return (
    <StatusScreen
      icon="!"
      iconClass="text-destructive bg-destructive/10"
      title="Algo deu errado"
      description={message ?? "Ocorreu um erro inesperado. Tente novamente."}
      actions={[
        { label: "Tentar novamente", onClick: onRetry, variant: "primary" },
        { label: "Voltar", onClick: onBack, variant: "outline" },
      ]}
    />
  )
}