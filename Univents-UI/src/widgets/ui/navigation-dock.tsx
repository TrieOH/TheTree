import { useState, useRef } from 'react';
import { motion, useMotionValue, useSpring, useTransform } from 'motion/react';
import { Home, Search, Settings, Bell, type LucideIcon } from 'lucide-react';
import { Tooltip, TooltipTrigger, TooltipContent } from '@/shared/ui/shadcn/tooltip';
import { cn } from '@/shared/lib/utils';

export type NavItemType = {
  id: string;
  label: string;
  icon: LucideIcon;
  href?: string;
};

export type NavigationDockProps = {
  items?: NavItemType[];
  activeId?: string;
  onNavigate?: (id: string) => void;
  className?: string;
};

const defaultItems: NavItemType[] = [
  { id: 'home', label: 'Home', icon: Home },
  { id: 'search', label: 'Search', icon: Search },
  { id: 'notifications', label: 'Notifications', icon: Bell },
  { id: 'settings', label: 'Settings', icon: Settings },
];

const DesktopNavItem = ({
  item,
  isActive,
  onClick,
  mouseX,
}: {
  item: NavItemType;
  isActive: boolean;
  onClick: () => void;
  mouseX: ReturnType<typeof useMotionValue<number>>;
}) => {
  const ref = useRef<HTMLButtonElement>(null);
  const Icon = item.icon;

  const distance = useTransform(mouseX, (val) => {
    const b = ref.current?.getBoundingClientRect() ?? { x: 0, width: 0 };
    return val - (b.x + b.width / 2);
  });

  const sizeRaw = useTransform(distance, [-130, 0, 130], [40, 56, 40]);
  const size = useSpring(sizeRaw, { mass: 0.08, stiffness: 200, damping: 18 });

  const iconSizeRaw = useTransform(distance, [-130, 0, 130], [16, 22, 16]);
  const iconSize = useSpring(iconSizeRaw, { mass: 0.08, stiffness: 200, damping: 18 });

  return (
    <Tooltip>
      <TooltipTrigger
        render={
          <motion.button
            ref={ref}
            onClick={onClick}
            style={{ width: size, height: size }}
            className={cn(
              'relative flex items-center justify-center rounded-full outline-none transition-colors duration-200',
              isActive
                ? 'bg-primary text-primary-foreground shadow-lg shadow-primary/30'
                : 'bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground',
            )}
            aria-label={item.label}
            aria-current={isActive ? 'page' : undefined}
            whileTap={{ scale: 0.88 }}
          >
            {isActive && (
              <motion.div
                layoutId="active-ring"
                className="absolute inset-0 rounded-full ring-2 ring-primary ring-offset-2 ring-offset-background"
                initial={false}
                transition={{ type: 'spring', stiffness: 400, damping: 30 }}
              />
            )}

            <motion.div style={{ width: iconSize, height: iconSize }} className="flex items-center justify-center">
              <Icon
                style={{ width: '100%', height: '100%' }}
                strokeWidth={isActive ? 2.5 : 2}
              />
            </motion.div>
          </motion.button>
        }
      />
      <TooltipContent side="top" sideOffset={8}>
        {item.label}
      </TooltipContent>
    </Tooltip>
  );
};

const MobileNavItem = ({
  item,
  isActive,
  onClick,
}: {
  item: NavItemType;
  isActive: boolean;
  onClick: () => void;
}) => {
  const Icon = item.icon;

  return (
    <button
      onClick={onClick}
      className={cn(
        'relative flex flex-col items-center justify-center flex-1 py-3 gap-1.5 outline-none',
        isActive ? 'text-primary' : 'text-muted-foreground',
      )}
      aria-label={item.label}
      aria-current={isActive ? 'page' : undefined}
    >
      {isActive && (
        <motion.div
          layoutId="mobile-indicator"
          className="absolute top-1 left-1/2 h-1 w-8 -translate-x-1/2 rounded-full bg-primary"
          initial={{ opacity: 0, y: -4 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -4 }}
          transition={{ type: 'spring', stiffness: 500, damping: 35 }}
        />
      )}

      <motion.div
        animate={isActive ? { scale: 1.05, y: -2 } : { scale: 1, y: 0 }}
        transition={{ type: 'spring', stiffness: 400, damping: 25 }}
      >
        <Icon size={22} strokeWidth={isActive ? 2.4 : 2} />
      </motion.div>

      <motion.span
        animate={isActive ? { opacity: 1, y: 0 } : { opacity: 0.6, y: 0 }}
        transition={{ duration: 0.15 }}
        className="text-[10px] font-medium tracking-tight"
      >
        {item.label}
      </motion.span>
    </button>
  );
};

export function NavigationDock({
  items = defaultItems,
  activeId: controlledActiveId,
  onNavigate,
  className,
}: NavigationDockProps) {
  const [internalActiveId, setInternalActiveId] = useState(items[0]?.id ?? '');
  const activeId = controlledActiveId ?? internalActiveId;

  const handleNavigate = (id: string) => {
    if (controlledActiveId === undefined) setInternalActiveId(id);
    onNavigate?.(id);
  };

  const navItems = items.slice(0, 6);
  const mouseX = useMotionValue(0);

  return (
    <>
      {/* Desktop */}
      <nav
        role="navigation"
        className={cn('fixed bottom-8 left-1/2 -translate-x-1/2 z-50 hidden md:flex', className)}
        onMouseMove={(e) => mouseX.set(e.clientX)}
        onMouseLeave={() => mouseX.set(0)}
      >
        <motion.div
          initial={{ y: 20, opacity: 0, filter: 'blur(10px)' }}
          animate={{ y: 0, opacity: 1, filter: 'blur(0px)' }}
          transition={{ type: 'spring', stiffness: 260, damping: 24, delay: 0.05 }}
          className="flex items-center gap-2 px-3 py-3 rounded-full bg-background/80 backdrop-blur-2xl border border-border/60 shadow-lg shadow-black/5"
        >
          {navItems.map((item) => (
            <DesktopNavItem
              key={item.id}
              item={item}
              isActive={activeId === item.id}
              onClick={() => handleNavigate(item.id)}
              mouseX={mouseX}
            />
          ))}
        </motion.div>
      </nav>

      {/* Mobile */}
      <nav
        role="navigation"
        className={cn('fixed bottom-0 left-0 right-0 z-50 md:hidden', className)}
      >
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ type: 'spring', stiffness: 260, damping: 24 }}
          className="flex items-stretch justify-around px-4 pt-2 pb-safe bg-background/90 backdrop-blur-2xl border-t border-border/40"
        >
          {navItems.map((item) => (
            <MobileNavItem
              key={item.id}
              item={item}
              isActive={activeId === item.id}
              onClick={() => handleNavigate(item.id)}
            />
          ))}
        </motion.div>
        <div className="h-safe-area-inset-bottom bg-background/90 backdrop-blur-2xl" />
      </nav>
    </>
  );
}

export { DesktopNavItem, MobileNavItem };