import { useState, useRef, useEffect } from "react"
import { MercadoPagoForm } from "./MercadoPago"
import { PaymentMethodSelector } from "./checkout/PaymentMethodSelector"
import { PaymentFormSheet } from "./checkout/PaymentFormSheet"
import type { PaymentProviderI, SubmitPaymentPayloadI, PaymentMethodI } from "../model"

interface PaymentProviderSelectorProps {
  provider?: PaymentProviderI
  amount: number
  handleSubmit: (data: SubmitPaymentPayloadI) => void
  sellerPublicKey: string
  onCancel?: () => void
}

export function PaymentProviderSelector({
  provider = "mercadopago",
  amount,
  handleSubmit,
  sellerPublicKey,
  onCancel,
}: PaymentProviderSelectorProps) {
  const isTooLowForCreditCard = amount < 100
  const [method, setMethod] = useState<PaymentMethodI | null>(null)
  const [sheetOpen, setSheetOpen] = useState(false)
  const [formReady, setFormReady] = useState(false)
  const closeTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (closeTimeoutRef.current) {
        clearTimeout(closeTimeoutRef.current)
      }
    }
  }, [])

  const handleSelectMethod = (m: PaymentMethodI) => {
    if (m === "credit_card" && isTooLowForCreditCard) return

    // Clear any pending cleanup from a previous close
    if (closeTimeoutRef.current) {
      clearTimeout(closeTimeoutRef.current)
      closeTimeoutRef.current = null
    }

    setMethod(m)
    setFormReady(false)
    setSheetOpen(true)
  }

  const handleClose = () => {
    setSheetOpen(false)
    setFormReady(false)

    if (closeTimeoutRef.current) {
      clearTimeout(closeTimeoutRef.current)
    }

    closeTimeoutRef.current = setTimeout(() => {
      setMethod(null)
      closeTimeoutRef.current = null
    }, 400)
  }

  switch (provider) {
    default:
      return (
        <>
          <PaymentMethodSelector
            amountCents={amount}
            selectedMethod={method}
            onSelectMethod={handleSelectMethod}
            onCancel={onCancel}
            isTooLowForCreditCard={isTooLowForCreditCard}
          />

          <PaymentFormSheet
            open={sheetOpen}
            method={method}
            onClose={handleClose}
            onReady={() => { setFormReady(true) }}
          >
            {method && formReady && (
              <MercadoPagoForm
                amount={amount}
                method={method}
                handleSubmit={(data) => {
                  handleSubmit(data)
                  handleClose()
                }}
                sellerPublicKey={sellerPublicKey}
              />
            )}
          </PaymentFormSheet>
        </>
      )
  }
}
