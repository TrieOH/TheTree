import { useEffect, useRef, useState } from "react"
import { X, Loader2 } from "lucide-react"
import { motion, AnimatePresence } from "motion/react"
import type { PaymentMethodI } from "../../model"
import { cn } from "@/shared/lib/utils"
import { useIsMobile } from "@/shared/hooks/use-mobile"

interface PaymentFormSheetProps {
  open: boolean
  method: PaymentMethodI | null
  onClose: () => void
  onReady: () => void
  children: React.ReactNode
}

const methodLabel: Record<string, string> = {
  credit_card: "Cartão de crédito",
  pix: "Pix",
}

export function PaymentFormSheet({
  open,
  method,
  onClose,
  onReady,
  children,
}: PaymentFormSheetProps) {
  const isMobile = useIsMobile()
  const [ready, setReady] = useState(false)
  const onReadyRef = useRef(onReady)
  onReadyRef.current = onReady

  useEffect(() => {
    if (!open) {
      setReady(false)
      return
    }

    const timer = setTimeout(() => {
      setReady(true)
      onReadyRef.current()
    }, 500)
    return () => { clearTimeout(timer); }
  }, [open])

  const handleClose = () => {
    setReady(false)
    onClose()
  }

  return (
    <AnimatePresence>
      {open && (
        <>
          {/* Backdrop */}
          <motion.div
            key="backdrop"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            onClick={handleClose}
            className="fixed inset-0 z-40 bg-background/60 backdrop-blur-sm"
          />

          {/* Single Sheet - Condicional para evitar IDs duplicados no DOM */}
          {isMobile ? (
            <motion.div
              key="sheet-mobile"
              initial={{ y: "100%" }}
              animate={{ y: 0 }}
              exit={{ y: "100%" }}
              transition={{ type: "spring", damping: 30, stiffness: 300 }}
              className={cn(
                "fixed inset-x-0 bottom-0 z-50",
                "bg-background border-t border-border rounded-t-2xl shadow-2xl",
                "max-h-[90dvh] flex flex-col"
              )}
            >
              <div className="flex justify-center pt-3 pb-1 shrink-0">
                <div className="w-10 h-1 rounded-full bg-muted-foreground/20" />
              </div>

              <div className="flex items-center justify-between px-5 py-3 border-b border-border shrink-0">
                <span className="text-sm font-bold uppercase tracking-wide">
                  {method ? methodLabel[method] : ""}
                </span>
                <button onClick={handleClose} className="p-1 hover:bg-muted rounded-full transition-colors">
                  <X className="w-4 h-4 text-muted-foreground" />
                </button>
              </div>

              <div className="overflow-y-auto px-5 py-4 flex-1 min-h-75">
                {ready ? children : (
                  <div className="flex items-center justify-center h-48">
                    <Loader2 className="w-6 h-6 animate-spin text-primary" />
                  </div>
                )}
              </div>
            </motion.div>
          ) : (
            <motion.div
              key="sheet-desktop"
              initial={{ x: "100%" }}
              animate={{ x: 0 }}
              exit={{ x: "100%" }}
              transition={{ type: "spring", damping: 32, stiffness: 280 }}
              className={cn(
                "fixed top-0 right-0 bottom-0 z-50",
                "w-105 flex flex-col",
                "bg-background border-l border-border shadow-2xl"
              )}
            >
              <div className="flex items-center justify-between px-6 py-5 border-b border-border shrink-0">
                <span className="text-sm font-bold uppercase tracking-wide">
                  {method ? methodLabel[method] : ""}
                </span>
                <button onClick={handleClose} className="p-1 hover:bg-muted rounded-full transition-colors">
                  <X className="w-5 h-5 text-muted-foreground" />
                </button>
              </div>

              <div className="overflow-y-auto px-6 py-5 flex-1">
                {ready ? children : (
                  <div className="flex items-center justify-center h-48">
                    <Loader2 className="w-6 h-6 animate-spin text-primary" />
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </>
      )}
    </AnimatePresence>
  )
}
