import { EnvironmentSelector } from "#/features/environment/ui/environment-selector"
import { env } from "#/env"

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
      <div className="max-w-7xl mx-auto px-4 flex h-(--header-height) items-center justify-between">
        <div className="flex items-center gap-4">
          <span className="font-bold hidden sm:inline-block">
            {env.VITE_APP_TITLE || "SpiceDB UI"}
          </span>
        </div>
        <div className="flex items-center gap-2">
          <EnvironmentSelector />
        </div>
      </div>
    </header>
  )
}
