import * as React from 'react'
import { Link } from '@tanstack/react-router'
import {
  LayoutDashboard,
  Key,
  Webhook,
  Layers,
  LogOut,
} from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { useAuth } from '@trieoh/node-auth-sdk/react'

export function AdminLayout({ children }: { children: React.ReactNode }) {
  const { auth } = useAuth()

  const navItems = [
    { to: '/admin', icon: Layers, label: 'Workspaces' },
    { to: '/admin/webhooks', icon: Webhook, label: 'Webhooks' },
    { to: '/admin/keys', icon: Key, label: 'API Keys' },
  ]

  return (
    <div className="flex min-h-screen bg-background font-sans selection:bg-primary/10">
      <div className="flex-1 flex flex-col min-w-0">
        {/* Header */}
        <header className="sticky top-0 z-40 flex h-16 items-center justify-between border-b border-border bg-background/80 backdrop-blur-md px-4 md:px-8 lg:px-10">
          <div className="flex items-center gap-8">
            <div className="flex items-center gap-2.5">
              <div className="w-8 h-8 rounded-none bg-primary flex items-center justify-center text-primary-foreground">
                <LayoutDashboard className="w-4 h-4" />
              </div>
              <span className="font-black text-xl tracking-tighter uppercase text-foreground">Trie</span>
            </div>

            {/* Desktop Navigation */}
            <nav className="hidden lg:flex items-center h-16">
              {navItems.map((item) => (
                <Link
                  key={item.to}
                  to={item.to}
                  activeOptions={item.to === '/admin' ? { exact: true } : undefined}
                  activeProps={{ className: 'text-primary border-b-2 border-primary h-full flex items-center' }}
                  inactiveProps={{ className: 'text-muted-foreground hover:text-foreground h-full flex items-center' }}
                  className="px-4 text-sm font-bold uppercase tracking-wider transition-all"
                >
                  {item.label}
                </Link>
              ))}
            </nav>
          </div>

          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="sm"
              className="hidden sm:flex items-center gap-2 text-muted-foreground hover:text-destructive hover:bg-destructive/5 font-bold uppercase tracking-wider"
              onClick={() => auth.logout()}
            >
              <LogOut className="w-4 h-4" />
              <span>Logout</span>
            </Button>
            
            {/* Mobile Logout (Icon only) */}
            <Button
              variant="ghost"
              size="icon"
              className="sm:hidden text-muted-foreground"
              onClick={() => auth.logout()}
            >
              <LogOut className="w-4 h-4" />
            </Button>
          </div>
        </header>

        <main className="flex-1 overflow-x-hidden">
          <div className="p-4 sm:p-6 md:p-8 lg:p-10 max-w-7xl mx-auto">
            {children}
          </div>
        </main>

        {/* Mobile Navigation (Bottom Bar) */}
        <nav className="lg:hidden sticky bottom-0 z-40 flex h-16 items-center justify-around border-t border-border bg-background/95 backdrop-blur-md px-4">
          {navItems.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              activeOptions={item.to === '/admin' ? { exact: true } : undefined}
              className="flex flex-col items-center gap-1 px-3 relative h-full justify-center group"
            >
              {({ isActive }) => (
                <>
                  <item.icon className={`w-5 h-5 transition-colors ${isActive ? 'text-primary' : 'text-muted-foreground group-hover:text-foreground'}`} />
                  <span className={`text-[9px] font-black uppercase tracking-tighter transition-colors ${isActive ? 'text-primary' : 'text-muted-foreground group-hover:text-foreground'}`}>
                    {item.label}
                  </span>
                  {isActive && (
                    <div className="absolute top-0 left-1/2 -translate-x-1/2 w-10 h-1 bg-primary rounded-b-full" />
                  )}
                </>
              )}
            </Link>
          ))}
        </nav>
      </div>
    </div>
  )
}
