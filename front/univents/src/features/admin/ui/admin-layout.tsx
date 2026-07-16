import { useLocation } from '@tanstack/react-router'
import {
  SidebarProvider,
  useSidebar,
} from '@/widgets/sidebar/hooks/use-sidebar'
import { MobileTopbar } from '@/widgets/sidebar/ui/mobile-topbar'
import { Sidebar } from '@/widgets/sidebar/ui/sidebar'
import { Breadcrumb } from '@/shared/ui/breadcrumb'

interface AdminLayoutProps {
  children: React.ReactNode
}

function AppShell({ children }: { children: React.ReactNode }) {
  const { collapsed } = useSidebar()
  const { pathname } = useLocation()
  const isCertificateEditor = pathname.endsWith('/certifications/editor')

  if (isCertificateEditor) {
    return <div className="h-dvh overflow-hidden bg-background">{children}</div>
  }

  return (
    <div className="min-h-dvh bg-background">
      <Sidebar />

      <div
        className={
          'flex min-h-dvh flex-col transition-[padding] duration-300 ease-in-out ' +
          (collapsed ? 'lg:pl-18' : 'lg:pl-72')
        }
      >
        <MobileTopbar />
        <div className="sticky top-0 z-30 hidden bg-card/95 shadow-sm shadow-black/5 lg:block">
          <Breadcrumb />
        </div>
        <main className="flex-1">{children}</main>
      </div>
    </div>
  )
}

export function AdminLayout({ children }: AdminLayoutProps) {
  return (
    <SidebarProvider>
      <AppShell>{children}</AppShell>
    </SidebarProvider>
  )
}
