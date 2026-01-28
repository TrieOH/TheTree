import { useMatches } from "@tanstack/react-router"
import { useState } from "react"
import { headerRegistry } from "../model/header.registry"
import MobileMenu from "./MobileMenu"
import { Menu, X } from "lucide-react"
import HeaderActionRenderer from "./HeaderActionRenderer"
import { cn } from "@/shared/lib/utils"

export default function Header() {
  const matches = useMatches()
  const headerVariant = [...matches]
    .reverse()
    .find(m => m.staticData.components.header)?.staticData.components.header

  const [isMenuOpen, setMenuOpen] = useState(false)

  if (!headerVariant || headerVariant === 'none') return null

  const header = headerRegistry[headerVariant]
  const mobileActions = [
    ...(header.leftActions ?? []).filter(a => a.visibleOn !== 'desktop'),
    ...(header.centerActions ?? []).filter(a => a.visibleOn !== 'desktop'),
    ...(header.rightActions ?? []).filter(a => a.visibleOn !== 'desktop'),
  ]

  return (
    <header className="sticky top-0 w-full z-10">
      <div 
        className={cn(
          "flex justify-between items-center border-b border-b-border px-6 py-4",
          "bg-background/80 backdrop-blur-sm select-none"
        )}
      >
        <div className="flex items-center gap-1.5">
          <button
            type="button"
            aria-label="Open menu" 
            className={cn(
              "md:hidden block active:scale-95 active:translate-y-px",
              "cursor-pointer transition-transform duration-100 ease-out"
            )}
            onClick={() => setMenuOpen(v => !v)}
          >
            { isMenuOpen ? <X size={24} /> : <Menu size={24} /> }
          </button>

          {/* left actions (e.g., back) */}
          <div className="hidden md:flex items-center gap-2">
            {header.leftActions?.filter(a => a.visibleOn !== "mobile").map((a, i) => (
              <HeaderActionRenderer key={`left-${i}`} action={a} />
            ))}
            {header.title && 
              <h2 className="text-2xl font-semibold text-foreground">
                {header.title}
              </h2>
            }
          </div>

          {/* mobile title next to hamburger */}
          <div className="md:hidden">
            {header.title && 
              <h2 className="text-lg font-semibold text-foreground">
                {header.title}
              </h2>
            }
          </div>
        </div>

        {/* center desktop nav */}
        <nav className="hidden md:flex gap-6">
          {header.centerActions?.filter(a => a.visibleOn !== "mobile").map((a, i) => (
            <HeaderActionRenderer key={`center-${i}`} action={a} />
          ))}
        </nav>

        {/* right actions */}
        <div className="flex items-center gap-2">
          {header.rightActions?.filter(a => a.visibleOn !== "mobile").map((a, i) => (
            <HeaderActionRenderer key={`right-${i}`} action={a} />
          ))}
        </div>
      </div>

      {/* mobile menu dropdown */}
      {isMenuOpen && <MobileMenu actions={mobileActions} onClose={() => setMenuOpen(false)} />}
    </header>
  )
}