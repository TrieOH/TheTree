import { cn } from '@/shared/lib/utils';
import type { HeaderAction } from '../model/header.types'
import HeaderActionRenderer from './HeaderActionRenderer';
import { useEffect } from 'react';

export default function MobileMenu({
  actions,
  onClose,
}: {
  actions: { left: HeaderAction[]; center: HeaderAction[]; right: HeaderAction[] }
  onClose: () => void
}) {

  useEffect(() => {
    // lock body scroll
    const prev = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => { document.body.style.overflow = prev }
  }, [])

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  const hasAny =
    (actions.left && actions.left.length) ||
    (actions.center && actions.center.length) ||
    (actions.right && actions.right.length)

  if (!hasAny) return null

  const hasLCDivider = actions.left.length > 0 && actions.center.length > 0;
  const hasCRivider = actions.center.length > 0 && actions.right.length > 0;
  return (
    <div
      role="menu"
      aria-label="Mobile navigation"
      className={cn(
        "absolute w-full flex flex-col md:hidden justify-center items-center gap-4",
        "text-lg border-b border-b-border py-4 bg-background/80 backdrop-blur-sm",
        "md:hidden top-full left-0"
      )}
    >
      {actions.left.map((a, i) => (
        <div key={i} onClick={onClose}>
          <HeaderActionRenderer action={a} />
        </div>
      ))}
      {hasLCDivider && <div className="w-11/12 h-px bg-border self-center" />}
      {actions.center.map((a, i) => (
        <div key={i} onClick={onClose}>
          <HeaderActionRenderer action={a} />
        </div>
      ))}
      {hasCRivider && <div className="w-11/12 h-px bg-border self-center" />}
      {actions.right.map((a, i) => (
        <div key={i} onClick={onClose}>
          <HeaderActionRenderer action={a} />
        </div>
      ))}
    </div>
  )
}
