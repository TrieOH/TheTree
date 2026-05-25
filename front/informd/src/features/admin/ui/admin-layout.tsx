import { Link } from '@tanstack/react-router'
import {
  LogOut,
  ChevronLeft,
  ChevronRight,
  FolderKanban,
  KeySquare,
  FileText
} from 'lucide-react'
import { useState } from 'react'
import { useAuthActions } from '#/features/auths/hooks/use-auth-actions'
import { cn } from '#/shared/lib/utils'
import { motion } from 'motion/react'
import { Button } from '#/shared/ui/shadcn/button'
import { Breadcrumb } from '#/shared/ui/breadcrumb'

export function AdminLayout({ children }: { children: React.ReactNode }) {
  const { handleLogout } = useAuthActions()
  const [isCollapsed, setIsCollapsed] = useState(false)

  const navItems = [
    {
      to: '/admin',
      icon: FolderKanban,
      label: 'Namespaces',
      exact: true
    },
    {
      to: '/admin/form',
      icon: FileText,
      label: 'Forms',
      exact: true
    },
    {
      to: '/admin/keys',
      icon: KeySquare,
      label: 'API Keys',
      exact: true
    },
  ]

  return (
    <div className="flex min-h-screen bg-background font-sans selection:bg-primary/10">
      {/* Desktop Sidebar Navigation */}
      <motion.aside
        initial={false}
        animate={{ width: isCollapsed ? 60 : 200 }}
        className="hidden lg:flex flex-col border-r border-border/60 sticky top-0 h-screen shrink-0 bg-background z-20"
      >
        <div className="p-4 flex items-center justify-between h-16 border-b border-border/60">
          {!isCollapsed && (
            <span className="text-xs font-bold truncate uppercase tracking-[0.3em] text-primary">
              Informd
            </span>
          )}
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setIsCollapsed(!isCollapsed)}
            className={cn('hover:bg-transparent', isCollapsed ? 'mx-auto' : 'ml-auto')}
          >
            {isCollapsed ? (
              <ChevronRight className="w-4 h-4" />
            ) : (
              <ChevronLeft className="w-4 h-4" />
            )}
          </Button>
        </div>

        <nav className="flex-1 py-4 flex flex-col">
          {navItems.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              activeOptions={{ exact: item.exact, includeSearch: false }}
              className="flex items-center gap-3 px-4 py-4 text-[10px] font-bold uppercase tracking-[0.2em] transition-colors relative group"
            >
              {({ isActive }) => (
                <>
                  <item.icon
                    className={cn(
                      'w-4 h-4 transition-colors duration-300',
                      isCollapsed && 'mx-auto',
                      isActive
                        ? 'text-primary'
                        : 'text-muted-foreground group-hover:text-foreground',
                    )}
                  />
                  {!isCollapsed && (
                    <span
                      className={cn(
                        'transition-colors duration-300 truncate',
                        isActive
                          ? 'text-foreground'
                          : 'text-muted-foreground group-hover:text-foreground',
                      )}
                    >
                      {item.label}
                    </span>
                  )}

                  {/* Desktop Indicator (Right) */}
                  <div
                    className={cn(
                      'absolute -right-px transition-all duration-300 ease-in-out bg-primary w-0.5',
                      isActive ? 'top-2 bottom-2 opacity-100' : 'top-1/2 bottom-1/2 opacity-0',
                    )}
                  />
                </>
              )}
            </Link>
          ))}

          <button
            onClick={handleLogout}
            className="mt-auto flex items-center gap-3 px-4 py-4 text-[10px] font-bold uppercase tracking-[0.2em] transition-colors relative group text-muted-foreground hover:text-destructive cursor-pointer"
          >
            <LogOut
              className={cn(
                'w-4 h-4 transition-colors duration-300',
                isCollapsed && 'mx-auto',
              )}
            />
            {!isCollapsed && <span className='truncate'>Logout</span>}
          </button>
        </nav>
      </motion.aside>

      {/* Main Content Area */}
      <div className="flex-1 min-w-0 w-full pb-24 lg:pb-0">
        <div className="sticky top-0 z-10">
          <Breadcrumb />
        </div>
        <main>{children}</main>
      </div>

      {/* Mobile Bottom Navigation */}
      <nav className="lg:hidden! fixed bottom-0 left-0 right-0 z-40 flex h-16 items-center justify-around border-t border-border bg-background/95 backdrop-blur-md px-4">
        {navItems.map((item) => (
          <Link
            key={item.to}
            to={item.to}
            activeOptions={{ exact: item.exact, includeSearch: false }}
            className="flex flex-col items-center gap-1 px-3 relative h-full justify-center group"
          >
            {({ isActive }) => (
              <>
                <item.icon
                  className={cn(
                    'w-5 h-5 transition-colors',
                    isActive
                      ? 'text-primary'
                      : 'text-muted-foreground group-hover:text-foreground',
                  )}
                />
                <span
                  className={cn(
                    'text-[9px] font-bold uppercase tracking-tighter transition-colors truncate',
                    isActive
                      ? 'text-primary'
                      : 'text-muted-foreground group-hover:text-foreground',
                  )}
                >
                  {item.label}
                </span>
                {/* Mobile Indicator (Top of the bar) */}
                <div
                  className={cn(
                    'absolute top-0 left-1/2 -translate-x-1/2 w-10 h-1 bg-primary rounded-b-full transition-all duration-300 ease-in-out',
                    isActive ? 'opacity-100 scale-x-100' : 'opacity-0 scale-x-0',
                  )}
                />
              </>
            )}
          </Link>
        ))}
        <button
          onClick={handleLogout}
          className="flex flex-col items-center gap-1 px-3 justify-center group cursor-pointer"
        >
          <LogOut className="w-5 h-5 text-muted-foreground group-hover:text-destructive transition-colors" />
          <span className="text-[9px] truncate font-bold uppercase tracking-tighter text-muted-foreground group-hover:text-destructive transition-colors">
            Logout
          </span>
        </button>
      </nav>
    </div>
  )
}
