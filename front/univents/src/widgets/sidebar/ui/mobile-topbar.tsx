import { ChevronLeft, Menu } from 'lucide-react'
import { Link, useRouterState } from '@tanstack/react-router'
import { useSidebar } from '../hooks/use-sidebar'
import { getAdminBackLink, getAdminShellLabel } from '../sidebar-menu'
import { cn } from '@/shared/lib/utils'

export function MobileTopbar() {
  const { setMobileOpen } = useSidebar()
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  })
  const label = getAdminShellLabel(pathname)
  const backLink = getAdminBackLink(pathname)

  return (
    <header className="sticky top-0 z-30 flex h-16 shrink-0 items-center border-b border-border/60 bg-card/95 px-3 shadow-sm shadow-black/5 backdrop-blur-xl lg:hidden!">
      <div className="flex min-w-0 flex-1 items-center gap-2">
        {backLink ? (
          <Link
            to={backLink.to as never}
            {...(backLink.params ? { params: backLink.params as never } : {})}
            preload="intent"
            className={cn(
              'inline-flex size-10 items-center justify-center rounded-xl text-foreground transition-colors hover:bg-muted/60',
            )}
            aria-label="Voltar"
          >
            <ChevronLeft className="h-5 w-5" />
          </Link>
        ) : (
          <button
            type="button"
            onClick={() => setMobileOpen(true)}
            className={cn(
              'inline-flex size-10 items-center justify-center rounded-xl text-foreground transition-colors hover:bg-muted/60',
            )}
            aria-label="Abrir menu"
          >
            <Menu className="h-5 w-5" />
          </button>
        )}

        <div className="min-w-0">
          <h1 className="truncate text-sm font-semibold text-foreground">
            {label.title}
          </h1>
        </div>
      </div>

      <button
        type="button"
        onClick={() => setMobileOpen(true)}
        className={cn(
          'inline-flex size-10 items-center justify-center rounded-xl text-foreground transition-colors hover:bg-muted/60',
          backLink ? 'opacity-100' : 'hidden',
        )}
        aria-label="Abrir menu"
      >
        <Menu className="h-5 w-5" />
      </button>
    </header>
  )
}
