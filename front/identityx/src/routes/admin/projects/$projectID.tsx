import { LayoutContext } from '@/shared/lib/hooks/layout-context'
import { cn } from '@/shared/lib/utils'
import { createFileRoute, Link, Outlet } from '@tanstack/react-router'
import { Workflow } from 'lucide-react'
import { useState } from 'react'
import z from 'zod'

const projectSearchSchema = z.object({
  organizationID: z.string().optional(),
})

export const Route = createFileRoute('/admin/projects/$projectID')({
  validateSearch: (search) => projectSearchSchema.parse(search),
  component: ProjectLayout,
})

function ProjectLayout() {
  const { projectID } = Route.useParams()
  const { organizationID } = Route.useSearch()

  const [headerSlot, setHeaderSlot] = useState<React.ReactNode>(null)

  const tabs = [
    {
      label: 'Main',
      to: '/admin/projects/$projectID',
      params: { projectID },
      icon: Workflow,
      exact: true,
    },
  ]

  return (
    <LayoutContext.Provider value={{ setHeader: setHeaderSlot }}>
      <div className="flex flex-col h-full">
        {/* Page Header Slot */}
        {/*
          Rendered only when a child page calls useLayoutHeader().
          Sits between the tab bar and the page content.
        */}
        {headerSlot && (
          <div className="border-b border-border/40 px-6 py-4 bg-background">
            {headerSlot}
          </div>
        )}

        {/* Tab Bar */}
        <div className="border-b border-border/60 bg-background/50 px-6 overflow-x-auto scrollbar-none">
          <div className="flex items-center gap-8 h-12 min-w-max">
            {tabs.map((tab) => (
              <Link
                key={tab.label}
                to={tab.to}
                params={tab.params}
                search={{ organizationID }}
                activeOptions={{ exact: tab.exact, includeSearch: false }}
                className="relative h-full flex items-center gap-2 text-[10px] font-bold uppercase tracking-widest transition-colors group"
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
                        'transition-colors whitespace-nowrap',
                        isActive
                          ? 'text-foreground'
                          : 'text-muted-foreground group-hover:text-foreground',
                      )}
                    >
                      {tab.label}
                    </span>
                    {isActive && (
                      <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary" />
                    )}
                  </>
                )}
              </Link>
            ))}
          </div>
        </div>

        {/* Page Content */}
        <div className="flex-1 p-6">
          <Outlet />
        </div>

      </div>
    </LayoutContext.Provider>
  )
}
