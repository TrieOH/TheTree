import type * as React from 'react'
import { Link } from '@tanstack/react-router'
import {
  LayoutDashboard,
  LogOut,
} from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { useAuth } from '@trieoh/node-auth-sdk/react'
import { cn } from '#/shared/lib/utils'

export function AdminLayout({ children }: { children: React.ReactNode }) {
  const { auth } = useAuth()

  return (
    <div className="flex min-h-screen bg-background font-sans selection:bg-primary/10">
      <div className="flex-1 flex flex-col min-w-0">
        {/* Header */}
        <header className="sticky top-0 z-40 flex h-16 items-center justify-between border-b border-border bg-background/80 backdrop-blur-md px-4 md:px-8 lg:px-10">
          <div className="flex items-center gap-8">
            <Link to="/admin" className="flex items-center gap-2.5 group transition-all hover:opacity-80">
              <div className="w-8 h-8 rounded-none bg-primary flex items-center justify-center text-primary-foreground group-hover:bg-primary/90 transition-all">
                <LayoutDashboard className="w-4 h-4" />
              </div>
              <span className="font-black text-xl tracking-tighter text-foreground">
                TriePayments
              </span>
            </Link>
          </div>

          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="sm"
              className={cn(
                "hidden sm:flex items-center gap-2 text-muted-foreground",
                "hover:text-destructive hover:bg-destructive/5 font-bold uppercase tracking-wider",
                "rounded-sm cursor-pointer py-4"
              )}
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
      </div>
    </div>
  )
}
