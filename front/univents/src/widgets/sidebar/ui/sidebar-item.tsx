import { Link } from '@tanstack/react-router'
import { cn } from '@/shared/lib/utils'
import type { SidebarMenuItem } from '../sidebar-menu'

interface SidebarItemProps {
  item: SidebarMenuItem
  collapsed: boolean
  pathname: string
}

function isActivePath(pathname: string, href: string, exact?: boolean) {
  const pattern = `^${href
    .split('/')
    .map((segment) => {
      if (!segment) return ''
      if (segment.startsWith('$')) return '[^/]+'
      return segment.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
    })
    .join('/')}${exact ? '$' : '(?:/.*)?$'}`

  return new RegExp(pattern).test(pathname)
}

export function SidebarItem({ item, collapsed, pathname }: SidebarItemProps) {
  const isActive = isActivePath(pathname, item.to, item.exact)

  const baseClasses =
    'group relative flex items-center gap-3 rounded-xl px-3 py-3 text-sm font-medium outline-none transition-colors duration-150 focus-visible:ring-2 focus-visible:ring-ring'

  const stateClasses = isActive
    ? 'text-primary'
    : 'text-muted-foreground hover:bg-muted/70 hover:text-foreground'

  const iconClasses = isActive
    ? 'text-primary'
    : 'text-muted-foreground group-hover:text-foreground'

  return (
    <Link
      to={item.to as never}
      {...(item.params ? { params: item.params as never } : {})}
      preload="intent"
      activeOptions={{ exact: item.exact, includeSearch: false }}
      className={cn(baseClasses, stateClasses)}
      aria-current={isActive ? 'page' : undefined}
    >
      <span
        className={cn(
          'absolute left-0 top-1/2 h-6 w-1 -translate-y-1/2 rounded-r-full bg-primary transition-all duration-200 ease-out',
          isActive ? 'scale-y-100 opacity-100' : 'scale-y-75 opacity-0',
        )}
      />

      <item.icon className={cn('h-4 w-4 shrink-0 transition-colors', iconClasses)} strokeWidth={2} />

      <span
        className={cn(
          'truncate transition-[opacity,width] duration-200',
          collapsed ? 'lg:w-0 lg:opacity-0' : 'w-auto opacity-100',
        )}
      >
        {item.label}
      </span>

      {collapsed && (
        <span className="pointer-events-none absolute left-full top-1/2 z-50 ml-3 hidden -translate-y-1/2 whitespace-nowrap rounded-md bg-primary px-2.5 py-1.5 text-xs font-medium text-primary-foreground opacity-0 shadow-md transition-opacity duration-150 group-hover:opacity-100 lg:block">
          {item.label}
        </span>
      )}
    </Link>
  )
}
