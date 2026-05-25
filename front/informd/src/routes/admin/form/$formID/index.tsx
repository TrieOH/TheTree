import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import { createFileRoute } from '@tanstack/react-router'
import { useMemo } from 'react'

export const Route = createFileRoute('/admin/form/$formID/')({
  component: RouteComponent,
})

function RouteComponent() {
  const header = useMemo(() => (
    <div className="flex items-start justify-between">
      <div>
        <h1 className="text-lg font-semibold tracking-tight">Steps</h1>
        <p className="text-sm text-muted-foreground">
          No steps yet in this form
        </p>
      </div>
    </div>
  ), [])

  useLayoutHeader(header)
  return <div>Hello "/admin/form/$formID/"!</div>
}
