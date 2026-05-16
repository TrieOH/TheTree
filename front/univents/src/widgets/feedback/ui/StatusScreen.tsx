interface StatusAction {
  label: string
  onClick: () => void
  variant: "primary" | "outline"
}

interface StatusScreenProps {
  icon: string
  iconClass: string
  title: string
  description: string
  actions: StatusAction[]
}

export default function StatusScreen({
  icon,
  iconClass,
  title,
  description,
  actions
}: StatusScreenProps) {
  return (
    <main className="w-full min-w-75 max-w-sm mx-auto px-3 py-16 flex flex-col items-center gap-5 text-center">
      <div className={`w-14 h-14 rounded-full flex items-center justify-center text-xl font-bold ${iconClass}`}>
        {icon}
      </div>
      <div className="space-y-1">
        <h1 className="text-lg font-bold text-foreground">{title}</h1>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
      {actions.length > 0 && (
        <div className="flex flex-col gap-2 w-full">
          {actions.map((action) => (
            <button
              key={action.label}
              onClick={action.onClick}
              className={
                action.variant === "primary"
                  ? "w-full rounded-md bg-primary text-primary-foreground px-4 py-2.5 text-sm font-medium hover:bg-primary/90 transition-colors"
                  : "w-full rounded-md border border-border px-4 py-2.5 text-sm text-muted-foreground hover:bg-muted/50 transition-colors"
              }
            >
              {action.label}
            </button>
          ))}
        </div>
      )}
    </main>
  )
}