import { cn } from '@/shared/lib/utils'

type TabValue = 'overview' | 'activities' | 'tickets'

interface EditionTabsProps {
  activeTab: TabValue
  onTabChange: (tab: TabValue) => void
  counts: {
    activities: number
    tickets: number
  }
}

const tabs = [
  { value: 'overview', label: 'Visão Geral' },
  { value: 'activities', label: 'Atividades', countKey: 'activities' as const },
  { value: 'tickets', label: 'Tickets', countKey: 'tickets' as const },
] as const

export function EditionTabs({ activeTab, onTabChange, counts }: EditionTabsProps) {
  return (
    <div className="border-b border-border">
      <div className="flex gap-1 overflow-x-auto scrollbar-hide">
        {tabs.map((tab) => (
          <button
            key={tab.value}
            onClick={() => { onTabChange(tab.value as TabValue); }}
            className={cn(
              "relative px-4 py-3 text-sm font-medium transition-colors whitespace-nowrap",
              "hover:text-foreground focus-visible:outline-none",
              activeTab === tab.value
                ? "text-foreground"
                : "text-muted-foreground"
            )}
          >
            <span className="flex items-center gap-2">
              {tab.label}
              {'countKey' in tab && (
                <span className={cn(
                  "px-1.5 py-0.5 rounded-full text-xs",
                  activeTab === tab.value
                    ? "bg-primary/10 text-primary"
                    : "bg-muted text-muted-foreground"
                )}>
                  {counts[tab.countKey]}
                </span>
              )}
            </span>
            {activeTab === tab.value && (
              <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary rounded-full" />
            )}
          </button>
        ))}
      </div>
    </div>
  )
}