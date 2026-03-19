// shared/ui.tsx - Shared components & design tokens

export const inputClass =
  'w-full px-3 py-2 text-sm border border-gray-200 rounded-lg bg-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent transition-all'

export const labelClass = 'block text-xs font-medium text-gray-500 uppercase tracking-wide mb-1'

export const btnPrimary =
  'inline-flex items-center gap-2 px-4 py-2 bg-gray-900 text-white text-sm font-medium rounded-lg hover:bg-gray-700 transition-colors disabled:opacity-40 disabled:cursor-not-allowed'

export const btnSecondary =
  'inline-flex items-center gap-2 px-4 py-2 bg-white text-gray-700 text-sm font-medium rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors'

export const btnDanger =
  'inline-flex items-center gap-2 px-4 py-2 bg-red-50 text-red-600 text-sm font-medium rounded-lg border border-red-100 hover:bg-red-100 transition-colors'

export const cardClass = 'bg-white border border-gray-100 rounded-xl p-4 shadow-sm hover:shadow-md transition-shadow'

export function PageHeader({
  title,
  subtitle,
  action,
}: {
  title: string
  subtitle?: string
  action?: React.ReactNode
}) {
  return (
    <div className="flex items-start justify-between mb-8">
      <div>
        <h1 className="text-2xl font-semibold text-gray-900 tracking-tight">{title}</h1>
        {subtitle && <p className="mt-1 text-sm text-gray-500">{subtitle}</p>}
      </div>
      {action && <div>{action}</div>}
    </div>
  )
}

export function FormField({
  label,
  children,
  required,
}: {
  label: string
  children: React.ReactNode
  required?: boolean
}) {
  return (
    <div>
      <label className={labelClass}>
        {label}
        {required && <span className="text-gray-400 ml-0.5">*</span>}
      </label>
      {children}
    </div>
  )
}

export function ErrorMsg({ msg }: { msg: string | null }) {
  if (!msg) return null
  return (
    <div className="flex items-center gap-2 text-sm text-red-600 bg-red-50 border border-red-100 rounded-lg px-3 py-2">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
        <circle cx="12" cy="12" r="10" />
        <line x1="12" y1="8" x2="12" y2="12" />
        <line x1="12" y1="16" x2="12.01" y2="16" />
      </svg>
      {msg}
    </div>
  )
}

export function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    draft: 'bg-amber-50 text-amber-700 border-amber-100',
    published: 'bg-green-50 text-green-700 border-green-100',
    active: 'bg-blue-50 text-blue-700 border-blue-100',
    archived: 'bg-gray-100 text-gray-500 border-gray-200',
  }
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 text-xs font-medium rounded-md border ${colors[status] ?? 'bg-gray-100 text-gray-600 border-gray-200'}`}
    >
      {status}
    </span>
  )
}

export function EmptyState({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) {
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <div className="w-12 h-12 rounded-xl bg-gray-100 flex items-center justify-center text-gray-400 mb-4">
        {icon}
      </div>
      <p className="text-sm font-medium text-gray-700">{title}</p>
      <p className="text-xs text-gray-400 mt-1 max-w-xs">{description}</p>
    </div>
  )
}

export function AdminShell({ children, breadcrumbs }: { children: React.ReactNode; breadcrumbs?: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-gray-50 font-sans">
      <nav className="bg-white border-b border-gray-100 px-6 py-3 flex items-center gap-3">
        <span className="text-sm font-semibold text-gray-900 tracking-tight">Admin</span>
        {breadcrumbs && (
          <>
            <span className="text-gray-300">/</span>
            {breadcrumbs}
          </>
        )}
      </nav>
      <main className="max-w-5xl mx-auto px-6 py-8">{children}</main>
    </div>
  )
}