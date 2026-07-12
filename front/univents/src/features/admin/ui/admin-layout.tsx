import { Breadcrumb } from '@/shared/ui/breadcrumb'

interface AdminLayoutProps {
  children: React.ReactNode
}

export function AdminLayout({ children }: AdminLayoutProps) {
  return (
    <div className="flex min-h-screen bg-background font-body selection:bg-primary/10">
      {/* Main Content Area */}
      <div className="flex-1 min-w-0 w-full pb-24">
        <div className="sticky top-0 z-10">
          <Breadcrumb />
        </div>
        <main>{children}</main>
      </div>
    </div>
  )
}
