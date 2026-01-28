import { cn } from '@/shared/lib/utils';
import { HeaderAction } from '../model/header.types'
import HeaderActionRenderer from './HeaderActionRenderer';

export default function MobileMenu({ actions = [], onClose }: { actions: HeaderAction[]; onClose: () => void }) {
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
      {actions.map((a, i) => (
        <div key={i} onClick={onClose}>
          <HeaderActionRenderer action={a} />
        </div>
      ))}
    </div>
  )
}
