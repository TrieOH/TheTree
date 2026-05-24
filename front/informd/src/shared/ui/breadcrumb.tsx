import { Link, useLocation } from '@tanstack/react-router'
import { ChevronRight } from 'lucide-react'
import { Fragment } from 'react'
import { cn } from '../lib/utils'

export function Breadcrumb() {
  const { pathname } = useLocation()

  const segments = pathname.split('/').filter(Boolean)

  return (
    <nav
      className={cn(
        'flex items-center space-x-2 text-muted-foreground',
        'font-bold uppercase tracking-[0.2em] text-[10px]',
        'px-6 h-16 border-b border-border/60',
        'bg-background/95 backdrop-blur-md',
      )}
    >
      {segments.map((segment, index) => {
        const isLast = index === segments.length - 1
        const path = `/${segments.slice(0, index + 1).join('/')}`

        // Capitalize and format label
        const label = segment.charAt(0).toUpperCase() + segment.slice(1)

        return (
          <Fragment key={path}>
            {index > 0 && <ChevronRight className="h-3 w-3 text-muted-foreground/40" />}
            {isLast ? (
              <span className="text-foreground">{label}</span>
            ) : (
              <Link to={path} className="hover:text-primary transition-colors">
                {label}
              </Link>
            )}
          </Fragment>
        )
      })}
    </nav>
  )
}
