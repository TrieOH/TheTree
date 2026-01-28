import { useMatches } from "@tanstack/react-router"
import { useEffect, useState } from "react"
import { headerRegistry } from "../model/header.registry"
import MobileMenu from "./MobileMenu"
import { Menu, X } from "lucide-react"
import HeaderActionRenderer from "./HeaderActionRenderer"
import { cn } from "@/shared/lib/utils"
import type { HeaderAction, VisibleOn } from "../model/header.types"

function visibilityClass(visibleOn?: VisibleOn) {
  if (!visibleOn || visibleOn === 'fixed') return ''
  if (visibleOn === 'desktop') return 'hidden md:flex'
  return 'flex md:hidden'
}

export default function Header() {
  const matches = useMatches()
  const headerVariant = [...matches]
    .reverse()
    .find(m => m.staticData.components.header)?.staticData.components.header

  const [isMenuOpen, setMenuOpen] = useState(false)

  if (!headerVariant || headerVariant === 'none') return null

  const header = headerRegistry[headerVariant]
  
  const shouldCollapse = (a: HeaderAction) => {
    if (header.disableMobileMenu) return false // doesn't collapse
    if (a.collapseToMenu === false) return false // doesn't collapse
    if (a.visibleOn === 'fixed') return false // doesn't collapse
    return true // if desktop-only or mobile-only collapse
  }

  const left = header.leftActions ?? []
  const center = header.centerActions ?? []
  const right = header.rightActions ?? []

  const mobileActions = [
    ...left.filter(a => shouldCollapse(a)),
    ...center.filter(a => shouldCollapse(a)),
    ...right.filter(a => shouldCollapse(a)),
  ]

  const showHamburger = !!mobileActions.length && !header.disableMobileMenu

  useEffect(() => {
    if (!showHamburger) setMenuOpen(false)
  }, [showHamburger])
  
  return (
    <header className="sticky top-0 w-full z-10">
      <div 
        className={cn(
          "flex justify-between items-center border-b border-b-border px-6 py-2",
          "bg-background/80 backdrop-blur-sm select-none min-h-16"
        )}
      >
        <div className="flex items-center gap-1">
          {showHamburger && (
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
          )}

          {/* left actions (e.g., back) */}
          <div className="flex items-center gap-2">
            {left.map((a, i) => (
              <div className={visibilityClass(a.visibleOn)}>
                <HeaderActionRenderer key={`left-${i}`} action={a} />
              </div>
            ))}
            {header.title && header.titlePosition === 'left' && (
              <h2 className="ml-2 text-2xl font-semibold hidden md:flex">{header.title}</h2>
            )}
          </div>

          <div className="md:hidden">
            {header.title && header.titlePosition === 'left' && (
              <h2 className="text-lg font-semibold text-foreground">
                {header.title}
              </h2>
            )}
          </div>
        </div>

        {/* center desktop nav */}
        <nav className="hidden md:flex gap-6 items-center">
          {center.map((a, i) => (
            <div className={visibilityClass(a.visibleOn)}>
              <HeaderActionRenderer key={`center-${i}`} action={a} />
            </div>
          ))}
        </nav>

        {/* right actions */}
        <div className="flex items-center gap-2">
          {right.map((a, i) => (
            <div className={visibilityClass(a.visibleOn)}>
              <HeaderActionRenderer key={`right-${i}`} action={a} />
            </div>
          ))}
        </div>
      </div>

      {/* mobile menu dropdown */}
      {isMenuOpen && showHamburger && (
        <MobileMenu
          actions={{
            left: left.filter(a => shouldCollapse(a)),
            center: center.filter(a => shouldCollapse(a)),
            right: right.filter(a => shouldCollapse(a)),
          }}
          onClose={() => setMenuOpen(false)}
        />
      )}
    </header>
  )
}