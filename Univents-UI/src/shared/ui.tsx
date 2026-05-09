export const labelClass = 'block text-xs font-medium text-gray-500 uppercase tracking-wide mb-1'

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