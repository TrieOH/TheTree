import { Link } from '@tanstack/react-router'
import {
  LayoutDashboard,
} from 'lucide-react'
import { cn } from '#/shared/lib/class-utils'

interface WorkspaceLayoutProps {
  children: React.ReactNode
}

export default function AdminLayout({ children }: WorkspaceLayoutProps) {
  const navItems = [
    {
      to: '/',
      icon: LayoutDashboard,
      label: 'Overview',
      exact: true
    },
    // {
    //   to: '/admin/',
    //   icon: Key,
    //   label: 'API Keys'
    // },
    // {
    //   to: '/admin/$name/webhooks',
    //   icon: Webhook,
    //   label: 'Webhooks'
    // },
    // {
    //   to: '/admin/$name/providers',
    //   icon: ArrowRightFromLine,
    //   label: 'Providers'
    // },
  ]

  return (
    <div className="flex flex-col lg:flex-row gap-8 items-start">
      {/* Desktop Sidebar Navigation */}
      <aside className="hidden lg:block w-48 sticky top-24 shrink-0">
        <nav className="flex lg:flex-col border-r border-border/60">
          {navItems.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              activeOptions={{ exact: item.exact }}
              className="flex items-center gap-3 px-4 py-4 text-[10px] font-black uppercase tracking-[0.2em] transition-colors relative group"
            >
              {({ isActive }) => (
                <>
                  <item.icon className={cn(
                    "w-3.5 h-3.5 transition-colors duration-300",
                    isActive ? "text-primary" : "text-muted-foreground group-hover:text-foreground"
                  )} />
                  <span className={cn(
                    "transition-colors duration-300",
                    isActive ? "text-foreground" : "text-muted-foreground group-hover:text-foreground"
                  )}>
                    {item.label}
                  </span>

                  {/* Desktop Indicator (Right) */}
                  <div className={cn(
                    "absolute -right-px transition-all duration-300 ease-in-out bg-primary w-0.5",
                    isActive ? "top-2 bottom-2 opacity-100" : "top-1/2 bottom-1/2 opacity-0"
                  )} />
                </>
              )}
            </Link>
          ))}
        </nav>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 min-w-0 w-full pb-24 lg:pb-0">
        {children}
      </div>

      {/* Mobile Bottom Navigation */}
      <nav className="lg:hidden fixed bottom-0 left-0 right-0 z-40 flex h-16 items-center justify-around border-t border-border bg-background/95 backdrop-blur-md px-4">
        {navItems.map((item) => (
          <Link
            key={item.to}
            to={item.to}
            activeOptions={{ exact: item.exact }}
            className="flex flex-col items-center gap-1 px-3 relative h-full justify-center group"
          >
            {({ isActive }) => (
              <>
                <item.icon className={cn(
                  "w-5 h-5 transition-colors",
                  isActive ? "text-primary" : "text-muted-foreground group-hover:text-foreground"
                )} />
                <span className={cn(
                  "text-[9px] font-black uppercase tracking-tighter transition-colors",
                  isActive ? "text-primary" : "text-muted-foreground group-hover:text-foreground"
                )}>
                  {item.label}
                </span>
                {/* Mobile Indicator (Top of the bar) */}
                <div className={cn(
                  "absolute top-0 left-1/2 -translate-x-1/2 w-10 h-1 bg-primary rounded-b-full transition-all duration-300 ease-in-out",
                  isActive ? "opacity-100 scale-x-100" : "opacity-0 scale-x-0"
                )} />
              </>
            )}
          </Link>
        ))}
      </nav>
    </div>
  )
}
