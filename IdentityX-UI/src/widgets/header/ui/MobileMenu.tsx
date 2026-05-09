import { cn } from '@/shared/lib/utils';
import type { HeaderAction } from '../model/header.types'
import HeaderActionRenderer from './HeaderActionRenderer';

export default function MobileMenu({
  actions,
  onClose,
}: {
  actions: { left: HeaderAction[]; center: HeaderAction[]; right: HeaderAction[] }
  onClose: () => void
}) {

  const hasAny = actions.left?.length || actions.center?.length || actions.right?.length
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
      {actions.left.map(a => (
        <button type="button" key={`leftM-${a.id}`} onClick={onClose}>
          <HeaderActionRenderer action={a} />
        </button>
      ))}
      {hasLCDivider && <div className="w-11/12 h-px bg-border self-center" />}
      {actions.center.map(a => (
        <button type="button" key={`centerM-${a.id}`} onClick={onClose}>
          <HeaderActionRenderer action={a} />
        </button>
      ))}
      {hasCRivider && <div className="w-11/12 h-px bg-border self-center" />}
      {actions.right.map(a => (
        <button type="button" key={`rightM-${a.id}`} onClick={onClose}>
          <HeaderActionRenderer action={a} />
        </button>
      ))}
    </div>
  )
}
