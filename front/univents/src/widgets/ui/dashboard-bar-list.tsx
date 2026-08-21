interface DashboardBarListItem {
  id: string;
  label: string;
  value: number;
  detail?: string;
  color?: string;
}

interface DashboardBarListProps {
  items: DashboardBarListItem[];
  maxValue?: number;
  emptyMessage?: string;
}

export function DashboardBarList({
  items,
  maxValue,
  emptyMessage = "Nenhum dado disponível.",
}: DashboardBarListProps) {
  const maximum =
    maxValue ?? Math.max(...items.map((item) => Number(item.value)), 1);

  if (items.length === 0) {
    return <p className="text-sm text-muted-foreground">{emptyMessage}</p>;
  }

  return (
    <div className="space-y-4">
      {items.map((item) => {
        const numericValue = Number(item.value);
        const width = numericValue
          ? Math.max((numericValue / maximum) * 100, 3)
          : 0;

        return (
          <div key={item.id} className="space-y-1.5">
            <div className="flex items-center justify-between gap-3 text-xs">
              <span className="min-w-0 truncate font-medium">{item.label}</span>
              <span className="shrink-0 text-muted-foreground">
                {item.detail ?? item.value}
              </span>
            </div>
            <div className="h-2 overflow-hidden rounded-full bg-muted">
              <div
                className={`h-full rounded-full ${item.color ?? "bg-primary"}`}
                style={{ width: `${width}%` }}
              />
            </div>
          </div>
        );
      })}
    </div>
  );
}
