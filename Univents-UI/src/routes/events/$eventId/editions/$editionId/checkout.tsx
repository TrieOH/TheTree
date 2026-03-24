import { useEffect, useRef, useState } from "react"
import { createFileRoute, useRouter } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import type { SubmitPaymentPayloadI } from "@/features/payments/model"
import { useCart } from "@/features/products/hooks/use-cart"
import { env } from "@/env"

// ── Phase components ──────────────────────────────────────────────────────────
import CheckoutConnecting from "@/features/payments/ui/checkout/CheckoutConnecting"
import CheckoutReservationFailed from "@/features/payments/ui/checkout/status/CheckoutReservationFailed"
import CheckoutPartialReservation from "@/features/payments/ui/checkout/status/CheckoutPartialReservation"
import CheckoutError from "@/features/payments/ui/checkout/status/CheckoutError"
import CheckoutOrderConfirmed from "@/features/payments/ui/checkout/status/CheckoutOrderConfirmed"
import CheckoutPixPending from "@/features/payments/ui/checkout/CheckoutPixPending"
import CheckoutPaymentPending from "@/features/payments/ui/checkout/status/CheckoutPaymentPending"
import { useCheckoutSocket } from "@/features/products/hooks/use-checkout-socket"
import CheckoutPaymentForm from "@/features/payments/ui/checkout/CheckoutPaymentForm"
import CheckoutCartChanged from "@/features/payments/ui/checkout/CheckoutCartChanged"
import { editionQueryOptions } from "@/features/editions/api"
import CheckoutPaymentFailed from "@/features/payments/ui/checkout/status/CheckoutPaymentFailed"

export const Route = createFileRoute(
  "/events/$eventId/editions/$editionId/checkout",
)({
  component: CheckoutPage,
})

interface SessionSnapshot {
  sessionId: string
  items: { product_id: string; quantity: number }[]
}

