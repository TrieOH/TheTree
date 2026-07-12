import { createFileRoute, Link, Outlet } from '@tanstack/react-router'
import { CalendarDays, LayoutGrid } from 'lucide-react'
import { cn } from '@/shared/lib/utils'

export const Route = createFileRoute('/admin/events/$eventId')({
  component: EventAdminLayout,
})

function EventAdminLayout() {
  const { eventId } = Route.useParams()

  const tabs = [
    {
      label: 'Visão Geral',
      to: '/admin/events/$eventId',
      params: { eventId },
      icon: LayoutGrid,
      exact: true,
    },
    {
      label: 'Edições',
      to: '/admin/events/$eventId/editions',
      params: { eventId },
      icon: CalendarDays,
      exact: true,
    },
  ]

  return (
    <div className="flex min-h-full flex-col">
      <div className="border-b border-border/60 bg-background/70 px-6 backdrop-blur">
        <div className="flex h-12 items-center gap-8">
          {tabs.map((tab) => (
            <Link
              key={tab.label}
              to={tab.to}
              params={tab.params}
              activeOptions={{ exact: tab.exact }}
              className="group relative flex h-full items-center gap-2 text-[10px] font-bold uppercase tracking-widest transition-colors"
            >
              {({ isActive }) => (
                <>
                  <tab.icon
                    className={cn(
                      'size-3.5 transition-colors',
                      isActive
                        ? 'text-primary'
                        : 'text-muted-foreground group-hover:text-foreground',
                    )}
                  />
                  <span
                    className={cn(
                      'transition-colors',
                      isActive
                        ? 'text-foreground'
                        : 'text-muted-foreground group-hover:text-foreground',
                    )}
                  >
                    {tab.label}
                  </span>
                  {isActive && <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary" />}
                </>
              )}
            </Link>
          ))}
        </div>
      </div>

      <div className="flex-1 p-6">
        <Outlet />
      </div>
    </div>
  )
}
