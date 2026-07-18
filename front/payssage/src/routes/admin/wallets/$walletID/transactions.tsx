import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/wallets/$walletID/transactions')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/wallets/$walletID/transactions"!</div>
}
