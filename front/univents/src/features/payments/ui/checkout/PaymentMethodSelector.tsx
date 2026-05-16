import { CreditCard, QrCode, Lock, AlertTriangle } from "lucide-react"
import { motion, AnimatePresence } from "motion/react"
import type { PaymentMethodI } from "../../model"
import { cn } from "@/shared/lib/utils"
import { Button } from "@/shared/ui/shadcn/button"

interface PaymentMethodSelectorProps {
  amountCents: number
  selectedMethod: PaymentMethodI | null
  onSelectMethod: (method: PaymentMethodI) => void
  onCancel?: () => void
  disabled?: boolean
  isTooLowForCreditCard?: boolean
}

const methods = [
  {
    id: "credit_card" as const,
    label: "Cartão de crédito",
    description: "Até 12 parcelas",
    icon: CreditCard,
  },
  {
    id: "pix" as const,
    label: "Pix",
    description: "Aprovação imediata",
    icon: QrCode,
  },
]

function formatCurrency(cents: number) {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(cents / 100)
}

export function PaymentMethodSelector({
  amountCents,
  selectedMethod,
  onSelectMethod,
  onCancel,
  disabled = false,
  isTooLowForCreditCard = false,
}: PaymentMethodSelectorProps) {
  return (
    <div className="w-full min-w-75 space-y-4">
      {/* Header */}
      <div className="flex items-center gap-2 pb-2 border-b border-border">
        <Lock className="w-4 h-4 text-primary" />
        <h2 className="text-xs font-bold uppercase tracking-wide text-muted-foreground">
          Pagamento seguro
        </h2>
      </div>

      {/* Methods */}
      <div className="space-y-2">
        {methods.map((method) => {
          const Icon = method.icon
          const isSelected = selectedMethod === method.id
          const isDisabled = disabled || (method.id === "credit_card" && isTooLowForCreditCard)

          return (
            <div key={method.id} className="space-y-1">
              <Button
                variant="picker"
                size="none"
                data-state={isSelected ? "active" : "inactive"}
                onClick={() => {
                  if (!isDisabled) onSelectMethod(method.id)
                }}
                disabled={isDisabled}
                className={cn(
                  "w-full flex items-center gap-3 p-3 text-left transition-all rounded-none border",
                  isSelected && "border-primary",
                  isDisabled && "opacity-50 cursor-not-allowed"
                )}
              >
                <Icon
                  className={cn(
                    "w-5 h-5 shrink-0 transition-colors duration-200",
                    isSelected ? "text-primary" : "text-muted-foreground"
                  )}
                />

                <div className="flex-1 min-w-0">
                  <p
                    className={cn(
                      "text-sm font-bold uppercase tracking-wide transition-colors duration-200",
                      isSelected ? "text-primary" : "text-foreground"
                    )}
                  >
                    {method.label}
                  </p>
                  <p className="text-xs text-muted-foreground font-medium normal-case tracking-normal">
                    {method.description}
                  </p>
                </div>

                <div
                  className={cn(
                    "w-4 h-4 border shrink-0 relative overflow-hidden transition-colors duration-200",
                    isSelected
                      ? "border-primary bg-primary"
                      : "border-muted-foreground/30"
                  )}
                >
                  <AnimatePresence initial={false}>
                    {isSelected && (
                      <motion.div
                        initial={{ scale: 0, opacity: 0 }}
                        animate={{ scale: 0.5, opacity: 1 }}
                        exit={{ scale: 0, opacity: 0 }}
                        className="absolute inset-0 bg-primary-foreground"
                      />
                    )}
                  </AnimatePresence>
                </div>
              </Button>

              {/* Aviso inline abaixo do botão de cartão */}
              <AnimatePresence initial={false}>
                {method.id === "credit_card" && isTooLowForCreditCard && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: "auto" }}
                    exit={{ opacity: 0, height: 0 }}
                    transition={{ duration: 0.2, ease: "easeInOut" }}
                    className="overflow-hidden"
                  >
                    <div className="flex items-start gap-2 px-3 py-2 bg-amber-50 border border-amber-200 dark:bg-amber-950/30 dark:border-amber-800">
                      <AlertTriangle className="w-3.5 h-3.5 text-amber-600 dark:text-amber-500 shrink-0 mt-0.5" />
                      <p className="text-[11px] text-amber-700 dark:text-amber-400 leading-relaxed">
                        Valor mínimo para cartão é{" "}
                        <span className="font-semibold">R$ 1,00</span>.
                        Use Pix para este valor.
                      </p>
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          )
        })}
      </div>

      {/* Total */}
      <div className="pt-4 border-t-2 border-primary/10 space-y-3">
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold uppercase tracking-wide text-foreground">
            Total
          </span>
          <span className="text-xl font-bold text-primary tabular-nums">
            {formatCurrency(amountCents)}
          </span>
        </div>

        <Button
          variant="ghost"
          size="xl"
          disabled={disabled}
          onClick={onCancel}
          className="w-full text-muted-foreground hover:text-foreground"
        >
          Cancelar e voltar
        </Button>
      </div>
    </div>
  )
}