function CheckoutPage() {
  const { eventId, editionId } = Route.useParams()
  const router = useRouter()
  const { items, clearCart, replaceCart } = useCart(editionId)
  const { data: edition } = useQuery(editionQueryOptions(eventId, editionId));

  const wsUrl = `${env.VITE_API_URL.replace("http", "ws")}events/${eventId}/editions/${editionId}/products/purchase`

  const sessionKey = `checkout_session_${editionId}`
  const didInitRef = useRef(false)
  const pendingSnapshotRef = useRef<SessionSnapshot | null>(null)

  // When this is set, we're waiting for the user to decide what to do
  // with their changed cart vs the existing session.
  const [cartChanged, setCartChanged] = useState(false)

  const cartSnapshot = items.map((i) => ({ product_id: i.id, quantity: i.quantity }))

  const cartMatchesSnapshot = (snapshot: SessionSnapshot) => {
    if (snapshot.items.length !== cartSnapshot.length) return false
    return snapshot.items.every((s) => {
      const current = cartSnapshot.find((c) => c.product_id === s.product_id)
      return current?.quantity === s.quantity
    })
  }

  const {
    state,
    buyRequest,
    resumeSession,
    confirmPartial,
    cancelReservation,
    submitPayment,
    reset,
    cancelPurchase
  } = useCheckoutSocket({
    url: wsUrl,
    onPartialReservation: (reserved) => {
      replaceCart(
        reserved.map((item) => ({
          id: item.product_id,
          name: item.name,
          price_cents: item.price_cents,
          quantity: item.quantity,
          inventory_remaining: item.quantity,
          has_inventory: true,
        }))
      )
    },
    onPixCreated: () => {
      sessionStorage.removeItem(sessionKey)
    },
  })

  // On mount: decide between resume, ask, or new purchase.
  useEffect(() => {
    if (didInitRef.current) return
    didInitRef.current = true

    const raw = sessionStorage.getItem(sessionKey)
    if (raw) {
      try {
        const snapshot = JSON.parse(raw) as SessionSnapshot
        if (cartMatchesSnapshot(snapshot)) {
          // Cart unchanged — resume silently.
          resumeSession(snapshot.sessionId)
          return
        }
        // Cart changed — hold the snapshot and ask the user.
        pendingSnapshotRef.current = snapshot
        setCartChanged(true)
        return
      } catch {
        // malformed — fall through to buyRequest
      }
      sessionStorage.removeItem(sessionKey)
    }

    // No session or malformed — start fresh.
    if (items.length > 0) {
      buyRequest(cartSnapshot)
    } else {
      router.history.back()
    }
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  // Persist session ID + cart snapshot whenever the session is established.
  useEffect(() => {
    if (!state.sessionId) return
    const snapshot: SessionSnapshot = { sessionId: state.sessionId, items: cartSnapshot }
    sessionStorage.setItem(sessionKey, JSON.stringify(snapshot))
  }, [state.sessionId]) // eslint-disable-line react-hooks/exhaustive-deps

  // When the server confirms the reservation, sync the cart with the
  // exact items/quantities/prices it reserved — backend is the source of truth.
  useEffect(() => {
    if (state.phase !== "reservation_confirmed") return
    replaceCart(
      state.reservedItems.map((item) => ({
        id: item.product_id,
        name: item.name,
        price_cents: item.price_cents,
        quantity: item.quantity,
        inventory_remaining: item.quantity,
        has_inventory: true,
      }))
    )
  }, [state.phase]) // eslint-disable-line react-hooks/exhaustive-deps

  // Clear session when order is confirmed.
  useEffect(() => {
    if (state.phase === "payment_confirmed") {
      sessionStorage.removeItem(sessionKey)
      clearCart()
    }
  }, [state.phase, sessionKey]) // eslint-disable-line react-hooks/exhaustive-deps

  // session_expired → drop stale session and immediately start a fresh purchase.
  useEffect(() => {
    if (state.phase !== "session_expired") return
    sessionStorage.removeItem(sessionKey)
    reset()
    if (items.length > 0) {
      buyRequest(cartSnapshot)
    } else {
      router.history.back()
    }
  }, [state.phase]) // eslint-disable-line react-hooks/exhaustive-deps

  const handleCancelReservation = () => {
    cancelReservation()
    router.history.back()
  }

  const handlePurchaseCancel = () => {
    cancelPurchase()
    router.history.back()
  }

  const handlePaymentSubmit = (data: SubmitPaymentPayloadI) => {
    submitPayment(data)
  }

  const handleRetry = () => {
    reset()
    if (items.length > 0) {
      buyRequest(cartSnapshot)
    } else {
      router.history.back()
    }
  }

  // ── Cart changed: user chose to resume the old session ────────────────────
  const handleResumeOldSession = () => {
    const snapshot = pendingSnapshotRef.current
    if (!snapshot) return
    setCartChanged(false)
    pendingSnapshotRef.current = null
    resumeSession(snapshot.sessionId)
  }

  // ── Cart changed: user chose to start fresh with the new cart ─────────────
  const handleUseNewCart = () => {
    sessionStorage.removeItem(sessionKey)
    pendingSnapshotRef.current = null
    setCartChanged(false)
    if (items.length > 0) {
      buyRequest(cartSnapshot)
    } else {
      router.history.back()
    }
  }

  const { phase } = state

  // ── Cart changed — ask the user ────────────────────────────────────────────
  if (cartChanged) {
    return (
      <CheckoutCartChanged
        onResume={handleResumeOldSession}
        onUseNew={handleUseNewCart}
        onCancel={() => { router.history.back(); }}
      />
    )
  }

  // ── Idle / connecting / awaiting reservation ───────────────────────────────
  if (phase === "idle" || phase === "connecting" || phase === "awaiting_reservation") {
    return <CheckoutConnecting />
  }

  // ── Reservation failed ─────────────────────────────────────────────────────
  if (phase === "reservation_failed") {
    return (
      <CheckoutReservationFailed
        message={state.errorMessage}
        onBack={() => { router.history.back() }}
      />
    )
  }

  // ── Partial reservation ────────────────────────────────────────────────────
  if (phase === "partial_reservation" && state.partialData) {
    return (
      <CheckoutPartialReservation
        reserved={state.partialData.reserved}
        unavailable={state.partialData.unavailable}
        confirmDeadline={state.partialData.confirmDeadline}
        totalCents={state.totalCents}
        onConfirm={confirmPartial}
        onCancel={handleCancelReservation}
      />
    )
  }

  if (phase === "payment_failed") {
    return (
      <CheckoutPaymentFailed
        message={state.errorMessage}
        onBack={() => { router.history.back() }}
      />
    )
  }

  // ── Generic error ──────────────────────────────────────────────────────────
  if (phase === "error") {
    return (
      <CheckoutError
        message={state.errorMessage}
        onRetry={handleRetry}
        onBack={() => { router.history.back() }}
      />
    )
  }

  // ── Order confirmed ────────────────────────────────────────────────────────
  if (phase === "payment_confirmed") {
    return <CheckoutOrderConfirmed />
  }

  // ── Pix pending ────────────────────────────────────────────────────────────
  if (phase === "pix_pending" && state.pixData) {
    return (
      <CheckoutPixPending
        qrCode={state.pixData.qrCode}
        qrCodeBase64={state.pixData.qrCodeBase64}
        totalCents={state.totalCents}
        onCancel={handlePurchaseCancel}
      />
    )
  }

  // ── Payment pending (card — webhook handling) ──────────────────────────────
  if (phase === "payment_pending") {
    return <CheckoutPaymentPending message={state.pendingMessage} />
  }

  // ── Main payment form ──────────────────────────────────────────────────────
  if (
    (phase === "reservation_confirmed" ||
      phase === "awaiting_payment" ||
      phase === "payment_processing") && edition?.trie_payments_provider_public_key
  ) {
    return (
      <CheckoutPaymentForm
        phase={phase}
        reservedItems={state.reservedItems}
        totalCents={state.totalCents}
        reservationExpiresAt={state.reservationExpiresAt}
        paymentIntentId={state.paymentIntentId}
        sellerPublicKey={edition.trie_payments_provider_public_key}
        onSubmit={handlePaymentSubmit}
        onCancel={handlePurchaseCancel}
        onExpire={() => {/* session_expired comes from server */ }}
      />
    )
  }

  return <CheckoutConnecting />
}