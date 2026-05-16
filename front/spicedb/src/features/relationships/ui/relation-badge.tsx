export function RelationBadge({ value }: { value: string }) {
  const colors: Record<string, string> = {
    viewer: 'bg-blue-50   text-blue-700   dark:bg-blue-950   dark:text-blue-300',
    editor: 'bg-violet-50 text-violet-700 dark:bg-violet-950 dark:text-violet-300',
    owner: 'bg-amber-50  text-amber-700  dark:bg-amber-950  dark:text-amber-300',
    member: 'bg-teal-50   text-teal-700   dark:bg-teal-950   dark:text-teal-300',
    admin: 'bg-rose-50   text-rose-700   dark:bg-rose-950   dark:text-rose-300',
  }
  return (
    <span
      className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${colors[value] ?? 'bg-muted text-muted-foreground'}`}
    >
      {value}
    </span>
  )
}
